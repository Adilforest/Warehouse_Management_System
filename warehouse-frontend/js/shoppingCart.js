// Базовый URL API
const apiUrl = "http://localhost:8080/cart";

// Функция для получения токена из localStorage
function getToken() {
    return localStorage.getItem("token");
}

// Функция для проверки авторизации пользователя
async function checkAuth() {
    const token = getToken();
    if (!token) {
        alert("You need to log in first.");
        window.location.href = "LogIn.html";
        return false;
    }
    return true;
}

// Функция для отображения товаров в корзине
async function displayCart() {
    if (!(await checkAuth())) return;

    try {
        const response = await fetch(apiUrl + "/", {
            method: "GET",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        if (!response.ok) {
            throw new Error("Failed to fetch cart data.");
        }

        const result = await response.json();
        const cart = result.data; // Извлекаем объект корзины из свойства data
        const cartContainer = document.getElementById("cart-container");

        // Очищаем контейнер перед отображением новых данных
        cartContainer.innerHTML = "";

        if (!cart || !cart.items || cart.items.length === 0) {
            cartContainer.innerHTML = "<p>Your cart is empty.</p>";
            return;
        }

        // Создаем элементы для каждого товара в корзине
        cart.items.forEach((item) => {
            const productCard = document.createElement("div");
            productCard.className = "cart-item";

            productCard.innerHTML = `
                <h3>${item.product.brand} ${item.product.model}</h3>
                <p>Type: ${capitalize(item.product.type)}</p>
                <p>Price: $${Number(item.product.price).toFixed(2)}</p>
                <p>Quantity: ${item.quantity}</p>
                <button class="remove-btn" data-product-id="${item.product._id}">Remove</button>
            `;

            // Добавляем обработчик события для кнопки "Remove" с подтверждением
            const removeButton = productCard.querySelector(".remove-btn");
            removeButton.addEventListener("click", () => {
                if (confirm("Are you sure you want to remove this product from your cart?")) {
                    removeFromCart(item.product._id);
                }
            });

            cartContainer.appendChild(productCard);
        });

        // Добавляем блок для отображения общей суммы корзины
        const totalDiv = document.createElement("div");
        totalDiv.className = "cart-total";
        totalDiv.innerHTML = `<h2>Total: $${Number(cart.total).toFixed(2)}</h2>`;
        cartContainer.appendChild(totalDiv);
    } catch (error) {
        console.error("Error fetching cart:", error);
        document.getElementById("cart-container").innerHTML =
            "<p>An error occurred while loading your cart. Please try again later.</p>";
    }
}



// Функция для удаления товара из корзины
async function removeFromCart(productId) {
    if (!(await checkAuth())) return;

    try {
        const response = await fetch(`${apiUrl}/${productId}/remove`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        if (!response.ok) {
            throw new Error("Failed to remove item from cart.");
        }

        alert("Product removed from cart successfully!");
        displayCart(); // Обновляем отображение корзины
    } catch (error) {
        console.error("Error removing item from cart:", error);
        alert("An error occurred while removing the item from the cart.");
    }
}

// Функция для очистки всей корзины
async function clearCart() {
    if (!(await checkAuth())) return;

    try {
        const response = await fetch(`${apiUrl}/clear`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        if (!response.ok) {
            throw new Error("Failed to clear cart.");
        }

        alert("Cart cleared successfully!");
        displayCart(); // Обновляем отображение корзины
    } catch (error) {
        console.error("Error clearing cart:", error);
        alert("An error occurred while clearing the cart.");
    }
}

// Вспомогательная функция для преобразования первой буквы строки в верхний регистр
function capitalize(text) {
    return text.charAt(0).toUpperCase() + text.slice(1);
}

// Инициализация при загрузке страницы
document.addEventListener("DOMContentLoaded", () => {
    displayCart();

    // Добавляем кнопку "Clear Cart" динамически
    const clearCartButton = document.createElement("button");
    clearCartButton.id = "clear-cart-btn";
    clearCartButton.innerText = "Clear Cart";
    clearCartButton.addEventListener("click", clearCart);

    const mainElement = document.querySelector("main");
    mainElement.appendChild(clearCartButton);
});