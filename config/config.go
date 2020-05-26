package config

import (
	"log"

	"github.com/dynastymasra/cartographer/infrastructure/provider"

	"github.com/dynastymasra/cookbook"
	"github.com/spf13/viper"
)

type Config struct {
	serverAddress string
	logger        LoggerConfig
	neo4j         provider.Neo4J
}

var config *Config

func Load() {
	viper.SetDefault(envServerPort, "8080")

	viper.AutomaticEnv()

	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	viper.AddConfigPath("../../../../")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	config = &Config{
		serverAddress: getString(envServerPort),
		logger: LoggerConfig{
			format: getString(envLoggerFormat),
			level:  getString(envLoggerLevel),
		},
		neo4j: provider.Neo4J{
			Address:     getString(envNeo4JAddress),
			Username:    getString(envNeo4JUsername),
			Password:    getString(envNeo4JPassword),
			MaxConnPool: getInt(envNeo4JMaxConnPool),
			Encrypted:   getBool(envNeo4JEncrypted),
			LogEnabled:  getBool(envNeo4JLogEnabled),
			LogLevel:    getInt(envNeo4JLogLevel),
		},
	}
}

func ServerAddress() string {
	return config.serverAddress
}

func Logger() LoggerConfig {
	return config.logger
}

func Neo4J() provider.Neo4J {
	return config.neo4j
}

func getString(key string) string {
	value, err := cookbook.StringEnv(key)
	if err != nil {
		log.Fatalf("%v env key is not set", key)
	}
	return value
}

func getInt(key string) int {
	value, err := cookbook.IntEnv(key)
	if err != nil {
		log.Fatalf("%v env key is not set", key)
	}
	return value
}

func getBool(key string) bool {
	value, err := cookbook.BoolEnv(key)
	if err != nil {
		log.Fatalf("%v env key is not set", key)
	}
	return value
}
