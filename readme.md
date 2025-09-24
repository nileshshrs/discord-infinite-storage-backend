# Infinite Storage Backend

Infinite Storage is a backend service for secure file storage and management. It leverages MongoDB and Discord to provide a unique approach to file uploads, allowing authenticated users to upload, store, and retrieve files efficiently.

## Features

- **Authenticated File Uploads**: Users can upload files securely via a protected API endpoint. Authentication is handled using JWT tokens and middleware.  
- **Chunked File Storage**: Large files are split into chunks and uploaded to Discord channels, while metadata is stored in MongoDB.  
- **User-specific File Retrieval**: Users can fetch all their uploaded files as JSON, including detailed chunk information.  
- **Secure & Scalable**: Only authorized users can upload and access files, ensuring data integrity and security.

## Folder Structure

```text

discord-infinite-storage/
│
├── main.go
├── application/
│   ├── app.go
│   └── route.go
│
├── bot/
│   └── bot.go
│
├── config/
│   └── config.go
│
├── db/
│   └── database.go
│
├── handler/
│   ├── auth_handler.go
│   ├── file_handler.go
│   ├── upload_handler.go
│   └── user_handler.go
│
├── middlewares/
│   └── middleware.go
│
├── model/
│   ├── file.go
│   ├── session.go
│   └── user.go
│
├── repository/
│   ├── file_repository.go
│   ├── session_repository.go
│   └── user_repository.go
│
├── service/
│   ├── auth_service.go
│   ├── file_service.go
│   ├── upload_service.go
│   └── user_service.go
│
└── utils/
    └── jwt.go
```

## API Endpoints

### Auth
- `POST /api/v1/auth/sign-up` – Register a new user
- `POST /api/v1/auth/sign-in` – Login and receive JWT
- `POST /api/v1/auth/refresh` – Refresh JWT token

### Users
- `GET /api/v1/users/` – Retrieve all users (protected route)

### Files
- `POST /api/v1/files/upload` – Upload a new file (protected route)
- `GET /api/v1/files/` – Retrieve all files uploaded by the authenticated user (protected route)

## Running the Server

### Using Go
```bash
go run main.go
```

### Using Air (Live Reload)
```bash
go install github.com/cosmtrek/air@latest
```

### Run Air
```bash
air
```

## Authentication

The API provides endpoints for user registration, login, and token refresh using JWT.

- **POST** `/api/v1/auth/sign-up`  
  Register a new user.

- **POST** `/api/v1/auth/sign-in`  
  Login and receive an access token.

- **POST** `/api/v1/auth/refresh`  
  Refresh JWT token.

---

## Users

- **GET** `/api/v1/users/` *(Protected)*  
  Retrieve all registered users.

---

## Files

Manage user-specific files with upload and retrieval functionality.

- **POST** `/api/v1/files/upload` *(Protected)*  
  Upload a file to Discord and save metadata.

- **GET** `/api/v1/files/` *(Protected)*  
  Retrieve all files uploaded by the authenticated user.

---

## Future Work

- Implement file merging and download functionality for users.
- Add pagination and filtering for file retrieval.
- Build a frontend interface for easier file management.