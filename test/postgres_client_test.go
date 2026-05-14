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
