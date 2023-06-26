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

func main() {
	urlExample := "postgres://user:password@localhost:5432/DelDB?sslmode=disable"
	db, err := sql.Open("postgres", urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	//func conn

	// Fiber setup and routes starting
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		res, err := queryToJson(db, "select * from warriors")
		if err != nil {
			log.Println("data error:", err)
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "something went wrong!",
			})
		}

		log.Printf("res: %#v \n", string(res))
		return c.SendString(string(res))
	})

	//simple endpoint that will scan values form db
	app.Get("/simple", func(c *fiber.Ctx) error {
		res, err := queryScanJson(db)
		if err != nil {
			log.Println("data error:", err)
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "something went wrong!",
			})
		}

		return c.SendString(string(res))
	})

	log.Fatal(app.Listen(":3000"))
}

type Warrior struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Teacher   string `json:"teacher"`
	IsActive  bool   `json:"isActive"`
}

func queryScanJson(db *sql.DB) ([]byte, error) {
	res := []Warrior{}
	q := "select first_name,last_name, teacher, is_active from warriors"

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		w := Warrior{}
		var lName sql.NullString
		// err := rows.Scan(&w.FirstName, &w.LastName, &w.Teacher, &w.IsActive)
		err := rows.Scan(&w.FirstName, &lName, &w.Teacher, &w.IsActive)
		if err != nil {
			return nil, err
		}
		if lName.Valid {
			w.LastName = lName.String
		}
		res = append(res, w)
	}

	return json.MarshalIndent(res, "", "\t")
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
