package main

import (
	"fmt"
	"net/http"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func moveHandler(w http.ResponseWriter, r *http.Request) {

	const MaxQueueLength = 3

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathBits := strings.Split(r.URL.Path, "/")
	direction := pathBits[len(pathBits)-1]
	fmt.Println("/direction/", direction)

	if len(player[0].targetQueue) > MaxQueueLength {
		return
	}

	var lastTarget Location
	if len(player[0].targetQueue) > 0 {
		lastTarget = player[0].targetQueue[len(player[0].targetQueue)-1]
	} else {
		lastTarget = player[0].lastLocation
	}

	switch direction {
	case "up":
		player[0].targetQueue = append(player[0].targetQueue, Location{lastTarget.x, lastTarget.y - 1})
	case "down":
		player[0].targetQueue = append(player[0].targetQueue, Location{lastTarget.x, lastTarget.y + 1})
	case "left":
		player[0].targetQueue = append(player[0].targetQueue, Location{lastTarget.x - 1, lastTarget.y})
	case "right":
		player[0].targetQueue = append(player[0].targetQueue, Location{lastTarget.x + 1, lastTarget.y})
	}

}

func startServer() {

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/direction/", moveHandler)

	err := http.ListenAndServe(":8081", http.DefaultServeMux)
	if err != nil {
		fmt.Println("Error:", err)
	}

}
