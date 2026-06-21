<template>
  <LayoutDefault title="System Operations">
    <div class="space-y-6">
      <section class="rounded-2xl border p-6" :class="status.maintenance_mode?'border-amber-300 bg-amber-50':'border-emerald-200 bg-emerald-50'">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
          <div><p class="text-xs font-bold uppercase tracking-widest" :class="status.maintenance_mode?'text-amber-700':'text-emerald-700'">System status</p><h2 class="text-2xl font-bold text-primary-800 mt-1">{{ status.maintenance_mode?'Maintenance mode active':'Platform operational' }}</h2><p class="text-sm text-gray-600 mt-2">{{ status.reason || 'Public and authenticated services are available.' }}</p></div>
          <button v-if="authStore.isSuperAdmin" @click="openMaintenance" class="rounded-xl px-5 py-3 font-bold text-white" :class="status.maintenance_mode?'bg-emerald-700':'bg-amber-600'">{{ status.maintenance_mode?'End maintenance':'Schedule maintenance' }}</button>
        </div>
      </section>
      <section><div class="flex items-center justify-between mb-3"><div><h2 class="font-bold text-primary-800">Live resources</h2><p class="text-xs text-gray-500">Updated every five seconds from the API container host.</p></div><span class="text-xs text-gray-400">{{ lastResourceUpdate }}</span></div><div class="grid sm:grid-cols-2 xl:grid-cols-4 gap-4"><article v-for="card in resourceCards" :key="card.label" class="bg-white border border-gray-100 rounded-2xl p-5 shadow-sm"><p class="text-xs uppercase tracking-wide text-gray-400 font-semibold">{{ card.label }}</p><p class="text-2xl font-bold text-primary-800 mt-2">{{ card.value }}</p><div class="h-1.5 rounded-full bg-gray-100 mt-4"><div class="h-full rounded-full bg-primary-600" :style="{width:Math.min(card.percent||0,100)+'%'}"></div></div></article></div></section>
      <section class="grid md:grid-cols-2 gap-4">
        <article class="bg-white border border-gray-100 rounded-2xl p-6"><h3 class="font-bold text-primary-800">Application cache</h3><p class="text-sm text-gray-500 mt-2">Invalidate memory-cached settings and force resolvers to reload.</p><p class="text-xs text-gray-400 mt-3">Generation {{ status.cache_generation || 0 }}</p><button v-if="authStore.isSuperAdmin" @click="flushCache" class="mt-4 rounded-xl bg-primary-700 text-white px-4 py-2.5 font-semibold">Flush cache</button></article>
        <article class="bg-white border border-red-100 rounded-2xl p-6"><h3 class="font-bold text-red-800">Emergency session revocation</h3><p class="text-sm text-gray-500 mt-2">Revoke every refresh token and reject all currently issued access tokens.</p><button v-if="authStore.isSuperAdmin" @click="revokeAll" class="mt-4 rounded-xl bg-red-700 text-white px-4 py-2.5 font-semibold">Force logout all users</button></article>
      </section>
    </div>
    <Modal :open="maintenanceModal" :title="status.maintenance_mode?'End maintenance mode':'Schedule maintenance'" @close="maintenanceModal=false">
      <div class="space-y-4"><label class="block text-sm font-semibold">Reason<textarea v-model="maintenanceReason" rows="3" class="mt-1 w-full rounded-xl border-gray-200" placeholder="Describe the approved maintenance work"></textarea></label><label v-if="!status.maintenance_mode" class="block text-sm font-semibold">Expected completion<input v-model="maintenanceEnd" type="datetime-local" class="mt-1 w-full rounded-xl border-gray-200"></label></div>
      <template #footer><button @click="maintenanceModal=false" class="px-4 py-2">Cancel</button><button @click="saveMaintenance" :disabled="saving" class="rounded-xl bg-primary-700 px-5 py-2 text-white font-bold">Confirm</button></template>
    </Modal>
  </LayoutDefault>
</template>
<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'; import LayoutDefault from '@/components/layout/LayoutDefault.vue'; import Modal from '@/components/ui/Modal.vue'; import apiClient from '@/api/client.js'; import { useAuthStore } from '@/stores/auth.js'
const authStore=useAuthStore(),status=ref({}),resources=ref({}),maintenanceModal=ref(false),maintenanceReason=ref(''),maintenanceEnd=ref(''),saving=ref(false),lastResourceUpdate=ref('Not loaded');let timer
const pct=v=>Number.isFinite(v)?v:0,fmtBytes=v=>{if(!v)return '0 B';const units=['B','KB','MB','GB','TB'];let i=0;while(v>=1024&&i<units.length-1){v/=1024;i++}return `${v.toFixed(1)} ${units[i]}`}
const resourceCards=computed(()=>[{label:'CPU',value:`${pct(resources.value.cpu_percent?.[0]).toFixed(1)}%`,percent:pct(resources.value.cpu_percent?.[0])},{label:'Memory',value:`${pct(resources.value.memory?.usedPercent).toFixed(1)}%`,percent:pct(resources.value.memory?.usedPercent)},{label:'Disk',value:`${pct(resources.value.disk?.usedPercent).toFixed(1)}%`,percent:pct(resources.value.disk?.usedPercent)},{label:'Heap',value:fmtBytes(resources.value.heap_alloc_bytes),percent:0}])
async function refresh(){const [s,r]=await Promise.all([apiClient.get('/api/v1/admin/system/status'),apiClient.get('/api/v1/admin/system/resources')]);status.value=s.data.data||{};resources.value=r.data.data||{};lastResourceUpdate.value=new Date().toLocaleTimeString()}
function openMaintenance(){maintenanceReason.value=status.value.maintenance_mode?status.value.reason||'Maintenance completed':'';maintenanceEnd.value='';maintenanceModal.value=true}
async function saveMaintenance(){saving.value=true;try{await apiClient.put('/api/v1/admin/system/maintenance',{enabled:!status.value.maintenance_mode,reason:maintenanceReason.value,expected_end:maintenanceEnd.value?new Date(maintenanceEnd.value).toISOString():null});maintenanceModal.value=false;await refresh()}finally{saving.value=false}}
async function flushCache(){if(!confirm('Invalidate all application caches?'))return;await apiClient.post('/api/v1/admin/system/cache/flush');await refresh()}
async function revokeAll(){if(prompt('Type REVOKE ALL to confirm')!=='REVOKE ALL')return;await apiClient.post('/api/v1/admin/system/sessions/revoke-all');await authStore.logout();location.assign('/login')}
onMounted(()=>{refresh();timer=setInterval(refresh,5000)});onBeforeUnmount(()=>clearInterval(timer))
</script>
