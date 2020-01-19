package main

import (
 "database/sql"
 "io/ioutil"
 "fmt"
 _ "github.com/mattn/go-sqlite3"
 "github.com/bmizerany/pat"
 "encoding/json"
 "log"
 "net/http"
 "time"
 "os"
 "encoding/base64"
)

var database *sql.DB

type serverInfo struct {
        ID      int64 `json:"id"`
        Name    string `json:"name"`
		Date    string `json:"date"`
		Comments string `json:"comments"`
}

type servers []serverInfo

func main(){
	db, err := sql.Open("sqlite3", "database.db")
	checkErr(err)
	database = db
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS server (id INTEGER PRIMARY KEY, name TEXT, date TEXT, comments TEXT)")
	statement.Exec()

	r := pat.New()
	r.Get("/v1/server", http.HandlerFunc(getservers))
	r.Post("/v1/server", http.HandlerFunc(insert))
	
	
	http.Handle("/", r)
	
	if (fileExists("/cert.pem") && fileExists("/key.pem")) {
		log.Print(" Found SSL certificate and private key...")
		log.Print(" Listening secure on port 8443")
		err = http.ListenAndServeTLS(":8443", "/cert.pem", "/key.pem", nil)
	} else {
		log.Print(" No SSL certificate or private key found")
		log.Print(" Listening unsecure on port 8080")
		err = http.ListenAndServe(":8080", nil)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

const unauth = http.StatusUnauthorized
const userPassword = "foo:bar"
func getservers(w http.ResponseWriter, r *http.Request){
	keys := r.URL.Query()["since"]
	auth := r.Header.Get("Authorization")
	if (len(auth) > 0){
		up, err := base64.StdEncoding.DecodeString(auth[6:])
		checkErr(err)
		if (string(up) != userPassword){
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}

		if (len(keys) > 0){
			sinceDate := keys[0]
			now := time.Now()
			timestamp := now.Format("2006-01-02T15:04")
			query := "SELECT id,name,date,comments FROM server WHERE datetime(date) BETWEEN datetime(\"" + sinceDate + "\") and datetime(\"" + timestamp + "\")"
			rows, err := database.Query(query)
			checkErr(err)
			
			var server servers

			for rows.Next() {
				var oneserver serverInfo
				err = rows.Scan(&oneserver.ID, &oneserver.Name, &oneserver.Date, &oneserver.Comments)
				checkErr(err)
				server = append(server, oneserver)
			}
			
			json, err := json.Marshal(server)
			checkErr(err)
			fmt.Fprintf(w, "%s", string(json))
		} else {
			rows, err := database.Query("SELECT id,name,date,comments FROM server")
			checkErr(err)

			var server servers

			for rows.Next() {
				var oneserver serverInfo
				err = rows.Scan(&oneserver.ID, &oneserver.Name, &oneserver.Date, &oneserver.Comments)
				checkErr(err)
				server = append(server, oneserver)
			}
			
			json, err := json.Marshal(server)
			checkErr(err)
			fmt.Fprintf(w, "%s", string(json))
		}
	} else {
		http.Error(w, http.StatusText(unauth), unauth)
		return
	}
}

func insert(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if (len(auth) > 0){
		up, err := base64.StdEncoding.DecodeString(auth[6:])
		checkErr(err)
		if (string(up) != userPassword){
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		checkErr(err)

		now := time.Now()
		timestamp := now.Format("2006-01-02T15:04")
		var oneserver serverInfo

		err = json.Unmarshal(body, &oneserver)
		checkErr(err)

		statement, err := database.Prepare("INSERT INTO server(name, date, comments) VALUES (?, ?, ?)")
		checkErr(err)

		result, err := statement.Exec(oneserver.Name, timestamp, oneserver.Comments)
		checkErr(err)

		newID, err := result.LastInsertId()
		checkErr(err)

		oneserver.ID = newID
		oneserver.Date = timestamp
		outjson, err := json.Marshal(oneserver)
		checkErr(err)

		fmt.Fprintf(w, "%s", string(outjson))
	} else {
		http.Error(w, http.StatusText(unauth), unauth)
		return
	}
}

func checkErr(err error) {
 if err != nil {
  panic(err)
 }
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}