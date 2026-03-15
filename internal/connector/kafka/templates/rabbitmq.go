package templates

import "maps"

var defaultRabbitMQConfig = map[string]string{
	"connector.class":                       "com.ibm.eventstreams.connect.rabbitmqsource.RabbitMQSourceConnector",
	"rabbitmq.topology.recovery.enabled":    "true",
	"tasks.max":                             "1",
	"rabbitmq.password":                     "",
	"rabbitmq.username":                     "",
	"rabbitmq.queue":                        "",
	"rabbitmq.network.recovery.interval.ms": "10000",
	"rabbitmq.virtual.host":                 "/import",
	"rabbitmq.prefetch.count":               "500",
	"rabbitmq.port":                         "5672",
	"rabbitmq.host":                         "",
	"kafka.topic":                           "",
	"rabbitmq.automatic.recovery.enabled":   "true",
}

func GetRabbitMQConnectorTemplate() map[string]string {
	configCopy := make(map[string]string)
	maps.Copy(configCopy, defaultRabbitMQConfig)
	return configCopy
}

// RabbitMQRequiredFields returns a list of mandatory fields for the RabbitMQ connector.
func RabbitMQRequiredFields() []string {
	return []string{
		"rabbitmq.username",
		"rabbitmq.password",
		"rabbitmq.queue",
		"rabbitmq.host",
		"kafka.topic",
	}
}
