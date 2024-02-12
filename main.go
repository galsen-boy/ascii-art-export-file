package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.HandleFunc("/", StartPage)
	http.HandleFunc("/ascii-art", SubmitTing)
	http.HandleFunc("/down", download)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
