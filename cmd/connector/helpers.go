package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sort"
	"strings"
	"time"

	"github.com/Maksim-Gr/kkon/internal/connector"
	"github.com/Maksim-Gr/kkon/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// configMapFromFile extracts a connector config map from a connector JSON file,
// accepting either the wrapped {"name":..,"config":{..}} form or a flat config
// map. For the wrapped form the top-level name is folded into the config, since
// the validate endpoint requires it. It returns nil if neither form parses, in
// which case validation is skipped.
func configMapFromFile(b []byte) map[string]string {
	var wrapper struct {
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
	}
	if err := json.Unmarshal(b, &wrapper); err == nil && len(wrapper.Config) > 0 {
		if wrapper.Name != "" && wrapper.Config["name"] == "" {
			cfg := make(map[string]string, len(wrapper.Config)+1)
			maps.Copy(cfg, wrapper.Config)
			cfg["name"] = wrapper.Name
			return cfg
		}
		return wrapper.Config
	}
	var flat map[string]string
	if err := json.Unmarshal(b, &flat); err == nil {
		return flat
	}
	return nil
}

// argOrEmpty returns the first positional arg, or "" if none was provided.
func argOrEmpty(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// errNonInteractiveJSON is returned when --output json is combined with a
// command invocation that would need to prompt.
var errNonInteractiveJSON = errors.New("--output json requires non-interactive usage (pass the connector name and --yes)")

// printActionJSON emits a single action result object to stdout for --output json mode.
func printActionJSON(name, action string) {
	b, _ := json.Marshal(map[string]string{"name": name, "action": action, "result": "ok"})
	fmt.Println(string(b))
}

// lifecycleRunE returns a RunE shared by pause and resume, which differ only
// in the API call and the verb used in messages.
func lifecycleRunE(verb string, op func(ctx context.Context, client *connector.Client, name string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if util.IsJSONOutput(cmd) && argOrEmpty(args) == "" {
			return errNonInteractiveJSON
		}

		client, err := util.NewKafkaConnectClient()
		if err != nil {
			return err
		}

		name, err := util.ResolveConnectorName(cmd.Context(), client, argOrEmpty(args))
		if err != nil {
			return err
		}

		if util.IsDryRun(cmd) {
			color.Yellow("[dry-run] Would %s connector %s\n", verb, name)
			return nil
		}

		if err := op(cmd.Context(), client, name); err != nil {
			return fmt.Errorf("failed to %s %s: %w", verb, name, err)
		}
		if util.IsJSONOutput(cmd) {
			printActionJSON(name, verb)
			return nil
		}
		color.Green("%s requested for %s\n", strings.ToUpper(verb[:1])+verb[1:], name)
		return nil
	}
}

// validateConfigOrConfirm runs server-side validation on a connector config.
// It returns true if the caller should proceed with submission. When the
// connector.class is unknown or validation can't run, it proceeds without
// blocking. When validation reports errors, it prints them and asks the user
// whether to submit anyway.
func validateConfigOrConfirm(ctx context.Context, client *connector.Client, cfg map[string]string) bool {
	class := cfg["connector.class"]
	if class == "" {
		return true
	}

	res, err := client.ValidateConnectorConfig(ctx, class, cfg)
	if err != nil {
		color.Yellow("Could not validate config: %v\n", err)
		return true
	}
	if res.ErrorCount == 0 {
		return true
	}

	color.Red("\nValidation found %d error(s):\n", res.ErrorCount)
	for _, c := range res.Configs {
		for _, e := range c.Value.Errors {
			color.Red("  %s: %s\n", c.Value.Name, e)
		}
	}

	var proceed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Submit anyway?",
		Default: false,
	}, &proceed); err != nil {
		return false
	}
	return proceed
}

// validateConfigStrict runs server-side validation and returns an error listing
// the failing fields. It is the non-interactive counterpart of
// validateConfigOrConfirm for contexts where prompting is not possible
// (--output json). Like the interactive variant, it proceeds when the
// connector.class is unknown or validation itself can't run.
func validateConfigStrict(ctx context.Context, client *connector.Client, cfg map[string]string) error {
	class := cfg["connector.class"]
	if class == "" {
		return nil
	}

	res, err := client.ValidateConnectorConfig(ctx, class, cfg)
	if err != nil || res.ErrorCount == 0 {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "validation found %d error(s):", res.ErrorCount)
	for _, c := range res.Configs {
		for _, e := range c.Value.Errors {
			fmt.Fprintf(&b, "\n  %s: %s", c.Value.Name, e)
		}
	}
	return errors.New(b.String())
}

// editConnectorConfig fetches the named connector's live config, lets the user
// edit fields interactively, shows a diff, validates, and applies the change.
// Applying restarts the connector, so it then verifies the connector comes back
// RUNNING and offers to revert to the original config if it does not.
// User cancels return util.ErrCanceled; with dryRun the diff is shown but the
// change is not applied.
func editConnectorConfig(ctx context.Context, client *connector.Client, selected string, dryRun bool) error {
	connectorConfig, err := client.GetConnectorConfigJSON(ctx, selected)
	if err != nil {
		return fmt.Errorf("failed to get connector config: %w", err)
	}

	// Snapshot the original config for diff display and potential revert.
	original := make(map[string]string, len(connectorConfig))
	for k, v := range connectorConfig {
		original[k] = v
	}

	for {
		pretty, err := util.ToPrettyJSON(connectorConfig)
		if err != nil {
			return fmt.Errorf("failed to format config: %w", err)
		}
		color.Cyan("\n Current config for %s:\n", selected)
		fmt.Println(pretty)

		fields := make([]string, 0, len(connectorConfig))
		for k := range connectorConfig {
			fields = append(fields, k)
		}
		sort.Strings(fields)

		var fieldToChange string
		if err := survey.AskOne(&survey.Select{
			Message: "Which field do you want to change?",
			Options: fields,
		}, &fieldToChange); err != nil {
			return util.ErrCanceled
		}

		var newValue string
		if err := survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("New value for %s (current: %v):", fieldToChange, connectorConfig[fieldToChange]),
		}, &newValue); err != nil {
			return util.ErrCanceled
		}
		connectorConfig[fieldToChange] = newValue

		var more bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Change another field?",
			Default: false,
		}, &more); err != nil {
			return util.ErrCanceled
		}
		if !more {
			break
		}
	}

	// Compute and display changed fields.
	var changedKeys []string
	for k, newV := range connectorConfig {
		if oldV, exists := original[k]; exists && oldV != newV {
			changedKeys = append(changedKeys, k)
		}
	}

	if len(changedKeys) == 0 {
		color.Yellow("No changes made\n")
		return nil
	}

	sort.Strings(changedKeys)
	maxKeyLen := 0
	for _, k := range changedKeys {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	color.Cyan("\nChanges:")
	for _, k := range changedKeys {
		fmt.Printf("  %-*s  %s  →  %s\n", maxKeyLen, k, original[k], connectorConfig[k])
	}

	if dryRun {
		color.Yellow("[dry-run] Would apply %d change(s) to %s\n", len(changedKeys), selected)
		return nil
	}

	var confirm bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Apply this config to " + selected + "?",
		Default: true,
	}, &confirm); err != nil || !confirm {
		return util.ErrCanceled
	}

	if !validateConfigOrConfirm(ctx, client, connectorConfig) {
		return util.ErrCanceled
	}

	if err := client.UpdateConnectorConfig(ctx, selected, connectorConfig); err != nil {
		return fmt.Errorf("failed to update connector: %w", err)
	}
	color.Green("Connector %s updated; restarting with new config...\n", selected)

	verifyConnectorOrRevert(ctx, client, selected, original)
	return nil
}

// verifyConnectorOrRevert polls the connector's status after a config change.
// Updating a config restarts the connector and its tasks, so this waits for it
// to settle and reports whether it is RUNNING. If it is not healthy, it offers
// to revert to the previous config.
func verifyConnectorOrRevert(ctx context.Context, client *connector.Client, name string, original map[string]string) {
	status, healthy := waitForConnectorRunning(ctx, client, name, 5, 2*time.Second)

	if healthy {
		color.Green("Connector %s is up and running: %s\n", name, util.ColorState(status.Connector.State))
		return
	}

	if status.Name == "" {
		color.Yellow("Could not verify connector status\n")
		return
	}

	color.Red("Connector %s did not come back RUNNING.\n", name)
	color.Red("  connector: %s\n", util.ColorState(status.Connector.State))
	for _, t := range status.Tasks {
		fmt.Printf("  task %d: %s\n", t.ID, util.ColorState(t.State))
	}

	var revert bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Connector %s is not RUNNING. Revert to previous config?", name),
		Default: true,
	}, &revert); err != nil || !revert {
		color.Yellow("Keeping new config in place\n")
		return
	}

	if err := client.UpdateConnectorConfig(ctx, name, original); err != nil {
		color.Red("Failed to revert connector: %v\n", err)
		return
	}
	color.Green("Reverted %s to previous config\n", name)
}

// waitForConnectorRunning polls the connector's status up to attempts times,
// sleeping delay between tries, and returns once it is healthy or attempts are
// exhausted. The returned bool reports whether it became healthy; the returned
// status is the last one observed (zero-valued, with empty Name, if every poll
// errored).
func waitForConnectorRunning(ctx context.Context, client *connector.Client, name string, attempts int, delay time.Duration) (connector.Status, bool) {
	var status connector.Status
	for i := 0; i < attempts; i++ {
		time.Sleep(delay)
		s, err := client.GetConnectorStatus(ctx, name)
		if err != nil {
			continue
		}
		status = s
		if connectorHealthy(status) {
			return status, true
		}
	}
	return status, false
}

// printConnectorStatus prints the state of a connector and its tasks.
// For FAILED tasks it fetches and shows the first 3 lines of the error trace.
func printConnectorStatus(ctx context.Context, client *connector.Client, name string, status connector.Status) {
	fmt.Printf("  %s  %s\n", name, util.ColorState(status.Connector.State))
	for _, t := range status.Tasks {
		fmt.Printf("    Task %d: %s\n", t.ID, util.ColorState(t.State))
		if t.State == "FAILED" {
			ts, err := client.GetConnectorTaskStatus(ctx, name, t.ID)
			if err == nil && ts.Trace != "" {
				lines := strings.Split(strings.TrimRight(ts.Trace, "\n"), "\n")
				shown := lines
				if len(lines) > 3 {
					shown = lines[:3]
				}
				for _, line := range shown {
					color.Yellow("      %s\n", line)
				}
				if len(lines) > 3 {
					color.Yellow("      ...\n")
					fmt.Printf("      To see full trace run: kkon task get --connector %s --id %d\n", name, t.ID)
				}
			}
		}
	}
}

// connectorHealthy reports whether the connector and all its tasks are RUNNING.
func connectorHealthy(status connector.Status) bool {
	if status.Connector.State != "RUNNING" {
		return false
	}
	for _, t := range status.Tasks {
		if t.State != "RUNNING" {
			return false
		}
	}
	return true
}
