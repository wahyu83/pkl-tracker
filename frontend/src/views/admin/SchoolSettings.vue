<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Pengaturan Sekolah</h2>
        <p class="text-sm text-gray-500 mt-0.5">Identitas sekolah untuk kop dan tanda tangan sertifikat PKL</p>
      </div>
      <button
        @click="save"
        :disabled="saving"
        class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark disabled:opacity-50"
      >
        <SaveIcon :size="16" /> {{ saving ? 'Menyimpan...' : 'Simpan' }}
      </button>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Logo -->
      <div class="bg-white rounded-2xl border border-gray-100 p-5">
        <h3 class="text-sm font-semibold text-gray-800 mb-3">Logo Sekolah</h3>

        <div class="border-2 border-dashed border-gray-200 rounded-xl p-4 flex items-center justify-center bg-gray-50/50 mb-3" style="min-height: 140px">
          <img v-if="logoPreview" :src="logoPreview" alt="Logo sekolah" class="max-h-28 max-w-full object-contain" />
          <div v-else class="text-center text-gray-400">
            <ImageIcon :size="28" class="mx-auto mb-2" />
            <p class="text-xs">Belum ada logo</p>
          </div>
        </div>

        <input ref="fileInput" type="file" accept="image/png,image/jpeg,image/webp" class="hidden" @change="onFileChange" />
        <button
          @click="$refs.fileInput.click()"
          :disabled="uploading"
          class="w-full py-2.5 rounded-xl text-sm font-medium border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-50"
        >
          {{ uploading ? 'Mengunggah...' : 'Pilih & Unggah Logo' }}
        </button>
        <p class="text-xs text-gray-400 mt-2">PNG/JPG/WEBP, maksimal 3 MB. Disarankan latar transparan.</p>
        <p v-if="logoMsg" :class="['text-xs mt-2', logoErr ? 'text-red-500' : 'text-accent']">{{ logoMsg }}</p>
      </div>

      <!-- Identitas -->
      <div class="bg-white rounded-2xl border border-gray-100 p-5 lg:col-span-2">
        <h3 class="text-sm font-semibold text-gray-800 mb-4">Identitas & Tanda Tangan</h3>

        <div class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Nama Sekolah</label>
            <input v-model="form.name" placeholder="SMKN 1 Arahan" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Alamat</label>
            <textarea v-model="form.address" rows="2" placeholder="Jl. Raya Arahan, Indramayu" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"></textarea>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Kota/Kabupaten</label>
              <input v-model="form.city" placeholder="Indramayu" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Nama Kepala Sekolah</label>
              <input v-model="form.headmaster_name" placeholder="Drs. H. ..." class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">NIP Kepala Sekolah</label>
              <input v-model="form.headmaster_nip" placeholder="1965..." class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Awalan Nomor Sertifikat</label>
            <input v-model="form.certificate_prefix" placeholder="421.5/PKL" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            <p class="text-xs text-gray-400 mt-1">
              Nomor sertifikat dibentuk otomatis: <span class="font-mono">{{ form.certificate_prefix || '421.5/PKL' }}/{{ periodeStartYear }}/&lt;NIS&gt;</span>
            </p>
          </div>
        </div>

        <p v-if="msg" :class="['text-xs mt-3', err ? 'text-red-500' : 'text-accent']">{{ msg }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { SaveIcon, ImageIcon } from 'lucide-vue-next'
import { get, put, upload } from '@/api'

const loading = ref(true)
const saving = ref(false)
const uploading = ref(false)
const msg = ref('')
const err = ref(false)
const logoMsg = ref('')
const logoErr = ref(false)
const logoPreview = ref('')
const fileInput = ref(null)

const form = reactive({
  name: '',
  address: '',
  city: '',
  headmaster_name: '',
  headmaster_nip: '',
  certificate_prefix: ''
})

const periodeStartYear = computed(() => new Date().getFullYear())

async function fetchProfile() {
  loading.value = true
  try {
    const res = await get('/school-profile')
    const d = res.data || {}
    Object.assign(form, {
      name: d.name || '',
      address: d.address || '',
      city: d.city || '',
      headmaster_name: d.headmaster_name || '',
      headmaster_nip: d.headmaster_nip || '',
      certificate_prefix: d.certificate_prefix || '421.5/PKL'
    })
    logoPreview.value = d.logo_url || ''
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function save() {
  msg.value = ''
  err.value = false

  if (!form.name.trim()) {
    msg.value = 'Nama sekolah wajib diisi'
    err.value = true
    return
  }

  saving.value = true
  try {
    await put('/school-profile', { ...form })
    msg.value = 'Pengaturan berhasil disimpan!'
    err.value = false
  } catch (e) {
    msg.value = e.message || 'Gagal menyimpan pengaturan'
    err.value = true
  } finally {
    saving.value = false
  }
}

async function onFileChange(event) {
  const file = event.target.files?.[0]
  if (!file) return

  logoMsg.value = ''
  logoErr.value = false

  if (file.size > 3 * 1024 * 1024) {
    logoMsg.value = 'Ukuran file maksimal 3 MB'
    logoErr.value = true
    return
  }

  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('logo', file)
    const res = await upload('/school-profile/logo', fd)
    logoPreview.value = res.data?.logo_url || ''
    logoMsg.value = 'Logo berhasil diunggah!'
    logoErr.value = false
  } catch (e) {
    logoMsg.value = e.message || 'Gagal mengunggah logo'
    logoErr.value = true
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

onMounted(fetchProfile)
</script>
