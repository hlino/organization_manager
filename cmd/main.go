package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"organization_manager/pkg/api"
	"organization_manager/pkg/database"
)

func main() {
	var envConfig EnvConfig
	err := envconfig.Process("", &envConfig)
	if err != nil {
		log.Fatalf("error loading environment variables: %v", err.Error())
	}

	err = database.Initialize(envConfig.PostgresUser, envConfig.PostgresPassword, envConfig.PostgresHost,
		envConfig.PostgresDb, envConfig.PostgresPort)
	if err != nil {
		log.Fatalf("error initializing the database: %v", err.Error())
	}

	// running migrations on startup
	err = database.RunMigrations(envConfig.MigrationsPath, envConfig.PostgresUser, envConfig.PostgresPassword,
		envConfig.PostgresHost, envConfig.PostgresDb, envConfig.PostgresPort)
	if err != nil {
		log.Fatalf("error running database migration on startup: %v", err.Error())
	}

	server := api.Server{}
	err = server.Initialize()
	if err != nil {
		log.Fatalf("error initializing database: %v", err.Error())
	}

	server.Run()
}

