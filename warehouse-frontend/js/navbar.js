document.addEventListener("DOMContentLoaded", async () => {
    // Получаем email и роль из localStorage
    const email = localStorage.getItem("email");
    const role = localStorage.getItem("role"); // Добавляем роль пользователя
    console.log("Email from localStorage:", email);
    console.log("Role from localStorage:", role);

    // Находим элементы навигации
    const loginLink = document.getElementById("login-link");
    const signupLink = document.getElementById("signup-link");
    const logoutLink = document.getElementById("logout-link");
    const adminLink = document.getElementById("admin-link");

    // Проверяем, существуют ли все необходимые элементы
    if (!loginLink || !signupLink || !logoutLink || !adminLink) {
        console.error("One or more navigation elements are missing.");
        return;
    }

    if (email) {
        console.log("User is logged in with email:", email);

        // Если пользователь авторизован
        loginLink.style.display = "none";
        signupLink.style.display = "none";
        logoutLink.style.display = "inline";

        // Проверяем роль пользователя
        if (role === "admin") {
            console.log("User is an admin. Showing Admin Panel link.");
            adminLink.style.display = "inline"; // Показываем ссылку на Admin Panel
        } else {
            console.log("User is not an admin. Hiding Admin Panel link.");
            adminLink.style.display = "none"; // Скрываем ссылку, если пользователь не админ
        }
    } else {
        console.log("User is not logged in.");

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
        localStorage.removeItem("token"); // Удаляем токен
        localStorage.removeItem("role");  // Удаляем роль
        window.location.href = "LogIn.html"; // Перенаправляем на страницу входа
    });
});