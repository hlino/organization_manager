package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(dbUser, dbPassword, dbHost, dbName string, dbPort int) error {
	gormDb := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		dbHost, dbUser, dbName, dbPassword, dbPort)
	log.Infof("Attempting to connect to database: %v", gormDb)
	var err error
	DB, err = gorm.Open(postgres.Open(gormDb), &gorm.Config{})
	if err != nil {
		log.Errorf("Unable to connect to postgres: %v", err)
		return err
	}
	return nil
}

func RunMigrations(migrationsPath, dbUser, dbPassword, dbHost, dbName string, dbPort int) error {
	dbUri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	migrations, err := migrate.New(migrationsPath, dbUri)
	if err != nil {
		return err
	}
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func InitializeTest() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	DB, err = gorm.Open(postgres.New(
		postgres.Config{Conn: db}),
		&gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}
