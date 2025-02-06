document.addEventListener("DOMContentLoaded", async () => {
    const email = localStorage.getItem("email");
    console.log("Email from localStorage:", email); // Логируем значение email

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
        console.log("User is logged in with email:", email); // Логируем, что пользователь авторизован

        // Если пользователь авторизован
        loginLink.style.display = "none";
        signupLink.style.display = "none";
        logoutLink.style.display = "inline";

        // Проверяем, является ли пользователь администратором
        if (email === "231441@astanait.edu.kz") {
            console.log("User is an admin. Showing Admin Panel link."); // Логируем, что пользователь админ
            adminLink.style.display = "inline"; // Показываем ссылку на Admin Panel
        } else {
            console.log("User is not an admin. Hiding Admin Panel link."); // Логируем, что пользователь не админ
            adminLink.style.display = "none"; // Скрываем ссылку, если пользователь не админ
        }
    } else {
        console.log("User is not logged in."); // Логируем, что пользователь не авторизован

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
        window.location.href = "LogIn.html"; // Перенаправляем на страницу входа
    });
});