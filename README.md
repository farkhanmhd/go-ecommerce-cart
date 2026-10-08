# Ecommerce Cart System

A practice RESTful API service for an ecommerce shopping cart system built with Go and PostgreSQL.

## Overview

This project provides backend API capabilities covering:

- **User Authentication & Management**: Registration, login, session/logout, and profile retrieval.
- **Product Catalog**: Browsing products and viewing product details.
- **Shopping Cart Operations**: Adding items, viewing cart contents, updating item quantities, and removing items.

API specifications and endpoint schemas can be referenced in [openapi.json](file:///Users/farkhanmhd/golang/go-ecommerce-cart/openapi.json).

---

## Prerequisites & Dependencies

To install the required packages and libraries:

```bash
# Request validation
go get github.com/go-playground/validator/v10

# HTTP routing
go get github.com/julienschmidt/httprouter

# Testing and assertions
go get github.com/stretchr/testify

# PostgreSQL driver and toolkit
go get github.com/jackc/pgx/v5

# Password hashing
go get golang.org/x/crypto/bcrypt

# JWT
go get github.com/golang-jwt/jwt/v5

# UUID
go get github.com/google/uuid

# Crypto
go get golang.org/x/crypto

# ENV
go get github.com/joho/godotenv
```

### Direct Links

- [validator](https://github.com/go-playground/validator) - Request and payload validation
- [httprouter](https://github.com/julienschmidt/httprouter) - High performance HTTP request multiplexer
- [testify](https://github.com/stretchr/testify) - Testing toolkit and assertions
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing utilities

---

## API Endpoints Summary

### Auth Module

| Method | Endpoint             | Description         |
| ------ | -------------------- | ------------------- |
| `POST` | `/api/auth/register` | Register new user   |
| `POST` | `/api/auth/login`    | Login user          |
| `POST` | `/api/auth/logout`   | Logout user         |

### Users Module

| Method | Endpoint             | Description         |
| ------ | -------------------- | ------------------- |
| `GET`  | `/api/users/current` | Get current user    |

### Products Module

| Method   | Endpoint                   | Description         |
| -------- | -------------------------- | ------------------- |
| `GET`    | `/api/products`            | Get list of products|
| `POST`   | `/api/products`            | Create product      |
| `GET`    | `/api/products/:productId` | Get product details |
| `PUT`    | `/api/products/:productId` | Update product      |
| `DELETE` | `/api/products/:productId` | Delete product      |

### Carts Module

| Method   | Endpoint                  | Description               |
| -------- | ------------------------- | ------------------------- |
| `GET`    | `/api/cart`               | Get user cart             |
| `DELETE` | `/api/cart`               | Clear cart                |
| `POST`   | `/api/cart/items`         | Add item to cart          |
| `PUT`    | `/api/cart/items/:itemId` | Update cart item quantity |
| `DELETE` | `/api/cart/items/:itemId` | Remove item from cart     |

