package main

import (
	"database/sql"
	"fmt"
	"log"
	"sales-api/api"
	"sales-api/config"
	db "sales-api/db/sqlc"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	if err = run(conf); err != nil {
		log.Fatalf("an error occurred: %v", err)
	}
}

func run(conf config.Config) error {
	conn, err := sql.Open(conf.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			conf.DBUser,
			conf.DBPassword,
			conf.DBHost,
			conf.DBPort,
			conf.DBName,
		),
	)
	if err != nil {
		return fmt.Errorf("unable to connect to db: %v", err)
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		return fmt.Errorf("database did not respond: %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, conf)
	if err != nil {
		return fmt.Errorf("unable to create api server: %v", err)
	}

	err = server.Start(conf.BindAddr)
	if err != nil {
		return fmt.Errorf("unable to start server: %v", err)
	}

	return nil
}
