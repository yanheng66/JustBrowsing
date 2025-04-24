package config

import (
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	MongoDB       MongoDBConfig       `yaml:"mongodb"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	Redis         RedisConfig         `yaml:"redis"`
	Kafka         KafkaConfig         `yaml:"kafka"`
	Logging       LoggingConfig       `yaml:"logging"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port     int    `yaml:"port"`
	BasePath string `yaml:"basePath"`
}

// MongoDBConfig represents the MongoDB configuration
type MongoDBConfig struct {
	URI      string        `yaml:"uri"`
	Database string        `yaml:"database"`
	PoolSize int           `yaml:"poolSize"`
	Timeout  time.Duration `yaml:"timeout"`
}

// ElasticsearchConfig represents the Elasticsearch configuration
type ElasticsearchConfig struct {
	Addresses   []string `yaml:"addresses"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	IndexPrefix string   `yaml:"indexPrefix"`
}

// RedisConfig represents the Redis configuration
type RedisConfig struct {
	Address  string        `yaml:"address"`
	Password string        `yaml:"password"`
	DB       int           `yaml:"db"`
	PoolSize int           `yaml:"poolSize"`
	TTL      time.Duration `yaml:"ttl"`
}

// KafkaConfig represents the Kafka configuration
type KafkaConfig struct {
	Brokers []string       `yaml:"brokers"`
	GroupID string         `yaml:"groupId"`
	Topics  KafkaTopicConfig `yaml:"topics"`
}

// KafkaTopicConfig represents the Kafka topic configuration
type KafkaTopicConfig struct {
	Product   string `yaml:"product"`
	Inventory string `yaml:"inventory"`
	Order     string `yaml:"order"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads the configuration from the given file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Replace environment variables in the YAML content
	expandedData := expandEnvVars(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, err
	}

	// Convert timeout from seconds to time.Duration
	config.MongoDB.Timeout = config.MongoDB.Timeout * time.Second
	config.Redis.TTL = config.Redis.TTL * time.Second

	return &config, nil
}

// expandEnvVars replaces ${VARIABLE} or $VARIABLE in the YAML with the corresponding environment variable.
func expandEnvVars(content string) string {
	// Replace ${var} format
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]
		content = strings.ReplaceAll(content, "${"+key+"}", value)
		content = strings.ReplaceAll(content, "$"+key, value)
	}
	return content
}