package main

import (
	sql "database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Table struct {
	TableName string
}

type Student struct {
	Code    string
	Name    string
	Program string
}

var tableInfo = &Table{}

func main() {
	os.Remove("sqlite-database.db")
	log.Println("Creating sqlite-database.db...")

	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	runHTTPServer()
}

// runHTTPServer starts the HTTP server and sets up routes for creating tables, inserting students, and displaying students.
//
// No parameters.
// No return values.
func runHTTPServer() {

	http.HandleFunc("/createTable/", handleCreateTable)

	http.HandleFunc("/insertStudent/", handleInsertStudent)

	http.HandleFunc("/displayStudents/", handleDisplayStudents)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)

}

// handleCreateTable handles the HTTP request to create a table in the SQLite database.
//
// Parameters:
//
//	w (http.ResponseWriter): the HTTP response writer
//	r (*http.Request): the HTTP request
//
// Return values:
//
//	None
func handleCreateTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, tableInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	if err != nil {
		log.Fatal(err.Error())
	}

	defer sqliteDatabase.Close() // Defer Closing the database

	err = createTable(tableInfo.TableName, sqliteDatabase)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// handleInsertStudent handles the HTTP request to insert a student into the SQLite database.
//
// Parameters:
//
//	w (http.ResponseWriter): the HTTP response writer
//	r (*http.Request): the HTTP request
//
// Return values:
//
//	None
func handleInsertStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	studentInfo := &Student{}

	err = json.Unmarshal(body, studentInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	if err != nil {
		log.Fatal(err.Error())
	}

	defer sqliteDatabase.Close() // Defer Closing the database

	err = insertStudent(sqliteDatabase, studentInfo.Code, studentInfo.Name, studentInfo.Program)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// handleDisplayStudents handles the HTTP request to display students from the SQLite database.
//
// Parameters:
//
//	w (http.ResponseWriter): the HTTP response writer
//	r (*http.Request): the HTTP request
//
// Return values:
//
//	None
func handleDisplayStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	if err != nil {
		log.Fatal(err.Error())
	}

	defer sqliteDatabase.Close() // Defer Closing the database

	student := displayStudents(sqliteDatabase, tableInfo.TableName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)

}

// createTable creates a table in the SQLite database.
//
// Parameters:
//
//	tabname (string): the name of the table to be created
//	db (*sql.DB): the SQLite database connection
//
// Return values:
//
//	error: any error that occurred during table creation
func createTable(tabname string, db *sql.DB) error {
	var createStudentTableSQL string
	var b strings.Builder

	b.WriteString(`CREATE TABLE `)
	b.WriteString(tabname)
	b.WriteString(` (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );`)
	createStudentTableSQL = b.String()

	log.Println("Create " + tabname + " table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println(tabname, "table created")

	return err
}

// insertStudent inserts a student record into the database.
//
// Parameters:
// - db: A pointer to a sql.DB object representing the database connection.
// - code: A string representing the student's code.
// - name: A string representing the student's name.
// - program: A string representing the student's program.
//
// Returns:
// - error: An error object if there was an error inserting the student record.
func insertStudent(db *sql.DB, code string, name string, program string) error {
	log.Println("Inserting student record ...")

	var b strings.Builder

	b.WriteString(`INSERT INTO `)
	b.WriteString(tableInfo.TableName)
	b.WriteString(`(code, name, program) VALUES (?, ?, ?)`)
	insertStudentSQL := b.String()

	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return err
}

// displayStudents retrieves and displays student records from the database.
//
// Parameters:
//
//	db (*sql.DB): the SQLite database connection
//	tabname (string): the name of the table to retrieve student records from
//
// Return values:
//
//	student (Student): a Student object containing the retrieved student record
func displayStudents(db *sql.DB, tabname string) (student Student) {
	row, err := db.Query("SELECT * FROM " + tabname + " ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	for row.Next() { // Iterate and fetch the records from result cursor
		var id string
		row.Scan(&id, &student.Code, &student.Name, &student.Program)
	}

	return student
}
