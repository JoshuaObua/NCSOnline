<template>
  <LayoutDefault title="Stores & Inventory Management">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-box text-xs"></i> Warehouse & Logistics
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Stores & Electronic Inventory Registry</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Real-time stock ledger, sports kits & equipment distributions, engineering spares, and Goods Received Notes (GRN).
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchData" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
        </div>
      </div>

      <!-- Macro Metric KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Catalogued SKUs</h5>
              <h2>{{ items.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-barcode"></i> Active Catalog</span>
            </div>
            <div class="banner-img">
              <i class="icofont-box"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Total Units In Stock</h5>
              <h2>{{ totalUnits }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Main Stores</span>
            </div>
            <div class="banner-img">
              <i class="icofont-cubes"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Low Stock Alerts</h5>
              <h2>{{ lowStockCount }}</h2>
              <span class="badge" :class="lowStockCount > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-warning-alt"></i> Reorder Trigger
              </span>
            </div>
            <div class="banner-img">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Inventory Valuation</h5>
              <h2>UGX {{ Number(totalValuation).toLocaleString() }}</h2>
              <span class="badge badge-info"><i class="icofont-money"></i> Stock Value</span>
            </div>
            <div class="banner-img">
              <i class="icofont-money"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div class="flex items-center gap-1.5 border-b border-slate-200 dark:border-slate-800 pb-2">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          @click="activeTab = tab.id"
          class="px-3.5 py-2 rounded-lg text-xs font-bold transition flex items-center gap-2"
          :class="activeTab === tab.id 
            ? 'bg-blue-600 text-white shadow-sm' 
            : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
        >
          <i :class="tab.icon"></i> {{ tab.label }}
        </button>
      </div>

      <!-- TAB 1: INVENTORY ITEMS LEDGER -->
      <div v-if="activeTab === 'items'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-3.5 px-5">
            <div>
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">Live Stores Bin Card Ledger</h4>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Physical bin locations, available balances, and reorder levels</p>
            </div>
          </div>

          <div class="card-body p-0 overflow-x-auto">
            <table class="table table-striped table-hover mb-0 text-left text-xs">
              <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 uppercase tracking-wider font-semibold">
                <tr>
                  <th class="py-3 px-4">SKU / Item Name</th>
                  <th class="py-3 px-4">Category</th>
                  <th class="py-3 px-4">Location</th>
                  <th class="py-3 px-4">Unit Cost</th>
                  <th class="py-3 px-4">Stock Level</th>
                  <th class="py-3 px-4">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                <tr v-for="item in items" :key="item.id">
                  <td class="py-3 px-4">
                    <div class="font-bold text-slate-900 dark:text-white">{{ item.item_name }}</div>
                    <div class="font-mono text-[11px] text-blue-600 dark:text-blue-400">{{ item.sku_code }}</div>
                  </td>
                  <td class="py-3 px-4 font-mono font-medium text-slate-600 dark:text-slate-400">{{ item.category }}</td>
                  <td class="py-3 px-4 text-slate-700 dark:text-slate-300">{{ item.warehouse_bin_location }}</td>
                  <td class="py-3 px-4 font-mono text-slate-700 dark:text-slate-300">UGX {{ Number(item.unit_cost_ugx).toLocaleString() }}</td>
                  <td class="py-3 px-4 font-mono font-bold text-slate-900 dark:text-white">{{ item.quantity_on_hand }} {{ item.unit_of_measure }}</td>
                  <td class="py-3 px-4">
                    <span class="badge" :class="item.status === 'LOW_STOCK' ? 'badge-warning' : 'badge-success'">
                      {{ item.status }}
                    </span>
                  </td>
                </tr>
                <tr v-if="!items.length">
                  <td colspan="6" class="text-center py-6 text-xs text-slate-400">
                    No store items registered in the inventory ledger.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- TAB 2: GRN -->
      <div v-if="activeTab === 'grn'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Goods Received Notes (GRN) Inward Logs</h4>
          </div>
          <div class="card-body p-5 space-y-3">
            <div v-for="g in grns" :key="g.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30">
              <div class="flex items-center justify-between">
                <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400">{{ g.grn_number }}</span>
                <span class="badge badge-success">Received</span>
              </div>
              <h5 class="text-sm font-bold text-slate-900 dark:text-white mt-1">{{ g.supplier_name }}</h5>
              <div class="text-xs text-slate-500 pt-2 mt-2 border-t border-slate-200 dark:border-slate-700/50">
                PO Reference: <strong class="text-slate-700 dark:text-slate-300">{{ g.po_reference }}</strong> · Received Value: <strong class="text-blue-600 dark:text-blue-400 font-mono">UGX {{ Number(g.total_received_value_ugx).toLocaleString() }}</strong>
              </div>
            </div>
            <div v-if="!grns.length" class="text-center py-6 text-xs text-slate-400">
              No Goods Received Notes recorded yet.
            </div>
          </div>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { apiGet } from '@/api/client'

const loading = ref(false)
const items = ref([])
const grns = ref([])
const activeTab = ref('items')

const tabs = [
  { id: 'items', label: 'Electronic Bin Cards & Stock Ledger', icon: 'icofont-box' },
  { id: 'grn', label: 'Goods Received Notes (GRN)', icon: 'icofont-delivery-time' }
]

const totalUnits = computed(() => {
  return items.value.reduce((acc, cur) => acc + (cur.quantity_on_hand || 0), 0)
})

const lowStockCount = computed(() => {
  return items.value.filter(i => (i.quantity_on_hand || 0) <= (i.minimum_reorder_level || 5)).length
})

const totalValuation = computed(() => {
  return items.value.reduce((acc, cur) => acc + ((cur.quantity_on_hand || 0) * (cur.unit_cost_ugx || 0)), 0)
})

async function fetchData() {
  loading.value = true
  try {
    const [invRes, grnRes] = await Promise.allSettled([
      apiGet('/stores/inventory'),
      apiGet('/stores/grn')
    ])

    if (invRes.status === 'fulfilled') {
      items.value = invRes.value.data?.inventory || invRes.value.inventory || []
    }
    if (grnRes.status === 'fulfilled') {
      grns.value = grnRes.value.data?.grns || grnRes.value.grns || []
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
