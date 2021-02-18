package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlConn *sql.DB

func init() {
	fmt.Println("MySQL initial")
	CreateConn()
}

func CreateConn() {
	// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/testdb") //@home
	db, err := sql.Open("mysql", "root:adminstrator@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	mysqlConn = db
	fmt.Println("Suceessfully Connected to MySQL")
}
func Weathersql(ID string) Response {
	results, err := mysqlConn.Query("Select text FROM Weather where ID = ?", ID)
	if err != nil {
		fmt.Println(err.Error())
		return Response{}
	}
	var weadb Response
	for results.Next() {

		err = results.Scan(&weadb.Text)
		if err != nil {
			fmt.Println(err.Error())
			return Response{}
		}
		defer results.Close()

	}
	return weadb
}
