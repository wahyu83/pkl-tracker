<template>
  <div class="print-root">
    <!-- Toolbar: tidak tercetak -->
    <div class="toolbar">
      <div class="toolbar-inner">
        <div>
          <p class="toolbar-title">Pratinjau Sertifikat</p>
          <p class="toolbar-sub">
            {{ certs.length }} sertifikat • A4 landscape • gunakan "Simpan sebagai PDF" pada dialog cetak
          </p>
        </div>
        <div class="toolbar-actions">
          <button class="btn btn-ghost" @click="goBack">Kembali</button>
          <button class="btn btn-primary" @click="print">Cetak / Simpan PDF</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="state">Memuat sertifikat...</div>
    <div v-else-if="error" class="state state-error">{{ error }}</div>
    <div v-else-if="certs.length === 0" class="state">Tidak ada data siswa untuk dicetak.</div>

    <div v-else class="pages">
      <div v-for="c in certs" :key="c.student.id" class="cert-page">
        <div class="frame-outer">
          <div class="frame-inner">
            <p class="title">SERTIFIKAT</p>
            <p class="subtitle">PRAKTIK KERJA LAPANGAN</p>
            <p class="school-name">{{ school.name || 'Nama Sekolah' }}</p>
            <p class="cert-number">Nomor: {{ c.certificate_number }}</p>

            <div class="body">
              <p class="lead">Yang bertanda tangan di bawah ini menerangkan bahwa:</p>

              <p class="student-name">{{ c.student.full_name }}</p>
              <p class="student-meta">
                NIS {{ c.student.nis_nip_nik }}<span v-if="c.student.jurusan"> • Kompetensi Keahlian {{ c.student.jurusan }}</span>
              </p>

              <p class="lead">
                telah melaksanakan Praktik Kerja Lapangan (PKL) di
                <strong>{{ c.dudi?.company_name || '.......................' }}</strong>
                <span v-if="c.dudi?.address">, {{ c.dudi.address }}</span><span v-if="periodRange">, pada periode {{ periodRange }}</span>.
              </p>

              <div class="score-box">
                <div class="score-item">
                  <p class="score-label">Nilai Akhir</p>
                  <p class="score-value">{{ c.score.assessed ? Number(c.score.final_score).toFixed(2) : '-' }}</p>
                </div>
                <div class="score-item">
                  <p class="score-label">Predikat</p>
                  <p class="score-value">{{ c.score.assessed ? (c.score.final_grade || '-') : '-' }}</p>
                </div>
                <div class="score-item">
                  <p class="score-label">Kehadiran</p>
                  <p class="score-value">
                    {{ c.score.attendance_score_auto ? Number(c.score.attendance_score_auto).toFixed(0) + '%' : '-' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Tanda tangan: Kepala Sekolah & PIC DUDI (Pembimbing) -->
            <div class="signatures">
              <div class="sig">
                <p class="sig-line">Pembimbing</p>
                <p class="sig-company">{{ c.dudi?.company_name || '.......................' }}</p>
                <div class="sig-space"></div>
                <p class="sig-name">{{ c.dudi?.pic_name || '................................' }}</p>
              </div>
              <div class="sig">
                <p class="sig-line">{{ school.city || '' }}<span v-if="school.city">, </span>{{ today }}</p>
                <p class="sig-line">Kepala Sekolah</p>
                <div class="sig-space"></div>
                <p class="sig-name">{{ school.headmaster_name || '................................' }}</p>
                <p v-if="school.headmaster_nip" class="sig-nip">NIP. {{ school.headmaster_nip }}</p>
              </div>
            </div>
          </div>

          <!-- Keterangan di luar kotak sertifikat -->
          <div class="cert-footer">
            <span>{{ c.student.full_name }} - {{ c.student.jurusan || '-' }}</span>
            <span>Dicetak dari PKL Tracker - {{ school.name || 'Nama Sekolah' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get } from '@/api'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref('')
const school = ref({})
const periode = ref(null)
const certs = ref([])

const today = computed(() =>
  new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
)

const periodRange = computed(() => {
  const p = periode.value
  if (!p?.start_date || !p?.end_date) return ''
  const opts = { day: 'numeric', month: 'long', year: 'numeric' }
  return `${new Date(p.start_date).toLocaleDateString('id-ID', opts)} – ${new Date(p.end_date).toLocaleDateString('id-ID', opts)}`
})

function goBack() {
  if (window.history.length > 1) router.back()
  else router.push('/')
}

function print() {
  window.print()
}

// Tunggu gambar (logo) selesai dimuat sebelum mencetak otomatis.
async function waitForImages() {
  const imgs = Array.from(document.images)
  await Promise.all(imgs.map(img => img.complete
    ? Promise.resolve()
    : new Promise(resolve => { img.onload = resolve; img.onerror = resolve })
  ))
}

async function fetchData() {
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams()
    if (route.query.periode_id) params.append('periode_id', route.query.periode_id)
    if (route.query.dudi_id) params.append('dudi_id', route.query.dudi_id)

    const studentId = route.query.student_id
    if (studentId) {
      const res = await get(`/certificate/${studentId}?${params.toString()}`)
      school.value = res.school || {}
      periode.value = res.periode || null
      certs.value = res.data ? [res.data] : []
    } else {
      const qs = params.toString()
      const res = await get(`/certificate${qs ? '?' + qs : ''}`)
      school.value = res.school || {}
      periode.value = res.periode || null
      certs.value = res.data || []
    }
  } catch (e) {
    error.value = e.message || 'Gagal memuat data sertifikat'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await fetchData()
  await nextTick()
  await waitForImages()
  setTimeout(() => window.print(), 400)
})
</script>

<style scoped>
.print-root {
  background: #e5e7eb;
  min-height: 100vh;
}

/* --- Toolbar (layar saja) --- */
.toolbar {
  position: sticky;
  top: 0;
  z-index: 10;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
}
.toolbar-inner {
  max-width: 1100px;
  margin: 0 auto;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}
.toolbar-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}
.toolbar-sub {
  font-size: 12px;
  color: #6b7280;
}
.toolbar-actions {
  display: flex;
  gap: 8px;
}
.btn {
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid transparent;
}
.btn-primary {
  background: #1e3a8a;
  color: #fff;
}
.btn-ghost {
  background: #f3f4f6;
  color: #374151;
  border-color: #e5e7eb;
}

.state {
  text-align: center;
  padding: 48px 16px;
  color: #6b7280;
  font-size: 14px;
}
.state-error {
  color: #dc2626;
}

.pages {
  padding: 24px 16px 48px;
}

/* --- Halaman sertifikat A4 landscape --- */
.cert-page {
  width: 297mm;
  height: 210mm;
  margin: 0 auto 24px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  page-break-after: always;
  break-after: page;
  overflow: hidden;
}
.cert-page:last-child {
  page-break-after: auto;
  break-after: auto;
}

.frame-outer {
  width: 100%;
  height: 100%;
  padding: 7mm 8mm 5mm;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}
.frame-inner {
  width: 100%;
  flex: 1;
  min-height: 0;
  box-sizing: border-box;
  border: 3px double #1e3a8a;
  padding: 8mm 12mm;
  display: flex;
  flex-direction: column;
  font-family: "Times New Roman", Georgia, serif;
  color: #111827;
}

/* Keterangan kecil di luar kotak sertifikat */
.cert-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 6mm;
  padding-top: 2.5mm;
  font-family: "Times New Roman", Georgia, serif;
  font-size: 8pt;
  color: #6b7280;
}

.title {
  margin-top: 6mm;
  text-align: center;
  font-size: 26pt;
  font-weight: 700;
  letter-spacing: 6px;
  color: #1e3a8a;
}
.subtitle {
  text-align: center;
  font-size: 13pt;
  letter-spacing: 3px;
  margin-top: 2px;
}
.school-name {
  text-align: center;
  font-size: 13pt;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-top: 3mm;
  color: #111827;
}
.cert-number {
  text-align: center;
  font-size: 10.5pt;
  color: #374151;
  margin-top: 4px;
}

.body {
  margin-top: 7mm;
  text-align: center;
  flex: 1;
}
.lead {
  font-size: 12pt;
  line-height: 1.6;
}
.student-name {
  font-size: 22pt;
  font-weight: 700;
  margin-top: 3mm;
  text-decoration: underline;
  text-underline-offset: 3px;
}
.student-meta {
  font-size: 11pt;
  color: #374151;
  margin-top: 1mm;
}

.score-box {
  display: inline-flex;
  gap: 10mm;
  border: 1.5px solid #1e3a8a;
  border-radius: 6px;
  padding: 3mm 8mm;
  margin-top: 6mm;
}
.score-item {
  text-align: center;
  min-width: 26mm;
}
.score-label {
  font-size: 9pt;
  color: #4b5563;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.score-value {
  font-size: 16pt;
  font-weight: 700;
  color: #1e3a8a;
}

.signatures {
  display: flex;
  justify-content: space-between;
  margin-top: 6mm;
  padding: 0 10mm;
}
.sig {
  text-align: center;
  font-size: 11pt;
  line-height: 1.4;
}
.sig-line {
  line-height: 1.5;
}
.sig-company {
  font-weight: 600;
  line-height: 1.4;
  text-transform: uppercase;
}
.sig-space {
  height: 18mm;
}
.sig-name {
  font-weight: 700;
  text-decoration: underline;
  text-underline-offset: 3px;
}
.sig-nip {
  font-size: 10pt;
  color: #374151;
}

/* --- Cetak --- */
@page {
  size: A4 landscape;
  margin: 0;
}

@media print {
  .toolbar {
    display: none !important;
  }
  .print-root {
    background: #fff;
  }
  .pages {
    padding: 0;
  }
  .cert-page {
    margin: 0;
    box-shadow: none;
    width: 297mm;
    height: 210mm;
  }
}
</style>
