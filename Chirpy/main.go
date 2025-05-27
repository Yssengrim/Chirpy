package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()

	serveMux.Handle("/", http.FileServer(http.Dir("."))) // <-- Look at this line!

	chirpyServer := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	chirpyServer.ListenAndServe()
}
