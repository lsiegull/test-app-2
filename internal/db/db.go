package db

import (
    "database/sql"
    "errors"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
)

type DB struct{
    conn *sql.DB
}

type User struct{
    ID int64 `json:"id"`
    Email string `json:"email"`
    FullName string `json:"full_name"`
}

type Class struct{
    ID int64 `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    Capacity int `json:"capacity"`
}

type Activity struct{
    ID int64 `json:"id"`
    UserID int64 `json:"user_id"`
    Type string `json:"type"`
    DurationMinutes int `json:"duration_minutes"`
    Notes string `json:"notes"`
    CreatedAt time.Time `json:"created_at"`
}

func New(dsn string) (*DB, error) {
    conn, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, err
    }
    return &DB{conn: conn}, nil
}

func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) ExecSchema(schema string) error {
    _, err := db.conn.Exec(schema)
    return err
}

func (db *DB) CreateUser(email, password, fullName string) (int64, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil { return 0, err }
    res, err := db.conn.Exec(`INSERT INTO users (email, password_hash, full_name) VALUES (?, ?, ?)`, email, string(hashed), fullName)
    if err != nil { return 0, err }
    return res.LastInsertId()
}

func (db *DB) GetUserByEmail(email string) (User, string, error) {
    var u User
    var pw string
    row := db.conn.QueryRow(`SELECT id, email, password_hash, full_name FROM users WHERE email = ?`, email)
    var id int64; var e, ph, fn sql.NullString
    if err := row.Scan(&id, &e, &ph, &fn); err != nil {
        if err == sql.ErrNoRows { return u, "", errors.New("not found") }
        return u, "", err
    }
    u.ID = id
    u.Email = e.String
    u.FullName = fn.String
    pw = ph.String
    return u, pw, nil
}

func (db *DB) Authenticate(email, password string) (User, error) {
    u, pw, err := db.GetUserByEmail(email)
    if err != nil { return u, err }
    if bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)) != nil {
        return u, errors.New("invalid credentials")
    }
    return u, nil
}

func (db *DB) CreateClass(name, description string, capacity int) (int64, error) {
    res, err := db.conn.Exec(`INSERT INTO classes (name, description, capacity) VALUES (?, ?, ?)`, name, description, capacity)
    if err != nil { return 0, err }
    return res.LastInsertId()
}

func (db *DB) ListClasses(limit, offset int) ([]Class, error) {
    rows, err := db.conn.Query(`SELECT id, name, description, capacity FROM classes ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    out := []Class{}
    for rows.Next() {
        var c Class
        if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Capacity); err != nil { return nil, err }
        out = append(out, c)
    }
    return out, nil
}

func (db *DB) SignupClass(userID, classID int64) error {
    _, err := db.conn.Exec(`INSERT INTO class_signups (user_id, class_id) VALUES (?, ?)`, userID, classID)
    return err
}

func (db *DB) LogActivity(userID int64, typ string, duration int, notes string) (int64, error) {
    res, err := db.conn.Exec(`INSERT INTO activities (user_id, type, duration_minutes, notes) VALUES (?, ?, ?, ?)`, userID, typ, duration, notes)
    if err != nil { return 0, err }
    return res.LastInsertId()
}

func (db *DB) ListActivitiesByUser(userID int64, limit, offset int) ([]Activity, error) {
    rows, err := db.conn.Query(`SELECT id, user_id, type, duration_minutes, notes, created_at FROM activities WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    out := []Activity{}
    for rows.Next() {
        var a Activity
        var created string
        if err := rows.Scan(&a.ID, &a.UserID, &a.Type, &a.DurationMinutes, &a.Notes, &created); err != nil { return nil, err }
        a.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
        out = append(out, a)
    }
    return out, nil
}
