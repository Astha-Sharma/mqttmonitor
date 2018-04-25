package mysql

type SqlClient interface {
}

func NewClient(config *Config) SqlClient {
	return &mySqlClient{
		config: config,
	}
}

type Config struct {
	DBConnection string
}

type mySqlClient struct {
	config *Config
}
