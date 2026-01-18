# Cacto CMS

<div align="center">

**Performance-focused, Enterprise-grade CMS built with Go + Templ + HTMX**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 📋 Table of Contents

- [Features](#-features)
- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Usage](#-usage)
- [Architecture](#-architecture)
- [API Endpoints](#-api-endpoints)
- [Development](#-development)
- [Deployment](#-deployment)
- [Troubleshooting](#-troubleshooting)
- [Changelog](#-changelog)
- [Contributing](#-contributing)

---

## ✨ Features

### 🏗️ Architecture
- ✅ **Clean Architecture** - Domain-Driven Design (DDD) approach
- ✅ **Layered Architecture** - Domain, Application, Infrastructure, Interface layers
- ✅ **Dependency Inversion** - Interface-based design
- ✅ **Testable** - Each layer can be tested independently

### 🔐 Security
- ✅ **JWT Authentication** - Secure token-based authentication
- ✅ **Argon2id Password Hashing** - Modern password hashing
- ✅ **Role-Based Access Control (RBAC)** - Admin, Editor, Author, Viewer roles
- ✅ **Permission System** - Granular permission management
- ✅ **Input Validation** - Comprehensive validation system

### 📝 Content Management
- ✅ **Page Management** - Dynamic page creation and management
- ✅ **Component System** - Reusable, database-driven components
- ✅ **Media Management** - File upload and media library
- ✅ **SEO Optimization** - Centralized SEO management
- ✅ **Sitemap Generation** - Automatic sitemap.xml generation

### 🛠️ Developer Experience
- ✅ **Artisan CLI** - Migration and seeding management
- ✅ **Structured Logging** - Go 1.21+ slog integration
- ✅ **Error Handling** - Comprehensive error management
- ✅ **Hot Reload** - Development mode with Air
- ✅ **Type-Safe Templates** - Templ for type-safe HTML

### 🚀 Performance
- ✅ **SQLite Database** - Lightweight, fast database
- ✅ **Component Caching** - Efficient component rendering
- ✅ **Static File Serving** - Optimized asset delivery

---

## 🚀 Quick Start

### Requirements

- Go 1.23 or higher
- Make (optional, for commands)
- SQLite3 (automatically included)
- Node.js 18+ (for Tailwind CSS v4, optional - only needed for CSS build)

### Installation in 3 Steps

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/cacto-cms.git
cd cacto-cms

# 2. Install dependencies
make install

# 3. Setup database and seed data
./artisan migrate:fresh --seed

# 4. Run in development mode (hot reload)
make dev
```

Server: **http://localhost:8080**

### First Login

Login with the admin user created from seed data:

- **URL**: `http://localhost:8080/admin/login`
- **Email**: `admin@cacto-cms.local`
- **Password**: `admin123`

Available pages:
- Home: `/` (Hero + About sections)
- About: `/about` (from seed data)
- Contact: `/contact` (from seed data)
- Admin Login: `/admin/login`
- Admin Dashboard: `/admin/dashboard` (after login)

---

## 📦 Installation

### 1. Install Dependencies

```bash
make install
```

This command installs:
- Go module dependencies
- Templ CLI tool
- Air (for hot reload)

### 2. Database Setup

```bash
# Migration only (without seed)
./artisan migrate:fresh

# Migration + Seed (recommended)
./artisan migrate:fresh --seed
```

### 3. Environment Variables

Create a `.env` file (copy from `.env.example`):

```bash
cp .env.example .env
```

Or set as environment variables:

```bash
export PORT=8080
export BASE_URL=http://localhost:8080
export JWT_SECRET=your-secret-key
export ENV=development
export USE_HTTPS=false  # Set to true if using HTTPS
```

**Important**: Always use a strong `JWT_SECRET` in production!

---

## 💻 Usage

### Makefile Commands

```bash
make help          # Show all commands
make install       # Install dependencies (Go + Node.js)
make dev           # Run with hot reload
make run           # Run normally
make build         # Build (CSS + Go binaries)
make templ         # Generate templ files
make css           # Tailwind CSS v4 build (production, only used classes)
make css-watch     # Tailwind CSS v4 watch (development, JIT)
make clean         # Clean (binary, db, templ files, CSS)
make tidy          # go mod tidy
make test          # Run tests
make artisan       # Run Artisan CLI (ARGS="migrate:fresh --seed")
```

**Tailwind CSS v4.1:**
- Development: `make css-watch` (ayrı terminal, JIT compilation)
- Production: `make build` (otomatik CSS build, only used classes)

### Artisan CLI Commands

```bash
# Migration operations
./artisan migrate:fresh              # Reset database and run migrations
./artisan migrate:fresh --seed       # Migration + seed data

# Or with make
make artisan ARGS="migrate:fresh --seed"
```

### Running the Server

```bash
# Development (hot reload)
make dev

# Production
make build
./cacto-cms
```

---

## 🏛️ Architecture

Cacto CMS is designed according to **Clean Architecture** and **Domain-Driven Design** principles.

### Layers

```
┌─────────────────────────────────────┐
│   Interfaces Layer (HTTP, CLI)      │  ← Communication with external world
├─────────────────────────────────────┤
│   Application Layer (Services)     │  ← Business logic
├─────────────────────────────────────┤
│   Domain Layer (Entities)            │  ← Core business rules
├─────────────────────────────────────┤
│   Infrastructure Layer (DB, etc)    │  ← Technical implementations
└─────────────────────────────────────┘
```

### Proje Yapısı

```
cacto-cms/
├── cmd/
│   ├── server/              # HTTP server entry point
│   └── artisan/             # CLI tool (migrations, seeds)
│
├── app/
│   ├── domain/               # Domain Layer
│   │   ├── page/            # Page entity & repository interface
│   │   ├── component/        # Component entity
│   │   ├── user/            # User entity & roles
│   │   └── media/           # Media entity
│   │
│   ├── application/          # Application Layer
│   │   ├── page/            # Page service (business logic)
│   │   ├── component/       # Component service
│   │   ├── user/            # User service
│   │   ├── auth/            # Authentication service
│   │   └── media/           # Media service
│   │
│   ├── infrastructure/       # Infrastructure Layer
│   │   ├── database/        # DB connection, migrations & seeds
│   │   │   ├── migrations/  # SQL migration files
│   │   │   └── seeds/       # Database seeders
│   │   └── persistence/     # Repository implementations
│   │       ├── page/
│   │       ├── component/
│   │       ├── user/
│   │       └── media/
│   │
│   ├── interfaces/          # Interface Layer
│   │   ├── http/
│   │   │   ├── controller/  # HTTP controllers
│   │   │   ├── middleware/  # Auth, logging, error handling
│   │   │   └── router.go   # Route definitions
│   │   └── templates/       # Templ templates
│   │
│   └── shared/              # Shared utilities
│       ├── errors/          # Error handling
│       ├── logger/         # Structured logging
│       ├── validation/     # Input validation
│       ├── auth/           # JWT & password hashing
│       ├── seo/            # SEO management
│       └── sitemap/        # Sitemap generator
│
├── config/                  # Configuration
├── web/
│   ├── static/             # Static files (CSS, JS)
│   └── uploads/            # User uploads
└── Makefile
```

### Architecture Layers

#### 1. Domain Layer (`app/domain/`)
- **Purpose**: Business entities and domain logic
- **Content**: 
  - Entity definitions (Page, Component, User, Media)
  - Repository interfaces (domain contract)
  - Domain-specific business rules
- **Dependency**: Not dependent on any external layer

#### 2. Application Layer (`app/application/`)
- **Purpose**: Use cases and business logic orchestration
- **Content**:
  - Service implementations
  - Business logic
  - Use case coordination
- **Dependency**: Only depends on Domain layer

#### 3. Infrastructure Layer (`app/infrastructure/`)
- **Purpose**: Technical implementations
- **Content**:
  - Database connections
  - Repository implementations
  - External service integrations
- **Dependency**: Depends on Domain and Application layers

#### 4. Interface Layer (`app/interfaces/`)
- **Purpose**: Communication with external world
- **Content**:
  - HTTP handlers
  - CLI commands (Artisan)
  - Templates
- **Dependency**: Depends on Application layer

#### 5. Shared (`app/shared/`)
- **Purpose**: Common utilities
- **Content**:
  - Error handling
  - Logging
  - Validation
  - SEO utilities
  - Sitemap generator
- **Dependency**: May depend on Domain layer

### Dependency Flow

```
cmd/server/main.go
    ↓
app/interfaces/http (router, controllers)
    ↓
app/application (services)
    ↓
app/domain (entities, interfaces)
    ↑
app/infrastructure/persistence (implementations)
```

**Rule**: Outer layers can depend on inner layers, inner layers cannot depend on outer layers.

### Package Organization

#### Domain-Driven Design (DDD) Approach
- Each domain entity in its own package (`app/domain/page/`)
- Repository interface defined in domain
- Implementation in infrastructure

#### Clean Architecture Principles
- **Separation of Concerns**: Each layer focuses on its own responsibility
- **Dependency Inversion**: Domain layer is not dependent on anything
- **Testability**: Each layer can be tested independently

---

## 🔌 API Endpoints

Cacto CMS has an **API-first** architecture. All endpoints can return both JSON and HTML:
- **JSON**: With `Accept: application/json` header or `Content-Type: application/json`
- **HTML**: Returns HTML by default (for browser requests)

### Public Routes

| Method | Endpoint | Description | Response Type |
|--------|----------|-------------|---------------|
| GET | `/` | Home page | HTML |
| GET | `/{slug}` | Dynamic page by slug | HTML |
| GET | `/sitemap.xml` | Sitemap | XML |
| GET | `/static/*` | Static files | Static |
| GET | `/uploads/*` | Uploaded files | Static |

### Authentication Routes

| Method | Endpoint | Description | Auth Required | Response Type |
|--------|----------|-------------|---------------|---------------|
| POST | `/api/auth/login` | User login (API) | ❌ | JSON |
| POST | `/api/auth/register` | User registration | ❌ | JSON |
| POST | `/api/auth/logout` | User logout | ✅ | JSON |
| GET | `/admin/login` | Admin login page | ❌ | HTML |
| POST | `/admin/login` | Admin login (form/JSON) | ❌ | HTML/JSON |

### Protected Routes

Protected routes require authentication. Admin/Editor routes require specific roles.

| Method | Endpoint | Description | Roles | Response Type |
|--------|----------|-------------|-------|---------------|
| GET | `/admin/dashboard` | Admin dashboard | admin, editor | HTML/JSON |
| POST | `/admin/logout` | Admin logout | admin, editor | HTML/JSON |

### API-First Architecture

Cacto CMS has an **API-first** architecture. All controllers can return both HTML and JSON:

- **JSON Response**: When request is made with `Accept: application/json` header
- **HTML Response**: By default (for browser requests)

This allows:
- ✅ Existing HTML pages continue to work
- ✅ API endpoints return JSON from the same controller
- ✅ A single controller serves both web and API needs
- ✅ JSON API can be quickly used when needed

### Request/Response Examples

#### Admin Login (HTML - Tarayıcı)

Tarayıcıda: `http://localhost:8080/admin/login`

#### Admin Login (JSON - API)

```bash
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "email": "admin@cacto-cms.local",
    "password": "admin123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@cacto-cms.local",
    "name": "Admin User",
    "role": "admin",
    "is_active": true
  }
}
```

#### Admin Dashboard (JSON - API)

```bash
curl -X GET http://localhost:8080/admin/dashboard \
  -H "Accept: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Cookie: auth_token=YOUR_JWT_TOKEN"
```

Response:
```json
{
  "email": "admin@cacto-cms.local",
  "role": "admin"
}
```

#### API Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@cacto-cms.local",
    "password": "admin123"
  }'
```

#### Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "New User"
  }'
```

---

## 🛠️ Development

### Adding a New Domain Entity

1. `app/domain/{entity}/entity.go` - Entity definitions
2. `app/domain/{entity}/repository.go` - Repository interface
3. `app/application/{entity}/service.go` - Business logic
4. `app/infrastructure/persistence/{entity}/repository.go` - Implementation
5. `app/interfaces/http/controller/{entity}_controller.go` - HTTP controller

### Adding a New Component

1. Update component entity: `app/domain/component/entity.go`
2. Add to component renderer: `app/shared/component/renderer.go`
3. Create template: `app/interfaces/templates/components/new_component.templ`
4. Run `make templ`

### Adding a Migration

1. Create `app/infrastructure/database/migrations/003_new_migration.sql`
2. File name must be sequential (001_, 002_, etc.)
3. Write SQL file
4. Migration runs automatically

### Adding Seed Data

1. Create `app/infrastructure/database/seeds/03_new_data.go`
2. Write seed function
3. Add to `app/infrastructure/database/seeds/seeder.go`
4. Run `./artisan migrate:fresh --seed`

### Updating Templates

1. Edit `.templ` file
2. Run `make templ` or `templ generate`
3. `*_templ.go` file is automatically generated

### Updating Tailwind CSS v4.1

1. Edit `web/static/css/input.css` file
2. Development: Auto-updates if `make css-watch` is running (JIT)
3. Production: `make build` automatically builds CSS (only used classes)
4. See [TAILWIND.md](TAILWIND.md) for detailed information

### Database Operations

```bash
# Open SQLite shell
sqlite3 cacto.db

# List tables
.tables

# View page data
SELECT * FROM pages;

# Add new page
INSERT INTO pages (slug, title, content, status) 
VALUES ('test', 'Test Page', 'Test content', 'published');

# Exit
.quit
```

### Testing

```bash
# Run all tests
make test

# Test a specific package
go test ./app/application/page/...
```

---

## 🚀 Deployment

### Production Build

```bash
# Build
make build

# Environment variables ayarla
export ENV=production
export JWT_SECRET=your-production-secret
export BASE_URL=https://yourdomain.com

# Çalıştır
./cacto-cms
```

### Systemd Service

#### 1. Server Preparation

```bash
# Connect to server via SSH
ssh root@your-server-ip

# Update system
apt update && apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install required tools
apt install -y git nginx sqlite3
```

#### 2. Deploy Project

```bash
# Clone project
git clone https://github.com/youruser/cacto-cms.git /home/cacto/apps/cacto-cms
cd /home/cacto/apps/cacto-cms

# Install dependencies
make install

# Build
make build

# Environment variables
cat > .env << 'EOF'
PORT=8080
BASE_URL=https://yourdomain.com
DB_PATH=/home/cacto/apps/cacto-cms/cacto.db
UPLOAD_DIR=/home/cacto/apps/cacto-cms/web/uploads
ENV=production
JWT_SECRET=your-production-secret
EOF

# Set directory permissions
chmod +x cacto-cms
mkdir -p web/uploads logs
chmod 755 web/uploads
```

#### 3. Create Systemd Service

```bash
sudo nano /etc/systemd/system/cacto-cms.service
```

```ini
[Unit]
Description=Cacto CMS
After=network.target

[Service]
Type=simple
User=cacto
WorkingDirectory=/home/cacto/apps/cacto-cms
ExecStart=/home/cacto/apps/cacto-cms/cacto-cms
Restart=on-failure
RestartSec=5s
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"

# Logs
StandardOutput=append:/home/cacto/apps/cacto-cms/logs/app.log
StandardError=append:/home/cacto/apps/cacto-cms/logs/error.log

[Install]
WantedBy=multi-user.target
```

```bash
# Activate service
systemctl daemon-reload
systemctl enable cacto-cms
systemctl start cacto-cms

# Check status
systemctl status cacto-cms

# Watch logs
journalctl -u cacto-cms -f
```

#### 4. Nginx Configuration

```nginx
# HTTP -> HTTPS Redirect
server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com www.yourdomain.com;
    
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }
    
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS Server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL Sertifikaları (Let's Encrypt ile alınacak)
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # SSL Ayarları
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Gzip Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript 
               application/x-javascript application/xml+rss 
               application/javascript application/json;

    # Static files
    location /static/ {
        alias /home/cacto/apps/cacto-cms/web/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location /uploads/ {
        alias /home/cacto/apps/cacto-cms/web/uploads/;
        expires 30d;
        add_header Cache-Control "public";
    }

    # Sitemap
    location = /sitemap.xml {
        alias /home/cacto/apps/cacto-cms/web/static/sitemap.xml;
        expires 1d;
    }

    # Go App (reverse proxy)
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        proxy_cache_bypass $http_upgrade;
        
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    access_log /var/log/nginx/cacto-cms-access.log;
    error_log /var/log/nginx/cacto-cms-error.log;
}
```

---

## 🐛 Troubleshooting

### 1. "templ: command not found"

```bash
go install github.com/a-h/templ/cmd/templ@latest
export PATH=$PATH:$(go env GOPATH)/bin
# or Makefile already handles it: make install
```

### 2. "pattern migrations/*.sql: no matching files found"

✅ FIXED: 
- Migrations: `app/infrastructure/database/migrations/`
- Seeds: `app/infrastructure/database/seeds/` (same location as migrations)

### 3. "go.sum error"

```bash
make tidy
```

### 4. Port already in use

```bash
# Create .env file
cp .env.example .env
# Change PORT=8081 etc.
```

### 5. Database error

```bash
# Reset database
make clean
./artisan migrate:fresh --seed
```

### 6. Build error

```bash
# Update dependencies
make tidy
make install
make build
```

---

## 🔒 Security

### Security Features

Cacto CMS implements enterprise-level security measures:

#### ✅ Authentication & Authorization
- **JWT (JSON Web Token)** based authentication
- **Argon2id Password Hashing** - Modern, secure password hashing
- **Role-Based Access Control (RBAC)** - Admin, Editor, Author, Viewer roles
- **Cookie-based** - Automatic cookie management for browsers
- **Header-based** - `Authorization: Bearer TOKEN` for API requests
- **Secure Cookies** - `Secure` flag automatically enabled when HTTPS is detected

#### ✅ Input Validation & Sanitization
- **HTML Sanitization** - All user-generated content is sanitized using bluemonday
- **Input Validation** - Comprehensive struct-based validation
- **SQL Injection Protection** - Prepared statements throughout
- **Path Traversal Prevention** - All file paths are validated

#### ✅ Security Headers
- **X-Content-Type-Options: nosniff** - Prevents MIME type sniffing
- **X-Frame-Options: DENY** - Prevents clickjacking
- **X-XSS-Protection: 1; mode=block** - XSS protection
- **Content-Security-Policy** - Restricts resource loading
- **Strict-Transport-Security** - HSTS for HTTPS connections
- **Referrer-Policy** - Controls referrer information

#### ✅ CORS Protection
- **Development**: Allows all origins (`*`)
- **Production**: Only allows configured origins (from `BASE_URL` or `ALLOWED_ORIGINS`)
- **Configurable**: Set `ALLOWED_ORIGINS` environment variable for custom origins

#### ✅ Rate Limiting
- **Auth Endpoints**: 5 requests per minute (login, register)
- **API Endpoints**: 60 requests per minute
- **IP-based**: Rate limiting by client IP address

#### ✅ Error Handling
- **Production Mode**: Generic error messages (no internal details exposed)
- **Development Mode**: Detailed error messages for debugging
- **Structured Errors**: Consistent error response format

#### ✅ File Upload Security
- **MIME Type Validation** - Validates file types
- **File Size Limits** - Configurable max file size (default: 10MB)
- **Filename Sanitization** - Prevents path traversal and malicious filenames
- **Content Validation** - Magic number detection for file types
- **Path Validation** - Prevents directory traversal attacks

#### ✅ CSRF Protection
- **CSRF Middleware** - Available for state-changing requests
- **Double Submit Cookie** - Token-based CSRF protection
- **Configurable** - Can be enabled per route group

### Default Admin Credentials

Default users created from seed data:

- **Admin**: `admin@cacto-cms.local` / `admin123`
- **Editor**: `editor@cacto-cms.local` / `admin123`

**⚠️ IMPORTANT**: Always change passwords in production!

### Production Security Checklist

Before deploying to production, ensure:

1. ✅ **JWT Secret**: Set a strong, random `JWT_SECRET` (minimum 32 characters)
2. ✅ **HTTPS**: Set `USE_HTTPS=true` or use `https://` in `BASE_URL`
3. ✅ **Environment**: Set `ENV=production`
4. ✅ **CORS**: Configure `ALLOWED_ORIGINS` if needed (defaults to `BASE_URL`)
5. ✅ **Passwords**: Change default admin/editor passwords
6. ✅ **Database**: Use secure database path (not world-readable)
7. ✅ **Upload Directory**: Set proper permissions on upload directory
8. ✅ **Error Logging**: Monitor error logs for security issues
9. ✅ **Rate Limiting**: Verify rate limiting is working
10. ✅ **Security Headers**: Verify all security headers are present

### Environment Variables

**Required for Production:**

```bash
ENV=production
JWT_SECRET=<strong-random-secret-min-32-chars>
BASE_URL=https://yourdomain.com
USE_HTTPS=true  # Or set BASE_URL to https://
```

**Optional:**

```bash
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
PORT=8080
DB_PATH=/secure/path/to/cacto.db
UPLOAD_DIR=/secure/path/to/uploads
```

### Security Best Practices

1. **JWT Secret**: Generate using: `openssl rand -base64 32`
2. **HTTPS**: Always use HTTPS in production (required for secure cookies)
3. **Password Policy**: Implement strong password requirements
4. **Regular Updates**: Keep dependencies updated
5. **Monitoring**: Set up error logging and monitoring
6. **Backup**: Regular database backups
7. **Access Control**: Limit server access to authorized personnel only
8. **Firewall**: Configure firewall rules appropriately
9. **SSL/TLS**: Use strong SSL/TLS configuration
10. **Security Headers**: Verify all security headers are present (use securityheaders.com)

### Web Directory Security

**Public Directories:**
- ✅ `web/static/` - Public static files (CSS, JS, images)
- ✅ `web/uploads/` - Public uploaded files

**Protected:**
- ✅ Application code - Not accessible via web
- ✅ Database files - Not in web directory
- ✅ Configuration files - Not in web directory
- ✅ Source code - Not accessible

**Path Traversal Protection:**
- ✅ All file paths are validated
- ✅ `..` sequences are blocked
- ✅ Absolute paths are rejected
- ✅ Null bytes are filtered

---

## 📝 Changelog

### [1.0.0] - 2024-01-16

#### Added
- ✅ Clean Architecture implementation with Domain-Driven Design
- ✅ Page management system with dynamic routing
- ✅ Component-based architecture with database-driven components
- ✅ JWT Authentication system with Argon2id password hashing
- ✅ Role-Based Access Control (RBAC) with Admin, Editor, Author, Viewer roles
- ✅ Media management system with file validation
- ✅ SEO management with centralized meta tag handling
- ✅ Automatic sitemap generation
- ✅ Artisan CLI for migrations and seeding
- ✅ Structured error handling system
- ✅ Structured logging with Go 1.21+ slog
- ✅ Input validation system
- ✅ Environment-based configuration
- ✅ Component renderer with registry pattern
- ✅ Seed management system
- ✅ Protected routes with authentication middleware
- ✅ Error handling middleware
- ✅ CORS middleware
- ✅ Request logging middleware

#### Architecture
- Domain layer with entities and repository interfaces
- Application layer with business logic services
- Infrastructure layer with database and persistence implementations
- Interface layer with HTTP controllers and templates
- Shared utilities for common functionality

#### Security
- JWT token-based authentication
- Argon2id password hashing
- Role-based permission system
- Input validation
- CORS configuration

#### Developer Experience
- Hot reload with Air
- Type-safe templates with Templ
- Makefile with common commands
- Artisan CLI for database operations
- Comprehensive documentation

### [Unreleased]

#### Planned
- Admin dashboard UI
- Media upload handler
- File upload endpoint
- Content versioning
- Full-text search
- Caching layer
- API documentation (Swagger/OpenAPI)
- Unit and integration tests
- Docker support
- Rate limiting
- Email notifications

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards

- Use Go fmt
- Write tests
- Follow Clean Architecture principles
- Add documentation

---

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Templ](https://templ.guide) - Type-safe HTML templates
- [Chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [HTMX](https://htmx.org) - Modern web interactions
- [Go Community](https://go.dev/community) - Amazing Go ecosystem

---

<div align="center">

**Cacto CMS** - Performance-focused, modern, enterprise-ready CMS

Made with ❤️ using Go

</div>
