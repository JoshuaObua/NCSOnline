<template>
  <section class="sitemap-panel">
    <div class="cms-panel-head"><h2>Sitemap</h2><button type="button" @click="loadAll">{{ loading ? 'Refreshing…' : 'Refresh' }}</button></div>
    <p class="cms-hint">Every public URL currently registered in the system — static routes wired into the app, and the live URL for every published piece of content.</p>

    <div class="cms-panel">
      <div class="cms-panel-head"><h3>Static Routes <span class="sitemap-count">{{ staticRoutes.length }}</span></h3></div>
      <table class="sitemap-table">
        <thead><tr><th>Page</th><th>URL</th></tr></thead>
        <tbody>
          <tr v-for="r in staticRoutes" :key="r.path">
            <td>{{ r.title }}</td>
            <td><a :href="r.path" target="_blank" rel="noopener">{{ r.path }}</a></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-for="group in dynamicGroups" :key="group.key" class="cms-panel">
      <div class="cms-panel-head"><h3>{{ group.label }} <span class="sitemap-count">{{ group.items.length }}</span></h3></div>
      <table class="sitemap-table">
        <thead><tr><th>Title</th><th>URL</th></tr></thead>
        <tbody>
          <tr v-for="item in group.items" :key="item.url">
            <td>{{ item.title }}</td>
            <td><a :href="item.url" target="_blank" rel="noopener">{{ item.url }}</a></td>
          </tr>
          <tr v-if="!group.items.length"><td colspan="2" class="cms-empty">No published items yet.</td></tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import * as cms from '@/api/cms.js'

const emit = defineEmits(['message', 'error'])
const router = useRouter()
const loading = ref(false)

const staticRoutes = computed(() => {
  const seen = new Set()
  return router.getRoutes()
    .filter(r => r.name && typeof r.path === 'string' && !r.path.includes(':') && r.path !== '/login' && r.path !== '/cms')
    .filter(r => (seen.has(r.path) ? false : (seen.add(r.path), true)))
    .map(r => ({ path: r.path || '/', title: String(r.name).replace(/([a-z])([A-Z])/g, '$1 $2') }))
    .sort((a, b) => a.path.localeCompare(b.path))
})

const posts = ref([])
const pages = ref([])
const events = ref([])
const careers = ref([])
const facilities = ref([])

const dynamicGroups = computed(() => [
  { key: 'news', label: 'News Articles', items: posts.value.filter(p => p.status === 'published').map(p => ({ title: p.title, url: `/news/${p.slug}` })) },
  { key: 'pages', label: 'Static Pages', items: pages.value.filter(p => p.status === 'published').map(p => ({ title: p.title, url: `/pages/${p.slug}` })) },
  { key: 'events', label: 'Events', items: events.value.filter(e => e.status === 'published').map(e => ({ title: e.title, url: `/events/${e.slug}` })) },
  { key: 'careers', label: 'Career Posts', items: careers.value.filter(c => c.status === 'published').map(c => ({ title: c.title, url: `/careers/${c.id}` })) },
  { key: 'facilities', label: 'Facilities', items: facilities.value.filter(f => f.is_active).map(f => ({ title: f.name, url: `/facilities/${f.slug}` })) },
])

function listData(res) {
  const value = res?.data?.data
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.items)) return value.items
  return []
}

async function loadAll() {
  loading.value = true
  try {
    const results = await Promise.allSettled([
      cms.adminListPosts({ category: 'blog', per_page: 200 }),
      cms.adminListPosts({ category: 'page', per_page: 200 }),
      cms.adminListEvents({ per_page: 200 }),
      cms.adminListCareers({ per_page: 200 }),
      cms.adminListFacilities(),
    ])
    posts.value = listData(results[0].value)
    pages.value = listData(results[1].value)
    events.value = listData(results[2].value)
    careers.value = listData(results[3].value)
    facilities.value = listData(results[4].value)
  } catch (err) {
    emit('error', err.response?.data?.error?.message || err.message || 'Could not load sitemap data')
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.cms-hint { color: var(--text-muted, #6c757d); margin-top: -.5rem; }
.sitemap-count { font-weight: 400; color: #98a6ad; font-size: .8rem; }
.sitemap-table { width: 100%; border-collapse: collapse; font-size: .85rem; }
.sitemap-table th { text-align: left; padding: .5rem .25rem; border-bottom: 2px solid #eef0f7; color: #98a6ad; font-weight: 600; }
.sitemap-table td { padding: .5rem .25rem; border-bottom: 1px solid #f4f6f9; }
.sitemap-table a { color: #6777ef; }
</style>
