package middleware

import (
	repository "SessionCookie/Repository"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var Rep = &repository.UserRepository{}

func ValidateSession(ctx *fiber.Ctx) error {

	tokenString := ctx.Cookies("jwt")
	if tokenString == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(
			&fiber.Map{"message": "Not logged in"},
		)
	}

	token, err := jwt.Parse(tokenString, verifyJWT)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Error parsing token"},
		)
	}

	if !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(
			&fiber.Map{"message": "Invalid token"},
		)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Problem claiming token"},
		)
	}

	floatExpiration, ok := claims["exp"].(float64)
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "Invalid expiration time",
			},
		)
	}

	expirationTime := time.Unix(int64(floatExpiration), 0)
	currentTime := time.Now()
	if currentTime.After(expirationTime) {
		return ctx.Status(http.StatusUnauthorized).JSON(
			&fiber.Map{"message": "Logged out, token expired"},
		)
	}

	user, err := Rep.FindUserByName(claims["user"].(string))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Error finding user"},
		)
	}

	ctx.Locals("user", user)

	return ctx.Next()
}

func verifyJWT(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing error")
	}

	return []byte(os.Getenv("JWT_KEY")), nil
}
