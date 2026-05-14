Berikut adalah **Product Requirements Document (PRD)** untuk aplikasi **PKL (Praktik Kerja Lapangan)** dengan fitur Jurnal, Absensi (timestamp + lokasi + foto ke Google Drive), dan Penilaian oleh DUDI.

---

# PRD Aplikasi PKL: Jurnal, Absensi, & Penilaian DUDI

## 1. Pendahuluan

### 1.1 Tujuan Produk
Membangun sistem manajemen PKL yang terintegrasi untuk siswa, guru pembimbing, dan DUDI (Dunia Usaha/Dunia Industri). Aplikasi ini menggantikan proses manual (kertas) dengan sistem digital yang mencatat lokasi, waktu, dan bukti kehadiran, serta memudahkan penilaian akhir PKL.

### 1.2 Target Pengguna
| Role | Deskripsi |
|------|------------|
| **Siswa** | Melakukan absensi, mengisi jurnal harian, melihat nilai. |
| **Guru Pembimbing** | Memantau absensi & jurnal siswa, membimbing. |
| **DUDI (Pembimbing Lapangan)** | Memberikan nilai PKL berdasarkan performa siswa di tempat magang. |
| **Admin** | Mengelola data siswa, DUDI, periode PKL, dan konfigurasi sistem. |

---

## 2. Arsitektur Teknologi

| Komponen | Teknologi |
|----------|------------|
| Backend | Golang (dengan Gin/Echo framework, GORM) |
| Frontend | Vue.js + Tailwind CSS |
| Database | PostgreSQL |
| Penyimpanan Foto | Google Drive API |
| Autentikasi | JWT (Access + Refresh Token) |
| Real-time (opsional) | WebSocket (untuk notifikasi) |

---

## 3. Core Features

### 3.1 Manajemen Pengguna & Role
- Registrasi/login dengan NIS/NIP/NIK.
- Verifikasi email/WhatsApp OTP.
- Setiap role memiliki dashboard berbeda.

### 3.2 Absensi (dengan Timestamp, Lokasi, Foto ke Google Drive)
**Alur Absensi Siswa:**
1. Siswa membuka halaman absensi.
2. Sistem mengambil **timestamp** otomatis (server-side).
3. Siswa mengambil **foto selfie dengan latar tempat PKL** (via kamera browser/mobile).
4. Sistem mengambil **lokasi GPS** (longitude, latitude) melalui browser.
5. Foto diupload ke **Google Drive** (folder per siswa/bulan).
6. Sistem menyimpan:  
   - `id_absensi`, `user_id`, `timestamp`, `latitude`, `longitude`, `google_drive_image_url`, `status` (hadir, terlambat, izin, sakit).
7. Validasi: Apakah lokasi siswa berada dalam radius yang ditentukan oleh DUDI? (jika ya → valid, jika tidak → warning).

### 3.3 Jurnal Harian PKL
- Siswa mengisi:  
  - Tanggal kegiatan (default hari ini).  
  - Uraian kegiatan.  
  - Dokumentasi (foto opsional → upload ke Google Drive juga).  
  - Komentar/refleksi.  
- Guru/DUDI bisa memberikan komentar balik pada jurnal.

### 3.4 Penilaian PKL oleh DUDI
- DUDI melihat daftar siswa yang magang di perusahaannya.
- DUDI mengisi form penilaian dengan kriteria:
  - Kehadiran (otomatis dari absensi).  
  - Kedisiplinan, tanggung jawab, kerjasama, inisiatif (skala 1-5).  
  - Nilai akhir PKL (angka + huruf).  
- Nilai otomatis terlihat oleh siswa dan guru pembimbing.

### 3.5 Monitoring & Laporan (Guru/Admin)
- Rekap absensi per siswa (export PDF/Excel).
- Rekap jurnal (waktu baca, komentar).
- Rekap nilai dari DUDI.
- Peta lokasi absensi (heatmap siswa selama PKL).

### 3.6 Notifikasi
- Pengingat absensi (setiap pagi via email/Web push).
- DUDI mengisi nilai → notif ke guru.
- Komentar pada jurnal → notif ke siswa.

---

## 4. App Flow

### 4.1 Login & Otorisasi
```
Landing Page → Login (pilih role: Siswa/Guru/DUDI/Admin) → Dashboard sesuai role
```

### 4.2 Flow Absensi (Siswa)
```
Dashboard Siswa → Menu Absensi → Sistem ambil timestamp & lokasi → 
Ambil foto (camera) → Upload ke Google Drive → 
Simpan data absensi → Sukses → History absensi
```

### 4.3 Flow Jurnal (Siswa)
```
Dashboard → Tulis Jurnal Baru → Isi form + upload dokumen → 
Simpan → Daftar jurnal (editable sebelum dikomentari) → 
Lihat komentar dari DUDI/Guru
```

### 4.4 Flow Penilaian (DUDI)
```
Login as DUDI → Daftar siswa magang → Pilih siswa → 
Form penilaian (otomatis tarik data kehadiran) → 
Submit nilai → Final (terkunci untuk diedit) → 
Guru & siswa bisa lihat
```

### 4.5 Flow Monitoring (Guru)
```
Dashboard Guru → Pilih periode PKL → Pilih siswa → 
Lihat ringkasan (total hadir, jurnal terisi, nilai DUDI) → 
Klik detail absensi/jurnal → Beri komentar
```

---

## 5. Detail Database Schema (PostgreSQL)

### users
```sql
id UUID PK
full_name VARCHAR(255)
email VARCHAR(255) UNIQUE
password_hash VARCHAR(255)
role ENUM('student', 'teacher', 'dudi', 'admin')
nis_nip_nik VARCHAR(50) UNIQUE
dudi_id UUID FK (if role='student')
created_at TIMESTAMP
```

### dudi
```sql
id UUID PK
company_name VARCHAR(255)
address TEXT
latitude DECIMAL(10,8)
longitude DECIMAL(11,8)
radius_allowed INTEGER (meter)
pic_name VARCHAR(255)
phone VARCHAR(20)
```

### absensi
```sql
id UUID PK
student_id UUID FK
timestamp TIMESTAMP
latitude DECIMAL(10,8)
longitude DECIMAL(11,8)
photo_url TEXT (Google Drive link)
status ENUM('hadir', 'terlambat', 'izin', 'sakit')
is_verified BOOLEAN (radius check)
created_at TIMESTAMP
```

### jurnal
```sql
id UUID PK
student_id UUID FK
date DATE
activity TEXT
documentation_url TEXT (Google Drive, optional)
reflection TEXT
teacher_comment TEXT
dudi_comment TEXT
created_at TIMESTAMP
updated_at TIMESTAMP
```

### penilaian
```sql
id UUID PK
student_id UUID FK
dudi_id UUID FK
attendance_score_auto DECIMAL(5,2) (dari absensi)
discipline INTEGER (1-5)
responsibility INTEGER (1-5)
teamwork INTEGER (1-5)
initiative INTEGER (1-5)
final_score DECIMAL(5,2)
final_grade CHAR(2) (A/B/C/D)
notes TEXT
submitted_at TIMESTAMP
```

---

## 6. Non-Functional Requirements

| Aspek | Deskripsi |
|-------|------------|
| **Keamanan** | Semua endpoint JWT-protected, upload foto hanya via token valid. |
| **Privasi Lokasi** | Lokasi hanya direkam saat absensi, tidak ada tracking terus-menerus. |
| **Ketersediaan** | Target uptime 99.5% karena PKL harian. |
| **Responsif** | Tailwind CSS + Vue.js memastikan mobile-friendly (siswa gunakan HP). |
| **Backup Data** | PostgreSQL backup harian, Google Drive untuk penyimpanan media. |

---

## 7. User Interface Guidelines (Tailwind + Vue)

- **Warna tema**: Biru (#1E3A8A) untuk institusi, aksen hijau (#10B981) untuk aksi sukses.
- **Komponen**:
  - Form absensi dengan tombol kamera (menggunakan `navigator.mediaDevices`).
  - Map (Leaflet/Google Maps) untuk menampilkan lokasi absensi siswa.
  - Table/tailwind bisa dengan daisyUI atau headlessui untuk modal/filter.
- **Layout**: Sidebar untuk navigasi (menu: Dashboard, Absensi, Jurnal, Nilai). Mobile → hamburger menu.

---

## 8. API Endpoints (Golang) – Contoh

| Method | Endpoint | Deskripsi |
|--------|----------|-------------|
| POST | /api/login | Login JWT |
| POST | /api/absensi | Siswa melakukan absensi + upload foto ke Google Drive |
| GET | /api/absensi/history | Riwayat absensi siswa |
| POST | /api/jurnal | Tambah jurnal |
| GET | /api/jurnal/{student_id} | Lihat jurnal (guru/DUDI) |
| POST | /api/jurnal/comment | Beri komentar |
| POST | /api/nilai | DUDI input nilai |
| GET | /api/nilai/{student_id} | Lihat nilai (siswa/guru) |
| GET | /api/report/absensi?periode=... | Export rekap (guru) |

---

## 9. Integrasi Google Drive (untuk foto absensi & jurnal)

- Gunakan **Service Account** Google Drive API.
- Setiap folder dibuat per user: `pkl_photos/{user_id}/{YYYY-MM}/`
- Upload file dengan nama `absensi_{timestamp}.jpg`
- Kembalikan `webViewLink` atau `thumbnailLink` untuk ditampilkan di frontend.
- Akses publik terbatas (hanya user yang login bisa lihat fotonya).

---

## 10. Prioritas Pengembangan (MVP)

1. **Minggu 1**: Setup Golang + PostgreSQL, auth JWT, CRUD user/dudi.
2. **Minggu 2**: Absensi dengan location + timestamp + upload foto ke Google Drive.
3. **Minggu 3**: Jurnal sederhana (tanpa komentar dulu).
4. **Minggu 4**: Form penilaian oleh DUDI.
5. **Minggu 5**: Dashboard guru (monitoring absensi & jurnal).
6. **Minggu 6**: Notifikasi peran, export laporan, testing + deployment.

---

## 11. Risiko & Mitigasi

| Risiko | Mitigasi |
|--------|-----------|
| Siswa memalsukan lokasi GPS | Gunakan kombinasi foto + timestamp + validasi radius dari server. |
| Google Drive quota penuh | Setiap user dapat upload maks 10MB per foto, kompres sebelum upload. |
| Siswa lupa absensi | Notifikasi push setiap jam 08:00, 12:00. |
| DUDI tidak aktif memberikan nilai | Guru dapat mengingatkan via sistem notifikasi, admin dapat override dengan nilai sementara. |

---

Dokumen ini siap digunakan sebagai acuan tim developer (backend Golang, frontend Vue.js + Tailwind, database PostgreSQL, dan integrasi Google Drive API).
