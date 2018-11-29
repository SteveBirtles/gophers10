package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func moveHandler(w http.ResponseWriter, r *http.Request) {

	/* Code goes here */

}

func startServer() {

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/direction/", moveHandler)

	err := http.ListenAndServe(":8081", http.DefaultServeMux)
	if err != nil {
		fmt.Println("Error:", err)
	}

}
