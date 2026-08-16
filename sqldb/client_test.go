package sqldb

import (
	"database/sql"
	"os"
	"strconv"
	"testing"
)

// getTestDNS constructs SqlDNS from environment variables or sensible defaults.
func getTestDNS() SqlDNS {
	host := os.Getenv("PGHOST")
	if host == "" {
		host = "localhost"
	}

	port := 5432
	if portStr := os.Getenv("PGPORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	user := os.Getenv("PGUSER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("PGPASSWORD")
	if password == "" {
		password = "postgres"
	}

	database := os.Getenv("PGDATABASE")
	if database == "" {
		database = "postgres"
	}

	sslmode := os.Getenv("PGSSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	return SqlDNS{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslmode,
	}
}

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func setupUsersTable(t *testing.T, client *Client) {
	t.Helper()

	// Ensure users table exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE
	);`

	_, err := client.Execute(createTableSQL)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	// Clean table contents before each test for isolation
	_, err = client.Execute("TRUNCATE TABLE users RESTART IDENTITY")
	if err != nil {
		t.Fatalf("failed to truncate users table: %v", err)
	}
}

func TestPostgresClient_UsersTable(t *testing.T) {
	dns := getTestDNS()
	client := NewPostgresClient(dns)
	if client == nil {
		t.Fatal("expected non-nil Client from NewPostgresClient")
	}
	defer client.Close()

	if client.GetDB() == nil {
		t.Fatal("expected non-nil *sql.DB from GetDB()")
	}

	setupUsersTable(t, client)

	t.Run("Execute Insert and Query", func(t *testing.T) {
		// Test Execute (INSERT)
		res, err := client.Execute("INSERT INTO users (name, email) VALUES ($1, $2)", "Alice", "alice@example.com")
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			t.Fatalf("failed to get rows affected: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("expected 1 row affected, got %d", rowsAffected)
		}

		// Insert second user
		_, err = client.Execute("INSERT INTO users (name, email) VALUES ($1, $2)", "Bob", "bob@example.com")
		if err != nil {
			t.Fatalf("failed to insert second user: %v", err)
		}

		// Test Query raw rows
		rows, err := client.Query("SELECT id, name, email FROM users ORDER BY id ASC")
		if err != nil {
			t.Fatalf("failed to query users: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}
			count++
			if count == 1 && (u.Name != "Alice" || u.Email != "alice@example.com") {
				t.Errorf("unexpected first user row: %+v", u)
			}
			if count == 2 && (u.Name != "Bob" || u.Email != "bob@example.com") {
				t.Errorf("unexpected second user row: %+v", u)
			}
		}

		if count != 2 {
			t.Errorf("expected 2 user rows, got %d", count)
		}
	})

	t.Run("Execute Update and Delete", func(t *testing.T) {
		setupUsersTable(t, client)

		_, err := client.Execute("INSERT INTO users (name, email) VALUES ($1, $2)", "Charlie", "charlie@example.com")
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		// Update
		res, err := client.Execute("UPDATE users SET name = $1 WHERE email = $2", "Charlie Updated", "charlie@example.com")
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}
		affected, _ := res.RowsAffected()
		if affected != 1 {
			t.Errorf("expected 1 row updated, got %d", affected)
		}

		// Query typed row via QueryRow helper
		user, err := QueryRow[User](client, "SELECT id, name, email FROM users WHERE email = $1", "charlie@example.com")
		if err != nil {
			t.Fatalf("failed to query updated row: %v", err)
		}
		if user.Name != "Charlie Updated" {
			t.Errorf("expected updated name 'Charlie Updated', got '%s'", user.Name)
		}

		// Delete
		res, err = client.Execute("DELETE FROM users WHERE email = $1", "charlie@example.com")
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}
		affected, _ = res.RowsAffected()
		if affected != 1 {
			t.Errorf("expected 1 row deleted, got %d", affected)
		}

		// Verify deletion
		_, err = QueryRow[User](client, "SELECT id, name, email FROM users WHERE email = $1", "charlie@example.com")
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows after deletion, got %v", err)
		}
	})

	t.Run("QueryRows Integration Helper", func(t *testing.T) {
		setupUsersTable(t, client)

		_, err := client.Execute("INSERT INTO users (name, email) VALUES ($1, $2), ($3, $4)", "Dave", "dave@example.com", "Eve", "eve@example.com")
		if err != nil {
			t.Fatalf("failed to insert users: %v", err)
		}

		users, err := QueryRows[User](client, "SELECT id, name, email FROM users ORDER BY id ASC")
		if err != nil {
			t.Fatalf("failed to query rows using mapper helper: %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("expected 2 users mapped, got %d", len(users))
		}
		if users[0].Name != "Dave" || users[1].Name != "Eve" {
			t.Errorf("unexpected mapped users: %+v", users)
		}
	})
}
