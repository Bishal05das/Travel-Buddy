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
RATE_LIMIT_PER_MINUTE=30          # requests per client IP per minute
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

# 🖥 Web apps

The customer site and agency portal are separate Next.js applications. They can
run on different servers and both call the same Go API.

| Application | Folder | Local URL | Users |
|---|---|---|---|
| Customer site | [`frontend/`](frontend/README.md) | http://localhost:3001 | Tour browsing, registration, customer bookings |
| Agency portal | [`agency-portal/`](agency-portal/README.md) | http://localhost:3003 | Agency staff dashboard, guest bookings, platform admin onboarding |
| Go API | `cmd/api/` | `http://localhost:${HTTP_PORT}` | Shared backend |

Start the API, then start each frontend in its own terminal:

```bash
cd frontend
cp .env.example .env.local
npm ci
npm run dev
```

```bash
cd agency-portal
cp .env.example .env.local
npm ci
npm run dev
```

Both apps read `NEXT_PUBLIC_API_URL` at build time. Set it to your API URL
in each app's `.env.local`; for example, use `http://localhost:3002` if
`HTTP_PORT=3002`. The frontends use separate browser sessions because they
run on different origins. Set `NEXT_PUBLIC_AGENCY_PORTAL_URL` on the customer
site and `NEXT_PUBLIC_CUSTOMER_URL` on the portal when deploying them to
different hosts.

---

# 🐳 Running with Docker

Build and run using Docker Compose:

```
docker compose up --build
```

This will start:

* Go API service
* PostgreSQL database

Both Compose files use `postgres:17-alpine` with the `tour_data_pg17` data
volume. The previous PostgreSQL 15 `tour_data` volume is retained for rollback.
Existing PostgreSQL 15 data needs a dump and restore into the new volume;
starting Compose alone does not migrate it.

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
| GET | `/search?q=&min_price=&max_price=&start_date=&end_date=&page=&limit=` | public |
| GET | `/tours/{tour_id}` | public |
| GET | `/agency/{agency_id}/tours/list?page=&limit=` | public |
| POST | `/agency/{agency_id}/tours` (multipart, `image` file) | `tour:create` |
| PUT | `/agency/{agency_id}/tours/{tour_id}` | `tour:update` |
| PUT | `/agency/{agency_id}/tours/{tour_id}/image` (multipart, `image` file) | `tour:update`; own agency, non-cancelled tours |
| PATCH | `/tours/{tour_id}/tour-status` (body: `"open"`/`"closed"`/`"cancelled"`) | `tour:update` |
| DELETE | `/tours/{tour_id}` | `tour:delete` |
| POST | `/users` | public |
| POST | `/users/login` | public |
| PUT, DELETE | `/users/{user_id}` | the user themselves, or super |
| POST | `/bookings/{tour_id}` | role `user` |
| POST | `/admin/bookings/{tour_id}` (guest booking) | member with `booking:create` |
| GET | `/agency/{agency_id}/bookings?status=&tour_id=&page=&limit=` | `booking:read` |
| GET | `/agency/{agency_id}/bookings/{booking_id}` | `booking:read` |
| PATCH | `/agency/{agency_id}/bookings/{booking_id}/status` (body: `{"status": "confirmed"}`) | `booking:update` |
| GET | `/me/bookings?status=&tour_id=&page=&limit=` | role `user` (own bookings) |
| GET | `/me/bookings/{booking_id}` | role `user` (own bookings) |
| POST | `/me/bookings/{booking_id}/cancel` | role `user` (own bookings) |
| POST | `/agency` (multipart, `image` file) | super |
| PUT | `/agency/{agency_id}` | `agency:update` |
| PUT | `/agency/{agency_id}/image` (multipart, `image` file) | `agency:update` |
| DELETE | `/agency/{agency_id}` | `agency:delete`; blocked when the agency has an owner |
| POST | `/members/{agency_id}` | `member:create` |
| GET | `/members/{agency_id}` | `member:read` |
| PUT | `/members/{member_id}/permissions` | `member:update`; owner access is protected |
| DELETE | `/members/{member_id}` | `member:delete`; owners cannot be deleted |
| POST | `/members/login` | public |
| GET | `/members/me` | signed-in member; no team-read permission required |
| POST, DELETE | `/permissions`, `/permissions/{id}` | super |
| GET | `/images/{path}` | public (files only, no directory listing) |

Search supports case-insensitive, typo-tolerant matching across tour names,
descriptions and agency names. Every search word must match; words can appear in
any order or across these fields. Exact names rank before partial and fuzzy
matches. Price filters use the discounted booking price. Only active agencies
and their tours appear. Queries accept up to 200 characters and 8 words;
`page` defaults to 1 and `limit` defaults to 20 (maximum 50). Responses include
tour images and pagination totals in `Meta`.

The booking `total_price` must equal `(price - price * discount / 100) * number_of_people`.

**Booking lifecycle:** a booking starts `pending`. The agency moves it to
`confirmed` once it has verified the payment, which also marks the payment
`success`. From `confirmed` it can become `completed`. The agency can cancel a
`pending` or `confirmed` booking at any time, and the customer can cancel their
own until the tour's start date. Cancelling returns the seats to the tour and
marks an unverified (`pending`) payment `failed`. A verified payment stays
`success`, and the refund is handled outside the system. `cancelled` and
`completed` are final. Invalid status changes return `409 Conflict`. A booking
outside the caller's agency or account returns `404`.

Each booking response includes `tour_status` and, when cancelled, a
`cancellation_reason`: `customer`, `agency` or `tour_cancelled`.

**Tour status:** `open` and `closed` can be switched freely, and closing a tour
only stops new bookings. Setting a tour to `cancelled` is final (reopening it
returns `409`). It also cancels every `pending` or `confirmed` booking on the
tour with reason `tour_cancelled`, returns their seats, and fails their
unverified payments, all in one transaction. The response reports how many
bookings were cancelled. Verified payments stay `success` for an offline
refund.

**Seats:** a tour's `total_seat` is its capacity and is set when the tour is
created or updated (tour creation still accepts `available_seat` as an alias).
`available_seat` is read-only: it starts at `total_seat` and bookings reduce it.
Updating `total_seat` keeps existing bookings, so `available_seat` becomes
`total_seat - booked`. Setting a capacity below the booked seats returns
`409 Conflict`.

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
  creates each agency's owner. Owner deletion and access changes are blocked
  even for this role.
* **owner** (`agency_members.is_owner = true`): logs in as a member and has
  full access within their own agency, including adding and removing staff.
  Ownership is separate from editable role titles and permission grants. No
  current account can delete an owner or change their access; a system-creator
  override is not implemented. Agency/role deletion cannot cascade to an owner.
* **member**: an agency staff account, limited to its own agency. Each management route
  requires a permission (`tour:create`, `member:update`, …) granted through the
  member's role. Permissions are checked against the database on every request,
  so revoking one applies immediately. A member cannot grant permissions they
  do not hold.
* **user**: a customer. Can book tours and manage only their own account.

The permissions the routes use are seeded by migration `000015`. List them with
`SELECT permission_id, name FROM permissions;` to find the ids to pass as
`permissions` when creating members.

The portal shows the signed-in member's name and owner/staff status in its
header, with a **My profile** link for email, phone, role, and agency details.

Agency owners and staff can upload or replace the agency image under
**Settings** and tour covers under **Tours → Edit** in the portal. Both controls
preview JPG, PNG and WebP images up to 10 MB. Staff need `agency:update` for
agency images and `tour:update` for tour images; owners have both permissions.

Restart the API after updating to apply migration `000019`. For existing
agencies, it marks the earliest member with an `Owner` role as the owner, or
the earliest member overall if no role is named `Owner` (case-insensitive).
Review existing role titles before applying this migration if your owner was
not the first member. New owners created through the API are marked explicitly.

### Bootstrapping a new environment

1. Register an account with `POST /users`.
2. Promote it once in the database:
   `UPDATE users SET role = 'super' WHERE email = 'admin@example.com';`
3. Log in, create an agency with `POST /agency`, then create its owner with
   `POST /members/{agency_id}` using `role_name: "Owner"` and the permission
   ids. The portal onboarding form sets this automatically. Only a platform
   admin may create this owner account; each agency has at most one owner.

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
