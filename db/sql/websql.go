package sql

import (
	"fmt"

	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func Websql(ID string) Response {
	results, err := mysqlConn.Query("Select name,text FROM page3 where ID = ?", ID)
	if err != nil {
		fmt.Println(err.Error())
		return Response{}
	}
	var pagedb Response
	for results.Next() {

		err = results.Scan(&pagedb.Name, &pagedb.Text)
		if err != nil {
			fmt.Println(err.Error())
			return Response{}
		}
		defer results.Close()

	}
	return pagedb
}
