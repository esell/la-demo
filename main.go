package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type conf struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var prePopulateDatabase = flag.Bool("p", false, "Pre-populate the database with sample values")

func main() {
	flag.Parse()

	var parsedconfig conf
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("unable to read config file, exiting...")
	}
	if err := json.Unmarshal(file, &parsedconfig); err != nil {
		log.Fatal("unable to marshal config file, exiting...")
	}

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=true", parsedconfig.Username, parsedconfig.Password, parsedconfig.Host, parsedconfig.Database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	if *prePopulateDatabase {
		prePopulate(db)
	}
	// Vendor #1
	http.Handle("/vendor1/listsubmissions", submissionsListHandler(db))
	// Vendor #1
	http.Handle("/vendor1/newsubmission", vendor1SpeakerSubmission(db))
	// Vendor #2
	http.Handle("/vendor2/newsubmission", vendor2SpeakerSubmission(db))
	http.Handle("/", http.FileServer(http.Dir("html")))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func prePopulate(db *sql.DB) {
	// Drop previous table of same name if one exists.
	_, err := db.Exec("DROP TABLE IF EXISTS submissions;")
	if err != nil {
		panic(err.Error())
	}
	log.Println("Finished dropping table (if existed).")

	// Create table.
	_, err = db.Exec("CREATE TABLE submissions (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, speaker_name VARCHAR(150), speaker_email VARCHAR(150), speaker_topic VARCHAR(150), speaker_status VARCHAR(50));")

	if err != nil {
		panic(err.Error())
	}
	log.Println("Finished creating table.")

	// Insert some data into table.
	sqlStatement, err := db.Prepare("INSERT INTO submissions (speaker_name, speaker_email, speaker_topic, speaker_status) VALUES (?, ?, ?, ?);")
	res, err := sqlStatement.Exec("Jane Doe", "jane@blah.com", "Some really neat stuff", "COMPLETE")
	if err != nil {
		panic(err.Error())
	}
	rowCount, err := res.RowsAffected()
	log.Printf("Inserted %d row(s) of data.\n", rowCount)

	res, err = sqlStatement.Exec("John Doe", "john@blah.com", "Things you should know", "COMPLETE")
	if err != nil {
		panic(err.Error())
	}
	rowCount, err = res.RowsAffected()
	log.Printf("Inserted %d row(s) of data.\n", rowCount)

	log.Println("Done.")
}

// Vendor #1
func submissionsListHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}

		var (
			id     int
			name   string
			email  string
			topic  string
			status string
		)

		rows, err := db.Query("SELECT id, speaker_name, speaker_email, speaker_topic, speaker_status from submissions;")
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to retrieve rows", http.StatusMethodNotAllowed)
			return
		}
		defer rows.Close()
		subListFinal := SubmitterList{}
		subListFinal.Submissions = make([]Submitter, 1)
		for rows.Next() {
			err := rows.Scan(&id, &name, &email, &topic, &status)
			if err != nil {
				log.Println(err)
			}
			tempSub := Submitter{ID: id, Name: name, Email: email, Topic: topic, Status: status}
			subListFinal.Submissions = append(subListFinal.Submissions, tempSub)
		}

		subListFinalJSON, err := json.Marshal(subListFinal)
		if err != nil {
			http.Error(w, "marshal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(subListFinalJSON)
		return
	})
}

// Vendor #1
func vendor1SpeakerSubmission(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var speakerSubmission Submitter
		err := decoder.Decode(&speakerSubmission)
		if err != nil {
			http.Error(w, "decode error", http.StatusInternalServerError)
			return
		}
		log.Println(speakerSubmission)

		// Insert some data into table.
		sqlStatement, err := db.Prepare("INSERT INTO submissions (speaker_name, speaker_email, speaker_topic, speaker_status) VALUES (?, ?, ?, ?);")
		res, err := sqlStatement.Exec(speakerSubmission.Name, speakerSubmission.Email, speakerSubmission.Topic, speakerSubmission.Status)
		if err != nil {
			http.Error(w, "sql error", http.StatusInternalServerError)
			return
		}
		rowCount, err := res.RowsAffected()
		log.Printf("Inserted %d row(s) of data.\n", rowCount)

		if err != nil {
			http.Error(w, "sql error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})
}

// Vendor #2
func vendor2SpeakerSubmission(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var speakerSubmission Submitter
		err := decoder.Decode(&speakerSubmission)
		if err != nil {
			http.Error(w, "decode error", http.StatusInternalServerError)
			return
		}
		log.Println(speakerSubmission)

		// change status
		speakerSubmission.Status = "COMPLETE"
		// Insert some data into table.
		sqlStatement, err := db.Prepare("UPDATE submissions set speaker_status=? where speaker_email=? AND speaker_topic=?")
		res, err := sqlStatement.Exec(speakerSubmission.Status, speakerSubmission.Email, speakerSubmission.Topic)
		if err != nil {
			http.Error(w, "sql error", http.StatusInternalServerError)
			return
		}
		rowCount, err := res.RowsAffected()
		log.Printf("Inserted %d row(s) of data.\n", rowCount)

		if err != nil {
			http.Error(w, "sql error", http.StatusInternalServerError)
			return
		}
		speakerSubmissionJSON, err := json.Marshal(speakerSubmission)
		if err != nil {
			http.Error(w, "marshal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(speakerSubmissionJSON)
		return
	})
}
