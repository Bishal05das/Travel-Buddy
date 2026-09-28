package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetConnectionString(cnf *config.DBConfig) string {
	//return "user=ecommerce_user password=ecommerce_password dbname=ecommerce_db host=localhost port=5434 sslmode=disable"
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",
		quoteDSN(cnf.User), quoteDSN(cnf.Password), quoteDSN(cnf.Host), cnf.Port, quoteDSN(cnf.Name))
	if !cnf.EnableSSLMode {
		connString += " sslmode=disable"
	}
	return connString
}

// quoteDSN quotes a key/value connection-string value so passwords
// containing spaces, quotes or backslashes are passed through intact.
func quoteDSN(v string) string {
	v = strings.ReplaceAll(v, `\`, `\\`)
	v = strings.ReplaceAll(v, `'`, `\'`)
	return "'" + v + "'"
}

func NewConnection(cnf *config.Config) (*sqlx.DB, error) {
	dbSource := GetConnectionString(cnf.DB)
	var db *sqlx.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("postgres", dbSource)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			fmt.Println("Database connected!")
			break
		}

		fmt.Println("Waiting for DB, retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}
