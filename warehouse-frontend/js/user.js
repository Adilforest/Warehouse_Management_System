function createUser() {
    const name = document.getElementById("create-name").value.trim();
    const email = document.getElementById("create-email").value.trim();
    const password = document.getElementById("create-password").value.trim();

    fetch("http://localhost:8080/users/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email, password }),
    })
        .then((response) => response.json())
        .then((data) => {
            document.getElementById("create-user-response").innerText = JSON.stringify(data, null, 2);
        })
        .catch((error) => console.error(error));
}

document.getElementById("create-user-form").addEventListener("submit", (e) => {
    e.preventDefault();
    createUser();
});

function getAllUsers() {
    fetch("http://localhost:8080/users/", {
        method: "GET",
        headers: {Authorization: `Bearer ${localStorage.getItem("token")}`},
    })
        .then((response) => response.json())
        .then((data) => {
            document.getElementById("all-users-response").innerText = JSON.stringify(data, null, 2);
        })
        .catch((error) => console.error(error));
}
function getUser() {
    const id = document.getElementById("get-id").value.trim();
    fetch(`http://localhost:8080/users/${id}`, {
        method: "GET",
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    })
        .then((response) => response.json())
        .then((data) => displayResponse("get-user-response", data))
        .catch((error) => console.error(error));
}

function updateUser() {
    const id = document.getElementById("update-id").value.trim();
    const name = document.getElementById("update-name").value.trim();
    const email = document.getElementById("update-email").value.trim();

    fetch(`http://localhost:8080/users/${id}`, {
        method: "PUT",
        headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, email }),
    })
        .then((response) => response.json())
        .then((data) => displayResponse("update-user-response", data))
        .catch((error) => console.error(error));
}

function deleteUser() {
    const id = document.getElementById("delete-id").value.trim();
    fetch(`http://localhost:8080/users/${id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    })
        .then((response) => response.json())
        .then((data) => displayResponse("delete-user-response", data))
        .catch((error) => console.error(error));
}

function deleteAllUsers() {
    fetch("http://localhost:8080/users/deleteAll", {
        method: "DELETE",
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    })
        .then((response) => response.json())
        .then((data) => displayResponse("delete-all-users-response", data))
        .catch((error) => console.error(error));
}

document.getElementById("get-user-form").addEventListener("submit", (e) => {
    e.preventDefault();
    getUser();
});

document.getElementById("update-user-form").addEventListener("submit", (e) => {
    e.preventDefault();
    updateUser();
});

document.getElementById("delete-user-form").addEventListener("submit", (e) => {
    e.preventDefault();
    deleteUser();
});
