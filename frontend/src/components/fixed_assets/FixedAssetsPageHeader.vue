<template>
  <div class="space-y-4">
    <!-- Top Header & Breadcrumb -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white p-6 rounded-xl border border-gray-100 shadow-sm">
      <div>
        <div class="flex items-center space-x-2 text-sm text-gray-500 mb-1">
          <span>Financial & Assets Ledger</span>
          <span>/</span>
          <span class="text-emerald-700 font-medium">Fixed Assets Register</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight flex flex-wrap items-center gap-2">
          {{ title }}
          <span v-if="badge" class="text-xs bg-emerald-100 text-emerald-800 font-semibold px-2.5 py-0.5 rounded-full border border-emerald-200">
            {{ badge }}
          </span>
        </h1>
        <p class="text-xs text-gray-500 mt-1">
          IPSAS 17 Master Fixed Asset Register & Statutory Revaluation Ledger (UGX 31.02B FB_COST baseline).
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <router-link
          to="/fixed-assets/depreciation"
          class="inline-flex items-center px-4 py-2 bg-indigo-50 text-indigo-700 rounded-lg hover:bg-indigo-100 text-xs font-semibold transition"
        >
          <i class="icofont-clock-time mr-1.5"></i>
          Run Monthly Depreciation
        </router-link>
        <router-link
          to="/fixed-assets/new"
          class="inline-flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 text-xs font-semibold shadow-sm transition"
        >
          <i class="icofont-plus-circle mr-1.5"></i>
          New Fixed Asset
        </router-link>
      </div>
    </div>

  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Fixed Assets Master Register' },
  badge: { type: String, default: '' },
})
</script>
