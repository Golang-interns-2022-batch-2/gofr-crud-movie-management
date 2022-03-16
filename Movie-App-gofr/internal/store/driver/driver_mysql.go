package driver

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Initializing the Driver to register the mysql driver
)

const (
	dbDriver string = "mysql"
	dbUser   string = "root"
	dbPass   string = "password"
)

func DBConn(dbConfigName string) (*sql.DB, error) {
	dbName := dbConfigName
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		return nil, err
	}

	return db, nil
}
