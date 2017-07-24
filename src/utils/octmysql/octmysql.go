package octmysql

import (
	"database/sql"
	"fmt"
	"octlink/mirage/src/utils/octlog"

	_ "github.com/go-sql-driver/mysql"
)

type OctMysql struct {
	conn *sql.DB
}

var logger *octlog.LogConfig

func InitLog(level int) {
	logger = octlog.InitLogConfig("octmysql.log", level)
}

const (
	DB_NAME     = "dbmirage"
	DB_USER     = "root"
	DB_PASSWORD = "123456"
	DB_SERVER   = "127.0.0.1"
	DB_PORT     = 3306
)

func (octmysql *OctMysql) Open() error {

	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		DB_USER, DB_PASSWORD, DB_SERVER, DB_PORT, DB_NAME)

	db, err := sql.Open("mysql", ds)
	if err != nil {
		fmt.Println("failed to open database:", err.Error())
		return err
	}

	octmysql.conn = db

	octlog.Debug("open db OK %p\n", octmysql.conn)

	return err
}

func (octmysql *OctMysql) Close() {
	if octmysql.conn != nil {
		octmysql.conn.Close()
	}
}

func (octmysql *OctMysql) Query(query string, args ...interface{}) (*sql.Rows, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}

	return octmysql.conn.Query(query, args...)
}

func (octmysql *OctMysql) Prepare(query string) (*sql.Stmt, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}

	return octmysql.conn.Prepare(query)
}

func (octmysql *OctMysql) QueryRow(query string, args ...interface{}) *sql.Row {
	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil
		}
	}
	return octmysql.conn.QueryRow(query, args...)
}

func (octmysql *OctMysql) Begin() (*sql.Tx, error) {
	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}
	return octmysql.conn.Begin()
}

func (octmysql *OctMysql) Exec(query string, args ...interface{}) (sql.Result, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}

	return octmysql.conn.Exec(query, args...)
}

func (octmysql *OctMysql) Count(table string, cond string, args ...interface{}) (int, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return 0, err
		}
	}

	query := "SELECT COUNT(*) FROM " + table + " " + cond

	var count int
	row := octmysql.conn.QueryRow(query, args...)

	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, err
}
