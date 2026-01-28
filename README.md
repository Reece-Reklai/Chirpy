# Chirpy - Go Backend API for Social Networking

A production-ready RESTful API backend for a Twitter-like social platform, built in Go.

## 🚀 Live Demo
[Link if deployed, or "Local installation instructions below"]

## ⚡ Features
- **JWT Authentication**: Secure token-based auth with refresh tokens
- **User Management**: Registration, login, profile updates
- **CRUD Operations**: Create, read, delete chirps (social posts)
- **Authorization**: Protected endpoints with token validation
- **RESTful API**: Clean HTTP endpoints following REST principles

## 🛠️ Tech Stack
- **Backend**: Go (stdlib `net/http` or framework)
- **Authentication**: JWT (access tokens + refresh tokens)
- **Storage**: JSON file-based database (or PostgreSQL if you upgraded)
- **API Design**: RESTful architecture

## 📊 Architecture
[Insert simple architecture diagram: Client → API Server → JSON DB]

## 🔑 API Endpoints

### User Endpoints
- `POST /api/users` - Register new user
- `PUT /api/users` - Update user profile (auth required)

### Authentication
- `POST /api/login` - User login (returns JWT)
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token (auth required)

### Chirp Endpoints
- `POST /api/chirps` - Create chirp (auth required)
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp (auth required)

## 💻 Local Setup
```bash
git clone https://github.com/Reece-Reklai/chirpy.git
cd chirpy
cp .env.example .env  # Add your JWT secret
go mod download
go run main.go
