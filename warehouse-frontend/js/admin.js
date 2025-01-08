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

function validateCreateForm() {
    // Пример логики валидации
    const brand = document.getElementById("create-brand").value.trim();
    const model = document.getElementById("create-model").value.trim();
    const price = document.getElementById("create-price").value.trim();
    const quantity = document.getElementById("create-quantity").value.trim();
    const link = document.getElementById("create-link").value.trim();

    // Проверка обязательных полей
    if (!brand) {
        displayResponse("server-response", { error: "Поле 'Бренд' обязательно для заполнения." });
        return false;
    }

    if (!model) {
        displayResponse("server-response", { error: "Поле 'Модель' обязательно для заполнения." });
        return false;
    }

    if (!price || isNaN(price)) {
        displayResponse("server-response", { error: "Цена должна быть числом." });
        return false;
    }

    if (!quantity || isNaN(quantity)) {
        displayResponse("server-response", { error: "Количество должно быть числом." });
        return false;
    }

    if (!link || !isValidURL(link)) {
        displayResponse("server-response", { error: "Укажите корректную ссылку на фото." });
        return false;
    }

    // Если все проверки пройдены
    return true;
}

// Функция для проверки корректности URL
function isValidURL(url) {
    try {
        new URL(url);
        return true;
    } catch (error) {
        return false;
    }
}

function createProduct() {
    if (!validateCreateForm()) {
        return; // Остановка, если валидация не пройдена
    }

    const linkField = document.getElementById("create-link").value.trim();

    if (!linkField || !isValidURL(linkField)) {
        highlightErrorField("create-link");
        displayResponse("server-response", { error: "Please provide a valid URL for the photo." });
        return;
    }

    const productData = {
        type: document.getElementById("create-type").value,
        brand: document.getElementById("create-brand").value.trim(),
        model: document.getElementById("create-model").value.trim(),
        color: document.getElementById("create-color").value.trim(),
        price: parseFloat(document.getElementById("create-price").value),
        quantity: parseInt(document.getElementById("create-quantity").value),
        warranty: parseInt(document.getElementById("create-warranty").value) || 0,
        link: linkField,
        description: document.getElementById("create-description")?.value.trim() || "",
    };

    // Необязательные поля
    const optionalFields = {
        processor: document.getElementById("create-processor").value.trim(),
        ram: document.getElementById("create-ram").value.trim(),
        storage: document.getElementById("create-storage").value.trim(),
    };

    for (const [key, value] of Object.entries(optionalFields)) {
        if (value) {
            productData[key] = value;
        }
    }

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
            displayResponse("server-response", { error: "Something went wrong while creating the product." });
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
        processor: document.getElementById("update-processor").value.trim() || null,
        ram: document.getElementById("update-ram").value.trim() || null,
        storage: document.getElementById("update-storage").value.trim() || null,
        color: document.getElementById("update-color").value.trim() || null,
        price: parseFloat(document.getElementById("update-price").value) || null,
        quantity: parseInt(document.getElementById("update-quantity").value) || null,
        warranty: parseInt(document.getElementById("update-warranty").value) || null,
        description: document.getElementById("update-description")?.value.trim() || null,
        link: document.getElementById("update-link").value.trim() || null, // новое поле
    };

    // Удаляем пустые поля (равные null)
    for (const key in updatedData) {
        if (updatedData[key] === null || updatedData[key] === "") {
            delete updatedData[key];
        }
    }

    // Проверяем, что хотя бы одно поле для изменения заполнено
    if (Object.keys(updatedData).length === 0) {
        displayResponse("server-response", { error: "Please fill at least one field to update." });
        return;
    }

    // Отправляем запрос на сервер
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