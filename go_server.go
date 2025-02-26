package main

// I think main here defines the code as an executable, as opposed to a library.

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq" // The funky underscore at the beginning allows me to import a package without using it. In this case we just want to "register" the postgres driver. Compiler will complain if it's missing.
)

// The handler must be a method. Simplest way looks to be to attach to an empty struct.
type GetStudentsHandler struct {
	db *sql.DB
}

func (h GetStudentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// syntax above: "h" is like self in python class definition. The stuff in parenthesis attaches the function to the struct as a method. In go, structs can have methods.
	// "w" is an instance of the http.ResponseWriter type
	// Likewise with r but in this case is a pointer, as required by the package.
	fmt.Fprintf(w, "Hello!")
}

func main() {

	db, _ := sql.Open("postgres", "user=postgress password=dbzsuper dbname=school_records sslmode=disable")
	// first argument is the driver name, hence postgres.
	// Normally I wouldn't write a password straight into the code like this but secrets management feels
	// out of scope for the assignment.

	defer db.Close()
	handler := GetStudentsHandler{db: db}
	http.Handle("/GetStudents", handler)
	http.ListenAndServe(":8080", nil)
}
