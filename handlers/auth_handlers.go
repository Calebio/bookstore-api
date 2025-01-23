package handlers

import (
    "bookstore-api/data"
    "bookstore-api/models"
    "bookstore-api/utils"
    "encoding/json"
    "net/http"
    "time"

    "golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    _, err = data.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)

    row := data.DB.QueryRow("SELECT id, password FROM users WHERE username = $1", user.Username)
    var storedUser models.User
    err := row.Scan(&storedUser.ID, &storedUser.Password)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    token, expiresAt, err := utils.GenerateToken(storedUser.ID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    _, err = data.DB.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", storedUser.ID, token, expiresAt)
    if err != nil {
        http.Error(w, "Error storing token", http.StatusInternalServerError)
        return
    }

    // Send the token and its expiration to the client
    json.NewEncoder(w).Encode(map[string]interface{}{
        "token":      token,
        "expires_at": expiresAt.Format(time.RFC3339),
    })
}
