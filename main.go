package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/microsoft/go-mssqldb"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysqlemployee.database.windows.net"
	dbUser := "root1"
	dbPass := "Nitya12$$"
	dbName := "mysql-employee"
	port := 1433
	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", dbUser, dbPass, dbDriver, dbName)
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		dbDriver, dbUser, dbPass, port, dbName)
	db, err := sql.Open("sqlserver", connectionString)
	//db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("forms/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	ctx := context.Background()
	sqlstat := fmt.Sprintf("SELECT * FROM Employee WHERE id=@nId")
	selDB, err := db.QueryContext(ctx, sqlstat, sql.Named("nId", nId))
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db := dbConn()
	nId := r.URL.Query().Get("id")
	sqlstat := fmt.Sprintf("SELECT * FROM Employee WHERE id=@nId")
	selDB, err := db.QueryContext(ctx, sqlstat, sql.Named("nId", nId))
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("id")
		sqlstat := "INSERT INTO Employee(id, name, city) VALUES(@id, @name, @city)"
		ctx := context.Background()
		insForm, err := db.Prepare(sqlstat)
		if err != nil {
			panic(err.Error())
		}
		row := insForm.QueryRowContext(ctx, sql.Named("id", id), sql.Named("name", name), sql.Named("city", city))
		err = row.Scan()
		if err != nil {
			log.Println("Error is", err)
		}
		log.Println("INSERT: Name: " + name + " | City: " + city + "| Id:" + id)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("id")
		sqlstat := fmt.Sprintf("UPDATE Employee SET name=@name, city=@city WHERE id=@id")
		// insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
		ctx := context.Background()
		_, err := db.ExecContext(ctx, sqlstat, sql.Named("name", name), sql.Named("city", city), sql.Named("id", id))
		if err != nil {
			panic(err.Error())
		}
		log.Println("UPDATE: Name: " + name + " | City: " + city + "| Id: " + id)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	sqlstat := "DELETE FROM Employee WHERE id=@id"
	delForm, err := db.Prepare(sqlstat)
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	row := delForm.QueryRowContext(ctx, sql.Named("id", emp))
	err = row.Scan()
	if err != nil {
		log.Println("Error is", err)
	}
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
