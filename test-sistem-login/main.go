package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"test-sistem-login/models"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()
var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Tidak dapat terhubung ke redis: %v", err)
	}
	fmt.Println("Berhasil terhubung ke redis:", pong)
}

func main() {
	initRedis()

	app := fiber.New()

	app.Post("/register", registerUser)
	app.Post("/login", loginUser)

	log.Fatal(app.Listen(":3000"))
}

func registerUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Permintaan tidak valid"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal hashing password"})
	}
	user.Password = string(hashedPassword)

	// simpan data ke redis
	userKey := fmt.Sprintf("login_%s", user.Email)
	userData, _ := json.Marshal(user)
	if err := redisClient.Set(ctx, userKey, userData, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan data ke redis"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Register berhasil, silahkan login"})
}

func loginUser(c *fiber.Ctx) error {
	payload := new(models.LoginPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Permintaan tidak valid"})
	}

	// ambil data user dari redis
	userKey := fmt.Sprintf("login_%s", payload.Email)
	userData, err := redisClient.Get(ctx, userKey).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau password salah"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data dari redis"})
	}

	user := new(models.User)
	if err := json.Unmarshal([]byte(userData), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal parsing data user"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau password salah"})
	}

	// login sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"response": fiber.Map{
			"realname": user.RealName,
		},
	})
}
