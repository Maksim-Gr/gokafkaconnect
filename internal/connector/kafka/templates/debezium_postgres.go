package templates

import "maps"

var defaultDebeziumPostgresConfig = map[string]string{
	"connector.class":      "io.debezium.connector.postgresql.PostgresConnector",
	"tasks.max":            "1",
	"database.hostname":    "",
	"database.port":        "5432",
	"database.user":        "",
	"database.password":    "",
	"database.dbname":      "",
	"database.server.name": "",
	"plugin.name":          "pgoutput",
	"table.include.list":   "",
	"topic.prefix":         "",
}

func GetDebeziumPostgresConnectorTemplate() map[string]string {
	configCopy := make(map[string]string)
	maps.Copy(configCopy, defaultDebeziumPostgresConfig)
	return configCopy
}

// DebeziumPostgresRequiredFields returns a list of mandatory fields for the Debezium PostgreSQL connector.
func DebeziumPostgresRequiredFields() []string {
	return []string{
		"database.hostname",
		"database.user",
		"database.password",
		"database.dbname",
		"database.server.name",
		"table.include.list",
	}
}
