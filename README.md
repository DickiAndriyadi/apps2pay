# Cinema Ticket System – Backend Test APPS2PAY 2025

**Nama**:  Dicki Andriyadi

**Tanggal**: 08 Oktober 2025

**Email**: dicki.andriyadi14@gmail.com

**Nomor Telepon**: +62851158011846

**Repository**: https://github.com/DickiAndriyadi/apps2pay

## 📌 Deskripsi
Sistem pembelian tiket bioskop nasional yang memungkinkan:
- Customer memilih kursi tanpa takut double booking
- Restok otomatis jika tidak dikonfirmasi dalam 15 menit
- Pembatalan jadwal oleh bioskop dengan refund massal

Dibangun dengan **Go (Fiber)**, **PostgreSQL**, dan **Redis**.

## ▶️ Cara Menjalankan

### Prasyarat
- Go 1.21+
- PostgreSQL
- Redis
- Git

### Langkah 1: Clone Repository
```bash
git clone https://github.com/DickiAndriyadi/apps2pay.git
```

```bash
cd apps2pay
```

### Langkah 2: Setup Database 

Buat database baru dengan nama `cinema_db`

```bash
createdb cinema_db
```

```bash
psql -d cinema_db -f db/schema.sql
```

### Langkah 3: Jalankan Redis

- Pastikan Redis telah di install
- Lalu jalankan

```bash
redis-server
```

### Langkah 4: Jalankan Api

- Sesuaikan .env yang ada dengan config database dan redis yang ada di local
- Lalu jalankan

```bash
go mod tidy
```

```bash
go run main.go
```

- Server sudah berjalan di `http://localhost:3000`
- Atau anda bisa jalankan via mode debugging yang sudah saya siapkan di launch.json

### Langkah 5: Test API

- Import Collection Postman
- Login → salin token → uji endpoint


## 🌐 Daftar Endpoint API

| Endpoint | Method | Deskripsi |
|:----------|:--------:|:-----------|
||||
||||
| `/login` | POST | Login |
| `/api/schedules` | GET, POST, PUT, DELETE | CRUD Schedule |
| `/api/schedules/:id` | GET | Get Schedule by ID |
| `/api/schedules/:id/seats` | GET | Get Seats for Schedule |
| `/api/schedules/:id/seats/lock` | POST | Lock Seat |
| `/api/schedules/:id/seats/release` | POST | Release Seat |
| `/api/schedules/:id/seats/confirm` | POST | Confirm Seat Purchase |
| `/api/schedules/:id/cancel` | POST | Cancel Schedule |

## 🧪 Akun Demo

| Email | Password |
|:-------|:----------|
||||
||||
| admin@example.com | passwordadmin |
| user@example.com | passworduser |

## 🧪 Note

- Untuk system design ada di folder `system-design`
- Untuk diagram database ada di folder `db`
