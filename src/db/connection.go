package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Describes the database connection specifications
type connectionSpecs struct {
	username   string
	password   string
	socketPath string
	database   string
}

// Returns a MySQL connection string
func getConnStr(specs *connectionSpecs) string {
	connStr := specs.username + ":" + specs.password +
		"@unix(" + specs.socketPath + ")" +
		"/" + specs.database +
		"?charset=utf8"
	return connStr
}

// Creates a database connection and returns its instance
func Connection() (*sql.DB, error) {

	// Mysql connection specifications
	connConfig := connectionSpecs{
		username:   os.Getenv("MYSQL_USERNAME"),
		password:   os.Getenv("MYSQL_PASSWORD"),
		socketPath: os.Getenv("MYSQL_SOCKETPATH"),
		database:   os.Getenv("MYSQL_DATABASE_NAME"),
	}

	var connStr = getConnStr(&connConfig)
	conn, err := sql.Open("mysql", connStr)
	return conn, err
}
