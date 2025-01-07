const apiUrl = "http://localhost:8080/products";

function displayResponse(targetId, responseText) {
    const targetElement = document.getElementById(targetId);
    if (targetElement) {
        targetElement.innerText = JSON.stringify(responseText, null, 2);
    } else {
        console.error(`Element with ID '${targetId}' not found.`);
    }
}

function createProduct() {
    const productData = {
        type: document.getElementById("product-type").value,
        brand: document.getElementById("create-brand").value.trim(),
        model: document.getElementById("create-model").value.trim(),
        specifications: document.getElementById("create-specifications").value.trim(),
        color: document.getElementById("create-color").value.trim(),
        price: parseFloat(document.getElementById("create-price").value),
        quantity: parseInt(document.getElementById("create-quantity").value),
        warranty: parseInt(document.getElementById("create-warranty").value),
    };

    fetch(apiUrl + "/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(productData),
    })
        .then((response) => response.json())
        .then((data) => displayResponse("server-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("server-response", error);
        });
}

function getProductById() {
    const productId = document.getElementById("get-id").value;

    // Проверяем, является ли ID корректным
    if (!productId || productId.length !== 24) {
        document.getElementById("product-response").innerText =
            "Invalid Product ID. Please enter a valid 24-character ObjectID.";
        return;
    }

    // Замените относительный URL на полный URL API
    const apiUrl = `http://localhost:8080/products/${productId}`;

    fetch(apiUrl)
        .then((response) => {
            // Проверяем успешность статуса (200-299)
            if (!response.ok) {
                throw new Error("Product not found. Status: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            document.getElementById("product-response").innerText = JSON.stringify(data, null, 2);
        })
        .catch((error) => {
            document.getElementById("product-response").innerText = `Error: ${error.message}`;
            console.error(error);
        });
}
function updateProduct() {
    const id = document.getElementById("update-id").value.trim();

    // Проверяем, введен ли корректный ID
    if (!id || id.length !== 24) {
        displayResponse("server-response", { error: "Invalid Product ID. Please enter a valid 24-character ObjectID." });
        return;
    }

    const updatedData = {
        type: document.getElementById("update-type").value,
        brand: document.getElementById("update-brand").value.trim() || null,
        model: document.getElementById("update-model").value.trim() || null,
        specifications: document.getElementById("update-specifications").value.trim() || null,
        color: document.getElementById("update-color").value.trim() || null,
        price: parseFloat(document.getElementById("update-price").value) || null,
        quantity: parseInt(document.getElementById("update-quantity").value) || null,
        warranty: parseInt(document.getElementById("update-warranty").value) || null,
    };

    // Проверки на неотправку пустого объекта
    if (Object.values(updatedData).every((value) => value === null)) {
        displayResponse("server-response", { error: "Please fill at least one field to update." });
        return;
    }

    // URL для API
    const apiUrl = "http://localhost:8080/products";

    fetch(`${apiUrl}/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedData),
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Error: ${response.status} - ${response.statusText}`);
            }
            return response.json();
        })
        .then((data) => {
            displayResponse("server-response", data); // Обрабатываем успешный ответ
        })
        .catch((error) => {
            console.error(error);
            displayResponse("server-response", { error: error.message }); // Выводим сообщение об ошибке
        });
}

function getAllProducts() {
    fetch(apiUrl + "/")
        .then((response) => response.json())
        .then((data) => displayResponse("all-products-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("all-products-response", error);
        });
}

function deleteProduct() {
    const id = document.getElementById("delete-id").value.trim();

    // Проверяем, является ли ID корректным
    if (!id || id.length !== 24) {
        document.getElementById("delete-response").innerText =
            "Invalid Product ID. Please enter a valid 24-character ObjectID.";
        return;
    }

    // URL для API
    const apiUrl = "http://localhost:8080/products";

    fetch(`${apiUrl}/${id}`, {
        method: "DELETE",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Error: ${response.status} - ${response.statusText}`);
            }
            return response.json();
        })
        .then((data) => {
            document.getElementById("delete-response").innerText = JSON.stringify(data, null, 2);
        })
        .catch((error) => {
            console.error(error);
            document.getElementById("delete-response").innerText = `Error: ${error.message}`;
        });
}

function deleteAllProducts() {
    fetch(apiUrl + "/deleteAll", {
        method: "DELETE",
    })
        .then((response) => response.json())
        .then((data) => displayResponse("server-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("server-response", error);
        });
}