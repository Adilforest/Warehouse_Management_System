const apiUrl = "http://localhost:8080/products";

// Toggle visibility of extra fields based on product type
function toggleFields(context) {
    const type = document.getElementById(`${context}-type`).value;
    const extraFields = document.getElementById(`extra-fields-${context}`);
    if (type === 'headphones') {
        extraFields.style.display = 'none';  // Hide for headphones
    } else {
        extraFields.style.display = '';  // Show for laptops and smartphones
    }
}

// Initialize visibility on page load
document.addEventListener('DOMContentLoaded', () => {
    toggleFields('create');
    toggleFields('update');
});

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
        processor: document.getElementById("create-processor").value.trim(), // добавлено
        ram: document.getElementById("create-ram").value.trim(), // добавлено
        storage: document.getElementById("create-storage").value.trim(), // добавлено
        color: document.getElementById("create-color").value.trim(),
        price: parseFloat(document.getElementById("create-price").value),
        quantity: parseInt(document.getElementById("create-quantity").value),
        warranty: parseInt(document.getElementById("create-warranty").value),
        description: document.getElementById("create-description")?.value.trim() || "", // добавлено описание
    };

    fetch(`${apiUrl}/create`, {
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

    if (!productId || productId.length !== 24) {
        document.getElementById("product-response").innerText =
            "Invalid Product ID. Please enter a valid 24-character ObjectID.";
        return;
    }

    fetch(`${apiUrl}/${productId}`)
        .then((response) => {
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

    // Проверяем корректность ID
    if (!id || id.length !== 24) {
        displayResponse("server-response", { error: "Invalid Product ID. Please enter a valid 24-character ObjectID." });
        return;
    }

    const updatedData = {
        type: document.getElementById("update-type").value,
        brand: document.getElementById("update-brand").value.trim() || null,
        model: document.getElementById("update-model").value.trim() || null,
        processor: document.getElementById("update-processor").value.trim() || null, // новое поле
        ram: document.getElementById("update-ram").value.trim() || null, // новое поле
        storage: document.getElementById("update-storage").value.trim() || null, // новое поле
        color: document.getElementById("update-color").value.trim() || null,
        price: parseFloat(document.getElementById("update-price").value) || null,
        quantity: parseInt(document.getElementById("update-quantity").value) || null,
        warranty: parseInt(document.getElementById("update-warranty").value) || null,
        description: document.getElementById("update-description")?.value.trim() || null, // описание
    };

    // Проверяем, что хотя бы одно поле для изменения заполнено
    if (Object.values(updatedData).every((value) => value === null)) {
        displayResponse("server-response", { error: "Please fill at least one field to update." });
        return;
    }

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

    if (!id || id.length !== 24) {
        document.getElementById("delete-response").innerText =
            "Invalid Product ID. Please enter a valid 24-character ObjectID.";
        return;
    }

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
    fetch(`${apiUrl}/deleteAll`, {
        method: "DELETE",
    })
        .then((response) => response.json())
        .then((data) => displayResponse("server-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("server-response", error);
        });
}