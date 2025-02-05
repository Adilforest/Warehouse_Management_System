document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");

    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;

            const response = await fetch("http://localhost:8080/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            });

            const result = await response.json();
            const loginResponse = document.getElementById("login-response");

            if (response.ok) {
                localStorage.setItem("email", result.email); // Сохраняем email
                loginResponse.textContent = "Login successful!";
                window.location.href = "index.html"; // Перенаправляем на главную страницу
            } else {
                loginResponse.textContent = `Error: ${result.error}`;
            }
        });
    }
});