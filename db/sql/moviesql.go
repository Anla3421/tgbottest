package sql

import (
	"fmt"
)

//撈資料放DB
func Moviesql(Rank int, IdResql string, TitleResql string) {
	results, err := mysqlConn.Query("INSERT INTO movielist (ID,idre,moviename) VALUES(?,?,?) ON DUPLICATE KEY UPDATE idre=?,moviename=?",
		Rank, IdResql, TitleResql, IdResql, TitleResql)
	if err != nil {
		panic(err)

	}
	defer results.Close()
	//fmt.Println("Suceessfully Inserted into users table")
}

//撈DB資料出來
func Moviesqlget(Rank int) Response {
	results, err := mysqlConn.Query("Select idre,moviename FROM movielist where ID = ?", Rank)
	if err != nil {
		fmt.Println(err.Error())
		return Response{}
	}
	var moviedb Response
	for results.Next() {

		err = results.Scan(&moviedb.Idre, &moviedb.Moviename)
		if err != nil {
			fmt.Println(err.Error())
			return Response{}
		}
		defer results.Close()

	}
	return moviedb
}
