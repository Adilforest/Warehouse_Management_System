package controllers

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
	"time"
	"warehouse-backend/database"
	"warehouse-backend/middleware"
	"warehouse-backend/models"
)

// CreateUser создает нового пользователя в базе данных
func CreateUser(name, email, password string) (models.User, error) {
	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.New("failed to hash password: " + err.Error())
	}

	// Создаем пользователя
	user := models.User{
		ID:                primitive.NewObjectID(),
		Name:              name,
		Email:             email,
		Password:          string(hashedPassword),
		Verified:          false,
		VerificationToken: generateVerificationToken(),
		Role:              "user", // По умолчанию роль пользователя — "user"
	}

	// Получаем коллекцию MongoDB
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, существует ли пользователь с таким email
	existingUser := collection.FindOne(ctx, bson.M{"email": email})
	if existingUser.Err() == nil {
		return models.User{}, errors.New("user with this email already exists")
	}

	// Вставляем нового пользователя
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, errors.New("failed to create user: " + err.Error())
	}

	// Отправляем email с токеном верификации
	err = sendVerificationEmail(email, user.VerificationToken)
	if err != nil {
		return models.User{}, errors.New("failed to send verification email: " + err.Error())
	}

	return user, nil
}

// generateVerificationToken генерирует случайный токен верификации
func generateVerificationToken() string {
	// Здесь можно использовать UUID или другой способ генерации уникального токена
	return primitive.NewObjectID().Hex()
}

// sendVerificationEmail отправляет email с ссылкой для верификации
func sendVerificationEmail(email, token string) error {
	verificationLink := "http://localhost:8080/auth/verify?token=" + token
	// Временно выводим ссылку в консоль (для тестирования)
	fmt.Println("Verification link:", verificationLink)

	// Настройка SMTP
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("failed to parse SMTP_PORT: %w", err)
	}
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Создаем новое сообщение
	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify your email")
	m.SetBody("text/html", fmt.Sprintf(`
        <h1>Thank you for registering!</h1>
        <p>Please verify your email by clicking the link below:</p>
        <a href="%s">Verify Email</a>
    `, verificationLink))

	// Настраиваем SMTP-клиент
	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

	// Отправляем письмо
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// LoginUser выполняет вход пользователя
func LoginUser(email, password string) (string, error) {
	// Находим пользователя по email
	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}
	// Проверяем, подтвержден ли email
	if !user.Verified {
		return "", errors.New("email not verified")
	}
	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}
	// Генерируем JWT-токен
	token, err := middleware.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	// Возвращаем токен и роль пользователя
	return token, nil
}
