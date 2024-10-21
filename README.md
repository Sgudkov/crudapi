 # crudapi
 
 This is a simple CRUD (Create, Read, Update, Delete) API written in Go that interacts with a SQLite database.
 
 ## Functionality
 
 The API provides the following endpoints:
 
 * **Create Table**: Creates a new table in the SQLite database with the given table name.
 * **Insert Student**: Inserts a new student record into the database with the given student code, name, and program.
 * **Display Students**: Retrieves and displays all student records from the database.
 
 ## API Endpoints
 
 * `POST /createTable/`: Creates a new table in the SQLite database.
 * `POST /insertStudent/`: Inserts a new student record into the database.
 * `GET /displayStudents/`: Retrieves and displays all student records from the database.
 
 ## Database
 
 The API uses a SQLite database to store student records. The database is created automatically when the API starts, and the table is created when the [createTable](cci:1://file:///f:/GO/crudapi/main.go:175:0-208:1) endpoint is called.
 
 ## Code Structure
 
 The code is organized into several functions:
 
 * [main](cci:1://file:///f:/GO/crudapi/main.go:27:0-39:1): The entry point of the program, which sets up the HTTP server and database.
 * [runHTTPServer](cci:1://file:///f:/GO/crudapi/main.go:41:0-56:1): Starts the HTTP server and sets up routes for the API endpoints.
 * [handleCreateTable](cci:1://file:///f:/GO/crudapi/main.go:58:0-98:1), [handleInsertStudent](cci:1://file:///f:/GO/crudapi/main.go:100:0-142:1), [handleDisplayStudents](cci:1://file:///f:/GO/crudapi/main.go:144:0-173:1): Handle the API requests for creating tables, inserting students, and displaying students, respectively.
 * [createTable](cci:1://file:///f:/GO/crudapi/main.go:175:0-208:1), [insertStudent](cci:1://file:///f:/GO/crudapi/main.go:210:0-241:1), [displayStudents](cci:1://file:///f:/GO/crudapi/main.go:243:0-266:1): Perform the actual database operations for creating tables, inserting students, and displaying students, respectively.
 
  

