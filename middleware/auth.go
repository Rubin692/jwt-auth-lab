package middleware

import (
    "fmt"
    "net/http"
    "strings"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "jwt-auth-service/config"
    "jwt-auth-service/models"
)

func GenerateToken(client *models.Client) (string, time.Time, error) {
    expirationTime := time.Now().Add(time.Duration(client.TokenTTL) * time.Minute)
    
    claims := &models.TokenClaims{
        ClientID: client.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "jwt-auth-service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(client.SecretKey))
    
    return tokenString, expirationTime, err
}

func ValidateToken(tokenString string) (*models.TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        
        if claims, ok := token.Claims.(*models.TokenClaims); ok {
            client := config.GetClientByID(claims.ClientID)
            if client != nil {
                return []byte(client.SecretKey), nil
            }
        }
        
        return nil, fmt.Errorf("invalid token")
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/login" {
            next.ServeHTTP(w, r)
            return
        }

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }

        claims, err := ValidateToken(tokenParts[1])
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        client := config.GetClientByID(claims.ClientID)
        if client == nil {
            http.Error(w, "Client not found", http.StatusForbidden)
            return
        }

        methodAllowed := false
        for _, method := range client.AllowedMethods {
            if method == r.Method {
                methodAllowed = true
                break
            }
        }

        if !methodAllowed {
            http.Error(w, "Method not allowed for this client", http.StatusForbidden)
            return
        }

        pathAllowed := false
        for _, path := range client.AllowedPaths {
            if strings.HasPrefix(r.URL.Path, "/api/proxy"+path) {
                pathAllowed = true
                break
            }
        }

        if !pathAllowed {
            http.Error(w, "Path not allowed for this client", http.StatusForbidden)
            return
        }

        r.Header.Set("X-Client-ID", claims.ClientID)
        
        next.ServeHTTP(w, r)
    })
}
EOF
