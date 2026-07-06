<template>
  <div class="bg-white min-h-screen">
    <section class="bg-cream py-16 px-4 text-center"><span class="section-tag">Leadership</span><h1 class="text-4xl md:text-5xl font-bold text-darken mt-3">Governing Council</h1><p class="mt-4 text-gray-500">The council overseeing Uganda's sports development.</p></section>
    <section class="max-w-7xl mx-auto px-4 py-16"><div v-if="loading" class="text-center text-gray-400">Loading council…</div><div v-else class="grid sm:grid-cols-2 lg:grid-cols-4 gap-6"><article v-for="member in council" :key="member.id" class="rounded-2xl overflow-hidden border border-gray-100 shadow-sm"><img v-if="member.image_url" :src="mediaUrl(member.image_url)" :alt="member.full_name" class="w-full aspect-square object-cover"/><div v-else class="aspect-square bg-gray-100 flex items-center justify-center"><i class="icofont-user text-5xl text-gray-300"></i></div><div class="p-5"><h2 class="font-bold text-darken">{{ member.full_name }}</h2><p class="text-accent text-sm font-semibold mt-1">{{ member.designation }}</p><div v-if="member.bio"><div v-if="expandedIds.has(member.id)" class="text-gray-500 text-sm mt-3 prose-bio" v-html="sanitizeRichHtml(member.bio)"></div><button type="button" class="text-accent text-sm font-semibold mt-3 hover:underline" @click="toggleBio(member.id)">{{ expandedIds.has(member.id) ? 'View less' : 'View more' }}</button></div></div></article><p v-if="!council.length" class="col-span-full text-center text-gray-400">Council profiles will appear here when published.</p></div></section>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { listCouncil } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import { sanitizeRichHtml } from '@/utils/sanitize.js'
const council=ref([]),loading=ref(true)
const expandedIds=ref(new Set())
function toggleBio(id){
  const next=new Set(expandedIds.value)
  next.has(id)?next.delete(id):next.add(id)
  expandedIds.value=next
}
onMounted(async()=>{try{const r=await listCouncil();council.value=r.data?.data||[]}catch{council.value=[]}finally{loading.value=false}})
</script>
<style scoped>
.prose-bio :deep(p) { margin: 0 0 0.5rem; }
.prose-bio :deep(p:last-child) { margin-bottom: 0; }
</style>
