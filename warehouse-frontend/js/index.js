const API_URL = 'http://localhost:8080/products/';

let products = [];
let currentPage = 1;
const itemsPerPage = 9;

let productList, pagination, sortBy;
let ramFilterGroup, storageFilterGroup, processorFilterGroup;

function initDOM() {
    productList = document.getElementById("product-list");
    pagination = document.getElementById("pagination");
    sortBy = document.getElementById("sort-by");
    ramFilterGroup = document.getElementById("ram-filter-group");
    storageFilterGroup = document.getElementById("storage-filter-group");
    processorFilterGroup = document.getElementById("processor-filter-group");
}

async function fetchProducts(filters = {}, sortBy = "", order = "asc", page = 1, minPrice = 0, maxPrice = 0) {
    const loadingIndicator = document.getElementById("loading");
    loadingIndicator.style.display = "block";

    // Удаляем пустые параметры из фильтров
    const cleanedFilters = {};
    for (const key in filters) {
        if (filters[key]) {
            cleanedFilters[key] = filters[key];
        }
    }

    // Добавляем параметры запроса
    const params = new URLSearchParams({
        ...cleanedFilters,
        sortBy,
        order,
        limit: itemsPerPage,
        offset: (page - 1) * itemsPerPage,
        minPrice, // Передаем minPrice
        maxPrice, // Передаем maxPrice
    });

    try {
        const response = await fetch(`${API_URL}?${params.toString()}`);
        const result = await response.json();

        if (result.status === "success") {
            // Извлекаем массив товаров из result.data.data
            const products = result.data.data;
            const total = result.data.total;

            renderProducts(products);
            renderPagination(total);
        } else {
            console.error("Failed to load products:", result.message);
            productList.innerHTML = '<p>Impression error while there are no such products in the Database. Try to find out later.</p>';
        }
    } catch (error) {
        console.error("Error fetching products:", error);
        productList.innerHTML = '<p>Impression error while there are no such products in the Database. Try to find out later.</p>';
    } finally {
        loadingIndicator.style.display = "none";
    }
}

function renderProducts(productsToRender) {
    productList.innerHTML = "";

    if (!Array.isArray(productsToRender)) {
        console.error("Expected an array of products, but got:", productsToRender);
        productList.innerHTML = '<p>No products found.</p>';
        return;
    }

    if (productsToRender.length === 0) {
        productList.innerHTML = '<p>No products found.</p>';
        return;
    }

    productsToRender.forEach(product => {
        const productCard = document.createElement("div");
        productCard.className = "product-card";

        productCard.innerHTML = `
            <img src="${product.link || 'https://via.placeholder.com/250x150'}" alt="${product.type}">
            <div class="card-content">
                <h3>${product.brand} ${product.model}</h3>
                <p>Type: ${capitalize(product.type)}</p>
                <p>${product.ram ? `RAM: ${product.ram}` : ""}</p>
                <p>${product.storage ? `Storage: ${product.storage}` : ""}</p>
                <p>${product.color ? `Color: ${product.color}` : ""}</p>
                <p>${product.processor ? `Processor: ${product.processor}` : ""}</p>
                <p>${product.description ? `Description: ${product.description}` : ""}</p>
                <p class="price">$${product.price.toFixed(2)}</p>
            </div>
            <div class="card-footer">
                <button>Details</button>
            </div>
        `;
        productList.appendChild(productCard);
    });
}

function renderPagination(totalItems) {
    pagination.innerHTML = "";

    const totalPages = Math.ceil(totalItems / itemsPerPage);

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement("button");
        button.textContent = i;
        button.className = i === currentPage ? "active" : "";
        button.addEventListener("click", () => {
            currentPage = i;
            updateView(); // Отправляем запрос на сервер при изменении страницы
        });
        pagination.appendChild(button);
    }
}

function getSelectedValues(name) {
    const checkboxes = document.querySelectorAll(`input[name="${name}"]:checked`);
    return Array.from(checkboxes).map(checkbox => checkbox.value);
}

function updateView() {
    const filters = {
        type: getSelectedValues("type").join(","),
        brand: getSelectedValues("brand").join(","),
        price: getSelectedValues("price").join(","), // Передаем выбранные диапазоны цен
        ram: getSelectedValues("ram").join(","),
        storage: getSelectedValues("storage").join(","),
        processor: getSelectedValues("processor").join(","),
        color: getSelectedValues("color").join(","),
    };

    const sortByValue = sortBy.value.split("-")[0]; // Убираем -asc или -desc
    const order = sortBy.value.includes("asc") ? "asc" : "desc";

    // Получаем выбранные значения для minPrice и maxPrice
    const selectedPrices = getSelectedValues("price");
    let minPrice = 0;
    let maxPrice = 0;

    if (selectedPrices.length > 0) {
        const priceRange = selectedPrices[0];
        if (priceRange === "2000-0") {
            // Обработка диапазона "$2000+"
            minPrice = 2000;
            maxPrice = 0; // maxPrice не задан
        } else {
            const [min, max] = priceRange.split("-").map(Number);
            minPrice = min;
            maxPrice = max;
        }
    }

    fetchProducts(filters, sortByValue, order, currentPage, minPrice, maxPrice);
}

function capitalize(text) {
    return text.charAt(0).toUpperCase() + text.slice(1);
}

function toggleFilters() {
    const selectedTypes = getSelectedValues("type");

    if (selectedTypes.includes("headphones")) {
        ramFilterGroup.style.display = "none";
        storageFilterGroup.style.display = "none";
        processorFilterGroup.style.display = "none";
    } else {
        ramFilterGroup.style.display = "block";
        storageFilterGroup.style.display = "block";
        processorFilterGroup.style.display = "block";
    }

    if (selectedTypes.includes("headphones")) {
        document.querySelectorAll('input[name="type"]').forEach(checkbox => {
            if (checkbox.value !== "headphones") {
                checkbox.checked = false;
            }
        });
    } else {
        document.querySelector('input[name="type"][value="headphones"]').checked = false;
    }
}

document.addEventListener("DOMContentLoaded", () => {
    initDOM();
    fetchProducts(); // Загружаем данные при загрузке страницы

    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.addEventListener("change", () => {
            currentPage = 1;
            toggleFilters();
            updateView(); // Отправляем запрос на сервер при изменении фильтров
        });
    });

    sortBy.addEventListener("change", () => {
        currentPage = 1;
        updateView(); // Отправляем запрос на сервер при изменении сортировки
    });
});
