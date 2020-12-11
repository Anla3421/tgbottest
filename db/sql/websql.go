package sql

import (
	"fmt"
)

type Response struct {
	Name      string `json:"name"`
	Text      string `json:"text"`
	Idre      string `json:"idre"`
	Moviename string `json:"moviename"`
	ID        int    `json:"id"`
	Who       string `json:"who"`
	Drink     string `json:"drink"`
	Sugar     string `json:"sugar"`
	Ice       string `json:"ice"`
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
