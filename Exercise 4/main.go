// Client/Server AJAX JSON Communication using golang web-server and JQuery
// Visit: http://127.0.0.1:8080
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"time"
)

//ReqDate is request struct
type ReqDate struct {
	Date string
}

//RespDate is request struct
type RespDate struct {
	Date           string
	LastDayOfMonth string
}

// Default Request Handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("template", "aj-json.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AJAX Request Handler
func ajaxHandler(w http.ResponseWriter, r *http.Request) {
	//parse request to struct
	var d ReqDate
	var resp RespDate
	fmt.Println("Came inside ajax handler")
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp.Date = d.Date

	yourDate, _ := time.Parse("2006-01-02", resp.Date)
	currentYear, currentMonth, _ := yourDate.Date()
	currentLocation := yourDate.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	_, _, Lastday := lastOfMonth.Date()

	resp.LastDayOfMonth = strconv.Itoa(Lastday)

	// create json response from struct
	a, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(string(a))
	w.Write([]byte("<h2>" + string(a) + "<h2>"))
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/ajax", ajaxHandler)
	http.ListenAndServe(":8000", nil)
}
