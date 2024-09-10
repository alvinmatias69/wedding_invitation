package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alvinmatias69/wedding_invitation/internal/controller"
	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/alvinmatias69/wedding_invitation/internal/handler"
	"github.com/alvinmatias69/wedding_invitation/internal/repository"
	"github.com/alvinmatias69/wedding_invitation/internal/resource"
	"github.com/alvinmatias69/wedding_invitation/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	defaultConfigPath   = "./config.toml"
	dbUserVar           = "POSTGRES_USER"
	dbNameVar           = "POSTGRES_DB"
	dbPassVar           = "POSTGRES_PASSWORD"
	dbHostVar           = "DB_HOST"
	connectionStringfmt = "postgres://%s:%s@%s:5432/%s"
)

func main() {
	var config entities.Config
	_, err := toml.DecodeFile(defaultConfigPath, &config)
	if err != nil {
		log.Fatalf("Error while decoding config file: %v\n", err)
	}

	pgpool, err := pgxpool.New(context.Background(), getDbConnectionString())
	if err != nil {
		log.Fatalf("Error while creating db connection: %v\n", err)
	}

	var (
		jwtResource        = resource.NewJwtResource(config)
		exifResource       = resource.NewExifResource(config)
		tokenRepository    = repository.New(pgpool)
		messageRepository  = repository.NewMessageRepository(pgpool)
		controllerInstance = controller.New(config, jwtResource, exifResource, tokenRepository, messageRepository)
		handlerInstance    = handler.New(config, controllerInstance)
		serverInstance     = server.New(config, handlerInstance)
	)

	serverInstance.Start()
}

func getDbConnectionString() string {
	var (
		dbUser = os.Getenv(dbUserVar)
		dbName = os.Getenv(dbNameVar)
		dbPass = os.Getenv(dbPassVar)
		dbHost = os.Getenv(dbHostVar)
	)

	return fmt.Sprintf(connectionStringfmt, dbUser, dbPass, dbHost, dbName)
}
