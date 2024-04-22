package controllers

import (
	middleware "SessionCookie/Middleware"
	repository "SessionCookie/Repository"
	service "SessionCookie/Service"
	"SessionCookie/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	App *fiber.App
}

var userService *service.UserService

func InitializeUserController(DataBase *gorm.DB, App *fiber.App) {

	c := &UserController{}
	CurrentRepository := &repository.UserRepository{DB: DataBase}
	userService = &service.UserService{Rep: CurrentRepository}
	middleware.Rep = CurrentRepository
	c.App = App
	c.setRoutes()
}

func (c *UserController) setRoutes() {
	app := c.App
	app.Get("/", c.list)
	app.Post("/login", c.login)
	app.Get("/validate", middleware.ValidateSession, c.validate)
	app.Post("/logout", middleware.ValidateSession, c.logout)
}

func (c *UserController) list(ctx *fiber.Ctx) error {

	users, err := userService.ListAllUsers()
	if err != nil {
		return ctx.Status(http.StatusBadGateway).JSON(
			&fiber.Map{
				"message": "Error fetching the users",
			},
		)
	}

	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Users fetched successfully",
			"users":   users,
		},
	)

	return nil
}

func (c *UserController) login(ctx *fiber.Ctx) error {

	user := &models.User{}
	err := ctx.BodyParser(user)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "Unable to parse user",
			},
		)
	}

	err = userService.AuthenticateUser(user)
	if err == fiber.ErrNotFound {
		return ctx.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message": "User is not registered",
			},
		)
	} else if err == fiber.ErrUnauthorized {
		return ctx.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message": "User is not registered",
			},
		)
	}

	token, err := userService.CreateJWT(user.UserName)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "error creating token",
			},
		)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Minute * 5),
		HTTPOnly: false,
		SameSite: "strict",
		Secure:   true,
	})

	ctx.Locals("user", user)

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Logged in successfully",
	})

	return nil
}

func (c *UserController) validate(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(*models.User)
	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"user":    user,
			"message": "You are logged",
		},
	)

	return nil
}

func (c *UserController) logout(ctx *fiber.Ctx) error {

	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: false,
		SameSite: "strict",
		Secure:   true,
	})

	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Logged out successfully",
		},
	)

	return nil
}
