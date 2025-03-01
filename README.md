<p align="center">
  <img src="./static/logo.png" alt="thaym api Logo" width="280" >
</p>

<p align="center">
  <strong>Thāyṃ</strong> is a robust and scalable e-commerce API built with GoLang, leveraging Hexagonal Architecture. This design pattern ensures that the application is maintainable and adaptable over time, with clear separation of concerns. The API provides essential features for managing e-commerce operations, including secure JWT authentication, email confirmation, advanced filtering capabilities, inventory management, role-based access control and sure payment gateway integration. Thāyṃ is designed to support multiple user roles such as Admin, User, and Store Owner, and allows for seamless product management, order tracking, and payment integrations. With its flexible and modular structure, Thāyṃ is ready to power a wide range of e-commerce applications while maintaining high performance and security.
</p>

---

## Setup Guide


### Environment Setup
1. Configure Your .env File

Create and configure a .env file in your project directory with the following variables:

```bash
# Server
PUBLIC_HOST=
PORT=8085

# Database
DB_USER=root
DB_PASSWORD=root
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=ecom

# JWT
JWT_SECRET=

# SMTP for Gmail
FROM_EMAIL=
FROM_EMAIL_PASSWORD= # using passkey
FROM_EMAIL_SMTP="smtp.gmail.com"
SMTP_ADDR="smtp.gmail.com:587"

# Stripe
SECRET_KEY_STRIPE=
WEBHOOK_SECRET_STRIPE=

#Payment Variable
ORDER_STATUS_PENDING="pending"  
ORDER_STATUS_PROCESSING="processing"  
ORDER_STATUS_SHIPPED="shipped"  
ORDER_STATUS_COMPLETED="completed"  
ORDER_STATUS_CANCELLED="cancelled"  
ORDER_STATUS_REFUNDED="refunded"  
PAYMENT_STATUS_PENDING="pending"  
PAYMENT_STATUS_PAID="paid"  
PAYMENT_STATUS_REFUNDED="refunded" 
```


### Using Containerized Environment

1. Build and Run the Project Using Docker Compose

```bash
docker-compose up --build
```

Once the project is built and the containers are running, you can access the API at:

> API URL: http://localhost:8085

2. Stop Containers

```bash
docker-compose down
```

### Running Locally Using `Make`


**Build**

Compile the application and output the binary in the `bin` directory:

```bash
make build
```

**Run**

Build and run the application:

```bash
make run
```

Once the project is built and the containers are running, you can access the API at:

> API URL: http://localhost:8085

### Database Migrations

**Generate a new migration file in the `cmd/migrate/migrations` directory:**

```bash
make migration <your_migration_name>
```

**Apply Migrations**

```bash
make migrate-up
```

**Rollback Migrations**

```bash
make migrate-down
```

--- 

## Features

### Authentication and Authorization
- **JWT Authentication:** Secure authentication using JSON Web Tokens.
- **Role-based Access Control:** Supports roles such as Admin, User, and Store Owner.

### User Management
- Email confirmation using custom HTML templates.
- Complete CRUD (Create, Read, Update, Delete) operations for user accounts.

### E-commerce Functionalities
- Advanced filtering by tags and categories.
- Search functionality for products and orders.
- Inventory management:
  - Control product quantity
  - Product stocking
  - Activation and deactivation of products.
- Order management:
  - Seamless integration with payment gateways.
  - Tracking and updating order statuses.

### Payment Gateway Integration
- Supports mutliple payment providers (e.g., Stripe, Banks[can be configure])  
- Secure transaction handling with encryption.
- Payment successfull Alerting mechanism

---

## Technologies Used

### Core Language and Architecture
- **Go 1.23.4:** The core programming language.
- **Hexagonal Architecture:** Ensures separation of concerns and maintainability.

### Dependencies
- **[go-playground/validator](https://github.com/go-playground/validator):** Input validation.
- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql):** MySQL database driver.
- **[golang-jwt/jwt](https://github.com/golang-jwt/jwt):** JWT authentication.
- **[golang-migrate/migrate](https://github.com/golang-migrate/migrate):** Database migrations.
- **[gorilla/mux](https://github.com/gorilla/mux):** HTTP request router.
- **[joho/godotenv](https://github.com/joho/godotenv):** Environment variable management.
- **[stretchr/testify](https://github.com/stretchr/testify):** Testing utilities.
- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto):** Cryptographic functions.

---

