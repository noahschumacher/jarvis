package main

import (
	"fmt"
	"net/http"
)

// Write a simply program that will fetch the price information from Jupiter API
// every 5 seconds and print the price of the asset in the console.
func main() {
	fmt.Println("Jarvis, yeet me coin")

	fartFetcher := newPriceModel("fart", fartAddr)
	go fartFetcher.run()

	http.HandleFunc("/events", chartHandler(fartFetcher))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/index.html") // Serve the HTML file
	})
	fmt.Println("Server running on http://localhost:3333")
	http.ListenAndServe(":3333", nil)
}
