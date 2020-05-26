package config

const (
	ServiceName = "Cartographer"
	Version     = "1.0.0"

	// Service port address
	envServerPort = "SERVER_PORT"

	// Logger config
	envLoggerFormat = "LOGGER_FORMAT"
	envLoggerLevel  = "LOGGER_LEVEL"

	// Neo4J config
	envNeo4JAddress     = "NEO4J_ADDRESS"
	envNeo4JUsername    = "NEO4J_USERNAME"
	envNeo4JPassword    = "NEO4J_PASSWORD"
	envNeo4JMaxConnPool = "NEO4J_MAX_CONN_POOL"
	envNeo4JEncrypted   = "NEO4J_ENCRYPTED"
	envNeo4JLogEnabled  = "NEO4J_LOG_ENABLED"
	envNeo4JLogLevel    = "NEO4J_LOG_LEVEL"

	Limit  = 25
	Offset = 0
)
