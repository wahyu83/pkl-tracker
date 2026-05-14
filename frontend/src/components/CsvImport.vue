<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="$emit('close')">
    <div class="fixed inset-0 bg-black/50" />
    <div class="relative bg-white rounded-2xl shadow-xl p-6 w-full max-w-lg">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-bold text-gray-800">{{ title }}</h3>
        <button @click="$emit('close')" class="p-1.5 rounded-lg hover:bg-gray-100 text-gray-400">
          <XIcon :size="20" />
        </button>
      </div>

      <p class="text-sm text-gray-500 mb-1">Upload file CSV untuk import {{ label }}.</p>
      <p class="text-xs text-gray-400 mb-4">
        Format CSV: <code class="bg-gray-100 px-1.5 py-0.5 rounded text-xs">{{ csvColumns }}</code>
      </p>

      <button
        v-if="sampleCsv"
        @click="downloadSample"
        class="text-xs text-primary font-medium hover:underline mb-4 block"
      >
        Download template CSV
      </button>

      <div
        class="border-2 border-dashed border-gray-300 rounded-xl p-8 text-center mb-4 cursor-pointer hover:border-primary transition-colors"
        :class="file ? 'border-accent bg-accent/5' : ''"
        @click="triggerInput"
        @dragover.prevent="dragover = true"
        @dragleave="dragover = false"
        @drop.prevent="handleDrop"
      >
        <input ref="fileInput" type="file" accept=".csv" class="hidden" @change="handleFile" />
        <template v-if="!file">
          <UploadIcon :size="32" class="text-gray-400 mx-auto mb-2" />
          <p class="text-sm text-gray-500">Klik atau drag & drop file CSV</p>
          <p class="text-xs text-gray-400 mt-1">Hanya file .csv</p>
        </template>
        <template v-else>
          <FileTextIcon :size="32" class="text-accent mx-auto mb-2" />
          <p class="text-sm font-medium text-gray-800">{{ file.name }}</p>
          <p class="text-xs text-gray-400 mt-1">{{ (file.size / 1024).toFixed(1) }} KB</p>
        </template>
      </div>

      <div v-if="result" :class="['rounded-xl p-4 mb-4', result.errors?.length ? 'bg-warning/10' : 'bg-accent/10']">
        <p class="text-sm font-semibold mb-2">{{ result.message }}</p>
        <div class="text-xs space-y-1">
          <p class="text-accent font-medium">Berhasil: {{ result.imported }} data</p>
          <p v-if="result.skipped" class="text-gray-500">Dilewati/duplikat: {{ result.skipped }}</p>
          <div v-if="result.errors?.length" class="mt-2">
            <p class="text-warning font-medium">Error:</p>
            <ul class="list-disc list-inside text-gray-600">
              <li v-for="(e, i) in result.errors" :key="i">{{ e }}</li>
            </ul>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="$emit('close')"
          class="flex-1 py-2.5 rounded-xl text-sm font-medium border border-gray-200 text-gray-600 hover:bg-gray-50 transition-colors"
        >
          {{ result ? 'Tutup' : 'Batal' }}
        </button>
        <button
          v-if="!result"
          @click="handleUpload"
          :disabled="!file || uploading"
          class="flex-1 py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
        >
          <LoaderIcon v-if="uploading" :size="18" class="animate-spin" />
          <span>{{ uploading ? 'Mengimport...' : 'Import' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { XIcon, UploadIcon, FileTextIcon, LoaderIcon } from 'lucide-vue-next'
import { post } from '../api'

const props = defineProps({
  show: Boolean,
  title: String,
  label: String,
  endpoint: String,
  csvColumns: String,
  sampleCsv: { type: String, default: '' }
})

const emit = defineEmits(['close', 'done'])

const fileInput = ref(null)
const file = ref(null)
const uploading = ref(false)
const result = ref(null)
const dragover = ref(false)

function triggerInput() {
  fileInput.value?.click()
}

function handleFile(e) {
  const f = e.target.files?.[0]
  if (f) file.value = f
}

function handleDrop(e) {
  const f = e.dataTransfer?.files?.[0]
  if (f?.name.endsWith('.csv')) file.value = f
}

async function handleUpload() {
  if (!file.value) return
  uploading.value = true
  result.value = null
  try {
    const formData = new FormData()
    formData.append('file', file.value)
    const res = await fetch('/api' + props.endpoint, {
      method: 'POST',
      headers: { Authorization: 'Bearer ' + (localStorage.getItem('token') || '') },
      body: formData
    })
    result.value = await res.json()
    if (res.ok) emit('done')
  } catch (e) {
    result.value = { message: 'Gagal import', errors: [e.message], imported: 0, skipped: 0 }
  } finally {
    uploading.value = false
  }
}

function downloadSample() {
  const blob = new Blob([props.sampleCsv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `template-${props.label.toLowerCase().replace(/\s/g, '-')}.csv`
  a.click()
  URL.revokeObjectURL(url)
}
</script>
