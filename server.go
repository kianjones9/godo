// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var existsSQL = "SELECT name FROM sqlite_master WHERE type='table' AND name='tasks';"
var createSQL = `CREATE TABLE tasks ("taskID" integer NOT NULL PRIMARY KEY AUTOINCREMENT, "taskName" TEXT);`

func handler(w http.ResponseWriter, r *http.Request) {

	db, _ := sql.Open("sqlite3", "./tasks.db")
	defer db.Close()

	fmt.Println(db, existsSQL)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(body))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":6000", nil))
}

func execSQL(db *sql.DB, query string) sql.Result {

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := statement.Exec()
	if err != nil {
		log.Fatal(err.Error())
	}

	return rows
}

// We are passing db reference connection from main to our method with other parameters
func insertStudent(db *sql.DB, code string, name string, program string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Student: ", code, " ", name, " ", program)
	}
}
