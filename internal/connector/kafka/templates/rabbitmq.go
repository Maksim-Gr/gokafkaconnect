package templates

var defaultRedisConfig = map[string]string{
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

func GetRedisConnectorTemplate() map[string]string {
	configCopy := make(map[string]string)
	for k, v := range defaultRedisConfig {
		configCopy[k] = v
	}
	return configCopy
}

// RequiredFields returns a list of mandatory fields for rabbit connector
func RequiredFields() []string {
	return []string{
		"rabbitmq.username",
		"rabbitmq.password",
		"rabbitmq.queue",
		"rabbitmq.host",
		"kafka.topic",
	}
}
