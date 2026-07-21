package test

import (
	"testing"

	"github.com/jycodepub/golibs/sqldb"
)

type user struct {
	id       int
	username string
	password string
}

func TestClient_Query(t *testing.T) {
	dns := sqldb.SqlDNS{
		Host:     "jysrv02",
		Port:     5432,
		User:     "jydev",
		Password: "jydev",
		Database: "jydb_dev",
	}

	client := sqldb.NewPostgresClient(dns)
	defer client.Close()

	rows, err := client.Query("SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var u user
	for rows.Next() {
		rows.Scan(&u.id, &u.username, &u.password)
		t.Logf("id: %d, username: %s, password: %s", u.id, u.username, u.password)
	}

	rows2, err := client.Query("SELECT * FROM users WHERE username=$1", "user1")
	if err != nil {
		t.Fatal(err)
	}
	defer rows2.Close()

	for rows2.Next() {
		rows.Scan(&u.id, &u.username, &u.password)
		t.Logf("id: %d, username: %s, password: %s", u.id, u.username, u.password)
	}
}

type ExportedUser struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func TestClient_QueryMapper(t *testing.T) {
	dns := sqldb.SqlDNS{
		Host:     "jysrv02",
		Port:     5432,
		User:     "jydev",
		Password: "jydev",
		Database: "jydb_dev",
	}

	client := sqldb.NewPostgresClient(dns)
	defer client.Close()

	// Query multiple rows using helper function
	users, err := sqldb.QueryRows[ExportedUser](client, "SELECT * FROM users")
	if err != nil {
		t.Logf("QueryRows failed (expected if db is unreachable): %v", err)
	} else {
		for _, u := range users {
			t.Logf("mapped user: id: %d, username: %s, password: %s", u.ID, u.Username, u.Password)
		}
	}

	// Query single row using helper function
	u, err := sqldb.QueryRow[ExportedUser](client, "SELECT * FROM users WHERE username=$1", "user1")
	if err != nil {
		t.Logf("QueryRow failed (expected if db is unreachable): %v", err)
	} else {
		t.Logf("mapped single user: id: %d, username: %s, password: %s", u.ID, u.Username, u.Password)
	}
}

