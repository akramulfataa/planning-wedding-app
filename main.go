package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"wedding-planner/handlers"
)

func main() {
	http.HandleFunc("/hitung", handlers.HandleHitung)
	fmt.Println("Server started at :8080")
	slog.Info("server berjalan di port: ", http.ListenAndServe(":8080", nil))
}
