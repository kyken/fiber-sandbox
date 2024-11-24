package user

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kyken/fiber-sandbox/lib/db"
	"github.com/kyken/fiber-sandbox/lib/model"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	UserName     string    `db:"username" json:"username"` // タグを修正
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserHandler struct {
	db db.Database
}

func NewUserHandler(db db.Database) *UserHandler {
	return &UserHandler{db: db}
}

// GET /users
func (h *UserHandler) GetUsersHandler(c *fiber.Ctx) error {
	users := make([]User, 0)

	findAllSql := "select id, username, email, created_at from users"

	err := h.db.Select(&users, findAllSql)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(users)
}

// GET /user/id
func (h *UserHandler) GetUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(model.ErrorResponse{
			Message: "Invalid user ID",
			Status:  400,
		})
	}

	user := User{}

	findSql := "select id, username, email, created_at from users where id = ?"

	err = h.db.Get(&user, findSql, id)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(user)
}

// PUT /user
func (h *UserHandler) PutUserHandler(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	insertSql := "insert into users (username, email, password_hash) values (?, ?, ?)"

	result, err := h.db.ExecContext(c.Context(), insertSql, user.UserName, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	// 挿入されたレコードのIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get inserted ID")
	}
	return c.Status(200).JSON(id)
}

// POST /user/:id
func (h *UserHandler) PostUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(model.ErrorResponse{
			Message: "Invalid user ID",
			Status:  400,
		})
	}
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	updateSql := "update users set username = ?, email = ?, password_hash = ? where id = ?"
	_, err = h.db.ExecContext(c.Context(), updateSql, user.UserName, user.Email, user.PasswordHash, id)
	if err != nil {
		return err
	}
	return c.Status(200).JSON("ok")
}

// DELETE /user/:id
func (h *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(model.ErrorResponse{
			Message: "Invalid user ID",
			Status:  400,
		})
	}
	ctx := context.Background()

	deleteSql := "delete from users where id = ?"
	_, err = h.db.ExecContext(ctx, deleteSql, id)
	if err != nil {
		return err
	}
	return c.Status(200).JSON("ok")
}
