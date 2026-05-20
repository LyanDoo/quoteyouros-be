# Authentication Implementation Analysis

## 🎯 Overview

We need to implement a complete JWT-based authentication system with:
- **Register** - Create new admin user account
- **Login** - Authenticate and return JWT token
- **Validate** - Protected routes verify JWT tokens

---

## 📋 Detailed Analysis

### 1. Data Model & Database

**Current User Table:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**Status:** ✅ Already exists, no changes needed

**DTOs to Create:**
```go
type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token     string `json:"token"`
    ExpiresIn int    `json:"expires_in"`
    User      struct {
        ID    string `json:"id"`
        Email string `json:"email"`
    } `json:"user"`
}
```

---

### 2. Security Architecture

#### Password Hashing
```go
// Use golang.org/x/crypto/bcrypt
// Hash password before storing: bcrypt.GenerateFromPassword([]byte(password), 10)
// Verify on login: bcrypt.CompareHashAndPassword(storedHash, []byte(inputPassword))
```

**Why bcrypt:**
- Adaptive - slows down as computers get faster
- Salted automatically
- Industry standard
- Already in go.mod via crypto package

#### JWT Token Structure
```
Header: {
  "alg": "HS256",
  "typ": "JWT"
}

Payload: {
  "user_id": "uuid",
  "email": "user@example.com",
  "iat": 1234567890,
  "exp": 1234567890 + 24h
}

Signature: HMAC-SHA256(secret)
```

**Configuration:**
- Secret: `JWT_SECRET` from `.env`
- Expiration: `JWT_EXPIRATION` from `.env` (default 24h)
- Algorithm: HS256

---

### 3. Implementation Components

### A. Password Utilities (New File)
**Location:** `pkg/crypto/password.go`

```go
package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plaintext password
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), err
}

// VerifyPassword compares a hashed password with plaintext
func VerifyPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### B. JWT Utilities (New File)
**Location:** `pkg/jwt/token.go`

```go
package jwt

import (
    "time"
    jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwtlib.RegisteredClaims
}

// GenerateToken creates a JWT token
func GenerateToken(userID, email, secret string, expiration time.Duration) (string, error) {
    expirationTime := time.Now().Add(expiration)
    claims := &Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwtlib.RegisteredClaims{
            ExpiresAt: jwtlib.NewNumericDate(expirationTime),
            IssuedAt:  jwtlib.NewNumericDate(time.Now()),
        },
    }
    
    token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates a token
func ValidateToken(tokenString, secret string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err != nil || !token.Valid {
        return nil, err
    }
    
    return claims, nil
}
```

### C. Auth Repository (Already Exists)
**Location:** `internal/repository/user/postgres.go`

Status: ✅ Exists, methods already available:
- `CreateUser(ctx context.Context, user *domain.User) error`
- `GetUserByEmail(ctx context.Context, email string) (*domain.User, error)`
- `GetUserByID(ctx context.Context, id string) (*domain.User, error)`

### D. Auth Use Case (To Create)
**Location:** `internal/usecase/auth/usecase.go`

```go
type AuthUseCase struct {
    userRepo domain.UserRepository
    jwtSecret string
    jwtExpiration time.Duration
}

// Register creates a new user
func (u *AuthUseCase) Register(ctx context.Context, req *domain.RegisterRequest) error {
    // 1. Validate request (email format, password strength, match)
    // 2. Check if email already exists
    // 3. Hash password
    // 4. Create user
    // 5. Return success
}

// Login authenticates user and returns JWT
func (u *AuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (string, error) {
    // 1. Find user by email
    // 2. Verify password
    // 3. Generate JWT token
    // 4. Return token
}

// GetCurrentUser retrieves user from JWT claims
func (u *AuthUseCase) GetCurrentUser(ctx context.Context, userID string) (*domain.User, error) {
    // 1. Get user by ID
    // 2. Return user (without password)
}
```

### E. Auth Handlers (To Create)
**Location:** `internal/handler/auth.go`

```go
type AuthHandler struct {
    usecase domain.AuthUseCase
}

// Register POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // 1. Parse request body
    // 2. Call usecase.Register()
    // 3. Return success/error response
}

// Login POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // 1. Parse request body
    // 2. Call usecase.Login()
    // 3. Return token response
}

// GetCurrentUser GET /api/auth/me (protected)
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
    // 1. Get userID from c.Locals("user_id")
    // 2. Call usecase.GetCurrentUser()
    // 3. Return user data
}
```

---

### 4. Database Schema Updates

**Add to migrations/001_init_schema.sql:**

Potentially add these optional fields:
```sql
ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT TRUE;
```

**Status:** Optional - Current schema works, these are enhancements

---

### 5. Request/Response Flow

#### Register Flow
```
POST /api/auth/register
{
  "email": "admin@example.com",
  "password": "SecurePass123",
  "confirm_password": "SecurePass123"
}

↓
1. Validate email format
2. Check password strength (8+ chars)
3. Verify passwords match
4. Check email not already registered
5. Hash password with bcrypt
6. Save user to database
7. Return 201 Created

{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "admin@example.com"
  }
}
```

#### Login Flow
```
POST /api/auth/login
{
  "email": "admin@example.com",
  "password": "SecurePass123"
}

↓
1. Find user by email
2. Verify password (bcrypt compare)
3. If valid, generate JWT token (24h expiration)
4. Return token

{
  "success": true,
  "data": {
    "token": "eyJhbGc...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "email": "admin@example.com"
    }
  }
}
```

#### Protected Route Flow (GET /api/auth/me)
```
GET /api/auth/me
Header: Authorization: Bearer eyJhbGc...

↓
1. JWT middleware validates token
2. Extract user_id from claims
3. Get user from database
4. Return user data

{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "admin@example.com",
    "created_at": "2026-05-20T10:00:00Z"
  }
}
```

---

### 6. Error Scenarios to Handle

| Scenario | Status | Response |
|----------|--------|----------|
| Email already exists | 409 Conflict | `"Email already registered"` |
| Invalid email format | 400 Bad Request | `"Invalid email format"` |
| Password too weak | 400 Bad Request | `"Password must be 8+ chars"` |
| Passwords don't match | 400 Bad Request | `"Passwords do not match"` |
| User not found | 401 Unauthorized | `"Invalid email or password"` |
| Wrong password | 401 Unauthorized | `"Invalid email or password"` |
| Invalid JWT | 401 Unauthorized | `"Invalid or expired token"` |
| Token expired | 401 Unauthorized | `"Token expired"` |

---

### 7. Validation Requirements

**Email:**
- Valid email format (RFC 5322)
- Unique in database

**Password:**
- Minimum 8 characters
- Optional: Should include uppercase, lowercase, number, special char
- Confirm password matches

---

### 8. Implementation Order

1. ✅ Create password hashing utility (`pkg/crypto/password.go`)
2. ✅ Create JWT utility (`pkg/jwt/token.go`)
3. ✅ Update domain entities with RegisterRequest
4. ✅ Create auth use case (`internal/usecase/auth/usecase.go`)
5. ✅ Create auth handler (`internal/handler/auth.go`)
6. ✅ Update main.go routes to use handlers
7. ✅ Update middleware to use new JWT utilities
8. ✅ Add validation middleware
9. ✅ Test endpoints with Postman

---

### 9. Files to Create/Modify

**New Files:**
- [ ] `pkg/crypto/password.go` - Password hashing utilities
- [ ] `pkg/jwt/token.go` - JWT token utilities
- [ ] `internal/usecase/auth/usecase.go` - Auth business logic
- [ ] `internal/handler/auth.go` - Auth HTTP handlers

**Modified Files:**
- [ ] `internal/domain/entities.go` - Add RegisterRequest
- [ ] `internal/domain/usecases.go` - Already has AuthUseCase interface
- [ ] `internal/middleware/auth.go` - Update to use JWT utils
- [ ] `cmd/main.go` - Update route handlers

---

### 10. Configuration Needed

From `.env`:
- `JWT_SECRET` - Used to sign/verify tokens
- `JWT_EXPIRATION` - Token lifetime (default 24h)

Validation:
- Check JWT_SECRET is at least 32 characters for production
- Warn if default secret is used in production

---

## 🚀 Next Steps

Ready to start implementation? I can:

1. **Create the utility packages** (password & JWT)
2. **Build the auth use case** with all business logic
3. **Create the auth handlers** for register/login endpoints
4. **Update routes** in main.go
5. **Test with Postman** collection

Which would you like me to start with?

