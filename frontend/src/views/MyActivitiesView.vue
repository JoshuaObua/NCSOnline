<template>
  <LayoutDefault title="My Activities">
    <div class="p-6 max-w-7xl mx-auto">
      <header class="mb-5">
        <h2 class="text-2xl font-semibold text-gray-900">My Activities</h2>
        <p class="text-xs text-gray-500 mt-1">A chronological ledger of everything your account has done in the system.</p>
      </header>

      <div class="bg-white border border-gray-200 rounded-xl overflow-hidden">
        <div v-if="loading" class="p-6 space-y-2">
          <div v-for="i in 6" :key="i" class="h-10 bg-gray-100 rounded animate-pulse"/>
        </div>
        <div v-else-if="!logs.length" class="p-10 text-center text-sm text-gray-500">
          No activity recorded yet.
        </div>
        <table v-else class="w-full text-sm">
          <thead class="bg-gray-50 text-xs uppercase tracking-wider text-gray-500">
            <tr>
              <th class="px-4 py-3 text-left">When (UTC)</th>
              <th class="px-4 py-3 text-left">Action</th>
              <th class="px-4 py-3 text-left">Target</th>
              <th class="px-4 py-3 text-left">Status</th>
              <th class="px-4 py-3 text-left">IP / Device</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="l in logs" :key="l.id">
              <td class="px-4 py-2.5 font-mono text-xs text-gray-600 whitespace-nowrap">{{ formatUTC(l.created_at) }}</td>
              <td class="px-4 py-2.5">
                <div class="text-gray-900 font-medium">{{ l.event_type || l.action || '—' }}</div>
                <div v-if="l.endpoint" class="text-[10px] text-gray-400 font-mono">{{ l.method }} {{ l.endpoint }}</div>
              </td>
              <td class="px-4 py-2.5 text-xs text-gray-700">{{ targetLabel(l) }}</td>
              <td class="px-4 py-2.5">
                <span :class="statusClass(l)" class="text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-full">{{ l.event_status || statusFromCode(l.response_code) }}</span>
              </td>
              <td class="px-4 py-2.5 text-xs text-gray-500">
                <div>{{ l.ip_address || '—' }} <span v-if="l.geo_country" class="text-gray-400">· {{ l.geo_city || l.geo_country }}</span></div>
                <div class="text-[10px] text-gray-400">{{ l.browser || '' }}{{ l.os_name ? ' · ' + l.os_name : '' }}</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="total > perPage" class="flex items-center justify-between mt-4 text-xs text-gray-500">
        <span>Showing {{ ((page - 1) * perPage) + 1 }}–{{ Math.min(page * perPage, total) }} of {{ total }}</span>
        <div class="flex gap-2">
          <button :disabled="page === 1" @click="page--; load()" class="px-2 py-1 border border-gray-200 rounded disabled:opacity-40">← Prev</button>
          <button :disabled="page * perPage >= total" @click="page++; load()" class="px-2 py-1 border border-gray-200 rounded disabled:opacity-40">Next →</button>
        </div>
      </div>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { getMyActivities } from '@/api/security.js'

const logs = ref([])
const loading = ref(false)
const page = ref(1)
const perPage = ref(25)
const total = ref(0)

function formatUTC(d) {
  if (!d) return ''
  const dt = new Date(d)
  return dt.toISOString().replace('T', ' ').slice(0, 19) + 'Z'
}

function targetLabel(l) {
  const r = l.resource || ''
  const id = l.resource_id ? ` · ${l.resource_id}` : ''
  return r ? `${r}${id}` : '—'
}

function statusFromCode(code) {
  if (!code) return ''
  return code >= 200 && code < 400 ? 'SUCCESS' : 'FAILURE'
}

function statusClass(l) {
  const s = (l.event_status || statusFromCode(l.response_code) || '').toUpperCase()
  if (s === 'SUCCESS') return 'bg-green-100 text-green-700'
  if (s === 'FAILURE') return 'bg-red-100 text-red-700'
  return 'bg-gray-100 text-gray-700'
}

async function load() {
  loading.value = true
  try {
    const r = await getMyActivities({ page: page.value, per_page: perPage.value })
    logs.value = r.data || []
    total.value = r.meta?.total || 0
  } finally { loading.value = false }
}

onMounted(load)
</script>
