# Travel-Buddy

A **Tour Management Backend API** built with **Go** following **Clean Architecture** and **SOLID principles**.
The system allows travel agencies to manage tours, bookings, members, and permissions while customers can explore and book tours.

This project demonstrates a **production-style backend architecture** including authentication, middleware, testing, Docker support, database migrations, and CI/CD with GitHub Actions.

---

# 🚀 Features

* User authentication using **JWT**
* Tour management
* Booking system
* Travel agency management
* Member and permission management
* Search functionality
* Tour booking
* Payment support
* Rate limiting middleware
* Logging middleware
* Database migrations
* Dockerized environment
* CI/CD using GitHub Actions
* Unit testing with mocks

---

# 🧱 Architecture

This project follows **Clean Architecture**:

```
HTTP Layer (Handlers / Router)
        │
        ▼
Usecases (Business Logic)
        │
        ▼
Repository Interfaces (Ports)
        │
        ▼
Infrastructure (PostgreSQL / External Services)
```

Advantages of this architecture:

* Separation of concerns
* Highly testable business logic
* Independent infrastructure layer
* Easier scalability and maintainability

---

# 🛠 Tech Stack

* **Go**
* **PostgreSQL**
* **SQLX**
* **Docker**
* **JWT Authentication**
* **REST API**
* **Clean Architecture**
* **SOLID Principles**
* **GitHub Actions**
* **Go Testing**

---

# 📂 Project Structure

```
.
├── cmd
│   └── api
│       └── main.go
├── config
│   └── config.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── adapter
│   │   └── http
│   ├── domain
│   │   ├── agency.go
│   │   ├── agencyMember.go
│   │   ├── booking.go
│   │   ├── customer.go
│   │   ├── home.go
│   │   ├── payment.go
│   │   ├── permission.go
│   │   ├── Review.go
│   │   ├── role.go
│   │   ├── search.go
│   │   ├── tour.go
│   │   └── user.go
│   ├── infrastructure
│   │   └── postgres
│   ├── mocks
│   │   ├── repository
│   │   └── usecase
│   ├── usecase
│   │   ├── agency
│   │   ├── agencyMember
│   │   ├── booking
│   │   ├── home
│   │   ├── permission
│   │   ├── port
│   │   ├── search
│   │   ├── tour
│   │   └── user
│   └── validation
│       └── validator.go
├── migrations
├── utils
└── makefile
```

---

# ⚙️ Environment Variables

Create a `.env` file in the root directory.

```
VERSION=1.0.0
SERVICE_NAME=Tour App
HTTP_PORT=3000

JWT_SECRET_KEY=my_secret_key

DBHOST=localhost
DBPORT=5432
DBNAME=travelbuddy
DBUSER=postgres
DBPASSWORD=password
ENABLE_SSL_MODE=false
```

⚠️ Do not commit `.env` files to version control.

---

# ▶️ Running the Application

Run the application locally:

```
go run cmd/api/main.go
```

---

# 🐳 Running with Docker

Build and run using Docker Compose:

```
docker compose up --build
```

This will start:

* Go API service
* PostgreSQL database

---

# 🗄 Database Migrations

Database migrations are located in:

```
/migrations
```

They define schema changes for:

* users
* customers
* travel agencies
* permissions
* roles
* agency members
* tours
* bookings
* payments
* reviews

---

# 📡 API Endpoints

## Home

```
GET /home
```

## Search

```
GET /search
```

## Tours

```
POST   /tours
GET    /tours/{tour_id}
GET    /tours/list/{agency_id}
PUT    /tours/{tour_id}
DELETE /tours/{tour_id}
```

## Users

```
POST   /users
POST   /users/login
DELETE /users/{user_id}
PUT    /users/{user_id}
```

## Bookings

```
POST /bookings/{tour_id}
```

## Agencies

```
POST   /agency
PUT    /agency/{agency_id}
DELETE /agency/{agency_id}
```

## Members

```
POST   /members
DELETE /members/{member_id}
GET    /members/{agency_id}
PUT    /members/{member_id}/permissions
POST   /members/login
```

## Permissions

```
POST   /permissions
DELETE /permissions/{id}
```

---

# 🧪 Running Tests

Run all tests:

```
go test ./...
```

Mocks are used for testing repositories and usecases.

---

# 🔐 Authentication

Authentication is handled using **JWT tokens**.

Utilities for authentication are located in:

```
utils/
```

Including:

* JWT token generation
* password hashing
* payload extraction

---

# 📈 Middleware

The API includes middleware for:

* **Logging**
* **Rate Limiting**
* **Authentication**

---

# 🧑‍💻 Development

This project follows best practices such as:

* Clean architecture
* SOLID principles
* Dependency injection
* Layer separation
* Interface-driven development
* Testable use cases

---

# 📦 CI/CD

GitHub Actions is used for:

* running tests
* validating builds
* maintaining code quality

---

# 📄 License

This project is licensed under the MIT License.

---

# 👨‍💻 Author

Developed by **Bishal Das**

Backend Engineer | Go Developer
