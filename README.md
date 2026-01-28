# Chirpy - Production-Ready Social Media API

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=flat&logo=postgresql)

*A scalable, secure REST API for Twitter-like social networking built in Go*

</div>

## Overview

Chirpy is a production-ready backend API that demonstrates enterprise-grade Go development practices. This project showcases advanced backend engineering concepts including JWT authentication, database design, API security, and scalable architecture patterns. Built as a comprehensive example for technical evaluation and learning purposes.

## Key Features

### Security & Authentication
- **JWT Token System**: Secure access tokens with configurable expiration
- **Refresh Tokens**: Long-lived token rotation for enhanced security
- **Password Hashing**: Argon2ID implementation for industry-standard password security
- **API Key Authentication**: Webhook integration with secure API key validation
- **Content Moderation**: Automated profanity filtering with word censorship

### Architecture & Design
- **Clean Architecture**: Separated concerns with organized package structure
- **Database Integration**: PostgreSQL with SQLC for type-safe database operations
- **Migration Management**: Schema versioning with migration scripts
- **RESTful Design**: Compliant API design following REST principles
- **Middleware Pattern**: Request tracking, metrics collection, and authentication middleware

### Core Functionality
- **User Management**: Registration, authentication, profile updates with premium tier support
- **Social Features**: Create, read, delete chirps with user authorization
- **Content Operations**: Advanced querying with sorting, filtering by author
- **Admin Dashboard**: Metrics collection and administrative operations
- **Webhook Integration**: External service integration for user upgrades

## Technical Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Runtime** | Go 1.24+ | High-performance backend |
| **Database** | PostgreSQL 15+ | Relational data storage |
| **Authentication** | JWT + Argon2ID | Secure user authentication |
| **Database Access** | SQLC | Type-safe SQL operations |
| **HTTP Server** | Go stdlib `net/http` | Lightweight HTTP handling |
| **Configuration** | Godotenv | Environment variable management |

## Project Structure

```
chirpy/
├── main.go                    # Application entry point
├── router.go                  # HTTP routing and handlers
├── go.mod/go.sum             # Go dependency management
├── internal/
│   ├── auth/                 # Authentication utilities
│   └── database/             # Database models and queries
├── sql/
│   ├── schema/               # Database migrations
│   └── queries/              # SQLC query definitions
├── public/                   # Static assets
├── test/                     # Test suites
└── docs/                     # Documentation
```

## Quick Start

### Prerequisites
- Go 1.24 or higher
- PostgreSQL 15 or higher
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/Reece-Reklai/chirpy.git
cd chirpy
```

2. **Set up environment variables**
```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
# Required: DB_URL, SECRET, PLATFORM, POLKA_KEY
```

3. **Install dependencies**
```bash
go mod download
```

4. **Set up database**
```bash
# Run database migrations
psql $DB_URL -f sql/schema/001_users.sql
psql $DB_URL -f sql/schema/002_chirps.sql
psql $DB_URL -f sql/schema/003_refresh_tokens.sql
```

5. **Generate database code**
```bash
sqlc generate
```

6. **Start the server**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Docs

### Authentication Endpoints

#### Register User
```http
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com", 
  "password": "securepassword123"
}
```

#### Refresh Token
```http
POST /api/refresh
Authorization: Bearer <refresh_token>
```

#### Revoke Token
```http
POST /api/revoke
Authorization: Bearer <refresh_token>
```

### User Management

#### Update Profile
```http
PUT /api/users
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

### Chirp Operations

#### Create Chirp
```http
POST /api/chirps
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "body": "Hello, Chirpy world!"
}
```

#### Get All Chirps
```http
GET /api/chirps?author_id=<uuid>&sort=asc|desc
```

#### Get Specific Chirp
```http
GET /api/chirps/{chirpID}
```

#### Delete Chirp
```http
DELETE /api/chirps/{chirpID}
Authorization: Bearer <access_token>
```

### Admin Endpoints

#### Get Metrics
```http
GET /admin/metrics
```

#### Reset Development Environment
```http
POST /admin/reset
```

### Webhook Integration

#### User Upgrade Webhook
```http
POST /api/polka/webhooks
Authorization: ApiKey <polka_key>
Content-Type: application/json

{
  "event": "user.upgraded",
  "data": {
    "user_id": "user-uuid"
  }
}
```

## Implementation Goals

### Code Quality
- **Type Safety**: Full type safety with SQLC and strong Go typing
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Security**: Input validation, SQL injection prevention, secure authentication
- **Testing**: Unit tests for critical authentication and business logic

### Performance Considerations
- **Database Indexing**: Optimized queries with proper indexing
- **Connection Pooling**: Efficient database connection management
- **Middleware Metrics**: Request tracking for performance monitoring
- **Resource Management**: Proper resource cleanup and memory management

### Security Implementation
- **Password Security**: Argon2ID with proper salt and parameters
- **Token Security**: JWT with proper expiration and refresh rotation
- **Input Validation**: Request validation and sanitization
- **Authorization**: Role-based access control and resource ownership checks

## Architecture

### Database Design
- **Normalized Schema**: Proper relational design with foreign key constraints
- **UUID Primary Keys**: Globally unique identifiers for distributed systems
- **Audit Fields**: Created/updated timestamps for data tracking
- **Cascade Deletion**: Referential integrity maintenance

### Authentication Architecture
- **JWT Implementation**: Custom JWT handling with proper claims validation
- **Token Rotation**: Secure refresh token mechanism preventing token fixation
- **Multi-factor Security**: Bearer token and API key authentication patterns
- **Session Management**: Stateless authentication design for scalability

### API Design Patterns
- **RESTful Compliance**: Proper HTTP method usage and status codes
- **Content Negotiation**: JSON request/response handling
- **Error Responses**: Consistent error format with proper HTTP semantics
- **Versioning Ready**: Structure supports future API versioning

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

Run specific test packages:
```bash
go test ./test/auth_test.go
```
