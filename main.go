package main // main defines executable, as opposed to a library.

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // The funky underscore at the beginning allows me to import a package without using it. In this case we just want to "register" the postgres driver. Compiler will complain if it's missing.
)

type student struct {
	StudentID   int       `json:"student_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// Architecture note: The way this library works the handler functions for each endpoint need to be methods on a struct.
// Despite being identical, I need two separate types because each differentiates Get vs Post.
type GetHelper struct {
	db *sql.DB
}
type PostHelper struct {
	db *sql.DB
}

// The struct should be referenced by pointer, for efficiency reasons. Otherwise we copy the struct.
func (h GetHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// syntax above: "h" is like self in python class definition. The stuff in parenthesis attaches the function to the struct as a method. In go, structs can have methods.
	// "w" is an instance of the http.ResponseWriter type
	// Likewise with r but in this case is a pointer, as required by the package.

	rows, err := h.db.Query("SELECT student_id, first_name, last_name, birth_date FROM students")
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []student // "slice" syntax, simliar to python list.

	for rows.Next() {
		var s student
		rows.Scan(&s.StudentID, &s.FirstName, &s.LastName, &s.DateOfBirth)
		// Assigns s, go passes value, need ampersand to write into s.
		students = append(students, s)
	}

	StudentData, err := json.Marshal(students)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(StudentData)
}

func (h PostHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var s student
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("INSERT INTO students (first_name, last_name, birth_date) VALUES ($1, $2, $3)", s.FirstName, s.LastName, s.DateOfBirth)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	slog.Info("A new student has been added. ")

}

func main() {

	db, err := sql.Open("postgres", "user=postgres password=dbzsuper dbname=school_records sslmode=disable")
	// first argument is the driver name, hence postgres.
	// Normally I wouldn't write a password straight into the code like this but secrets management feels
	// out of scope for the assignment.
	if err != nil || db == nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer db.Close()

	http.Handle("/GetStudents", GetHelper{db: db})
	http.Handle("/AddStudent", PostHelper{db: db})

	slog.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
