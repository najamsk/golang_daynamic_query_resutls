package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

//NOTE: https://forum.golangbridge.org/t/database-rows-scan-unknown-number-of-columns-json/7378/2
//above is the reference for code sample.

// TODO: should use gendric sql driver interface
// and https://pkg.go.dev/github.com/lib/pq
func main() {
	urlExample := "postgres://user:password@localhost:5432/DelDB?sslmode=disable"
	db, err := sql.Open("postgres", urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	//func conn
	res, err := queryToJson(db, "select * from warriors")
	if err != nil {
		fmt.Fprintf(os.Stderr, "queryToJson failed: %v\n", err)
		os.Exit(1)
	}
	log.Printf("res: %#v \n", string(res))
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(string(res))
	})

	log.Fatal(app.Listen(":3000"))

}
func queryToJson(db *sql.DB, query string, args ...interface{}) ([]byte, error) {
	// an array of JSON objects
	// the map key is the field name
	var objects []map[string]interface{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		// figure out what columns were returned
		// the column names will be the JSON object field keys
		columns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			// object[column.Name()] = reflect.New(column.ScanType()).Interface()
			object[column.Name()] = new(*string)
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	// indent because I want to read the output
	return json.MarshalIndent(objects, "", "\t")
}
