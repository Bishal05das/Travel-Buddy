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

# Optional
JWT_TTL=24h                       # access token lifetime (Go duration)
TRUST_PROXY_HEADERS=false         # true only behind a proxy that sets X-Real-IP / X-Forwarded-For
MIGRATIONS_PATH=file://migrations # the Docker image uses file:///migrations
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

Protected routes need `Authorization: Bearer <token>` from `POST /users/login`
or `POST /members/login`. The **Access** column lists who may call each route
(see [Roles and permissions](#-roles-and-permissions)).

| Method | Path | Access |
|---|---|---|
| GET | `/home` | public |
| GET | `/search?q=&min_price=&max_price=&start_date=&end_date=` | public |
| GET | `/tours/{tour_id}` | public |
| GET | `/agency/{agency_id}/tours/list?page=&limit=` | public |
| POST | `/agency/{agency_id}/tours` (multipart, `image` file) | `tour:create` |
| PUT | `/agency/{agency_id}/tours/{tour_id}` | `tour:update` |
| PATCH | `/tours/{tour_id}/tour-status` (body: `"open"`/`"closed"`/`"cancelled"`) | `tour:update` |
| DELETE | `/tours/{tour_id}` | `tour:delete` |
| POST | `/users` | public |
| POST | `/users/login` | public |
| PUT, DELETE | `/users/{user_id}` | the user themselves, or super |
| POST | `/bookings/{tour_id}` | role `user` |
| POST | `/admin/bookings/{tour_id}` (guest booking) | member with `booking:create` |
| POST | `/agency` (multipart, `image` file) | super |
| PUT | `/agency/{agency_id}` | `agency:update` |
| DELETE | `/agency/{agency_id}` | `agency:delete` |
| POST | `/members/{agency_id}` | `member:create` |
| GET | `/members/{agency_id}` | `member:read` |
| PUT | `/members/{member_id}/permissions` | `member:update` |
| DELETE | `/members/{member_id}` | `member:delete` |
| POST | `/members/login` | public |
| POST, DELETE | `/permissions`, `/permissions/{id}` | super |
| GET | `/images/{path}` | public (files only, no directory listing) |

The booking `total_price` must equal `(price - price * discount / 100) * number_of_people`.

---

# 🧪 Running Tests

Run all tests:

```
go test ./...
```

Mocks are used for testing repositories and usecases.

---

# 🔐 Authentication

Access tokens are HS256 JWTs signed with `JWT_SECRET_KEY`. They expire after
`JWT_TTL` and carry the caller's id, role and, for members, their agency.

## 🛡 Roles and permissions

* **super** (`users.role = 'super'`): the platform administrator. Creates
  agencies, manages the permission catalogue, can act on any agency, and
  creates each agency's first member.
* **member**: an agency staff account, limited to its own agency. Each route
  requires a permission (`tour:create`, `member:update`, …) granted through the
  member's role. Permissions are checked against the database on every request,
  so revoking one applies immediately. A member cannot grant permissions they
  do not hold.
* **user**: a customer. Can book tours and manage only their own account.

The permissions the routes use are seeded by migration `000015`. List them with
`SELECT permission_id, name FROM permissions;` to find the ids to pass as
`permissions` when creating members.

### Bootstrapping a new environment

1. Register an account with `POST /users`.
2. Promote it once in the database:
   `UPDATE users SET role = 'super' WHERE email = 'admin@example.com';`
3. Log in, create an agency with `POST /agency`, then create its first
   member (for example a manager holding every permission) with
   `POST /members/{agency_id}`. That member can manage the agency from then on.

---

# 📈 Middleware

The API includes middleware for:

* **Logging**
* **Rate Limiting**: 30 requests/minute per client IP
* **Authentication**: token verification
* **Authorization**: role, self-or-super and permission checks
* **Request size limit**: 11 MB bodies

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

GitHub Actions (`.github/workflows/ci-cd.yml`):

* runs `gofmt`, `go vet` and `go test -race` on every push to `main`,
  `pre-release`, `release/**` and `feature/**`, and on pull requests
* builds and pushes the Docker image on manual dispatch, only after tests pass

---

# 📄 License

This project is licensed under the MIT License.

---

# 👨‍💻 Author

Developed by **Bishal Das**

Backend Engineer | Go Developer
