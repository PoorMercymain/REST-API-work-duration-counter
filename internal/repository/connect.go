package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

type db struct {
	conn *pgx.Conn
}

func New() *db {
	db := new(db)

	var (
		hostname     = os.Getenv("REST_API_HOSTNAME")
		databaseName = os.Getenv("REST_API_DATABASE_NAME")
		username     = os.Getenv("REST_API_USERNAME")
		password     = os.Getenv("REST_API_PASSWORD")
		port         = os.Getenv("REST_API_PORT")
	)

	dsn := db.getDsnString(hostname, databaseName, username, password, port)

	log.Println(hostname, databaseName, username, password)

	conn, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Error occured while connecting to database - %v", err)
	}

	db.conn = conn
	return db
}

func (d *db) getDsnString(hostname, databaseName, username, password, port string) string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable ", hostname, databaseName, username, password, port)
}
