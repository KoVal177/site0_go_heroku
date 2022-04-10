package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "golang",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// установка данных
	//insert, err := db.Query("INSERT INTO `users` (`name`, `age`) VALUES ('Bob', 35)")
	//if err != nil {
	//	panic(err)
	//}
	//defer insert.Close()

	// Выборка данных
	res, _ := db.Query("SELECT `name`, `age` FROM `users`")
	defer res.Close()
	var user User
	for res.Next() {
		err = res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("User: %s with age %d", user.Name, user.Age))

	}

}
