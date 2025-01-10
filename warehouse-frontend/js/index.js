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

async function fetchProducts() {
    const loadingIndicator = document.getElementById("loading");
    loadingIndicator.style.display = "block";

    try {
        const response = await fetch(API_URL);
        const result = await response.json();

        if (result.status === "success") {
            products = result.data;
            updateView();
        } else {
            console.error("Failed to load products:", result.message);
            productList.innerHTML = '<p>Failed to load products.</p>';
        }
    } catch (error) {
        console.error("Error fetching products:", error);
        productList.innerHTML = '<p>Failed to load products. Please try again later.</p>';
    } finally {
        loadingIndicator.style.display = "none";
    }
}

function renderProducts(productsToRender) {
    productList.innerHTML = "";

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
            updateView();
        });
        pagination.appendChild(button);
    }
}

function getSelectedValues(name) {
    const checkboxes = document.querySelectorAll(`input[name="${name}"]:checked`);
    return Array.from(checkboxes).map(checkbox => checkbox.value);
}

function filterProducts() {
    let filteredProducts = products;

    const selectedTypes = getSelectedValues("type");
    if (selectedTypes.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedTypes.includes(product.type));
    }

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

    const selectedBrands = getSelectedValues("brand");
    if (selectedBrands.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedBrands.includes(product.brand));
    }

    const selectedRAMs = getSelectedValues("ram");
    if (selectedRAMs.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedRAMs.includes(product.ram.toString()));
    }

    const selectedStorages = getSelectedValues("storage");
    if (selectedStorages.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedStorages.includes(product.storage.toString()));
    }

    const selectedProcessors = getSelectedValues("processor");
    if (selectedProcessors.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedProcessors.includes(product.processor));
    }

    const selectedColors = getSelectedValues("color");
    if (selectedColors.length > 0) {
        filteredProducts = filteredProducts.filter(product => selectedColors.includes(product.color));
    }

    return filteredProducts;
}

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

function updateView() {
    const filteredProducts = filterProducts();
    const sortedProducts = sortProducts(filteredProducts);

    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const paginatedProducts = sortedProducts.slice(startIndex, endIndex);

    renderProducts(paginatedProducts);
    renderPagination(filteredProducts.length);
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
    fetchProducts();

    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.addEventListener("change", () => {
            currentPage = 1;
            toggleFilters();
            updateView();
        });
    });

    sortBy.addEventListener("change", () => {
        currentPage = 1;
        updateView();
    });
});