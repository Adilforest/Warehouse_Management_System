// URL API для загрузки данных из базы данных
const API_URL = 'http://localhost:8080/products/';

// Текущее состояние
let products = []; // Сюда будет загружен массив продуктов
let currentPage = 1; // Текущая страница
const itemsPerPage = 9; // Количество элементов на странице

// Получение ссылок на DOM-элементы
let productList, pagination, sortBy;

// Функция для инициализации DOM-элементов
function initDOM() {
    productList = document.getElementById("product-list");
    pagination = document.getElementById("pagination");
    sortBy = document.getElementById("sort-by"); // Инициализация элемента сортировки
}

// Fetch данных с сервера
async function fetchProducts() {
    // Показываем индикатор загрузки
    const loadingIndicator = document.getElementById("loading");
    loadingIndicator.style.display = "block";

    try {
        const response = await fetch(API_URL); // Запрос к API
        const result = await response.json();

        if (result.status === "success") {
            products = result.data; // Сохраняем продукты
            updateView(); // Обновляем отображение
        } else {
            console.error("Failed to load products:", result.message);
            productList.innerHTML = '<p>Failed to load products.</p>';
        }
    } catch (error) {
        console.error("Error fetching products:", error);
        productList.innerHTML = '<p>Failed to load products. Please try again later.</p>';
    } finally {
        // Скрываем индикатор загрузки
        loadingIndicator.style.display = "none";
    }
}

// Функция для отображения продуктов на странице
function renderProducts(productsToRender) {
    productList.innerHTML = ""; // Очистка списка продуктов

    if (productsToRender.length === 0) {
        productList.innerHTML = '<p>No products found.</p>';
        return;
    }

    productsToRender.forEach(product => {
        const productCard = document.createElement("div");
        productCard.className = "product-card";

        // Используем поле `link` для изображения
        productCard.innerHTML = `
            <img src="${product.link || 'https://via.placeholder.com/250x150'}" alt="${product.type}">
            <div class="card-content">
                <h3>${product.brand} ${product.model}</h3>
                <p>Type: ${capitalize(product.type)}</p>
                <p>${product.ram ? `RAM: ${product.ram}` : ""}</p>
                <p>${product.storage ? `Storage: ${product.storage}` : ""}</p>
                <p>${product.color ? `Color: ${product.color}` : ""}</p>
                <p class="price">$${product.price.toFixed(2)}</p>
            </div>
            <div class="card-footer">
                <button>Details</button>
            </div>
        `;
        productList.appendChild(productCard);
    });
}

// Функция для создания кнопок пагинации
function renderPagination(totalItems) {
    pagination.innerHTML = ""; // Очистка кнопок

    const totalPages = Math.ceil(totalItems / itemsPerPage);

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement("button");
        button.textContent = i;
        button.className = i === currentPage ? "active" : "";
        button.addEventListener("click", () => {
            currentPage = i;
            updateView();
        });
        pagination.appendChild(button);
    }
}

// Функция для получения выбранных значений из чекбоксов
function getSelectedValues(name) {
    const checkboxes = document.querySelectorAll(`input[name="${name}"]:checked`);
    return Array.from(checkboxes).map(checkbox => checkbox.value);
}

// Функция для фильтрации продуктов
function filterProducts() {
    let filteredProducts = products;

    // Фильтрация по типу
    const selectedTypes = getSelectedValues("type");
    if (selectedTypes.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedTypes.includes(product.type));
    }

    // Фильтрация по цене
    const selectedPrices = getSelectedValues("price");
    if (selectedPrices.length > 0) {
        filteredProducts = filteredProducts.filter(product => {
            return selectedPrices.some(range => {
                const [min, max] = range.split("-").map(Number);
                if (range.endsWith("+")) {
                    return product.price >= min;
                }
                return product.price >= min && product.price <= max;
            });
        });
    }

    // Фильтрация по бренду
    const selectedBrands = getSelectedValues("brand");
    if (selectedBrands.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedBrands.includes(product.brand));
    }

    // Фильтрация по RAM
    const selectedRAMs = getSelectedValues("ram");
    if (selectedRAMs.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedRAMs.includes(product.ram.toString()));
    }

    // Фильтрация по хранилищу
    const selectedStorages = getSelectedValues("storage");
    if (selectedStorages.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedStorages.includes(product.storage.toString()));
    }

    // Фильтрация по цвету
    const selectedColors = getSelectedValues("color");
    if (selectedColors.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedColors.includes(product.color));
    }

    return filteredProducts;
}

// Функция для сортировки продуктов
function sortProducts(productsToSort) {
    const selectedSort = sortBy.value;

    switch (selectedSort) {
        case "price-asc":
            return productsToSort.sort((a, b) => a.price - b.price);
        case "price-desc":
            return productsToSort.sort((a, b) => b.price - a.price);
        case "brand":
            return productsToSort.sort((a, b) => a.brand.localeCompare(b.brand));
        case "model":
            return productsToSort.sort((a, b) => a.model.localeCompare(b.model));
        default:
            return productsToSort;
    }
}

// Функция для обновления отображения продуктов и пагинации
function updateView() {
    const filteredProducts = filterProducts();
    const sortedProducts = sortProducts(filteredProducts);

    // Пагинация
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const paginatedProducts = sortedProducts.slice(startIndex, endIndex);

    renderProducts(paginatedProducts); // Отображаем продукты
    renderPagination(filteredProducts.length); // Отображаем кнопки для пагинации
}

// Утилита для приведения текста к виду с заглавной буквы
function capitalize(text) {
    return text.charAt(0).toUpperCase() + text.slice(1);
}

// Инициализация
document.addEventListener("DOMContentLoaded", () => {
    initDOM(); // Инициализация DOM-элементов
    fetchProducts(); // Загрузка продуктов

    // Слушатели событий для всех чекбоксов
    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.addEventListener("change", () => {
            currentPage = 1; // Сброс на первую страницу
            updateView();
        });
    });

    // Слушатель для сортировки
    sortBy.addEventListener("change", () => {
        currentPage = 1; // Сброс на первую страницу
        updateView();
    });
});