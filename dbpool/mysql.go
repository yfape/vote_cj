package dbpool

import (
	"strings"
	"vote_cj/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Mysqlpool *sqlx.DB

var Articlepool *sqlx.DB

func init() {
	var sbd strings.Builder
	sbd.WriteString(config.Mysql.Username)
	sbd.WriteString(":")
	sbd.WriteString(config.Mysql.Password)
	sbd.WriteString("@tcp(")
	sbd.WriteString(config.Mysql.Host)
	sbd.WriteString(":")
	sbd.WriteString(config.Mysql.Port)
	sbd.WriteString(")/")
	sbd.WriteString(config.Mysql.Database)
	sbd.WriteString("?multiStatements=true")
	database, err := sqlx.Open("mysql", sbd.String())
	if err != nil {
		panic(err)
	}
	Mysqlpool = database
	Mysqlpool.SetMaxIdleConns(config.Mysql.MaxIdle)
	Mysqlpool.SetMaxOpenConns(config.Mysql.MaxOpen)
}
