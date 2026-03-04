package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "jwt-auth-service/config"
    "jwt-auth-service/handlers"
    "jwt-auth-service/middleware"
)

func main() {
    err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    handlers.InitLogger()
    defer handlers.CloseLogger()

    r := mux.NewRouter()
    r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
    r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)
    api.HandleFunc("/proxy/{path:.*}", handlers.ProxyHandler)

    r.HandleFunc("/admin/logs", handlers.GetLogs).Methods("GET")

    serverAddr := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
    fmt.Printf("Server starting on port %s\n", config.AppConfig.ServerPort)
    log.Fatal(http.ListenAndServe(serverAddr, r))
}
