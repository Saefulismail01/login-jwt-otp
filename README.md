# Login JWT OTP

Aplikasi autentikasi dengan JWT dan OTP menggunakan Golang dan PostgreSQL.

## Fitur Utama

- Sistem autentikasi dengan JWT
- Registrasi pengguna dengan verifikasi OTP
- Login pengguna
- Enkripsi password menggunakan bcrypt
- Validasi email dan OTP

## Struktur Proyek

```
login-jwt-otp/
├── config/           # Konfigurasi aplikasi
├── delivery/         # Layer delivery (HTTP)
│   ├── controller/   # Controller HTTP
│   └── server/       # Server HTTP
├── model/           # Model data
├── repository/      # Repository database
├── usecase/         # Usecase bisnis
├── utils/           # Utilitas
└── main.go          # Entry point aplikasi
```

## Persyaratan

- Go 1.23.0 atau lebih tinggi
- PostgreSQL
- Git

## Instalasi

1. Clone repository
```bash
git clone <repository-url>
```

2. Masuk ke direktori proyek
```bash
cd login-jwt-otp
```

3. Install dependencies
```bash
go mod download
```

4. Buat file `.env` dengan isi:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_password
DB_NAME=login_jwt_otp
DB_DRIVER=postgres
API_PORT=8080

JWT_SECRET=your_jwt_secret
JWT_EXPIRY=2h
APP_NAME=Login JWT OTP
```

5. Jalankan migrasi database
```bash
go run migrations/migrate.go
```

## Penggunaan

### Menjalankan Server
```bash
go run main.go
```

### API Endpoints

#### Registrasi
- POST `/api/auth/register`
  ```json
  {
    "name": "string",
    "email": "string",
    "password": "string",
    "birth_year": integer,
    "phone": "string"
  }
  ```
  
  Response:
  ```json
  {
    "message": "OTP sent successfully",
    "email": "string"
  }
  ```

#### Verifikasi OTP
- POST `/api/auth/verify-otp`
  ```json
  {
    "email": "string",
    "otp": "string",
    "name": "string",
    "password": "string",
    "birth_year": integer,
    "phone": "string"
  }
  ```
  
  Response:
  ```json
  {
    "message": "Registration successful",
    "user": {
      "id": integer,
      "name": "string",
      "email": "string",
      "role": "string"
    }
  }
  ```

#### Login
- POST `/api/auth/login`
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
  
  Response:
  ```json
  {
    "token": "string",
    "user": {
      "id": integer,
      "name": "string",
      "email": "string",
      "role": "string"
    }
  }
  ```

## Teknologi

- Backend: Golang
- Framework: Gin
- Database: PostgreSQL
- Authentication: JWT
- OTP: Custom implementation
- Password Hashing: bcrypt

