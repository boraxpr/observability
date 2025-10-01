package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var (
	// Array of random HTTP code
	frequent_http_code = [6]int{301, 302, 400, 403, 404, 500}
)

func main() {
	// Serve Hello World! on /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling " + r.URL.Path + " ...")
		http_code, http_text := randomHTTPCode("Hello World!")
		w.WriteHeader(http_code)
		fmt.Fprintln(w, http_text)
		if http_code != 200 {
			log.Println("Request " + r.URL.Path + " failed with HTTP code " + strconv.Itoa(http_code))
		}
	})

	// Assign application port to run
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8081"
	}

	log.Println("Starting service on port " + appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}

// Function to random HTTP code
func randomHTTPCode(text string) (int, string) {
	http_code := 200
	http_text := text
	// You will have 20% chance to get non-200 http code
	if rand.Float64() < 0.2 {
		http_code = frequent_http_code[rand.Intn(len(frequent_http_code))]
		http_text = "HTTP Code: " + strconv.Itoa(http_code)
	}
	return http_code, http_text
}
