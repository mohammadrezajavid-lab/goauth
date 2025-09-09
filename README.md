# GoAuth - OTP-Based Authentication Service

**GoAuth** is a backend service written in **Golang** that provides a secure, OTP-based login and registration system,
along with essential user management features. The project is built following Clean Architecture principles to ensure
maintainability, scalability, and separation of concerns. It features a RESTful API, rate limiting, JWT-based security
for protected endpoints, and is fully containerized for easy setup and deployment.

### Installation & Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/mohammadrezajavid-lab/goauth.git
   cd goauth
   ```

2. **Configuration:**

    * The main configuration file is `config.yml` located at `deploy/goauth/development/config.yml`.
    * It includes settings for the HTTP server, PostgreSQL database, JWT, and logger.
    * You can override configurations using environment variables with the prefix `AUTH_`. For example,
      `postgres.password` in YAML becomes `AUTH_POSTGRES_PASSWORD` as an environment variable.

3. **Install Go Dependencies:**

   ```bash
   go mod tidy
   ```

### Development Workflow with Makefile

This project uses a `Makefile` as the single entry point for all common development tasks.

* To see a full list of available commands, run:
  ```bash
  make help
  ```

#### Running with Docker (Recommended)

This is the simplest way to run the entire application stack.

* **Build and Start the Application:**
  This command builds the Go application, starts the PostgreSQL container, and runs database migrations automatically.

  ```bash
  make up
  ```

* **Follow Logs:**

  ```bash
  make logs
  ```

* **Stop the Application:**
  This stops the running containers without deleting any data.

  ```bash
  make stop
  ```

* **Tear Down Everything:**
  This stops and removes all containers, networks, and data volumes.

  ```bash
  make down
  ```

#### Running Locally (for Go Development)

This workflow is ideal when you want to run the Go service directly on your machine for faster development and
debugging.

1. **Start the Database:**
   First, start only the PostgreSQL database using Docker.

   ```bash
   make up-db
   ```

2. **Run the Go Service:**
   This command starts the main application server locally and automatically applies any pending database migrations.

   ```bash
   make run
   ```

The HTTP server will start on `127.0.0.1:8080` by default.

### API Documentation & Endpoints

The API is documented using Swagger. Once the server is running, you can access the interactive documentation at:
**[http://localhost:8080/swagger/](http://localhost:8080/swagger/)**

Below are detailed examples for each endpoint.

-----

#### **Authentication Endpoints (Public)**

These endpoints are used for user login and registration.

* **`POST /v1/auth/generateotp`**: Initiates the login/registration process by generating an OTP.

    * **Rate Limit**: 3 requests per phone number every 10 minutes.

  **cURL Example:**

  ```bash
  curl -X POST http://localhost:8080/v1/auth/generateotp \
  -H "Content-Type: application/json" \
  -d '{
      "phone_number": "+989123456789"
  }'
  ```

  **How to get the OTP Code:**
  For testing purposes, the generated OTP is not sent via SMS. You must check the application's console logs to find the
  code. Look for a log entry similar to this:

  ```json
  {"time":"2025-09-09T23:52:07.7497003+03:30","level":"INFO","msg":"OTP code generated successfully","phone_number":"+989123456789","otp_code":"626249"}
  ```

  In this example, the OTP is `626249`.

  **Success Response (200 OK):**

  ```json
  {
      "message": "OTP code has been generated and printed to the console."
  }
  ```

* **`POST /v1/auth/verify`**: Verifies the OTP and returns a JWT token.
  **cURL Example:**

  ```bash
  # Replace "123456" with the OTP from your console logs
  curl -X POST http://localhost:8080/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d '{
      "phone_number": "+989123456789",
      "otp": "123456"
  }'
  ```

  **Success Response (200 OK):**

  ```json
  {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "is_new": true
  }
  ```

-----

#### **User Management Endpoints (Protected)**

These endpoints require a valid JWT token in the `Authorization: Bearer <token>` header.

* **`GET /v1/users/{id}`**: Retrieves details for a single user by their ID.

  **cURL Example:**

  ```bash
  curl -X GET http://localhost:8080/v1/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
  ```

  **Success Response (200 OK):**

  ```json
  {
      "id": 1,
      "phone_number": "+989123456789",
      "created_at": "2025-09-09T18:30:00Z",
      "updated_at": "2025-09-09T18:30:00Z"
  }
  ```

* **`GET /v1/users`**: Retrieves a paginated and searchable list of users.

    * **Query Parameters**: `page` (int), `pageSize` (int), `search` (string).

  **cURL Example (Pagination & Search):**

  ```bash
  curl -X GET "http://localhost:8080/v1/users?page=1&pageSize=5&search=912" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
  ```

  **Success Response (200 OK):**

  ```json
  {
      "users": [
          {
              "id": 1,
              "phone_number": "+989121112233",
              "created_at": "2025-09-09T18:30:00Z",
              "updated_at": "2025-09-09T18:30:00Z"
          }
      ],
      "metadata": {
          "currentPage": 1,
          "pageSize": 5,
          "totalRecords": 1,
          "totalPages": 1
      }
  }
  ```

-----

### Database Choice Justification

- **PostgreSQL (for Users)**: We chose PostgreSQL for storing user data due to its robustness, reliability, and strong
  support for data integrity through constraints (like `UNIQUE` for phone numbers). Its rich feature set and excellent
  performance make it a standard choice for production-ready applications.

- **In-Memory Store (for OTPs)**: OTP data is temporary and high-traffic. Using a simple, concurrent in-memory map is
  significantly faster than writing to a disk-based database. This approach reduces latency and minimizes unnecessary
  load on the primary PostgreSQL database. For a multi-node deployment, this could be easily swapped with a distributed
  cache like Redis.

### Technology Stack

- **Language**: Go
- **Framework**: Echo
- **Database**: PostgreSQL
- **DB Driver**: pgx
- **Migrations**: sql-migrate
- **Configuration**: koanft
- **Validation**: ozzo-validation
- **Containerization**: Docker, Docker Compose
- **API Documentation**: Swaggo