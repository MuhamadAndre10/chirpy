# App like twiiter | Chirpy

> Aplikasi ini di bangun dengan mengikuti petunjuk pembelajaran di [boot.dev](https://www.boot.dev/lessons/861ada77-c583-42c8-a265-657f2c453103). Di bangun  hanya untuk pembelajaran. 

## Apa yang Akan Dipelajari? 📖
- Dasar-dasar Server Web: Memahami cara kerja server web dan aplikasinya.
- Membangun Server Go: Membuat server HTTP siap produksi di Go tanpa framework.
- API RESTful: Berkomunikasi dengan klien menggunakan JSON, header, dan kode status.
- Manajemen Data: Menyimpan dan mengambil data dari Postgres menggunakan SQL yang aman.
- Sistem Keamanan: Mengimplementasikan sistem otentikasi/otorisasi yang kuat.
- Integrasi & Otomasi: Membangun dan memahami webhook serta kunci API.
- Dokumentasi: Mendokumentasikan REST API dengan markdown.

## Mau mencoba? / You wanna try?
### API Dokumentasi
| Method | Route | Keterangan | Autentikasi |
| :--- | :--- | :--- | :---: |
| **POST** | `/chirps` | Membuat chirp baru | **Ya** |
| **DELETE** | `/chirps/{id}` | Menghapus chirp berdasarkan ID | **Ya** |
| **GET** | `/chirps` | Mengambil semua chirp | Tidak |
| **GET** | `/chirps/{id}` | Mengambil chirp berdasarkan ID | Tidak |
| **POST** | `/users` | Membuat pengguna baru | Tidak |
| **PUT** | `/users` | Memperbarui kata sandi pengguna | **Ya** |
| **POST** | `/login` | Otentikasi dan login pengguna | Tidak |
| **POST** | `/refresh` | Mendapatkan token otentikasi baru | Tidak |
| **POST** | `/revoke` | Mencabut token refresh | Tidak |
| **POST** | `/polka/webhooks` | Memperbarui status pengguna menjadi anggota Chirpy Red | Tidak |

> _Untuk router yang menggunakan `authentikasi` gunakan `JWT Bearer Token` Khusus untuk webhooks gunakan `ApiKey Token` yang ada di fil .env_


### 1. Persiapan awal
-  Pastikan instalasi golang pada komputer kalian. _lihat petunjuk [disini](https://go.dev/doc/install) jika belum_
-  Clone Repository Chirpy
```bash
git clone https://github.com/MuhamadAndre10/chirpy.git
cd chirpy
```
-  Copy example.env ke .env
```bash
    cp example.env .env
```
### 2. Konfigurasi proyek
- **Inisialisasi Go Module**: Inisialisasi Go module dan unduh semua dependency yang dibutuhkan.
- **Setup Database:** Proyek ini menggunakan PostgreSQL. Pastikan Anda sudah menginstalnya. Buat database baru dengan nama `chirpy`
```sql
    CREATE DATABASE chirpy;
```
### 3. Migrasi Database
- **Instalasi Tools:** Instal `Goose` dan `SQLC`, yang digunakan untuk migrasi database dan menghasilkan kode SQL yang aman.
```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
- **Jalankan Migrasi :** Lakukan migrasi database untuk membuat tabel-tabel yang diperlukan.
```bash
goose up
```
### 4. Menjalankan aplikasi 
- **Jalankan Aplikasi:**
```bash
    go run .
```
**Note** Jika program tidak berjalan call [me](https://www.instagram.com/m_andrepriyanto/)
