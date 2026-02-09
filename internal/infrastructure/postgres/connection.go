package db

import (
	"fmt"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
)

func GetConnectionString(cnf *config.DBConfig) string {
	//return "user=ecommerce_user password=ecommerce_password dbname=ecommerce_db host=localhost port=5434 sslmode=disable"
	connString := fmt.Sprintf("user=%s password=%s  host=%s port=%d dbname=%s",
		cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Name)
	if !cnf.EnableSSLMode {
		connString += " sslmode=disable"
	}
	return connString
}

func NewConnection(cnf *config.Config) (*sqlx.DB, error) {
	dbSource := GetConnectionString(cnf.DB)
	db, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}
