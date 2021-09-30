package main

type EnvConfig struct {
	PostgresUser     string `default:"user" split_words:"true"`
	PostgresPassword string `split_words:"true"`
	PostgresHost     string `default:"localhost" split_words:"true"`
	PostgresPort     int    `default:"5432" split_words:"true"`
	PostgresDb       string `default:"organization_service" split_words:"true"`
	MigrationsPath   string `default:"file://pkg/database/migrations" split_words:"true"`
	DatabaseURL      string `default:"" split_words:"true"`
	Port             int    `default:"8082" split_words:"true"`	// env var used by heroku to assign the port for the deployed app
}
