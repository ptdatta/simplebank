package main

import (
	"context"
	"log"

	"github.com/ptdatta/simplebank/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ptdatta/simplebank/api"
	db "github.com/ptdatta/simplebank/db/sqlc"
)



func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

    store := db.NewStore(conn)
	server,err := api.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}