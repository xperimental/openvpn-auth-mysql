package main

type ServerConfig struct {
	Host     string
	Username string
	Password string
	Database string
}

func getServerConfig() ServerConfig {
	// modify to fit your server configuration
	return ServerConfig{"127.0.0.1:3306", "username", "password", "database"}
}
