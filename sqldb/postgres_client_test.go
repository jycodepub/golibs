package sqldb

import "testing"

func TestClient_Query(t *testing.T) {
	dns := SqlDNS{
		Host:     "jysrv02",
		Port:     5432,
		User:     "jydev",
		Password: "jydev",
		Database: "jydb_dev",
	}

	client := NewPostgresClient(dns)
	defer client.Close()

	rows, err := client.Query("SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username string
		var password string
		err = rows.Scan(&id, &username, &password)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("id: %d, username: %s, password: %s", id, username, password)
	}
}
