package templates

import "maps"

var defaultS3SinkConfig = map[string]string{
	"connector.class": "io.confluent.connect.s3.S3SinkConnector",
	"tasks.max":       "1",
	"s3.region":       "",
	"s3.bucket.name":  "",
	"s3.part.size":    "5242880",
	"topics":          "",
	"flush.size":      "1000",
	"storage.class":   "io.confluent.connect.s3.storage.S3Storage",
	"format.class":    "io.confluent.connect.s3.format.json.JsonFormat",
	"topics.dir":      "topics",
}

func GetS3SinkConnectorTemplate() map[string]string {
	configCopy := make(map[string]string)
	maps.Copy(configCopy, defaultS3SinkConfig)
	return configCopy
}

// S3SinkRequiredFields returns a list of mandatory fields for the S3 Sink connector.
func S3SinkRequiredFields() []string {
	return []string{
		"s3.region",
		"s3.bucket.name",
		"topics",
	}
}
