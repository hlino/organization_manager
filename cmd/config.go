package main

type EnvConfig struct {
	MigrationsPath   string `default:"file://pkg/database/migrations" split_words:"true"`
	DatabaseURL      string `default:"postgres://user:Password123!@localhost:5432/organization_service?sslmode=disable" split_words:"true"`
	Port             int    `default:"8082" split_words:"true"`	// env var used by heroku to assign the port for the deployed app
}
