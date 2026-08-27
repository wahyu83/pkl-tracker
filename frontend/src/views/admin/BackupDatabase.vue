<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Backup Database</h2>
        <p class="text-sm text-gray-500 mt-0.5">Buat dan unduh backup database (pg_dump) secara manual</p>
      </div>
      <button @click="createBackup" :disabled="saving"
        class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark disabled:opacity-50">
        <DatabaseBackupIcon :size="16" /> {{ saving ? 'Membuat Backup...' : 'Buat Backup' }}
      </button>
    </div>

    <!-- Info card -->
    <div class="bg-white rounded-2xl border border-gray-100 p-5 mb-6">
      <div class="flex items-start gap-3">
        <div class="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center flex-shrink-0">
          <DatabaseBackupIcon :size="20" class="text-primary" />
        </div>
        <div>
          <h3 class="text-sm font-semibold text-gray-800">Backup otomatis penuh (pg_dump)</h3>
          <p class="text-xs text-gray-500 mt-0.5 leading-relaxed">
            Backup menyimpan seluruh data (siswa, guru, DUDI, absensi, jurnal, nilai) ke file
            <span class="font-mono text-gray-600">.sql.gz</span> yang bisa diunduh dan dipulihkan kapan saja.
            Untuk backup terjadwal harian, gunakan <span class="font-mono text-gray-600">./backup.sh --install-cron</span> di server.
          </p>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <!-- Backup list -->
    <div v-else class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-gray-800">Daftar Backup</h3>
        <span v-if="backups.length" class="text-xs text-gray-500">{{ backups.length }} file</span>
      </div>

      <div v-if="backups.length === 0" class="text-center py-10">
        <div class="w-12 h-12 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-3">
          <DatabaseBackupIcon :size="22" class="text-gray-300" />
        </div>
        <p class="text-sm text-gray-400">Belum ada backup. Klik "Buat Backup" untuk membuat yang pertama.</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-gray-50/50">
          <tr class="text-left text-xs text-gray-500">
            <th class="px-5 py-3 font-medium">Nama File</th>
            <th class="px-5 py-3 font-medium">Tanggal</th>
            <th class="px-5 py-3 font-medium">Ukuran</th>
            <th class="px-5 py-3 font-medium text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-50">
          <tr v-for="b in backups" :key="b.filename" class="hover:bg-gray-50/50">
            <td class="px-5 py-3 font-mono text-xs text-gray-700">{{ b.filename }}</td>
            <td class="px-5 py-3 text-xs text-gray-500">{{ formatDate(b.created_at) }}</td>
            <td class="px-5 py-3 text-xs text-gray-500">{{ b.size_human }}</td>
            <td class="px-5 py-3 text-right">
              <div class="flex justify-end gap-1">
                <button @click="downloadBackup(b)" title="Unduh"
                  class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10">
                  <DownloadIcon :size="16" />
                </button>
                <button @click="confirmDelete(b)" title="Hapus"
                  class="p-1.5 rounded-lg text-gray-400 hover:text-danger hover:bg-danger/10">
                  <TrashIcon :size="16" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Delete Confirm -->
    <div v-if="showDelete" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="showDelete = null">
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-12 h-12 bg-danger/10 rounded-full flex items-center justify-center mx-auto mb-3">
          <AlertCircleIcon :size="24" class="text-danger" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-1">Hapus Backup?</h3>
        <p class="text-sm text-gray-500 mb-4 break-all">{{ deleteTarget?.filename }}</p>
        <div class="flex gap-2">
          <button @click="showDelete = null" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
          <button @click="deleteBackup" :disabled="saving" class="flex-1 py-2.5 bg-danger text-white rounded-xl text-sm font-medium hover:bg-red-600 disabled:opacity-50">{{ saving ? 'Menghapus...' : 'Hapus' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Database as DatabaseBackupIcon, DownloadIcon, TrashIcon, AlertCircleIcon } from 'lucide-vue-next'
import { get, post, del, downloadFile } from '@/api'

const backups = ref([])
const loading = ref(true)
const saving = ref(false)
const showDelete = ref(null)
const deleteTarget = ref(null)

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

async function fetchBackups() {
  loading.value = true
  try {
    const res = await get('/admin/backups')
    backups.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function createBackup() {
  saving.value = true
  try {
    const res = await post('/admin/backups')
    backups.value = [res.data, ...backups.value]
    alert('Backup berhasil dibuat: ' + res.data.filename)
  } catch (e) {
    alert('Backup gagal: ' + e.message)
  } finally {
    saving.value = false
  }
}

async function downloadBackup(b) {
  try {
    await downloadFile(b.download_url, b.filename)
  } catch (e) {
    alert('Unduh gagal: ' + e.message)
  }
}

function confirmDelete(b) { deleteTarget.value = b; showDelete.value = true }

async function deleteBackup() {
  saving.value = true
  try {
    await del('/admin/backups/' + deleteTarget.value.filename)
    showDelete.value = null
    backups.value = backups.value.filter(x => x.filename !== deleteTarget.value.filename)
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(fetchBackups)
</script>
