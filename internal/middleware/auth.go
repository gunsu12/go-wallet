package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil handler autorization

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ada pada header"})
			c.Abort()
			return
		}

		// Format: Bearer token
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah" + tokenString})
			c.Abort()
			return
		}

		// fmt.Println("HEADER :", authHeader)
		// fmt.Println("TOKEN STRING:", tokenString)
		// fmt.Println("ENV JWT_SECRET_SAAT_VERIF:", os.Getenv("JWT_SECRET"))

		// Verifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			// fmt.Println("JWT ERROR:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid", "debug": err.Error()})
			c.Abort()
			return
		}

		// simpan claim ke konteks
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// fmt.Println("CLAIMS:", claims)
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
		}

		c.Next()
	}
}
