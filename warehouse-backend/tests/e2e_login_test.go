package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	seleniumPath = "C:\\Windows\\System32\\chromedriver.exe" // Путь к ChromeDriver (проверь свой!)
	port         = 9515                                      // Порт для WebDriver
	baseURL      = "http://localhost:8080"                   // URL твоего фронтенда
	loginPage    = "/LogIn.html"                             // Страница логина
	email        = "231441@astanait.edu.kz"                  // Тестовый email
	password     = "1111"                                    // Тестовый пароль
)

func TestE2ELogin(t *testing.T) {
	// Настройки браузера
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{}
	caps.AddChrome(chromeCaps)

	// Подключение к WebDriver
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		t.Fatalf("Ошибка подключения к WebDriver: %v", err)
	}
	defer wd.Quit()

	// Открытие страницы логина
	err = wd.Get("http://localhost:63342/Warehouse_Management_System/warehouse-frontend/login.html")
	if err != nil {
		t.Fatalf("Ошибка загрузки страницы логина: %v", err)
	}
	time.Sleep(2 * time.Second)
	// Поиск и ввод email
	emailInput, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		t.Fatalf("Не найдено поле email: %v", err)
	}
	emailInput.SendKeys(email)

	// Поиск и ввод пароля
	passwordInput, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatalf("Не найдено поле пароля: %v", err)
	}
	passwordInput.SendKeys(password)

	// Нажатие на кнопку логина
	loginButton, err := wd.FindElement(selenium.ByID, "submit")
	if err != nil {
		t.Fatalf("Не найдена кнопка входа: %v", err)
	}
	loginButton.Click()

	// Ожидание редиректа (до 5 секунд)
	time.Sleep(5 * time.Second)

	// Проверка, что мы на странице профиля
	currentURL, err := wd.CurrentURL()
	if err != nil {
		t.Fatalf("Ошибка получения URL: %v", err)
	}
	if currentURL != baseURL+"/profile.html" {
		t.Errorf("Ожидался переход на профиль, но текущий URL: %s", currentURL)
	}
}
