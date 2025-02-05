document.addEventListener("DOMContentLoaded", async () => {
    const email = localStorage.getItem("email");
    const loginLink = document.getElementById("login-link");
    const signupLink = document.getElementById("signup-link");
    const logoutLink = document.getElementById("logout-link");
    const adminLink = document.getElementById("admin-link");

    if (email) {
        // Если пользователь авторизован
        loginLink.style.display = "none";
        signupLink.style.display = "none";
        logoutLink.style.display = "inline";

        // Проверяем, является ли пользователь администратором
        if (email === "231441@astanait.edu.kz") {
            adminLink.style.display = "inline"; // Показываем ссылку на Admin Panel
        }
    } else {
        // Если пользователь не авторизован
        loginLink.style.display = "inline";
        signupLink.style.display = "inline";
        logoutLink.style.display = "none";
        adminLink.style.display = "none";
    }

    // Обработка выхода из системы
    logoutLink.addEventListener("click", (e) => {
        e.preventDefault();
        localStorage.removeItem("email"); // Удаляем email
        window.location.href = "LogIn.html"; // Перенаправляем на страницу входа
    });
});