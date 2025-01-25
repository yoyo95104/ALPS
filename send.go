package main
import(
	"database/sql"
	"log"
) 

func send(grade int , email string , name string , pass string , channel chan bool){
	db , err := sql.Open("sqlite3" , "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt , err := db.Prepare("INSERT INTO userdata (name , email , grade , pass) VALUES (? , ? , ? , ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_ , err = stmt.Exec(name , email , grade , pass)
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT id FROM your_table ORDER BY id DESC LIMIT 1;")
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	id++
	stmt, err = db.Prepare("INSERT INTO userdata (id) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	channel <- true
}