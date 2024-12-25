const apiUrl = "http://localhost:8080/products";

function displayResponse(targetId, responseText) {
    document.getElementById(targetId).innerText = JSON.stringify(responseText, null, 2);
}

// Create Product
function createProduct() {
    const name = document.getElementById('create-name').value.trim();
    const description = document.getElementById('create-description').value.trim();
    const price = parseFloat(document.getElementById('create-price').value);
    const quantity = parseInt(document.getElementById('create-quantity').value);

    // Проверим типы данных
    console.log({
        name,
        description,
        price,
        quantity
    });

    const productData = {
        name,
        description,
        price,
        quantity
    };

    fetch("http://localhost:8080/products/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(productData)
    })
        .then(response => response.json())
        .then(data => {
            displayResponse("server-response", data);
            console.log(data); // Просмотр ответа на запрос
        })
        .catch(error => {
            displayResponse("server-response", error);
            console.error(error);
        });
}

// Get Product By ID
function getProductById() {
    const id = document.getElementById('get-id').value;

    fetch(apiUrl + "/" + id)
        .then(response => response.json())
        .then(data => displayResponse("product-response", data))
        .catch(error => displayResponse("product-response", error));
}

// Update Product
function updateProduct() {
    const id = document.getElementById('update-id').value; // Получаем ID из формы
    const name = document.getElementById('update-name').value.trim();
    const description = document.getElementById('update-description').value.trim();
    const price = parseFloat(document.getElementById('update-price').value);
    const quantity = parseInt(document.getElementById('update-quantity').value);

    // Создаем объект данных для отправки
    const updatedData = {
        name,
        description,
        price,
        quantity
    };

    console.log(updatedData); // Печатаем данные в консоль для проверки

    // Отправляем PUT запрос
    fetch(`http://localhost:8080/products/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedData)
    })
        .then(response => response.json())
        .then(data => {
            displayResponse("server-response", data); // Показываем ответ на экране
            console.log(data); // Логируем данные в консоль
        })
        .catch(error => {
            displayResponse("server-response", error); // Показываем ошибку на экране
            console.error(error); // Логируем ошибку в консоль
        });
}
// Get All Products
function getAllProducts() {
    fetch(apiUrl + "/")
        .then(response => response.json())
        .then(data => displayResponse("all-products-response", data))
        .catch(error => displayResponse("all-products-response", error));
}

// Delete Product By ID
function deleteProduct() {
    const id = document.getElementById('delete-id').value;

    fetch(apiUrl + "/" + id, {
        method: "DELETE"
    })
        .then(response => response.json())
        .then(data => displayResponse("server-response", data))
        .catch(error => displayResponse("server-response", error));
}

// Delete All Products
function deleteAllProducts() {
    fetch("http://localhost:8080/products/deleteAll", {
        method: "DELETE"
    })
        .then(response => response.json())
        .then(data => {
            displayResponse("server-response", data); // Отобразить ответ на фронте
            console.log(data); // Логируем данные в консоль
        })
        .catch(error => {
            displayResponse("server-response", error); // Отобразить ошибку
            console.error(error); // Логировать ошибку в консоль
        });
}
