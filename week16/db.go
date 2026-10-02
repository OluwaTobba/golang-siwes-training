package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const dsn = "postgres://postgres:password@localhost:5432/siwes_db?sslmode=disable"

type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id       SERIAL PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			email    VARCHAR(200) NOT NULL,
			active   BOOLEAN DEFAULT TRUE,
			created  TIMESTAMP DEFAULT NOW()
		)`)
	return err
}

func insertUser(db *sql.DB, u User) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO users(username, email) VALUES($1, $2) RETURNING id",
		u.Username, u.Email,
	).Scan(&id)
	return id, err
}

func getUser(db *sql.DB, id int) (User, error) {
	var u User
	err := db.QueryRow(
		"SELECT id, username, email, active FROM users WHERE id=$1", id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Active)
	return u, err
}

func listUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username, email, active FROM users WHERE active=TRUE")
	if err != nil { return nil, err }
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.Active)
		users = append(users, u)
	}
	return users, rows.Err()
}

// Transaction example
func transferBatch(db *sql.DB, users []User) error {
	tx, err := db.Begin()
	if err != nil { return err }
	defer tx.Rollback() // no-op if committed

	stmt, err := tx.Prepare("INSERT INTO users(username, email) VALUES($1, $2)")
	if err != nil { return err }
	defer stmt.Close()

	for _, u := range users {
		if _, err := stmt.Exec(u.Username, u.Email); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}
	return tx.Commit()
}

func main() {
	db, err := sql.Open("pgx", dsn)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	// Pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil { log.Fatal("ping:", err) }
	fmt.Println("Connected to PostgreSQL")

	createTable(db)

	id, _ := insertUser(db, User{Username: "david_o", Email: "david@example.com"})
	fmt.Println("Inserted user id:", id)

	u, _ := getUser(db, id)
	fmt.Printf("Fetched: %+v\n", u)

	batch := []User{
		{Username: "alice", Email: "alice@example.com"},
		{Username: "bob",   Email: "bob@example.com"},
	}
	transferBatch(db, batch)

	users, _ := listUsers(db)
	for _, u := range users { fmt.Printf("  %d | %s | %s\n", u.ID, u.Username, u.Email) }
}