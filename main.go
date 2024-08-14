package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

var connections = map[string]*TursoConnection{}

type TursoConnection struct {
	DB *sql.DB
}

type TableRow struct {
	ID  int    `json:"id"`
	Col string `json:"col"`
}

func getTursoConnection(name string) (*TursoConnection, error) {
	var conn *TursoConnection
	conn, exists := connections[name]
	if exists {
		return conn, nil
	}

	dbURL := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_DB_URL"), os.Getenv("TURSO_AUTH_TOKEN"))
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	conn = &TursoConnection{
		DB: db,
	}
	connections[name] = conn

	return conn, nil
}

func main() {
	http.HandleFunc("/some_table", func(w http.ResponseWriter, r *http.Request) {
		conn, err := getTursoConnection("some_table")
		if err != nil {
			http.Error(w, "error connecting to database: "+err.Error(), 500)
		}

		rows, err := conn.DB.Query(`
            SELECT * FROM some_table
        `)
		if err != nil {
			http.Error(w, "error querying db: "+err.Error(), 500)
			return
		}

		tableRows := []TableRow{}
		for rows.Next() {
			tableRow := TableRow{}
			rows.Scan(&tableRow.ID, &tableRow.Col)
			tableRows = append(tableRows, tableRow)
		}
		defer rows.Close()

		if len(tableRows) == 0 {
			w.Write([]byte("No rows found"))
			return
		}

		bytes, err := json.Marshal(tableRows)
		if err != nil {
			http.Error(w, "error marshalling rows"+err.Error(), 500)
			return
		}

		w.Header().Add("content-type", "application/json")
		w.Write(bytes)
	})

	http.HandleFunc("/some_other_table", func(w http.ResponseWriter, r *http.Request) {
		conn, err := getTursoConnection("some_other_table")
		if err != nil {
			http.Error(w, "error connecting to database: "+err.Error(), 500)
		}

		rows, err := conn.DB.Query(`
            SELECT * FROM some_other_table
        `)
		if err != nil {
			http.Error(w, "error querying db: "+err.Error(), 500)
			return
		}

		tableRows := []TableRow{}
		for rows.Next() {
			tableRow := TableRow{}
			rows.Scan(&tableRow.ID, &tableRow.Col)
			tableRows = append(tableRows, tableRow)
		}
		defer rows.Close()

		if len(tableRows) == 0 {
			w.Write([]byte("No rows found"))
			return
		}

		bytes, err := json.Marshal(tableRows)
		if err != nil {
			http.Error(w, "error marshalling rows"+err.Error(), 500)
			return
		}

		w.Header().Add("content-type", "application/json")
		w.Write(bytes)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
