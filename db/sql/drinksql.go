package sql

import (
	"fmt"
)

func Drinksql(Drinkid int, Who string, Arguments string, Sugar string, Ice string) {
	results, err := mysqlConn.Query("INSERT INTO drink (ID,who,drink,sugar,ice) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE who=?,drink=?,sugar=?,ice=?",
		Drinkid, Who, Arguments, Sugar, Ice, Who, Arguments, Sugar, Ice)
	if err != nil {
		panic(err)

	}
	defer results.Close()
}

func Drinksqlget(ID int) Response {
	results, err := mysqlConn.Query("SELECT * FROM drink where ID=?", ID)
	if err != nil {
		fmt.Println(err.Error())
		return Response{}
	}
	var drinkdb Response
	for results.Next() {

		err = results.Scan(&drinkdb.ID, &drinkdb.Who, &drinkdb.Drink, &drinkdb.Sugar, &drinkdb.Ice)
		if err != nil {
			fmt.Println(err.Error())
			return Response{}
		}
		defer results.Close()

	}
	return drinkdb
}

func Drinksqltruncate() {
	results, err := mysqlConn.Query("TRUNCATE table drink")
	if err != nil {
		panic(err.Error())

	}
	defer results.Close()
}
