cat > models/models.go << 'EOF'
package models

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Client struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    AllowedPaths []string `json:"allowed_paths"`
    AllowedMethods []string `json:"allowed_methods"`
    TokenTTL     int      `json:"token_ttl"`
    SecretKey    string   `json:"-"`
}

type TokenClaims struct {
    ClientID string `json:"client_id"`
    jwt.RegisteredClaims
}

type LoginRequest struct {
    ClientID  string `json:"client_id"`
    SecretKey string `json:"secret_key"`
}

type TokenResponse struct {
    Token     string `json:"token"`
    ExpiresAt string `json:"expires_at"`
    TokenType string `json:"token_type"`
}

type LogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    ClientID    string    `json:"client_id"`
    Method      string    `json:"method"`
    Path        string    `json:"path"`
    StatusCode  int       `json:"status_code"`
    ResponseTime string    `json:"response_time"`
}
EOF
