# 🚀 Smart File API

A production-ready REST API for intelligent file management with JWT authentication, Redis caching, and advanced file processing capabilities. Built with Go (Golang) and designed for scalability and performance.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-success)

## ✨ Features

### Core Features
- 🔐 **JWT Authentication** - Secure user registration and login
- 📁 **File Management** - Upload, retrieve, update, and delete files
- 🗄️ **SQLite Database** - Lightweight database with GORM ORM
- ⚡ **Redis Caching** - 5-minute cache for improved performance
- 🔄 **Background Processing** - Asynchronous file processing
- 🗑️ **Soft & Hard Delete** - Flexible file deletion with restore capability

### Advanced Features
- 📊 **Pagination & Filtering** - Query files with page, limit, type, status, and search
- 📈 **Statistics Dashboard** - Real-time metrics on files, storage, and activity
- 📝 **Logging & Monitoring** - JSON-formatted logs with request tracking
- 🔍 **Swagger Documentation** - Interactive API documentation
- 🎯 **Rate Limiting** - Prevent API abuse (100 requests/hour)
- 📉 **Performance Metrics** - System health monitoring

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin Web Framework
- **Database:** SQLite with GORM
- **Cache:** Redis
- **Authentication:** JWT (golang-jwt/jwt)
- **Documentation:** Swagger (swaggo)
- **Logging:** Logrus

## 📋 Prerequisites

- Go 1.21 or higher
- Redis server
- Git

## 🚀 Installation

### 1. Clone the repository
git clone https://github.com/yourusername/smart-file-api.git
cd smart-file-api


### 2. Install dependencies
go mod download


### 3. Install Swagger CLI (optional, for regenerating docs)
go install github.com/swaggo/swag/cmd/swag@latest


### 4. Start Redis server
Windows
redis-server

Linux/Mac
sudo service redis-server start

Docker
docker run -d -p 6379:6379 --name redis redis:alpine


### 5. Run the application
go run main.go


The server will start on `http://localhost:8080`

## 📚 API Documentation

Once the server is running, access the interactive Swagger documentation:

**Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Quick Start

#### 1. Register a new user
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
"name": "John Doe",
"email": "john@example.com",
"password": "password123"
}


#### 2. Login
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
"email": "john@example.com",
"password": "password123"
}


**Response:**
{
"status": "success",
"message": "Login successful",
"data": {
"user": {
"id": 1,
"name": "John Doe",
"email": "john@example.com"
},
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
}


#### 3. Upload a file (Protected)
POST http://localhost:8080/api/files/upload
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: multipart/form-data

file: [your file]


#### 4. Get all files with pagination
GET http://localhost:8080/api/files/?page=1&limit=10&type=image&sort=file_size&order=desc
Authorization: Bearer YOUR_JWT_TOKEN


## 🔑 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | ❌ |
| POST | `/api/auth/login` | User login | ❌ |

### Files
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/files/upload` | Upload file (max 10MB) | ✅ |
| GET | `/api/files/` | Get all files (with pagination) | ✅ |
| GET | `/api/files/:id` | Get file details | ✅ |
| DELETE | `/api/files/:id` | Soft delete file | ✅ |
| DELETE | `/api/files/:id/permanent` | Hard delete file | ✅ |
| GET | `/api/files/deleted` | Get deleted files | ✅ |
| POST | `/api/files/:id/restore` | Restore deleted file | ✅ |
| GET | `/api/files/statistics` | Get file statistics | ✅ |

### Monitoring
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | ❌ |
| GET | `/api/metrics` | System metrics | ✅ |
| GET | `/api/logs` | Application logs | ✅ |

## 🎯 Query Parameters

### Pagination & Filtering (GET /api/files/)

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `page` | int | Page number (default: 1) | `?page=2` |
| `limit` | int | Items per page (max: 100, default: 10) | `?limit=20` |
| `type` | string | Filter by file type | `?type=image` |
| `status` | string | Filter by status | `?status=completed` |
| `sort` | string | Sort by field | `?sort=file_size` |
| `order` | string | Sort order (asc/desc) | `?order=desc` |
| `search` | string | Search by filename | `?search=document` |

### Supported File Types
- `image` - jpg, jpeg, png, gif, bmp
- `audio` - mp3, wav, flac, m4a, ogg
- `video` - mp4, avi, mkv, mov
- `document` - pdf, doc, docx, txt
- `other` - all other types

## 📊 Response Format

### Success Response
{
"status": "success",
"message": "Operation successful",
"data": {
// Response data
}
}


### Error Response
{
"status": "error",
"message": "Error description"
}


### Paginated Response
{
"status": "success",
"message": "Files retrieved successfully",
"data": {
"files": [...],
"pagination": {
"page": 1,
"limit": 10,
"total_rows": 45,
"total_pages": 5
},
"filter": {
"type": "image",
"sort_by": "created_at",
"sort_order": "desc"
}
}
}


## 🔒 Security

- **JWT Authentication** - Secure token-based authentication
- **Password Hashing** - Bcrypt with cost factor 14
- **User Isolation** - Users can only access their own files
- **Rate Limiting** - Prevent API abuse
- **Input Validation** - Request validation with Gin binding
- **CORS** - Configurable cross-origin resource sharing

## 🚦 Caching Strategy

- **Cache Duration:** 5 minutes
- **Cache Key:** MD5 hash of (endpoint + user_id)
- **Cache Invalidation:** Automatic on POST/DELETE operations
- **Cache Headers:** `X-Cache: HIT/MISS` for monitoring

### Performance Impact
- **First Request (Cache MISS):** ~50ms (database query)
- **Cached Request (Cache HIT):** ~5ms (10x faster!) ⚡

## 📈 Monitoring & Metrics

### System Metrics Endpoint
GET http://localhost:8080/api/metrics
Authorization: Bearer YOUR_JWT_TOKEN


**Response includes:**
- System uptime
- Memory usage (Alloc, Total, Sys)
- Number of goroutines
- Garbage collection stats
- Database statistics
- Redis cache stats
- Storage usage

### Logging
- **Format:** JSON
- **File:** `app.log`
- **Fields:** timestamp, status_code, latency, method, path, user_id, client_ip

## 📁 Project Structure

smart-file-api/
├── config/
│ ├── database.go # Database configuration
│ ├── redis.go # Redis configuration
│ └── logger.go # Logger setup
├── controllers/
│ ├── auth.go # Authentication handlers
│ ├── file.go # File management handlers
│ └── monitoring.go # Monitoring endpoints
├── middleware/
│ ├── auth.go # JWT authentication middleware
│ ├── cache.go # Caching middleware
│ └── logger.go # Request logging middleware
├── models/
│ ├── user.go # User model
│ └── file.go # File model
├── routes/
│ └── api.go # Route definitions
├── utils/
│ ├── jwt.go # JWT utilities
│ ├── password.go # Password hashing
│ ├── response.go # Response helpers
│ └── pagination.go # Pagination utilities
├── uploads/ # File storage directory
├── docs/ # Swagger documentation
├── main.go # Application entry point
├── go.mod # Go module dependencies
└── README.md # This file


## 🧪 Testing

### Manual Testing with cURL

**Register:**
curl -X POST http://localhost:8080/api/auth/register
-H "Content-Type: application/json"
-d '{"name":"Test User","email":"test@example.com","password":"password123"}'


**Login:**
curl -X POST http://localhost:8080/api/auth/login
-H "Content-Type: application/json"
-d '{"email":"test@example.com","password":"password123"}'


**Upload File:**
curl -X POST http://localhost:8080/api/files/upload
-H "Authorization: Bearer YOUR_TOKEN"
-F "file=@/path/to/your/file.jpg"


## 🐛 Troubleshooting

### Redis Connection Failed
⚠️ Redis connection failed: dial tcp [::1]:6379: connect: connection refused

**Solution:** Start Redis server or use Docker:
docker run -d -p 6379:6379 --name redis redis:alpine


### Port Already in Use
Error: listen tcp :8080: bind: address already in use

**Solution:** Change port in `main.go` or kill the process:
Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

Linux/Mac
lsof -i :8080
kill -9 <PID>


## 🔄 Environment Variables (Optional)

Create `.env` file:
PORT=8080
JWT_SECRET=your-secret-key-here
REDIS_HOST=localhost:6379
DB_PATH=smart-file-api.db
LOG_LEVEL=info


## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Your Name**
- GitHub: [@yourusername](https://github.com/yourusername)
- LinkedIn: [Your Name](https://linkedin.com/in/yourprofile)
- Email: your.email@example.com

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Redis](https://redis.io/)
- [Swagger](https://swagger.io/)

## 📸 Screenshots

### Swagger API Documentation
![Swagger UI](screenshots/swagger.png)

### API Response Example
![API Response](screenshots/response.png)

---

⭐ If you find this project useful, please give it a star!
📝 Additional Files
1. Create .gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output
*.out

# Go workspace file
go.work

# Database
*.db
*.db-shm
*.db-wal

# Uploads
uploads/*
!uploads/.gitkeep

# Logs
*.log
app.log

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
2. Create uploads/.gitkeep
bash
# Windows
type nul > uploads/.gitkeep

# Linux/Mac
touch uploads/.gitkeep
3. Create LICENSE (MIT)
MIT License

Copyright (c) 2025 Your Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
