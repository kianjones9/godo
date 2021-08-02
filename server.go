package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
)

var existsSQL = "SELECT name FROM sqlite_master WHERE type='table' AND name='tasks';"
var createSQL = `CREATE TABLE IF NOT EXISTS tasks ("id" integer NOT NULL PRIMARY KEY AUTOINCREMENT, "name" TEXT);`

// task struct
// constructInsertQueryFromTask
// prepAndExec

type Task struct {
	Name string `json:"name"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	db, _ := sql.Open("sqlite3", "./tasks.db")
	defer db.Close()

	_, err := prepAndExec(db, createSQL)
	if err != nil {
		log.Printf("Error creating table: %v", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
	}

	insertSQL := fmt.Sprintf(`INSERT INTO tasks(name) VALUES ('%v')`, task.Name)
	_, err = prepAndExec(db, insertSQL)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "Success")

	if err != nil {
		log.Printf("Error constructing insertion query: %v", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":6000", nil))
}

func prepAndExec(db *sql.DB, query string) (sql.Result, error) {

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	result, err := statement.Exec()
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return result, nil
}


/*
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

 */