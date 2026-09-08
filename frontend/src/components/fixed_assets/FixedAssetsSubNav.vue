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

    <!-- Sub-navigation Tabs -->
    <div class="border-b border-gray-200 bg-white rounded-xl px-4 shadow-sm">
      <nav class="-mb-px flex flex-wrap gap-x-8" aria-label="Fixed Assets Navigation">
        <router-link
          to="/fixed-assets"
          exact
          :class="[
            $route.path === '/fixed-assets' ? 'border-emerald-600 text-emerald-700 font-bold' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium',
            'whitespace-nowrap py-3.5 px-1 border-b-2 text-sm flex items-center gap-2 transition'
          ]"
        >
          <i class="icofont-listine-dots text-base"></i>
          <span>Manage Assets</span>
        </router-link>

        <router-link
          to="/fixed-assets/value-adjustments"
          :class="[
            $route.path === '/fixed-assets/value-adjustments' ? 'border-emerald-600 text-emerald-700 font-bold' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium',
            'whitespace-nowrap py-3.5 px-1 border-b-2 text-sm flex items-center gap-2 transition'
          ]"
        >
          <i class="icofont-history text-base"></i>
          <span>Value Adjustments</span>
        </router-link>

        <router-link
          to="/fixed-assets/pivot-engine"
          :class="[
            $route.path === '/fixed-assets/pivot-engine' ? 'border-emerald-600 text-emerald-700 font-bold' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium',
            'whitespace-nowrap py-3.5 px-1 border-b-2 text-sm flex items-center gap-2 transition'
          ]"
        >
          <i class="icofont-chart-histogram text-base"></i>
          <span>Dynamic Pivot Engine</span>
        </router-link>
      </nav>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Fixed Assets Master Register' },
  badge: { type: String, default: '' },
})
</script>
