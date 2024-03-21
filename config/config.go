package config

import "fmt"

func New() (*Config, error) {
	config, err := NewFromEnv()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *MongoConfig) DSN() string {
	var dsn string
	if c.Username != "" && c.Password != "" {
		dsn = fmt.Sprintf("mongodb://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)
	} else {
		dsn = fmt.Sprintf("mongodb://%s:%d", c.Host, c.Port)
	}

	return dsn
}
