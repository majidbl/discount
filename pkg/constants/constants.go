package constants

const (
	SERVER_PORT  = "SERVER_PORT"
	GRPC_PORT    = "GRPC_PORT"
	METRICS_PORT = "METRICS_PORT"

	NATS_URL       = "NATS_URL"
	CLUSTER_ID     = "CLUSTER_ID"
	NATS_CLIENT_ID = "NATS_CLIENT_ID"

	REDIS_URL      = "REDIS_URL"
	REDIS_PASSWORD = "REDIS_PASSWORD"

	POSTGRES_HOST     = "POSTGRES_HOST"
	POSTGRES_PORT     = "POSTGRES_PORT"
	POSTGRES_USER     = "POSTGRES_USER"
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
	POSTGRES_DB       = "POSTGRES_DB"
	POSTGRES_SSL      = "POSTGRES_SSL"

	JAEGER_HOST = "JAEGER_HOST"
)

const (
	GrpcPort          = "GRPC_PORT"
	HttpPort          = "HTTP_PORT"
	ConfigPath        = "CONFIG_PATH"
	KafkaBrokers      = "KAFKA_BROKERS"
	JaegerHostPort    = "JAEGER_HOST"
	RedisAddr         = "REDIS_ADDR"
	MongoDbURI        = "MONGO_URI"
	PostgresqlHost    = "POSTGRES_HOST"
	PostgresqlPort    = "POSTGRES_PORT"
	ElasticUrl        = "ELASTIC_URL"
	TCP               = "tcp"
	MIGRATIONS_DB_URL = "MIGRATIONS_DB_URL"

	ReaderServicePort = "READER_SERVICE"

	Yaml          = "yaml"
	Redis         = "redis"
	Kafka         = "kafka"
	Postgres      = "postgres"
	MongoDB       = "mongo"
	ElasticSearch = "elasticSearch"

	GRPC     = "GRPC"
	SIZE     = "SIZE"
	URI      = "URI"
	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"

	Topic     = "topic"
	Partition = "partition"
	Message   = "message"
	WorkerID  = "workerID"
	Headers   = "headers"
	Offset    = "offset"
	Time      = "time"

	Page   = "page"
	Size   = "size"
	Search = "search"
	ID     = "id"

	EsAll = "$all"

	Validate        = "validate"
	FieldValidation = "field validation"
	RequiredHeaders = "required header"
	Base64          = "base64"
	Unmarshal       = "unmarshal"
	Uuid            = "uuid"
	Cookie          = "cookie"
	Token           = "token"
	Bcrypt          = "bcrypt"
	SQLState        = "sqlstate"

	MongoAggregateID = "aggregateID"

	MongoProjection   = "(MongoDB Projection)"
	ElasticProjection = "(Elastic Projection)"

	EventType   = "(EventType)"
	AggregateID = "(AggregateID)"
	Version     = "(Version)"
	TimeStamp   = "(TimeStamp)"
	Metadata    = "(Metadata)"

	MessageSize = "MessageSize"

	BankAccountIndex = "aggregateID"
	BankAccountId    = "BankAccountId"

	KafkaHeaders = "kafkaHeaders"

	Tcp = "tcp"
)
