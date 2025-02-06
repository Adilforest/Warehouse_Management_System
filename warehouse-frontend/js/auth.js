document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault(); // Предотвращаем стандартное поведение формы

            // Получаем значения email и password из формы
            const email = document.getElementById("email").value.trim();
            const password = document.getElementById("password").value.trim();

            // Проверяем, что поля не пустые
            if (!email || !password) {
                displayLoginResponse("Please fill in both email and password fields.");
                return;
            }

            try {
                // Отправляем запрос на сервер для аутентификации
                const response = await fetch("http://localhost:8080/auth/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password }),
                });

                // Парсим ответ от сервера
                const result = await response.json();

                // Обрабатываем успешный ответ
                if (response.ok) {
                    // Сохраняем токен и данные пользователя в localStorage
                    localStorage.setItem("token", result.token);
                    localStorage.setItem("email", email); // Сохраняем email
                    localStorage.setItem("isAdmin", email === "231441@astanait.edu.kz"); // Определяем роль администратора

                    // Выводим сообщение об успехе
                    displayLoginResponse("Login successful!");

                    // Перенаправляем пользователя в зависимости от роли
                    setTimeout(() => {
                        if (email === "231441@astanait.edu.kz") {
                            window.location.href = "admin.html"; // Перенаправление на админку
                        } else {
                            window.location.href = "index.html"; // Перенаправление на главную
                        }
                    }, 1000); // Задержка для отображения сообщения
                } else {
                    // Обрабатываем ошибку
                    displayLoginResponse(`Error: ${result.error}`);
                }
            } catch (error) {
                // Обработка сетевых ошибок
                console.error("Network error:", error);
                displayLoginResponse("An unexpected error occurred. Please try again later.");
            }
        });
    }
});

// Функция для отображения сообщений на странице
function displayLoginResponse(message) {
    const loginResponse = document.getElementById("login-response");
    if (loginResponse) {
        loginResponse.textContent = message;
    } else {
        console.error("Element with ID 'login-response' not found.");
    }
}