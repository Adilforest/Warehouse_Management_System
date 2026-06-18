# Warehouse Management System

A role-based inventory management web application built with Go (Gin), MongoDB, and vanilla HTML/CSS/JS — featuring JWT authentication, rate limiting, a separate payment microservice, and PDF receipt generation.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.10-00ADD8?logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-7.0-47A248?logo=mongodb&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-HS256-000000?logo=jsonwebtokens&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white)

---

## Overview

Built as a university group project at Astana IT University (SE-2311) by Adil Ormanov and Madi Kassymov.

The system lets different user roles manage a product inventory backed by MongoDB. A lightweight second service (`transaction-service`) handles the payment flow separately: the main backend creates a pending transaction and redirects the user to the payment form; the transaction service validates the card, updates the status, generates a PDF receipt, and emails it to the buyer.

---

## Features

- **Role-based access control** — three roles enforced by JWT middleware:
  - `admin` — full CRUD on products and users, manage roles
  - `cashier` — manage product quantities (assumed via the role system)
  - unregistered / guest — read-only product listing
- **User management** — register, email verification, login, profile CRUD, admin-controlled role updates
- **Product catalog** — full CRUD with rich attributes: type, brand, model, processor, RAM, storage, color, price, quantity, warranty; paginated and filterable listing (type, price range, brand, RAM, storage, processor, color) with sort support
- **Shopping cart** — add, remove, clear, and pay via cart; cart items persisted in MongoDB
- **Payment microservice** — create pending transaction, show payment form, validate card details (number format, CVV regex, expiry MM/YY), simulate approval (cards starting with `4` succeed), update status to `paid` or `declined`
- **PDF receipt generation** — on successful payment, generates a PDF receipt (`gofpdf`) and emails it to the buyer via SMTP (`gomail`)
- **Rate limiting** — `/api/contact` endpoint limited to 1 request per 15 seconds using `golang.org/x/time/rate`
- **Contact form** — email sending via SMTP with rate protection
- **Structured logging** — Logrus with log rotation (lumberjack), request logging middleware
- **CORS** — configurable allowed origins

---

## Tech Stack

| Layer | Technologies |
|---|---|
| Language | Go 1.23 |
| HTTP framework | Gin |
| Database | MongoDB (`go.mongodb.org/mongo-driver`) |
| Auth | JWT HS256 (`golang-jwt/jwt/v5`) |
| PDF generation | `jung-kurt/gofpdf` |
| Email | `gopkg.in/gomail.v2` |
| Rate limiting | `golang.org/x/time/rate` |
| Logging | Logrus + lumberjack |
| Frontend | HTML, CSS, JavaScript (no build step) |
| Containers | Docker, Docker Compose |

---

## Project Structure

```
Warehouse_Management_System/
├── warehouse-backend/
│   ├── main.go                      # Entry point: router setup, graceful shutdown
│   ├── config/config.go             # Config loading
│   ├── models/
│   │   ├── user.go                  # User schema (id, name, email, role, verified)
│   │   ├── product.go               # Product schema (type, brand, model, specs, price, qty)
│   │   └── cart.go                  # Cart schema
│   ├── controllers/
│   │   ├── auth_controller.go       # Register, login, email verification
│   │   ├── profile_controller.go    # Get/update profile
│   │   ├── user_controller.go       # Admin user management
│   │   ├── shoppingCart_controller.go # Cart operations + pay
│   │   └── contact_controller.go    # Contact form + email
│   ├── middleware/
│   │   ├── auth.go                  # JWT generation, validation, AuthMiddleware
│   │   └── rate_limit_middleware.go # Token bucket rate limiter
│   ├── routes/
│   │   ├── auth_routes.go
│   │   └── product_routes.go
│   ├── database/
│   │   ├── database.go              # Product CRUD with filtering + pagination
│   │   ├── mongo_connection.go      # MongoDB connect/disconnect
│   │   └── cart.go                  # Cart DB operations
│   └── logger/logger.go             # Logrus setup with lumberjack rotation
├── transaction-service/
│   └── main.go                      # Payment microservice (Gin :8081)
│                                    # Endpoints: POST /transactions,
│                                    #   GET  /transactions/:id/payment,
│                                    #   POST /transactions/:id/pay
│                                    # PDF receipt + email on successful payment
└── warehouse-frontend/
    ├── index.html                   # Product listing (guest access)
    ├── admin.html                   # Admin dashboard
    ├── LogIn.html / SignUp.html
    ├── profile.html
    ├── shoppingСart.html
    ├── ContactUs.html
    ├── js/                          # auth.js, index.js, admin.js, shoppingCart.js, etc.
    └── css/                         # Per-page stylesheets
```

---

## Getting Started

### Prerequisites

- Go 1.23+
- MongoDB (local or Atlas)
- Docker + Docker Compose (optional)
- SMTP credentials for email features

### Environment variables

Create a `.env` file in `warehouse-backend/` with:

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB=warehouse
JWT_SECRET=<your-secret>
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=<your-email>
SMTP_PASSWORD=<app-password>
```

Create a `.env` in `transaction-service/` with:

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB=warehouse
PORT=8081
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=<your-email>
SMTP_PASSWORD=<app-password>
```

### Run the backend

```bash
cd warehouse-backend
go mod tidy
go run main.go
# Server starts on :8080
```

### Run the transaction service

```bash
cd transaction-service
go mod tidy
go run main.go
# Service starts on :8081
```

### Run with Docker Compose

```bash
cd warehouse-backend
docker compose up --build
```

### Open the frontend

Open `warehouse-frontend/index.html` directly in a browser, or serve it with any static file server.

---

## API Endpoints

### Auth

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | — | Register new user |
| `POST` | `/auth/login` | — | Login, receive JWT |
| `GET` | `/auth/verify` | — | Verify email token |

### Products

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/products/` | — | List products (paginated, filterable) |
| `GET` | `/products/:id` | — | Get product by ID |
| `POST` | `/products/create` | admin | Create product |
| `PUT` | `/products/:id` | admin | Update product |
| `DELETE` | `/products/:id` | admin | Delete product |
| `DELETE` | `/products/deleteAll` | admin | Delete all products |

### Users (admin only)

| Method | Path | Description |
|---|---|---|
| `GET` | `/users/` | List all users |
| `GET` | `/users/:id` | Get user by ID |
| `PUT` | `/users/:id` | Update user |
| `PUT` | `/users/:id/role` | Change user role |
| `DELETE` | `/users/:id` | Delete user |

### Profile & Cart

| Method | Path | Description |
|---|---|---|
| `GET` | `/protected/profile` | Get current user profile |
| `PUT` | `/protected/profile` | Update profile |
| `POST` | `/cart/:product_id/add` | Add item to cart |
| `GET` | `/cart/` | View cart |
| `DELETE` | `/cart/:product_id/remove` | Remove item |
| `DELETE` | `/cart/clear` | Clear cart |
| `POST` | `/cart/pay` | Checkout (creates transaction) |

### Transaction service (:8081)

| Method | Path | Description |
|---|---|---|
| `POST` | `/transactions` | Create pending transaction |
| `GET` | `/transactions/:id/payment` | Show payment form |
| `POST` | `/transactions/:id/pay` | Submit payment (validates card, generates PDF receipt) |

### Other

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/contact` | Contact form (rate-limited: 1 req / 15 s) |

---

## Screenshot

![Homepage](warehouse-frontend/images/Screenshot%202024-12-20%20at%2007.19.04.png)

---

Adil Ormanov — [GitHub](https://github.com/Adilforest)
