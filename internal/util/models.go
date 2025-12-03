package util

type RestAPIConfig struct {
	KafkaConnect KafkaConnectConfig `yaml:"kafkaConnect"`
}

type KafkaConnectConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
