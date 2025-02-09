package controllers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"io"
	"net/http"
	"os"
	"warehouse-backend/logger"
)

// ContactController обрабатывает запросы от формы "Contact Us".
func ContactController(c *gin.Context) {
	// Логируем начало обработки запроса
	logger.LogOperationStart("ContactController", map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	})

	// Проверяем метод запроса (должен быть POST)
	if c.Request.Method != http.MethodPost {
		logger.Log.Error("ContactController", "Неподдерживаемый метод запроса", map[string]interface{}{
			"method": c.Request.Method,
		}, nil)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не поддерживается"})
		return
	}

	// Парсим multipart/form-data (максимальный размер файла — 10 MB)
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logger.Log.Error("ContactController", "Ошибка при парсинге формы", map[string]interface{}{
			"error": err.Error(),
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке формы"})
		return
	}

	// Извлекаем данные из формы
	name := c.PostForm("name")
	email := c.PostForm("email")
	message := c.PostForm("message")

	// Логируем полученные данные
	logger.LogInfo("ContactController", "Получены данные от пользователя", map[string]interface{}{
		"name":    name,
		"email":   email,
		"message": message,
	})

	// Обрабатываем прикрепленный файл (если есть)
	var filePath string
	file, handler, err := c.Request.FormFile("attachment")
	if err == nil {
		defer file.Close()

		// Создаем папку uploads, если она не существует
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			os.Mkdir("uploads", 0755)
		}

		// Сохраняем файл на сервере
		filePath = "uploads/" + handler.Filename
		f, err := os.Create(filePath)
		if err != nil {
			logger.Log.Error("ContactController", "Ошибка при создании файла", map[string]interface{}{
				"file_path": filePath,
				"error":     err.Error(),
			}, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
			return
		}
		defer f.Close()

		// Копируем содержимое файла
		_, err = io.Copy(f, file)
		if err != nil {
			logger.Log.Error("ContactController", "Ошибка при копировании файла", map[string]interface{}{
				"file_path": filePath,
				"error":     err.Error(),
			}, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
			return
		}

		// Логируем успешное сохранение файла
		logger.LogFileOperation("SaveFile", filePath, map[string]interface{}{
			"file_name": handler.Filename,
			"file_size": handler.Size,
		})
	}

	// Отправляем email
	err = SendEmail(name, email, message, filePath)
	if err != nil {
		logger.Log.Error("ContactController", "Ошибка при отправке email", map[string]interface{}{
			"sender_email": email,
			"error":        err.Error(),
		}, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке email"})
		return
	}

	// Логируем успешную отправку email
	logger.LogEmail("support@example.com", "Новое сообщение от "+name, map[string]interface{}{
		"sender_email": email,
		"message":      message,
	})

	// Логируем завершение обработки запроса
	logger.LogOperationEnd("ContactController", map[string]interface{}{
		"status": "success",
	})

	// Возвращаем успешный ответ клиенту
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SendEmail отправляет email с данными из формы.
func SendEmail(name, email, message, filePath string) error {
	// Создаем HTML-содержимое письма
	htmlContent := `
		<h1>Новое сообщение от ` + name + `</h1>
		<p><strong>Email:</strong> ` + email + `</p>
		<p><strong>Сообщение:</strong></p>
		<p>` + message + `</p>
	`

	// Создаем новое сообщение
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", os.Getenv("SMTP_RECEIVER"))
	m.SetHeader("Subject", "Новое сообщение от "+name)
	m.SetBody("text/html", htmlContent)

	// Прикрепляем файл, если он есть
	if filePath != "" {
		m.Attach(filePath)
	}

	// Настраиваем SMTP-клиент
	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// Отправляем письмо
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
