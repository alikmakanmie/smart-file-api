# 🚀 Smart File API

A production-ready REST API for intelligent file management with JWT authentication, Redis caching, and advanced file processing capabilities. Built with Go (Golang) and designed for scalability and performance.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-success)

---

## ✨ Features

### Core Features
- 🔐 **JWT Authentication** - Secure user registration and login
- 📁 **File Management** - Upload, retrieve, update, and delete files
- 🗄️ **SQLite Database** - Lightweight database with GORM ORM
- ⚡ **Redis Caching** - 5-minute cache for improved performance (10x faster!)
- 🔄 **Background Processing** - Asynchronous file processing
- 🗑️ **Soft & Hard Delete** - Flexible file deletion with restore capability

### Advanced Features
- 📊 **Pagination & Filtering** - Query files with page, limit, type, status, and search
- 📈 **Statistics Dashboard** - Real-time metrics on files, storage, and activity
- 📝 **Logging & Monitoring** - JSON-formatted logs with request tracking
- 🔍 **Swagger Documentation** - Interactive API documentation
- 🔒 **Security** - Password hashing, input validation, user isolation

---

## 🛠️ Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.21+ |
| **Framework** | Gin Web Framework |
| **Database** | SQLite with GORM |
| **Cache** | Redis |
| **Authentication** | JWT (golang-jwt/jwt) |
| **Documentation** | Swagger (swaggo) |
| **Logging** | Logrus |

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Redis server (optional, but recommended)
- Git

### Installation

```
# 1. Clone the repository
git clone https://github.com/alikmakanmie/smart-file-api.git
cd smart-file-api

# 2. Install dependencies
go mod download

# 3. Start Redis (optional)
docker run -d -p 6379:6379 --name redis redis:alpine

# 4. Run the application
go run main.go
```

The server will start on [**http://localhost:8080**](http://localhost:8080)

📖 **Swagger Documentation**: http://localhost:8080/swagger/index.html

---

## 📚 API Documentation

### Quick Example

#### 1️⃣ Register
```
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### 2️⃣ Login
```
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

#### 3️⃣ Upload File
```
POST http://localhost:8080/api/files/upload
Authorization: Bearer YOUR_TOKEN
Content-Type: multipart/form-data

file: [your file]
```

#### 4️⃣ Get Files with Pagination
```
GET http://localhost:8080/api/files/?page=1&limit=10&type=image
Authorization: Bearer YOUR_TOKEN
```

---

## 🔑 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/register` | Register new user | ❌ |
| POST | `/api/auth/login` | User login | ❌ |

### File Management
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/files/upload` | Upload file (max 10MB) | ✅ |
| GET | `/api/files/` | Get all files with pagination | ✅ |
| GET | `/api/files/:id` | Get file details | ✅ |
| DELETE | `/api/files/:id` | Soft delete file | ✅ |
| DELETE | `/api/files/:id/permanent` | Hard delete file | ✅ |
| GET | `/api/files/deleted` | Get deleted files | ✅ |
| POST | `/api/files/:id/restore` | Restore deleted file | ✅ |
| GET | `/api/files/statistics` | Get file statistics | ✅ |

### Monitoring
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | Health check | ❌ |
| GET | `/api/metrics` | System metrics | ✅ |
| GET | `/api/logs` | Application logs | ✅ |

---

## 🎯 Query Parameters

**Pagination & Filtering** (GET `/api/files/`)

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `page` | int | `?page=2` | Page number (default: 1) |
| `limit` | int | `?limit=20` | Items per page (max: 100) |
| `type` | string | `?type=image` | Filter by file type |
| `status` | string | `?status=completed` | Filter by status |
| `sort` | string | `?sort=file_size` | Sort by field |
| `order` | string | `?order=desc` | asc or desc |
| `search` | string | `?search=photo` | Search by filename |

**Supported File Types**: `image`, `audio`, `video`, `document`, `other`

---

## 📁 Project Structure

```
smart-file-api/
├── config/
│   ├── database.go          # Database configuration
│   ├── redis.go             # Redis configuration
│   └── logger.go            # Logger setup
├── controllers/
│   ├── auth.go              # Authentication handlers
│   ├── file.go              # File management handlers
│   └── monitoring.go        # Monitoring endpoints
├── middleware/
│   ├── auth.go              # JWT authentication middleware
│   ├── cache.go             # Caching middleware
│   └── logger.go            # Request logging middleware
├── models/
│   ├── user.go              # User model
│   └── file.go              # File model
├── routes/
│   └── api.go               # Route definitions
├── utils/
│   ├── jwt.go               # JWT utilities
│   ├── password.go          # Password hashing
│   ├── response.go          # Response helpers
│   └── pagination.go        # Pagination utilities
├── uploads/                 # File storage directory
├── docs/                    # Swagger documentation
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
└── README.md                # This file
```

---

## 📊 Response Format

### Success Response
```
{
  "status": "success",
  "message": "Operation successful",
  "data": {
    "files": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total_rows": 45,
      "total_pages": 5
    }
  }
}
```

### Error Response
```
{
  "status": "error",
  "message": "Error description"
}
```

---

## 🔒 Security Features

- ✅ JWT-based authentication
- ✅ Password hashing with Bcrypt (cost factor 14)
- ✅ User isolation (users can only access their own files)
- ✅ Input validation with Gin binding
- ✅ File type validation
- ✅ File size limits (10MB max)

---

## ⚡ Caching Strategy

- **Cache Duration**: 5 minutes
- **Cache Key**: MD5 hash of (endpoint + user_id)
- **Cache Invalidation**: Automatic on POST/DELETE operations
- **Performance**: 
  - First Request (MISS): ~50ms
  - Cached Request (HIT): ~5ms (**10x faster!** ⚡)

Check cache status via `X-Cache: HIT/MISS` response header.

---

## 🧪 Testing

### Using cURL

```
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Upload File
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/file.jpg"
```

### Using Swagger UI

Visit [**http://localhost:8080/swagger/index.html**](http://localhost:8080/swagger/index.html) for interactive API testing.

---

## 🐛 Troubleshooting

### Redis Connection Failed
```
⚠️ Redis connection failed (caching will be disabled)
```
**Solution**: Start Redis server
```
docker run -d -p 6379:6379 --name redis redis:alpine
```

### Port Already in Use
```
Error: listen tcp :8080: bind: address already in use
```
**Solution**: Change port in `main.go` or kill the process

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Malik**
- GitHub: [@alikmakanmie](https://github.com/alikmakanmie)
- Project: [Smart File API](https://github.com/alikmakanmie/smart-file-api)

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Redis](https://redis.io/)
- [Swagger](https://swagger.io/)

---

## 📸 Screenshots

### Swagger API Documentation
![Swagger UI](screenshots/swagger.png)

### API Response Example
![API Response](screenshots/response.png)

---

⭐ **If you find this project useful, please give it a star!**
```

***

## 💾 Commit & Push

```bash
# Update README
git add README.md
git commit -m "docs: Fix README formatting and improve structure"
git push
```

***

## ✅ Hasil Akhir

Setelah push, README Anda akan tampil **rapi dan professional** seperti ini:
- ✅ Tree structure ter-render dengan benar
- ✅ Table alignment sempurna
- ✅ Code blocks dengan syntax highlighting
- ✅ Emoji dan badges tampil dengan baik

**Refresh halaman GitHub Anda dalam 1-2 menit** untuk melihat hasil yang sudah rapi! 🎉
