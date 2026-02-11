package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// MySqlConnection represents a MySQL database connection
type MySqlConnection struct {
	db *sql.DB
}

// NewMySqlConnection creates a new MySQL connection using environment variables
// Required environment variables:
// - MYSQL_USER: MySQL username
// - MYSQL_PASSWORD: MySQL password
// - MYSQL_HOST: MySQL host (default: localhost)
// - MYSQL_PORT: MySQL port (default: 3306)
// - MYSQL_DATABASE: MySQL database name
func NewMySqlConnection() (*MySqlConnection, error) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	// Set defaults if not provided
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3306"
	}

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &MySqlConnection{db: db}, nil
}

// ValidateConnection runs a simple SELECT 1 query to validate the database connection
func (c *MySqlConnection) ValidateConnection() error {
	var result int
	err := c.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("failed to validate connection: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected result from validation query: got %d, expected 1", result)
	}

	return nil
}

// Close closes the database connection
func (c *MySqlConnection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// GetDB returns the underlying *sql.DB instance for advanced usage
func (c *MySqlConnection) GetDB() *sql.DB {
	return c.db
}
