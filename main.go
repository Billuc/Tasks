package main

import (
	"net/http"
)

func main() {
	// s := tasks.NewServer()
	// log.Println("Starting server on :8082")

	// s.ServeDir("./public", "/")
	// s.ServeRoute(tasks.Route{Method: tasks.GET, Path: "/tasks"}, tasksRoute.Test)

	// err := s.Start(8082)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	http.Handle("GET /", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":8082", nil)
}
