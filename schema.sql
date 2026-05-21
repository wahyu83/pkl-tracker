-- PKL Tracker Database Schema (PostgreSQL)
-- Run: psql -U pkl_user -d pkl_db -f schema.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. dudis
CREATE TABLE IF NOT EXISTS dudis (
    id UUID PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10,8) DEFAULT 0,
    longitude DECIMAL(11,8) DEFAULT 0,
    radius_allowed INT DEFAULT 500,
    pic_name VARCHAR(255) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'student' CHECK (role IN ('student','teacher','dudi','admin')),
    nis_nip_nik VARCHAR(50) NOT NULL UNIQUE,
    dudi_id UUID REFERENCES dudis(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. absensis
CREATE TABLE IF NOT EXISTS absensis (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ NOT NULL,
    latitude DECIMAL(10,8) DEFAULT 0,
    longitude DECIMAL(11,8) DEFAULT 0,
    photo_url TEXT DEFAULT '',
    type VARCHAR(10) NOT NULL DEFAULT 'masuk' CHECK (type IN ('masuk','pulang')),
    status VARCHAR(20) NOT NULL DEFAULT 'hadir' CHECK (status IN ('hadir','terlambat','izin','sakit')),
    is_verified BOOLEAN DEFAULT FALSE,
    ip_address VARCHAR(45) DEFAULT '',
    user_agent TEXT DEFAULT '',
    is_suspicious BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_absensis_student_id ON absensis(student_id);

-- 4. jurnals
CREATE TABLE IF NOT EXISTS jurnals (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    activity TEXT NOT NULL,
    documentation_url TEXT DEFAULT '',
    reflection TEXT DEFAULT '',
    teacher_comment TEXT DEFAULT '',
    dudi_comment TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_jurnals_student_id ON jurnals(student_id);

-- 5. penilaians
CREATE TABLE IF NOT EXISTS penilaians (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    dudi_id UUID NOT NULL REFERENCES dudis(id) ON DELETE CASCADE,
    attendance_score_auto DECIMAL(5,2) DEFAULT 0,
    discipline INT CHECK (discipline BETWEEN 1 AND 5) DEFAULT 0,
    responsibility INT CHECK (responsibility BETWEEN 1 AND 5) DEFAULT 0,
    teamwork INT CHECK (teamwork BETWEEN 1 AND 5) DEFAULT 0,
    initiative INT CHECK (initiative BETWEEN 1 AND 5) DEFAULT 0,
    final_score DECIMAL(5,2) DEFAULT 0,
    final_grade VARCHAR(2) DEFAULT '',
    notes TEXT DEFAULT '',
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_penilaians_student_id ON penilaians(student_id);

-- 6. periodes
CREATE TABLE IF NOT EXISTS periodes (
    id UUID PRIMARY KEY,
    tahun_pelajaran VARCHAR(20) NOT NULL,
    semester VARCHAR(10) NOT NULL CHECK (semester IN ('ganjil','genap')),
    is_active BOOLEAN DEFAULT FALSE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- SEED DATA (same as in backend/main.go seedDatabase + seedPeriode)
-- ============================================================

INSERT INTO dudis (id, company_name, address, latitude, longitude, radius_allowed, pic_name, phone) VALUES
  ('11111111-1111-1111-1111-111111111101', 'PT. Teknologi Maju',    'Jl. Sudirman No. 123, Jakarta Pusat',     -6.2088, 106.8456, 500, 'Hendra Gunawan', '021-5551234'),
  ('11111111-1111-1111-1111-111111111102', 'PT. Sejahtera Abadi',    'Jl. Gatot Subroto No. 45, Jakarta Selatan', -6.2297, 106.8243, 300, 'Ratna Dewi',     '021-5555678');

INSERT INTO users (id, full_name, email, password_hash, role, nis_nip_nik, dudi_id) VALUES
  ('22222222-2222-2222-2222-222222222201', 'Admin Utama',            'admin@pkl.local', 'password_hash_admin',  'admin',   'ADM-001',   NULL),
  ('22222222-2222-2222-2222-222222222202', 'Budi Santoso, S.Kom',    'budi@pkl.local',  'password_hash_teacher','teacher', '19850101',  NULL),
  ('22222222-2222-2222-2222-222222222203', 'Ahmad Rizky',            'ahmad@pkl.local', 'password_hash_student','student', '20230001',  '11111111-1111-1111-1111-111111111101'),
  ('22222222-2222-2222-2222-222222222204', 'Siti Nurhaliza',         'siti@pkl.local',  'password_hash_student','student', '20230002',  '11111111-1111-1111-1111-111111111102'),
  ('22222222-2222-2222-2222-222222222205', 'PT. Teknologi Maju',     'info@teknologimaju.id', 'password_hash_dudi','dudi', 'D-001',     '11111111-1111-1111-1111-111111111101');

INSERT INTO absensis (id, student_id, timestamp, latitude, longitude, status, is_verified) VALUES
  ('33333333-3333-3333-3333-333333333301', '22222222-2222-2222-2222-222222222203', CURRENT_DATE + TIME '07:45', -6.2088, 106.8456, 'hadir',    TRUE),
  ('33333333-3333-3333-3333-333333333302', '22222222-2222-2222-2222-222222222203', CURRENT_DATE - INTERVAL '1 day' + TIME '08:15', -6.2088, 106.8456, 'terlambat', TRUE),
  ('33333333-3333-3333-3333-333333333303', '22222222-2222-2222-2222-222222222203', CURRENT_DATE - INTERVAL '2 days' + TIME '07:30', -6.2088, 106.8456, 'hadir',    TRUE),
  ('33333333-3333-3333-3333-333333333304', '22222222-2222-2222-2222-222222222204', CURRENT_DATE - INTERVAL '1 day' + TIME '07:50', -6.2297, 106.8243, 'hadir',    TRUE);

INSERT INTO jurnals (id, student_id, date, activity, reflection, teacher_comment, dudi_comment) VALUES
  ('44444444-4444-4444-4444-444444444401', '22222222-2222-2222-2222-222222222203', CURRENT_DATE,
   'Mempelajari framework Laravel dan membuat CRUD untuk modul inventaris.',
   'Belajar banyak tentang MVC pattern.', NULL, NULL),
  ('44444444-4444-4444-4444-444444444402', '22222222-2222-2222-2222-222222222203', CURRENT_DATE - INTERVAL '1 day',
   'Debugging aplikasi internal, memperbaiki bug modul pelaporan.',
   NULL, 'Terus tingkatkan!', NULL),
  ('44444444-4444-4444-4444-444444444403', '22222222-2222-2222-2222-222222222204', CURRENT_DATE - INTERVAL '1 day',
   'Membantu tim network maintenance server.',
   NULL, NULL, 'Siswa menunjukkan antusiasme baik.');

INSERT INTO periodes (id, tahun_pelajaran, semester, is_active, start_date, end_date) VALUES
  ('55555555-5555-5555-5555-555555555501', '2025/2026', 'ganjil', TRUE,  '2025-07-14', '2025-12-20'),
  ('55555555-5555-5555-5555-555555555502', '2025/2026', 'genap',  FALSE, '2026-01-05', '2026-06-20');
