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
            // Если идентификатор продукта приходит как id, а не _id, используем его
            const productId = item.product._id || item.product.id;
            const productCard = document.createElement("div");
            productCard.className = "cart-item";

            productCard.innerHTML = `
                <h3>${item.product.brand} ${item.product.model}</h3>
                <p>Type: ${capitalize(item.product.type)}</p>
                <p>Price: $${Number(item.product.price).toFixed(2)}</p>
                <p>Quantity: ${item.quantity}</p>
                <button class="remove-btn" data-product-id="${productId}">Remove</button>
            `;

            // Добавляем обработчик события для кнопки "Remove" с подтверждением
            const removeButton = productCard.querySelector(".remove-btn");
            removeButton.addEventListener("click", () => {
                if (confirm("Are you sure you want to remove this product from your cart?")) {
                    removeFromCart(productId);
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

// Функция для оплаты корзины с использованием формы
async function submitPaymentForm() {
    // Получаем значения из формы
    const cardNumber = document.getElementById("cardNumber").value.trim();
    const expirationDate = document.getElementById("expirationDate").value.trim();
    const cvv = document.getElementById("cvv").value.trim();
    const name = document.getElementById("name").value.trim();
    const address = document.getElementById("address").value.trim();

    if (!cardNumber || !expirationDate || !cvv || !name || !address) {
        alert("Please fill in all payment fields.");
        return;
    }

    const paymentDetails = {
        cardNumber,
        expirationDate,
        cvv,
        name,
        address,
    };

    try {
        const response = await fetch(`${apiUrl}/pay`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${getToken()}`,
            },
            body: JSON.stringify(paymentDetails)
        });

        if (!response.ok) {
            throw new Error("Failed to process payment.");
        }

        const result = await response.json();
        alert(result.message || "Payment successful!");

        // После успешной оплаты скрываем форму и очищаем поля
        document.getElementById("payment-details-form").reset();
        document.getElementById("payment-form").style.display = "none";

        displayCart(); // Обновляем отображение корзины после оплаты
    } catch (error) {
        console.error("Error processing payment:", error);
        alert("An error occurred while processing payment.");
    }
}

// Функция для показа формы оплаты (при нажатии кнопки "Pay Cart")
function showPaymentForm() {
    document.getElementById("payment-form").style.display = "block";
}




// Функция для показа формы оплаты (при нажатии кнопки "Pay Cart")
function showPaymentForm() {
    const paymentForm = document.getElementById("payment-form");
    paymentForm.style.display = "block";
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

    // Добавляем кнопку "Pay Cart" динамически
    const payCartButton = document.createElement("button");
    payCartButton.id = "pay-cart-btn";
    payCartButton.innerText = "Pay Cart";
    payCartButton.addEventListener("click", showPaymentForm);

    const mainElement = document.querySelector("main");
    mainElement.appendChild(clearCartButton);
    mainElement.appendChild(payCartButton);

    // Добавляем обработчик отправки формы оплаты
    const paymentFormElement = document.getElementById("payment-details-form");
    paymentFormElement.addEventListener("submit", (e) => {
        e.preventDefault();
        submitPaymentForm();
    });
});
