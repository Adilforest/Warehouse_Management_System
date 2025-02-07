document.addEventListener("DOMContentLoaded", () => {
    // Обработка формы входа (Log In)
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
                    // Используем result.data.token, так как токен находится внутри data
                    localStorage.setItem("token", result.data.token);
                    localStorage.setItem("email", email); // Сохраняем email
                    // Сохраняем роль, возвращённую сервером (например, "admin" или "user")
                    localStorage.setItem("role", result.data.role);

                    // Выводим сообщение об успехе
                    displayLoginResponse("Login successful!");

                    // Перенаправляем пользователя в зависимости от его роли
                    setTimeout(() => {
                        if (result.data.role === "admin") {
                            window.location.href = "admin.html"; // Перенаправление для админа
                        } else {
                            window.location.href = "index.html"; // Перенаправление для обычного пользователя
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

    // Обработка формы регистрации (Sign Up)
    const signupForm = document.getElementById("signup-form");
    if (signupForm) {
        signupForm.addEventListener("submit", async (e) => {
            e.preventDefault(); // Предотвращаем стандартное поведение формы

            // Получаем значения name, email, password и confirm-password из формы
            const name = document.getElementById("name").value.trim();
            const email = document.getElementById("email").value.trim();
            const password = document.getElementById("password").value.trim();
            const confirmPassword = document.getElementById("confirm-password").value.trim();

            // Проверяем, что все поля заполнены
            if (!name || !email || !password || !confirmPassword) {
                displaySignupResponse("Please fill in all fields.");
                return;
            }

            // Проверяем совпадение паролей
            if (password !== confirmPassword) {
                displaySignupResponse("Passwords do not match.");
                return;
            }

            try {
                // Отправляем запрос на сервер для регистрации
                const response = await fetch("http://localhost:8080/auth/signup", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ name, email, password, confirm_password: confirmPassword }),
                });

                // Парсим ответ от сервера
                const result = await response.json();

                // Обрабатываем успешный ответ
                if (response.ok) {
                    displaySignupResponse("Registration successful! Please check your email to verify your account.");
                } else {
                    // Обрабатываем ошибку
                    displaySignupResponse(`Error: ${result.error}`);
                }
            } catch (error) {
                // Обработка сетевых ошибок
                console.error("Network error:", error);
                displaySignupResponse("An unexpected error occurred. Please try again later.");
            }
        });
    }
});

// Функция для отображения сообщений на странице входа
function displayLoginResponse(message) {
    const loginResponse = document.getElementById("login-response");
    if (loginResponse) {
        loginResponse.textContent = message;
    } else {
        console.error("Element with ID 'login-response' not found.");
    }
}

// Функция для отображения сообщений на странице регистрации
function displaySignupResponse(message) {
    const signupResponse = document.getElementById("signup-response");
    if (signupResponse) {
        signupResponse.textContent = message;
    } else {
        console.error("Element with ID 'signup-response' not found.");
    }
}
