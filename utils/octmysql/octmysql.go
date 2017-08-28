package octmysql

import (
	"database/sql"
	"fmt"
	"octlink/rstore/utils/octlog"
)

// OctMysql for oct basic mysql structure
type OctMysql struct {
	conn *sql.DB
}

var logger *octlog.LogConfig

// InitLog for log config init
func InitLog(level int) {
	logger = octlog.InitLogConfig("octmysql.log", level)
}

const (
	// DbName Basic Db name
	DbName = "dbrstore"

	// DbUser Default Database User
	DbUser = "root"

	// DbPassword Default Database Password
	DbPassword = "123456"

	// DbServer Default Database Server
	DbServer = "127.0.0.1"

	// DbPort Default Database Connection Port
	DbPort = 3306
)

// Open to open mysql db connection
func (octmysql *OctMysql) Open() error {

	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		DbUser, DbPassword, DbServer, DbPort, DbName)

	db, err := sql.Open("mysql", ds)
	if err != nil {
		fmt.Println("failed to open database:", err.Error())
		return err
	}

	octmysql.conn = db

	octlog.Debug("open db OK %p\n", octmysql.conn)

	return err
}

// Close to close database connection
func (octmysql *OctMysql) Close() {
	if octmysql.conn != nil {
		octmysql.conn.Close()
	}
}

// Query to do db queries
func (octmysql *OctMysql) Query(query string, args ...interface{}) (*sql.Rows, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}

	return octmysql.conn.Query(query, args...)
}

// Prepare to db db query preparation
func (octmysql *OctMysql) Prepare(query string) (*sql.Stmt, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}

	return octmysql.conn.Prepare(query)
}

// QueryRow to qeury my row
func (octmysql *OctMysql) QueryRow(query string, args ...interface{}) *sql.Row {
	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil
		}
	}
	return octmysql.conn.QueryRow(query, args...)
}

// Begin to start query
func (octmysql *OctMysql) Begin() (*sql.Tx, error) {
	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}
	return octmysql.conn.Begin()
}

// Exec to do db execution
func (octmysql *OctMysql) Exec(query string, args ...interface{}) (sql.Result, error) {

	if octmysql.conn == nil {
		err := octmysql.Open()
		if err != nil {
			return nil, err
		}
	}

	return octmysql.conn.Exec(query, args...)
}

// Count to return db count
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
