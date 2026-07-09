<template>
  <AccountShell title="My Audit Activities" subtitle="Searchable timeline of authenticated account activity.">
    <article class="account-panel">
      <div class="activity-toolbar">
        <input v-model="search" type="search" placeholder="Search activities, IP, endpoint, browser" />
        <button type="button" @click="loadActivities">Refresh</button>
      </div>
      <div class="timeline">
        <article v-for="item in filteredActivities" :key="item.id" class="timeline-item">
          <div class="activity-icon">{{ iconFor(item) }}</div>
          <div>
            <div class="activity-head">
              <strong>{{ categoryFor(item) }}</strong>
              <time>{{ formatDate(item.created_at) }}</time>
            </div>
            <p>{{ messageFor(item) }}</p>
            <small>{{ item.ip_address || 'Unknown IP' }} · {{ item.browser || item.device_info || 'Unknown browser' }} · {{ item.endpoint || item.resource }}</small>
          </div>
        </article>
        <p v-if="!filteredActivities.length" class="empty-copy">No matching activities found.</p>
      </div>
    </article>
  </AccountShell>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import AccountShell from './AccountShell.vue'
import { listMyAuditLogs } from '@/api/account.js'

const activities = ref([])
const search = ref('')

const filteredActivities = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return activities.value
  return activities.value.filter(item => JSON.stringify(item).toLowerCase().includes(q))
})

async function loadActivities() {
  const res = await listMyAuditLogs({ per_page: 100 })
  activities.value = res.data?.data || []
}

function categoryFor(item) {
  const event = String(item.event_type || item.action || '').toUpperCase()
  if (event.includes('AUTH')) return 'Authentication'
  if (event.includes('SECURITY') || event.includes('PIN') || event.includes('PASSWORD')) return 'Security Change'
  if (event.includes('CMS')) return 'Page Builder Deployment'
  if (event.includes('UPDATE') || event.includes('CREATE') || event.includes('DELETE')) return 'Data Update'
  return 'Account Activity'
}

function iconFor(item) {
  const category = categoryFor(item)
  return { Authentication: 'A', 'Security Change': 'S', 'Page Builder Deployment': 'P', 'Data Update': 'D' }[category] || 'L'
}

function messageFor(item) {
  const status = item.event_status ? `${item.event_status.toLowerCase()} ` : ''
  return `${status}${item.method || ''} ${item.endpoint || item.action || item.resource || 'account action'}`.trim()
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString() : 'Unknown time'
}

onMounted(loadActivities)
</script>

<style scoped>
.activity-toolbar{display:flex;gap:.75rem;margin-bottom:1rem}.activity-toolbar input{flex:1;height:42px;border:1px solid #cbd5e1;border-radius:6px;padding:0 .8rem}.activity-toolbar button{border:0;border-radius:6px;background:#112b4e;color:white;font-weight:800;padding:0 1rem}.timeline{display:grid;gap:.8rem}.timeline-item{display:grid;grid-template-columns:42px minmax(0,1fr);gap:.85rem;border:1px solid #e5e7eb;border-radius:8px;padding:.9rem}.activity-icon{width:42px;height:42px;border-radius:50%;display:grid;place-items:center;background:#eef2ff;color:#3730a3;font-weight:900}.activity-head{display:flex;justify-content:space-between;gap:.75rem}.activity-head strong{color:#112b4e}.activity-head time{color:#94a3b8;font-size:.8rem}.timeline-item p{margin:.2rem 0;color:#334155}.timeline-item small{color:#64748b}@media(max-width:680px){.activity-toolbar,.activity-head{flex-direction:column}.activity-toolbar{display:flex}}
</style>
