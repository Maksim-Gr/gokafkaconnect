package templates

import "maps"

var defaultJDBCSourceConfig = map[string]string{
	"connector.class":          "io.confluent.connect.jdbc.JdbcSourceConnector",
	"tasks.max":                "1",
	"connection.url":           "",
	"connection.user":          "",
	"connection.password":      "",
	"mode":                     "incrementing",
	"incrementing.column.name": "",
	"topic.prefix":             "",
	"poll.interval.ms":         "5000",
}

// GetJDBCSourceConnectorTemplate returns a copy of the default JDBC Source connector config.
func GetJDBCSourceConnectorTemplate() map[string]string {
	configCopy := make(map[string]string)
	maps.Copy(configCopy, defaultJDBCSourceConfig)
	return configCopy
}

// JDBCSourceRequiredFields returns a list of mandatory fields for the JDBC Source connector.
func JDBCSourceRequiredFields() []string {
	return []string{
		"connection.url",
		"connection.user",
		"connection.password",
		"topic.prefix",
	}
}

var defaultJDBCSinkConfig = map[string]string{
	"connector.class":     "io.confluent.connect.jdbc.JdbcSinkConnector",
	"tasks.max":           "1",
	"connection.url":      "",
	"connection.user":     "",
	"connection.password": "",
	"topics":              "",
	"auto.create":         "false",
	"insert.mode":         "insert",
}

// GetJDBCSinkConnectorTemplate returns a copy of the default JDBC Sink connector config.
func GetJDBCSinkConnectorTemplate() map[string]string {
	configCopy := make(map[string]string)
	maps.Copy(configCopy, defaultJDBCSinkConfig)
	return configCopy
}

// JDBCSinkRequiredFields returns a list of mandatory fields for the JDBC Sink connector.
func JDBCSinkRequiredFields() []string {
	return []string{
		"connection.url",
		"connection.user",
		"connection.password",
		"topics",
	}
}
