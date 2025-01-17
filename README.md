# Thāyṃ E-commerce API

**Thāyṃ** is an e-commerce API built in GoLang using Hexagonal Architecture. This API is designed to provide a scalable and maintainable structure for managing e-commerce operations. It supports features like JWT authentication, email confirmation, advanced filtering, inventory management, role-based access control, and much more.

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

### This project uses a Makefile to streamline common tasks. 

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

**Test**

Run all tests with verbose output:

```bash
make test
```

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

