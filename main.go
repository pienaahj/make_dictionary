package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type dict map[string][]string

var dData dict

type dictSQL struct {
	word     string
	meaning1 string
	meaning2 string
	meaning3 string
	meaning4 string
	meaning5 string
}

var dictSQLX []dictSQL

func convertToSQL(d dict) []dictSQL {
	dSQL := dictSQL{}
	dSQLX := make([]dictSQL, 0)
	for k, v := range d {
		dSQL.word = k
		dSQL.meaning1 = v[0]
		if len(v) == 2 {
			dSQL.meaning2 = v[1]
		} else if len(v) == 3 {
			dSQL.meaning3 = v[2]
		} else if len(v) == 4 {
			dSQL.meaning4 = v[3]
		} else if len(v) == 5 {
			dSQL.meaning5 = v[4]
		}

		for i := len(v); i < 4; i++ {
			if i == 1 {
				dSQL.meaning2 = " "
			} else if i == 2 {
				dSQL.meaning3 = " "
			} else if i == 3 {
				dSQL.meaning4 = " "
			} else if i == 4 {
				dSQL.meaning5 = " "
			}
		}
		dSQLX = append(dSQLX, dSQL)
		dSQL = dictSQL{}
	}
	return dSQLX
}

func findWord(s string, d dict) ([]string, error) {
	if def, ok := d[s]; ok {
		return def, nil
	} else if def, ok := d[strings.ToLower(s)]; ok {
		return def, nil
	}

	return nil, errors.New("could not find the word")
}

func read(fName string) {
	file, err := ioutil.ReadFile(fName)
	if err != nil {
		log.Fatalf("Could not open file: %v\n", err)
	}

	jErr := json.Unmarshal(file, &dData)
	if jErr != nil {
		log.Fatalf("Could not unmarshal the file: %v\n", err)
	}
}

// Read the command line
func getInput() string {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a word to search for")
	}
	var words string
	for _, v := range os.Args[1:] {
		words += v + " "
	}
	words = strings.TrimSuffix(words, " ")
	return words

}

func outPutSQL(ss string, s []dictSQL) {
	dictItem := dictSQL{}
	for _, v := range s {
		if v.word == ss {
			dictItem = v
			break
		}
	}
	fmt.Printf("You searched for the meaning of (%s):\n", ss)
	fmt.Println("Possible meanings: ")
	fmt.Printf(" %s\n %s\n %s\n %s\n %s\n", dictItem.meaning1, dictItem.meaning2, dictItem.meaning3, dictItem.meaning4, dictItem.meaning5)
}

func main() {
	fName := "data/data.json"
	read(fName)
	// connect to mysql
	connectStr := "root:Pinepine01#@tcp(127.0.0.1:3306)/dictionary?charset=utf8mb4"
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatalf("Could not connect to msql %v\n", err)
	}
	defer db.Close()
	fmt.Println("db is connected")
	// test connection
	err = db.Ping()
	if err != nil {
		log.Printf("Db not responding %v\n", err)
	}
	fmt.Println("db is available")
	// get work from user
	sWord := getInput()

	SQLData := convertToSQL(dData)
	outPutSQL(sWord, SQLData)
	// Create the db
	query := "INSERT INTO dict (id,word, meaning1, meaning2, meaning3, meaning4, meaning5) VALUES ( ?,?,?,?,?,?,? )"
	stat, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Cannot prepare statement %v\n", err)
	}
	recordC := 0
	for k, v := range SQLData {
		_, err := stat.Exec(k, v.word, v.meaning1, v.meaning2, v.meaning3, v.meaning4, v.meaning5)
		if err != nil {
			log.Printf("could not add record %v\n", err)
		}
		recordC += k
	}
	fmt.Printf("%d records added\n", recordC)

}
