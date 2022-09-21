package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type dsnBuilder struct {
	hostname     string
	databaseName string
	username     string
	password     string
}

func New() *dsnBuilder {
	dsn := new(dsnBuilder)

	dsn = &dsnBuilder{
		hostname:     "пока что тут ничего...",
		databaseName: "пока что тут ничего...",
		username:     "пока что тут ничего...",
		password:     "пока что тут ничего...",
	}

	return dsn
}

func (dsn *dsnBuilder) getDsnString() string {
	return fmt.Sprintf("host=%s;dbname=%s;user=%s;password=%s", dsn.hostname, dsn.databaseName, dsn.username, dsn.password)
}

func NewPg() *pgx.Conn {
	dsn := New()
	conn, err := pgx.Connect(context.Background(), os.Getenv(dsn.getDsnString()))

	if err != nil {
		fmt.Println("Error occured while connecting to database -", err.Error())
	}

	return conn
}
