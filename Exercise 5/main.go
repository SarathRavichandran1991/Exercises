package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
)

//GetPrime to check prime
func GetPrime(cs chan int, value int) {
	prime := func(value int) int {
		for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
			if value%i == 0 {
				return 0
			}
		}

		return value
	}
	cs <- prime(value)
}

//Home screen
func Home(w http.ResponseWriter, r *http.Request) {
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

func receiveAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postdata, err := strconv.Atoi(r.FormValue("postdata"))
		if err != nil {
			fmt.Println("Wrong input format")
			w.Write([]byte("<h2>Please enter valid number!<h2>"))
		}

		var primes []string

		cs := make(chan int)

		for i := 2; i <= postdata; i++ {
			go GetPrime(cs, i)
			printvalue := <-cs
			if printvalue > 0 {
				primes = append(primes, strconv.Itoa(printvalue))
			}
		}

		result := strings.Join(primes, ",")
		fmt.Println(result)
		w.Write([]byte("<h2>" + result + "<h2>"))
	}
}

func main() {
	// http.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/receive", receiveAjax)

	http.ListenAndServe(":8080", mux)
}
