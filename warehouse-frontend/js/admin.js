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
    const id = document.getElementById("get-id").value;

    fetch(apiUrl + "/" + id)
        .then((response) => response.json())
        .then((data) => displayResponse("product-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("product-response", error);
        });
}

function updateProduct() {
    const updatedData = {
        type: document.getElementById("update-type").value,
        brand: document.getElementById("update-brand").value.trim(),
        model: document.getElementById("update-model").value.trim(),
        specifications: document.getElementById("update-specifications").value.trim(),
        color: document.getElementById("update-color").value.trim(),
        price: parseFloat(document.getElementById("update-price").value),
        quantity: parseInt(document.getElementById("update-quantity").value),
        warranty: parseInt(document.getElementById("update-warranty").value),
    };

    const id = document.getElementById("update-id").value;
    fetch(apiUrl + "/" + id, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedData),
    })
        .then((response) => response.json())
        .then((data) => displayResponse("server-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("server-response", error);
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
    const id = document.getElementById("delete-id").value;
    fetch(apiUrl + "/" + id, {
        method: "DELETE",
    })
        .then((response) => response.json())
        .then((data) => displayResponse("server-response", data))
        .catch((error) => {
            console.error(error);
            displayResponse("server-response", error);
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