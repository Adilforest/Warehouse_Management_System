const userApiUrl = "http://localhost:8080/users"; // URL для работы с пользователями

// Функция для отображения ответа сервера
function displayUserResponse(targetId, responseText) {
    const targetElement = document.getElementById(targetId);
    if (targetElement) {
        targetElement.innerText = JSON.stringify(responseText, null, 2);
    } else {
        console.error(`Element with ID '${targetId}' not found.`);
    }
}

// Получение токена из localStorage
function getToken() {
    return localStorage.getItem("token");
}

// Проверка авторизации пользователя (админа)
async function checkAdminAccess() {
    const token = getToken();
    if (!token) {
        alert("You need to log in first.");
        window.location.href = "LogIn.html";
        return false;
    }

    try {
        const response = await fetch("http://localhost:8080/protected/profile", {
            method: "GET",
            headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Failed to load profile data.");
        }

        const user = await response.json();

        console.log("Server response:", user);

        // Проверяем, содержит ли ответ поле isAdmin
        if (typeof user.isAdmin === "undefined") {
            throw new Error("isAdmin is missing in the server response.");
        }

        // Проверяем роль пользователя
        if (!user.isAdmin) {
            alert("You do not have permission to access this page.");
            window.location.href = "index.html";
            return false;
        }

        return true; // Пользователь является администратором
    } catch (error) {
        console.error("Error checking admin access:", error);
        alert(`An error occurred: ${error.message}`);
        window.location.href = "LogIn.html";
        return false;
    }
}

// Создание нового пользователя через админ-панель
document.getElementById("create-user-form").addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!(await checkAdminAccess())) return;

    const name = document.getElementById("create-name").value.trim();
    const email = document.getElementById("create-email").value.trim();
    const password = document.getElementById("create-password").value.trim();
    const confirmPassword = document.getElementById("create-confirm-password").value.trim();

    if (!name || !email || !password || !confirmPassword) {
        displayUserResponse("create-user-response", { error: "All fields are required." });
        return;
    }

    if (password !== confirmPassword) {
        displayUserResponse("create-user-response", { error: "Passwords do not match." });
        return;
    }

    const userData = {
        name,
        email,
        password,
        confirm_password: confirmPassword  // если сервер ожидает именно такое имя поля
    };

    try {
        // Используем URL, который работает (как в cURL-запросе)
        const response = await fetch("http://localhost:8080/auth/signup", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${getToken()}`,
            },
            body: JSON.stringify(userData),
        });

        let data;
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
            data = await response.json();
        } else {
            data = await response.text();
        }

        if (!response.ok) {
            throw new Error(data.error || data || "Failed to create user.");
        }

        displayUserResponse("create-user-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("create-user-response", { error: error.message });
    }
});



// Получение всех пользователей
async function getAllUsers() {
    if (!(await checkAdminAccess())) return;

    try {
        const response = await fetch(`${userApiUrl}/`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || "Failed to fetch users.");
        }

        displayUserResponse("all-users-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("all-users-response", { error: error.message });
    }
}

// Получение пользователя по ID
document.getElementById("get-user-form").addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!(await checkAdminAccess())) return;

    const userId = document.getElementById("user-get-id").value.trim(); // Используем новый ID

    if (!userId || userId.length !== 24) {
        displayUserResponse("get-user-response", { error: "Invalid User ID. Please enter a valid 24-character ObjectID." });
        return;
    }

    try {
        const response = await fetch(`${userApiUrl}/${userId}`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || "User not found.");
        }

        displayUserResponse("get-user-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("get-user-response", { error: error.message });
    }
});

// Обновление пользователя по ID
document.getElementById("update-user-form").addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!(await checkAdminAccess())) return;

    const userId = document.getElementById("user-update-id").value.trim(); // Новый ID
    const name = document.getElementById("user-update-name").value.trim();
    const email = document.getElementById("user-update-email").value.trim();

    if (!userId || userId.length !== 24) {
        displayUserResponse("update-user-response", { error: "Invalid User ID. Please enter a valid 24-character ObjectID." });
        return;
    }

    const updatedData = {};
    if (name) updatedData.name = name;
    if (email) updatedData.email = email;

    if (Object.keys(updatedData).length === 0) {
        displayUserResponse("update-user-response", { error: "Please fill at least one field to update." });
        return;
    }

    try {
        const response = await fetch(`${userApiUrl}/${userId}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${getToken()}`,
            },
            body: JSON.stringify(updatedData),
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || "Failed to update user.");
        }

        displayUserResponse("update-user-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("update-user-response", { error: error.message });
    }
});

// Удаление пользователя по ID
document.getElementById("delete-user-form").addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!(await checkAdminAccess())) return;

    const userId = document.getElementById("user-delete-id").value.trim(); // Новый ID

    if (!userId || userId.length !== 24) {
        displayUserResponse("delete-user-response", { error: "Invalid User ID. Please enter a valid 24-character ObjectID." });
        return;
    }

    try {
        const response = await fetch(`${userApiUrl}/${userId}`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || "Failed to delete user.");
        }

        displayUserResponse("delete-user-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("delete-user-response", { error: error.message });
    }
});

// Удаление всех пользователей
async function deleteAllUsers() {
    if (!(await checkAdminAccess())) return;

    try {
        const response = await fetch(`${userApiUrl}/deleteAll`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || "Failed to delete all users.");
        }

        displayUserResponse("delete-all-users-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("delete-all-users-response", { error: error.message });
    }
}

// Изменение роли пользователя
document.getElementById("change-role-form").addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!(await checkAdminAccess())) return;

    const userId = document.getElementById("user-role-id").value.trim(); // Новый ID
    const role = document.getElementById("new-role").value;

    if (!userId || userId.length !== 24) {
        displayUserResponse("change-role-response", { error: "Invalid User ID. Please enter a valid 24-character ObjectID." });
        return;
    }

    try {
        const response = await fetch(`${userApiUrl}/${userId}/role`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${getToken()}`,
            },
            body: JSON.stringify({ role }),
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || "Failed to change user role.");
        }

        displayUserResponse("change-role-response", data);
    } catch (error) {
        console.error(error);
        displayUserResponse("change-role-response", { error: error.message });
    }
});
