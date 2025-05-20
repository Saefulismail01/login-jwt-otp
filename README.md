# Login JWT OTP

Aplikasi autentikasi dengan JWT dan OTP menggunakan Golang dan PostgreSQL. Aplikasi ini menyediakan sistem registrasi dengan verifikasi OTP melalui email dan login menggunakan JWT untuk autentikasi.

## Fitur Utama

- ✅ Registrasi pengguna dengan verifikasi OTP
- ✅ Verifikasi OTP untuk aktivasi akun
- ✅ Login menggunakan JWT
- ✅ Validasi email dan password
- ✅ Enkripsi password menggunakan bcrypt
- ✅ Manajemen sesi dengan JWT
- ✅ Rate limiting untuk mencegah spam OTP

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

### 1. Registrasi

Mengirimkan data registrasi dan menerima OTP melalui email.

- **URL**: `POST /api/auth/register`
- **Request Body**:
  ```json
  {
    "name": "Deny Caknan",
    "email": "deny@example.com",
    "password": "securePassword123",
    "birth_year": 2002,
    "phone": "081234567890"
  }
  ```
- **Response Success (200 OK)**:
  ```json
  {
    "message": "OTP sent successfully. Please check your email.",
    "email": "deny@example.com"
  }
  ```
- **Error Response**:
  - 400: Invalid request data
  - 409: Email already registered
  - 500: Internal server error

### 2. Verifikasi OTP

Memverifikasi OTP dan menyelesaikan proses registrasi.

- **URL**: `POST /api/auth/verify-otp`
- **Request Body**:
  ```json
  {
    "email": "deny@example.com",
    "otp": "123456",
    "name": "Deny Caknan",
    "password": "securePassword123",
    "birth_year": 2002,
    "phone": "081234567890"
  }
  ```
- **Response Success (200 OK)**:
  ```json
  {
    "message": "Registration successful",
    "data": {
      "user": {
        "id": 1,
        "name": "Deny Caknan",
        "email": "deny@example.com",
        "birth_year": 2002,
        "phone": "081234567890",
        "role": "USER"
      },
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```
- **Error Response**:
  - 400: Invalid OTP or request data
  - 401: OTP expired or invalid
  - 500: Internal server error

### 3. Login

Login pengguna dengan email dan password.

- **URL**: `POST /api/auth/login`
- **Request Body**:
  ```json
  {
    "email": "deny@example.com",
    "password": "securePassword123"
  }
  ```
- **Response Success (200 OK)**:
  ```json
  {
    "message": "Login successful",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": 1,
        "name": "Deny Caknan",
        "email": "deny@example.com",
        "role": "USER"
      }
    }
  }
  ```
- **Error Response**:
  - 400: Invalid request data
  - 401: Invalid credentials
  - 500: Internal server error

## Alur Registrasi

1. Pengguna mengirim data registrasi ke `/api/auth/register`
2. Sistem mengirim OTP ke email pengguna
3. Pengguna memverifikasi OTP dengan mengirim data lengkap ke `/api/auth/verify-otp`
4. Jika OTP valid, akun pengguna dibuat dan token JWT dikembalikan
5. Pengguna dapat login menggunakan email dan password yang telah didaftarkan

## Teknologi

- **Bahasa Pemrograman**: Go 1.23.0+
- **Framework Web**: Gin
- **Database**: PostgreSQL
- **Autentikasi**: JWT (JSON Web Tokens)
- **Keamanan**:
  - Bcrypt untuk hashing password
  - OTP dengan masa berlaku 10 menit
  - Rate limiting untuk mencegah spam
- **Lainnya**:
  - Viper untuk manajemen konfigurasi
  - GORM untuk ORM
  - Validator untuk validasi input

## Pengembangan

### Menjalankan Aplikasi

1. Pastikan PostgreSQL berjalan
2. Buat file `.env` dari `.env.example`
3. Jalankan migrasi database:
   ```bash
   go run migrations/migrate.go
   ```
4. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

### Testing

```bash
go test -v ./...
```

### Environment Variables

Salin `.env.example` ke `.env` dan sesuaikan konfigurasi:

```
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_password
DB_NAME=login_jwt_otp
DB_DRIVER=postgres

# Server
API_PORT=8080

# JWT
JWT_SECRET=your_secure_secret
JWT_EXPIRY=24h

# Aplikasi
APP_NAME=Login JWT OTP
```

## Kontribusi

1. Fork repository
2. Buat branch fitur (`git checkout -b fitur/namafitur`)
3. Commit perubahan (`git commit -m 'Menambahkan fitur baru'`)
4. Push ke branch (`git push origin fitur/namafitur`)
5. Buat Pull Request

