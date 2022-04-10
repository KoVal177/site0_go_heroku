package main

import (
	"database/sql"
	"fmt"
	//	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	//"reflect"
	//"strconv"
)

/*
var Cfg = mysql.Config{
	User:                 "valiok",
	Passwd:               "zaebalinah123456",
	Net:                  "tcp",
	Addr:                 "www.db4free.net:3306",
	DBName:               "golang_valiok",
	AllowNativePasswords: true,
}
*/

/*
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)
*/

const (
	host     = "ec2-3-230-122-20.compute-1.amazonaws.com"
	port     = 5432
	user     = "azcivxdazpjqxx"
	password = "a4a263001481fdc1b78d3f6dba95951922238714f80692905aa11901d86035a2"
	dbname   = "da1uletl8ugi64"
)

var psqlInfo string = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s", //sslmode=disable",
	host, port, user, password, dbname)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles"))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title") //название переменной может быть любое. просто берет из формы по названию в name
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все поля заполнены")
	} else {

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		res, _ := db.Query(fmt.Sprintf("SELECT max(id) FROM articles"))
		var id int
		res.Next()
		err = res.Scan(&id)
		id++

		insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (id, title, anons, full_text) VALUES (%d, '%s', '%s', '%s')", id, title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func handleRequest() {
	// статику логично прописывать в начале. вторая переменная убирает префекс из пути в html-файле и говорит где искать далее
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create/", create)
	http.HandleFunc("/save_article/", save_article)
	http.ListenAndServe(":90", nil)
}

func main() {
	//fmt.Println("starting...")
	handleRequest()

}
