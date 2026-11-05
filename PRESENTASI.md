---
marp: true
theme: default
paginate: true
size: 16:9
header: '**PKL Tracker** — Sistem Manajemen Praktik Kerja Lapangan'
footer: 'SMKN 1 Arahan • Indramayu'
style: |
  section {
    font-family: 'Inter', 'Segoe UI', sans-serif;
    font-size: 24px;
  }
  section h1 { color: #1E3A8A; }
  section h2 { color: #1E3A8A; }
  section strong { color: #10B981; }
  section.lead h1 { font-size: 56px; color: #1E3A8A; }
  table { font-size: 20px; }
  code { background: #f3f4f6; }
---

<!-- _class: lead -->
<!-- _paginate: false -->
<!-- _header: '' -->
<!-- _footer: '' -->

# PKL Tracker

### Sistem Manajemen Praktik Kerja Lapangan Terintegrasi

Absensi GPS • Jurnal Harian • Penilaian DUDI • Monitoring & Laporan

SMKN 1 Arahan — Indramayu

---

# Latar Belakang

Manajemen PKL secara **manual (kertas)** punya banyak masalah:

- Absensi **sulit dibuktikan** — rawan titip absen & pemalsuan lokasi.
- Jurnal siswa **tidak terdokumentasi** rapi, sulit dipantau guru.
- Nilai dari **DUDI (instruktur lapangan)** sering telat / tidak terpusat.
- Rekap laporan **memakan waktu** karena harus dihitung manual.

**Solusi:** sistem digital terintegrasi yang mencatat **lokasi, waktu, dan bukti kehadiran**, serta menyatukan alur jurnal & penilaian.

---

# Tujuan Produk

- Menggantikan proses manual dengan **sistem digital** yang terekam & bisa diaudit.
- Menghubungkan **3 pihak**: Siswa, Guru Pembimbing, dan DUDI dalam satu platform.
- Menyediakan **bukti kehadiran yang kuat** (server-side timestamp + GPS + foto).
- Mempercepat **penilaian akhir PKL** dan **rekap laporan** berbasis periode.

---

# Target Pengguna (5 Role)

| Role | Peran |
|------|-------|
| **Siswa** | Absen masuk/pulang, mengisi jurnal harian, melihat riwayat & nilai |
| **Guru Pembimbing** | Memantau absensi & jurnal siswa bimbingannya, verifikasi, komentar |
| **Instruktur DUDI** | Menilai performa siswa magang di perusahaannya, review jurnal |
| **Admin** | Kelola pengguna, DUDI, jurusan, periode, backup, laporan |
| **Admin Jurusan** | Sama seperti admin, tapi **hanya untuk jurusannya sendiri** |

Setiap role melihat **dashboard dan menu yang berbeda**.

---

# Arsitektur Teknologi

| Komponen | Teknologi |
|----------|-----------|
| **Backend** | Go 1.25 • Gin • GORM |
| **Frontend** | Vue 3 (Composition API) • Vite • Pinia • Vue Router |
| **UI / Styling** | Tailwind CSS v4 • Headless UI • Lucide Icons |
| **PWA** | `vite-plugin-pwa` — installable, mobile-first |
| **Peta** | Leaflet + OpenStreetMap |
| **Database** | PostgreSQL |
| **Autentikasi** | JWT (Access 24 jam + Refresh 7 hari) |
| **Deployment** | Single binary Go (frontend ter-embed) + Nginx |

---

# Fitur Utama — Siswa & Guru

**Siswa**
- Absen **masuk & pulang** dengan GPS + peta
- **Jurnal harian** (kegiatan, refleksi, dokumentasi)
- Riwayat absensi per bulan, lihat nilai PKL

**Guru Pembimbing**
- Monitoring absensi + **verifikasi** kehadiran
- Catat **Izin / Sakit** siswa
- Baca & beri **komentar jurnal**
- Lihat nilai dari DUDI
- **Rekap & laporan** (export CSV)

---

# Fitur Utama — DUDI & Admin

**Instruktur DUDI**
- Daftar siswa magang di perusahaannya
- **Form penilaian**: skor kehadiran otomatis + 4 kriteria manual
- Review jurnal + komentar
- Alert untuk siswa yang belum dinilai

**Admin / Admin Jurusan**
- Manajemen **pengguna, DUDI, jurusan, periode**
- **Import CSV** massal (siswa / guru / instruktur)
- **Rekap & laporan** dengan filter
- **Backup database** langsung dari aplikasi

---

# Alur Login & Otorisasi

```
Landing Page
     ↓
Login (pilih role + NIS/NIP/NIK + password)
     ↓
  POST /api/login  →  JWT access + refresh
     ↓
Redirect otomatis ke dashboard sesuai role
     ↓
Router guard memulihkan sesi saat halaman di-refresh
```

- Akses ke route role lain → **dialihkan** ke dashboard miliknya.
- Semua endpoint (kecuali login/register) **wajib JWT**.

---

# Alur Absensi (Siswa)

```
Dashboard Siswa
     ↓
Menu Absensi → deteksi GPS (browser geolocation)
     ↓
Peta menampilkan posisi + lokasi DUDI
     ↓
Server validasi radius (geofencing)
     ↓
[ Valid ] → Absen Masuk Berhasil
     ↓
Absen Pulang aktif setelah ±6–7 jam
```

Validasi **server-side**: waktu diambil dari server (WIB), bukan dari HP siswa.

---

# Alur Jurnal (Siswa)

```
Dashboard → Tulis Jurnal Baru
     ↓
Isi: tanggal (maks 10 hari ke belakang),
     uraian kegiatan, refleksi
     ↓
Simpan → masuk ke daftar jurnal
     ↓
Guru / DUDI membaca & memberi komentar
     ↓
Jurnal terkunci (tidak bisa diedit) setelah dikomentari
```

Satu jurnal per tanggal. Tidak boleh menulis kegiatan masa depan.

---

# Alur Penilaian (DUDI)

```
Login as Instruktur DUDI
     ↓
Daftar siswa magang → pilih siswa
     ↓
Skor kehadiran otomatis dari data absensi
     ↓
Beri nilai 4 kriteria (0–100):
   Kedisiplinan, Tanggung Jawab,
   Kerja Sama, Inisiatif
     ↓
Preview nilai akhir + grade (live)
     ↓
Submit → nilai terkunci & terlihat oleh Guru + Siswa
```

---

# Alur Monitoring (Guru)

```
Dashboard Guru
     ↓
Kartu ringkasan per siswa:
   % kehadiran • jumlah jurnal • nilai DUDI
     ↓
Monitoring Absensi → filter siswa/status/tanggal
     ↓
Verifikasi kehadiran • catat Izin/Sakit
     ↓
Monitoring Jurnal → baca & beri komentar
     ↓
Rekap & Laporan → export CSV
```

---

# Business Logic Kunci

**1. Geofencing (Haversine)** — absensi ditolak jika di luar `radius_allowed` dari lokasi DUDI.

**2. Waktu Server (WIB)** — timestamp tidak bisa dipalsukan dari sisi klien.

**3. Anti-Fraud (Shared Device)** — jika 1 perangkat (IP + User-Agent) dipakai **siswa berbeda** di hari yang sama → ditandai **"mencurigakan"** di dashboard Guru & DUDI.

**4. Aturan Pulang** — absen pulang hanya bisa ≥ 6 jam setelah absen masuk.

**5. Nilai Akhir Berbobot:**
```
final = (skor_kehadiran × 30%) + (rata_rata_manual × 70%)
grade: ≥90 A • ≥80 B+ • ≥70 B • ≥60 C • <60 D
```

---

# Struktur Database

| Tabel | Isi |
|-------|-----|
| **users** | semua pengguna (5 role), relasi ke DUDI & guru pembimbing |
| **dudis** | data perusahaan: alamat, GPS, **radius**, PIC, jurusan |
| **absensis** | kehadiran: timestamp, GPS, foto, **type** (masuk/pulang), status, anti-fraud |
| **jurnals** | jurnal harian + komentar guru & DUDI |
| **penilaians** | nilai DUDI: skor otomatis + 4 kriteria + nilai akhir & grade |
| **periodes** | tahun pelajaran & semester (periode aktif) |
| **jurusans** | daftar jurusan (kode + nama) |

Relasi kunci: **users → dudis**, **absensis/jurnals/penilaians → users** (cascade delete).

---

# API Endpoint (Ringkasan)

| Kategori | Contoh Endpoint |
|----------|-----------------|
| **Auth** | `POST /api/login` • `GET /api/me` • `POST /api/change-password` |
| **Absensi** | `POST /api/absensi` • `GET /api/absensi/history` • `PUT /api/absensi/:id/verify` |
| **Jurnal** | `POST /api/jurnal` • `POST /api/jurnal/comment` |
| **Nilai** | `POST /api/nilai` • `GET /api/nilai/:studentId` |
| **Laporan** | `GET /api/report/{absensi,jurnal,nilai}` • `GET /api/export/{...}` |
| **Admin** | CRUD `/api/admin/{users,dudi,periode,jurusan}` |
| **Import / Backup** | `POST /api/import/{siswa,guru,dudi}` • `GET /api/admin/backups` |

Semua endpoint `/api/*` **dilindungi JWT** dengan pembatasan akses per baris data (row-level scoping).

---

# Keamanan & Non-Functional

- **Keamanan** — password di-hash **bcrypt**, seluruh endpoint JWT-protected, upload file hanya lewat token valid.
- **Privasi Lokasi** — lokasi hanya direkam **saat absensi**, bukan pelacakan terus-menerus.
- **Responsif / Mobile-first** — siswa memakai HP; UI PWA installable dengan bottom navigation.
- **Ketersediaan** — target uptime tinggi karena dipakai harian.
- **Backup** — `pg_dump` harian (retention) + tombol backup di panel admin.

---

# Deployment

**Dua pendekatan:**

**A. Build di VPS** — Go + Node + PostgreSQL + Nginx
systemd service `pkl-tracker` + reverse proxy `/api` & `/uploads`.

**B. Single Binary (direkomendasikan)**
```
make build   →  frontend di-embed ke dalam binary Go (~32 MB)
```
VPS hanya perlu **PostgreSQL + Nginx**. Deploy:
```
update.sh user@vps /var/www/pkl-tracker
   → git pull → make build → scp binary → mv (atomic) → restart
```
Migrasi database otomatis (AutoMigrate) saat server start.

---

# Status & Rencana Lanjutan

**Sudah berjalan:**
- Absensi GPS + geofencing + anti-fraud, Jurnal, Penilaian, Laporan, Import/Export, Backup.

**Dalam pengembangan (roadmap):**
- **Penyimpanan foto absensi ke Google Drive** (infrastruktur sudah siap).
- **Notifikasi** (pengingat absen, notif komentar/nilai).
- **Halaman nilai khusus siswa** (saat ini nilai dilihat terbatas).
- Peta **heatmap** sebaran lokasi absensi siswa.

---

<!-- _class: lead -->
<!-- _paginate: false -->
<!-- _header: '' -->
<!-- _footer: '' -->

# Terima Kasih

### PKL Tracker
**Digitalisasi PKL: tercatat, terverifikasi, dan transparan.**

Pertanyaan & Diskusi
