package sqldb

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sql.DB
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) GetDB() *sql.DB {
	return c.db
}

func NewPostgresClient(dns SqlDNS) *Client {
	db, err := sql.Open("postgres", dns.getPostgresDNS())
	if err != nil {
		panic(err)
	}
	return &Client{db: db}
}

func (c *Client) Execute(statement string, args ...any) (sql.Result, error) {
	stmt, err := c.db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (c *Client) Query(statement string, args ...any) (*sql.Rows, error) {
	stmt, err := c.db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(args...)
}
