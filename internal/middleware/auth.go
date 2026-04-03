package middleware

import (
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_KEY"))},
		Extractor:    extractors.Chain(extractors.FromCookie("uusr"), extractors.FromAuthHeader("Bearer")),
		ErrorHandler: jwtError,
	})
}

func GenerateNewUserToken(userID uuid.UUID, tokenType string) (string, error) {
	var expiry int64
	if tokenType == "refresh" {
		expiry = time.Now().Add(time.Hour * 72).Unix()
	} else {
		expiry = time.Now().Add(time.Minute * 30).Unix()
	}
	claims := jwt.MapClaims{
		"id":  userID.String(),
		"exp": expiry,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func jwtError(c fiber.Ctx, err error) error {
	return c.JSON(map[string]any{"error": err.Error()})
}
