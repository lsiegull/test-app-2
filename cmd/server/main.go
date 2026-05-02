package main

import (
    "database/sql"
    "embed"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "time"

    "test-app-2/internal/db"

    "github.com/golang-jwt/jwt/v5"
)

//go:embed ../../internal/db/migrations/schema.sql
var schemaFS embed.FS

var jwtKey = []byte("secret-demo-key")

func main() {
    dsn := ":memory:?cache=shared"
    database, err := db.New(dsn)
    if err != nil { log.Fatal(err) }
    defer database.Close()

    schemaBytes, _ := schemaFS.ReadFile("../../internal/db/migrations/schema.sql")
    if err := database.ExecSchema(string(schemaBytes)); err != nil { log.Fatal(err) }

    // seed a class
    _, _ = database.CreateClass("Morning Yoga", "Gentle stretch and breathwork", 15)

    mux := http.NewServeMux()
    mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request){
        if r.Method != http.MethodPost { http.Error(w, "method", http.StatusMethodNotAllowed); return }
        var req struct{ Email, Password, FullName string }
        decodeJSON(r.Body, &req)
        id, err := database.CreateUser(req.Email, req.Password, req.FullName)
        if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        writeJSON(w, map[string]any{"id": id})
    })

    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request){
        if r.Method != http.MethodPost { http.Error(w, "method", http.StatusMethodNotAllowed); return }
        var req struct{ Email, Password string }
        decodeJSON(r.Body, &req)
        user, err := database.Authenticate(req.Email, req.Password)
        if err != nil { http.Error(w, "invalid credentials", http.StatusUnauthorized); return }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": user.ID, "exp": time.Now().Add(24*time.Hour).Unix()})
        tokStr, _ := token.SignedString(jwtKey)
        writeJSON(w, map[string]any{"token": tokStr})
    })

    mux.HandleFunc("/classes", func(w http.ResponseWriter, r *http.Request){
        switch r.Method {
        case http.MethodGet:
            limit := 20
            offset := 0
            classes, err := database.ListClasses(limit, offset)
            if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
            writeJSON(w, classes)
        case http.MethodPost:
            var c struct{ Name, Description string; Capacity int }
            decodeJSON(r.Body, &c)
            id, err := database.CreateClass(c.Name, c.Description, c.Capacity)
            if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
            writeJSON(w, map[string]any{"id": id})
        default:
            http.Error(w, "method", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/classes/", func(w http.ResponseWriter, r *http.Request){
        // POST /classes/{id}/signup
        if r.Method != http.MethodPost { http.Error(w, "method", http.StatusMethodNotAllowed); return }
        parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/classes/"), "/")
        if len(parts) < 2 || parts[1] != "signup" { http.NotFound(w, r); return }
        id, _ := strconv.ParseInt(parts[0], 10, 64)
        userID, ok := authUserID(r)
        if !ok { http.Error(w, "unauth", http.StatusUnauthorized); return }
        if err := database.SignupClass(userID, id); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        writeJSON(w, map[string]any{"ok": true})
    })

    mux.HandleFunc("/activities", func(w http.ResponseWriter, r *http.Request){
        userID, ok := authUserID(r)
        if !ok { http.Error(w, "unauth", http.StatusUnauthorized); return }
        switch r.Method {
        case http.MethodPost:
            var a struct{ Type string; DurationMinutes int; Notes string }
            decodeJSON(r.Body, &a)
            _, err := database.LogActivity(userID, a.Type, a.DurationMinutes, a.Notes)
            if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
            writeJSON(w, map[string]any{"ok": true})
        case http.MethodGet:
            activities, _ := database.ListActivitiesByUser(userID, 50, 0)
            writeJSON(w, activities)
        default:
            http.Error(w, "method", http.StatusMethodNotAllowed)
        }
    })

    addr := ":8080"
    log.Printf("listening %s", addr)
    log.Fatal(http.ListenAndServe(addr, logging(mux)))
}

func decodeJSON(r io.Reader, v any) {
    _ = json.NewDecoder(r).Decode(v)
}

func writeJSON(w http.ResponseWriter, v any) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(v)
}

func logging(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        start := time.Now()
        h.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}

func authUserID(r *http.Request) (int64, bool) {
    ah := r.Header.Get("Authorization")
    if ah == "" { return 0, false }
    parts := strings.SplitN(ah, " ", 2)
    if len(parts) != 2 { return 0, false }
    tok := parts[1]
    token, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error){ return jwtKey, nil })
    if err != nil || !token.Valid { return 0, false }
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        switch v := claims["sub"].(type) {
        case float64:
            return int64(v), true
        case string:
            if id, err := strconv.ParseInt(v, 10, 64); err == nil { return id, true }
        }
    }
    return 0, false
}
