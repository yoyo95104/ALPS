package main

import(
	// "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
)
func verify(email string , pass string , channel chan bool){
	db , err := sql.Open("sqlite3" , "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows , err := db.Query("SELECT * FROM userdata WHERE email = ? AND pass = ?" , email , pass)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if exists := rows.Next(); exists {
		channel <- true
	}else{
		channel <- false
	}
}