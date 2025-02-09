package config

import (
	"fmt"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     string
}

func (c Config) GetMongoDBURI() string {
	if c.DBUser != "" && c.DBPassword != "" {
		return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
			c.DBUser, c.DBPassword, c.DBHost, c.DBName)
	}
	return fmt.Sprintf("mongodb+srv://%s/%s?retryWrites=true&w=majority", c.DBHost, c.DBName)
}
