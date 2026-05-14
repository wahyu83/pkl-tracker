<template>
  <div class="bg-white rounded-2xl p-5 border border-gray-100 card-hover">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-sm text-gray-500 mb-1">{{ label }}</p>
        <p class="text-2xl font-bold text-gray-800">{{ value }}</p>
        <p
          v-if="change"
          :class="[
            'text-xs font-medium mt-1 inline-flex items-center gap-0.5',
            change > 0 ? 'text-accent' : 'text-danger'
          ]"
        >
          <TrendingUpIcon v-if="change > 0" :size="12" />
          <TrendingDownIcon v-else :size="12" />
          {{ Math.abs(change) }}%
        </p>
      </div>
      <div
        :class="[
          'w-11 h-11 rounded-xl flex items-center justify-center flex-shrink-0',
          iconBg
        ]"
      >
        <component :is="icon" :size="22" :class="iconColor" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { TrendingUpIcon, TrendingDownIcon } from 'lucide-vue-next'

defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], required: true },
  icon: { type: Object, required: true },
  iconBg: { type: String, default: 'bg-primary/10' },
  iconColor: { type: String, default: 'text-primary' },
  change: { type: Number, default: null }
})
</script>
