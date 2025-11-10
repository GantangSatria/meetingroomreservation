---

# Meeting Room Reservation

Sistem reservasi ruang meeting dengan fitur approval oleh admin, QR code untuk check-in, check-out, dan manajemen pengguna/ruangan.

## Kegunaan

* Pengguna biasa dapat membuat reservasi ruang meeting dengan waktu mulai & selesai.
* Admin dapat menyetujui atau menolak reservasi.
* Setelah disetujui, muncul QR code yang bisa digunakan saat datang untuk **check-in**.
* Pengguna dapat melakukan **check-out** setelah selesai menggunakan ruang.
* Admin/manager juga dapat melihat semua reservasi, mengelola ruangan (CRUD), dan mengelola pengguna.

## Instalasi & penggunaan

1. Pastikan memiliki Go (versi 1.x) + database PostgreSQL.
2. Clone repositori:

   ```bash
   git clone https://github.com/GantangSatria/meetingroomreservation.git
   cd meetingroomreservation
   ```
3. Buat file environment (.env) dengan variabel seperti:

   ```env
   DB_URL=postgres://user:password@host:port/dbname
   JWT_SECRET=your_secret_key
   ```
4. Jalankan migrasi & build aplikasi:

   ```bash
   go run cmd/main.go 
   ```


## Struktur Path & Endpoints

### Folder proyek

* `internal/controller/` — handler HTTP untuk tiap route
* `internal/services/` — logika bisnis (reservasi, check-in/out)
* `internal/repository/` — akses database via GORM
* `internal/models/` — definisi model Go + mapping ke tabel DB
* `pkg/dto/` — objek untuk transfer data (response API)
* `routes/` (atau file router) — definisi semua route HTTP

### Route utama (/api/v1)

| Method                                         | Path                                     | Keterangan |
| ---------------------------------------------- | ---------------------------------------- | ---------- |
| `POST /api/v1/register`                        | Register user baru                       |            |
| `POST /api/v1/login`                           | Login user & mendapatkan JWT             |            |
| `GET /api/v1/rooms`                            | Daftar semua ruangan publik              |            |
| `GET /api/v1/rooms/:id`                        | Detail ruangan                           |            |
| `POST /api/v1/rooms`                           | (Admin) Tambah ruangan                   |            |
| `PUT /api/v1/rooms/:id`                        | (Admin) Update ruangan                   |            |
| `DELETE /api/v1/rooms/:id`                     | (Admin) Hapus ruangan                    |            |
| `POST /api/v1/reservations`                    | Buat reservasi (User)                    |            |
| `GET /api/v1/reservations`                     | Daftar semua reservasi (tergantung role) |            |
| `GET /api/v1/reservations/:id`                 | Detail reservasi                         |            |
| `PUT /api/v1/reservations/:id`                 | Update reservasi (User/Admin)            |            |
| `DELETE /api/v1/reservations/:id`              | Hapus reservasi (User/Admin)             |            |
| `GET /api/v1/reservations/:id/qrcode`          | Ambil QR code untuk reservasi            |            |
| `PUT /api/v1/admin/reservations/:id/approve`   | (Admin) Approve reservasi                |            |
| `PUT /api/v1/admin/reservations/:id/reject`    | (Admin) Reject reservasi                 |            |
| `POST /api/v1/checkin/:reservation_id`         | (User) Check-in tanpa QR                 |            |
| `POST /api/v1/checkin/qr`                      | (User) Check-in via QR code              |            |
| `PUT /api/v1/checkin/:reservation_id/checkout` | (User) Check-out dari check-in           |            |

> **Catatan:** Semua route di `/api/v1/...` kecuali register/login dilindungi JWT.
> Untuk prefix `/admin`, perlu role = `admin`.

## Workflow singkat

1. User mendaftar → login → dapat token.
2. User membuat reservasi → status = `pending`.
3. Admin menyetujui → status = `approved`, QR code di-generate.
4. User datang ke ruangan → scan QR. Backend cek: valid, belum pernah check-in → buat record check-in.
5. Setelah selesai, user/pegawai melakukan check-out → status `checked_out` atau sejenis.

## Testing & Curl

Contoh curl check-in via QR code:

```bash
curl -X POST http://localhost:8080/api/v1/checkin/qr \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"qr_data":"reservation:1:1:2025-01-23T16:00:00+07:00"}'
```
