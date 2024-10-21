package main

import (
	"log"
	"net/http"
	"database/sql"
	"time"

	inertia "github.com/romsar/gonertia"
	_ "github.com/mattn/go-sqlite3"
)

type Tes struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
	db, err := sql.Open("sqlite3", "./fuji.db")
    if err != nil {
        log.Fatal(err)
    }

    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)

    defer db.Close()

	i := InitInertia()

	mux := http.NewServeMux()

	mux.Handle("/", i.Middleware(homeHandler(i, db)))
	mux.Handle("/about", i.Middleware(about(i)))
	mux.Handle("/build/", http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	http.ListenAndServe(":3000", mux)
}

func homeHandler(i *inertia.Inertia, db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tes (id INTEGER NOT NULL PRIMARY KEY, name TEXT)`)
    if err != nil {
        http.Error(w, "Error creating table", http.StatusInternalServerError)
        return
    }

    // stmt, err := db.Prepare("INSERT INTO tes(name) VALUES(?)")
    // if err != nil {
    //     http.Error(w, "Error preparing insert statement", http.StatusInternalServerError)
    //     return
    // }
    // defer stmt.Close()

    // _, err = stmt.Exec("fuji")
    // if err != nil {
    //     http.Error(w, "Error inserting data", http.StatusInternalServerError)
    //     return
    // }

    rows, err := db.Query("SELECT id, name FROM tes")
    if err != nil {
        http.Error(w, "Error querying data", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var results []Tes
    for rows.Next() {
        var tes Tes
        if err := rows.Scan(&tes.Id, &tes.Name); err != nil {
            http.Error(w, "Error scanning row", http.StatusInternalServerError)
            return
        }
        results = append(results, tes)
    }

		err = i.Render(w, r, "index", inertia.Props{
			"tes": results,
			// "test": results,
		})
		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func handleServerErr(w http.ResponseWriter, err error) {
	log.Printf("http error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}

func about(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "about")
		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}