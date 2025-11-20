package storage

type StorageConfig struct {
	Username, Password, Host, Port, Database string
	MaxConnections, ConnectAttempts          int
}
