package main // main defines top level executable, as opposed to a library.

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	school_db "github.com/czarv30/Go_simple_server_db"
)

// The way this library works the handler functions for each endpoint need to be methods on a struct.
// Despite being identical, I need two separate types because each differentiates Get vs Post.
type GetHelper struct {
	db *school_db.SchoolDb
}
type PostHelper struct {
	db *school_db.SchoolDb
}

// The struct should be referenced by pointer, for efficiency reasons. Otherwise we copy the struct.
func (h GetHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// syntax above: "h" is like self in python class definition. The stuff in parenthesis attaches the function to the struct as a method. In go, structs can have methods.
	// "w" is an instance of the http.ResponseWriter type
	// Likewise with r but in this case is a pointer, as required by the package.

	students, err := h.db.GetAllStudents()

	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}

	StudentData, err := json.Marshal(students)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(StudentData)
}

func (h PostHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var s school_db.Student
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.db.PostStudent(s)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	slog.Info("A new student has been added. ")

}

func main() {

	dbi, err := school_db.InitSchoolDb("user=postgres password=dbzsuper dbname=school_records sslmode=disable") // Initialize the database connection.

	if err != nil || dbi == nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer dbi.Close() // In Go, "defer" schedules a function call to be executed after the surrounding function returns, regardless of whether the function exits normally or via an error.

	http.Handle("/GetStudents", GetHelper{db: dbi})
	http.Handle("/AddStudent", PostHelper{db: dbi})

	slog.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
