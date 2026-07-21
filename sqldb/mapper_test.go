package sqldb

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"
)

// Define mock driver structures to simulate database queries.
type mockDriver struct{}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{}, nil
}

type mockConn struct{}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{query: query}, nil
}
func (c *mockConn) Close() error               { return nil }
func (c *mockConn) Begin() (driver.Tx, error) { return nil, nil }

type mockStmt struct {
	query string
}

func (s *mockStmt) Close() error                                    { return nil }
func (s *mockStmt) NumInput() int                                   { return -1 }
func (s *mockStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, nil }
func (s *mockStmt) Query(args []driver.Value) (driver.Rows, error) {
	if mockQueryHandler != nil {
		return mockQueryHandler(s.query, args)
	}
	return &mockRows{}, nil
}

type mockRows struct {
	cols []string
	data [][]driver.Value
	idx  int
}

func (r *mockRows) Columns() []string {
	return r.cols
}

func (r *mockRows) Close() error {
	return nil
}

func (r *mockRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.data) {
		return io.EOF
	}
	row := r.data[r.idx]
	for i := range dest {
		dest[i] = row[i]
	}
	r.idx++
	return nil
}

var mockQueryHandler func(query string, args []driver.Value) (driver.Rows, error)

func init() {
	sql.Register("mock_db_driver", &mockDriver{})
}

// Test structs for mapping verification.
type SimpleUser struct {
	ID       int    `db:"user_id"`
	Username string `db:"username"`
	Email    string // maps to "email" automatically
}

type EmbedAudit struct {
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type UserWithEmbed struct {
	ID int `db:"id"`
	EmbedAudit
	Role string
}

type UserWithEmbedPtr struct {
	ID int `db:"id"`
	*EmbedAudit
	Role string
}

type UserWithUnexported struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	secret string // unexported field, should be ignored
}

func TestMapRows_Simple(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			cols: []string{"user_id", "username", "email"},
			data: [][]driver.Value{
				{int64(1), "alice", "alice@example.com"},
				{int64(2), "bob", "bob@example.com"},
			},
		}, nil
	}

	users, err := QueryRows[SimpleUser](client, "SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if users[0].ID != 1 || users[0].Username != "alice" || users[0].Email != "alice@example.com" {
		t.Errorf("unexpected user 0: %+v", users[0])
	}
	if users[1].ID != 2 || users[1].Username != "bob" || users[1].Email != "bob@example.com" {
		t.Errorf("unexpected user 1: %+v", users[1])
	}
}

func TestMapRow_Single(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			cols: []string{"user_id", "username", "email"},
			data: [][]driver.Value{
				{int64(1), "alice", "alice@example.com"},
			},
		}, nil
	}

	user, err := QueryRow[SimpleUser](client, "SELECT * FROM users LIMIT 1")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 || user.Username != "alice" || user.Email != "alice@example.com" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestMapRow_NoRows(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			cols: []string{"user_id", "username", "email"},
			data: [][]driver.Value{},
		}, nil
	}

	_, err = QueryRow[SimpleUser](client, "SELECT * FROM users WHERE id = -1")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestMapRows_Embedded(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			cols: []string{"id", "created_at", "updated_at", "role"},
			data: [][]driver.Value{
				{int64(42), "2026-01-01", "2026-01-02", "admin"},
			},
		}, nil
	}

	users, err := QueryRows[UserWithEmbed](client, "SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	u := users[0]
	if u.ID != 42 || u.CreatedAt != "2026-01-01" || u.UpdatedAt != "2026-01-02" || u.Role != "admin" {
		t.Errorf("unexpected user: %+v", u)
	}
}

func TestMapRows_EmbeddedPtr(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			cols: []string{"id", "created_at", "updated_at", "role"},
			data: [][]driver.Value{
				{int64(42), "2026-01-01", "2026-01-02", "admin"},
			},
		}, nil
	}

	users, err := QueryRows[UserWithEmbedPtr](client, "SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	u := users[0]
	if u.ID != 42 || u.Role != "admin" {
		t.Errorf("unexpected outer values: %+v", u)
	}
	if u.EmbedAudit == nil {
		t.Fatal("expected EmbedAudit to be initialized")
	}
	if u.CreatedAt != "2026-01-01" || u.UpdatedAt != "2026-01-02" {
		t.Errorf("unexpected embedded values: %+v", u.EmbedAudit)
	}
}

func TestMapRows_UnmatchedColumnsAndUnexportedFields(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			// "extra_col" doesn't exist in struct.
			// "secret" is unexported, so should not be mapped.
			cols: []string{"id", "name", "secret", "extra_col"},
			data: [][]driver.Value{
				{int64(10), "john", "my-secret-value", "something-extra"},
			},
		}, nil
	}

	users, err := QueryRows[UserWithUnexported](client, "SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	u := users[0]
	if u.ID != 10 || u.Name != "john" {
		t.Errorf("unexpected user values: %+v", u)
	}
	if u.secret != "" {
		t.Errorf("secret field was mapped, but should have been ignored: %s", u.secret)
	}
}

func TestMapRows_PointerType(t *testing.T) {
	db, err := sql.Open("mock_db_driver", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	mockQueryHandler = func(query string, args []driver.Value) (driver.Rows, error) {
		return &mockRows{
			cols: []string{"user_id", "username", "email"},
			data: [][]driver.Value{
				{int64(1), "alice", "alice@example.com"},
			},
		}, nil
	}

	users, err := QueryRows[*SimpleUser](client, "SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	u := users[0]
	if u == nil {
		t.Fatal("expected non-nil pointer")
	}
	if u.ID != 1 || u.Username != "alice" || u.Email != "alice@example.com" {
		t.Errorf("unexpected user: %+v", u)
	}
}
