document.addEventListener("DOMContentLoaded", async () => {
    const token = localStorage.getItem("token");
    if (!token) {
        alert("You need to log in first.");
        window.location.href = "LogIn.html";
        return;
    }

    const updateProfileForm = document.getElementById("update-profile-form");

    // Загрузка текущих данных пользователя
    try {
        const response = await fetch("http://localhost:8080/protected/profile", {
            method: "GET",
            headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
            alert("Failed to load profile data.");
            window.location.href = "LogIn.html";
            return;
        }

        const user = await response.json();
        document.getElementById("name").value = user.name;
        document.getElementById("email").value = user.email;
    } catch (error) {
        console.error("Error loading profile:", error);
        alert("An error occurred while loading your profile.");
        window.location.href = "LogIn.html";
    }

    // Обработка обновления профиля
    updateProfileForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const name = document.getElementById("name").value;
        const email = document.getElementById("email").value;
        const currentPassword = document.getElementById("current-password").value;
        const newPassword = document.getElementById("new-password").value;
        const confirmNewPassword = document.getElementById("confirm-new-password").value;

        if (newPassword !== confirmNewPassword) {
            document.getElementById("profile-response").textContent =
                "New passwords do not match.";
            return;
        }

        const requestBody = {
            name,
            email,
            current_password: currentPassword,
            new_password: newPassword || undefined, // Не отправляем пустой пароль
        };

        try {
            const response = await fetch("http://localhost:8080/protected/profile", {
                method: "PUT",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(requestBody),
            });

            const result = await response.json();
            const profileResponse = document.getElementById("profile-response");

            if (response.ok) {
                profileResponse.textContent = "Profile updated successfully!";
            } else {
                profileResponse.textContent = `Error: ${result.error}`;
            }
        } catch (error) {
            console.error("Error updating profile:", error);
            document.getElementById("profile-response").textContent =
                "An error occurred while updating your profile.";
        }
    });
});