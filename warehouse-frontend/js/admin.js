const apiUrl = "http://localhost:8080/products";

function toggleFields(context) {
    // Получаем суффикс: "create" или "update"
    const prefix = context || "create"; // Если не передан, по умолчанию "create"

    // Элементы для управления
    const productType = document.getElementById(`${prefix}-type`);
    const processorField = document.getElementById(`${prefix}-processor`);
    const ramField = document.getElementById(`${prefix}-ram`);
    const storageField = document.getElementById(`${prefix}-storage`);
    const extraFields = document.getElementById(`extra-fields-${prefix}`);

    // Проверяем выбранный тип продукта
    if (productType && productType.value === "headphones") {
        // Скрыть лишние поля для "Headphones"
        if (processorField) processorField.parentElement.style.display = "none";
        if (ramField) ramField.parentElement.style.display = "none";
        if (storageField) storageField.parentElement.style.display = "none";
        if (extraFields) extraFields.style.display = "none";
    } else {
        // Показать поля для остальных типов продуктов
        if (processorField) processorField.parentElement.style.display = "block";
        if (ramField) ramField.parentElement.style.display = "block";
        if (storageField) storageField.parentElement.style.display = "block";
        if (extraFields) extraFields.style.display = "block";
    }
}

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
        type: document.getElementById("create-type").value,
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