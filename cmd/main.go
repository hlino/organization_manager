package main

import (
	"fmt"
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

	var dbURI string
	if envConfig.DatabaseURL == "" {
		dbURI = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			envConfig.PostgresUser, envConfig.PostgresPassword, envConfig.PostgresHost, envConfig.PostgresPort,
			envConfig.PostgresDb)
	} else {
		dbURI = envConfig.DatabaseURL
	}
	err = database.Initialize(dbURI)
	if err != nil {
		log.Fatalf("error initializing the database: %v", err.Error())
	}

	// running migrations on startup
	err = database.RunMigrations(envConfig.MigrationsPath, dbURI)
	if err != nil {
		log.Fatalf("error running database migration on startup: %v", err.Error())
	}

	server := api.Server{}
	err = server.Initialize()
	if err != nil {
		log.Fatalf("error initializing server: %v", err.Error())
	}

	log.Infof("successfully started server")

	server.Run(envConfig.Port)
}

