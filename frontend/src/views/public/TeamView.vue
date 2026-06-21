<template>
  <div class="bg-white min-h-screen">
    <section class="bg-cream py-16 px-4 text-center"><span class="section-tag">Leadership</span><h1 class="text-4xl md:text-5xl font-bold text-darken mt-3">Current NCS Membership</h1><p class="mt-4 text-gray-500">The people serving Uganda's sports community.</p></section>
    <section class="max-w-7xl mx-auto px-4 py-16"><div v-if="loading" class="text-center text-gray-400">Loading team…</div><div v-else class="grid sm:grid-cols-2 lg:grid-cols-4 gap-6"><article v-for="member in team" :key="member.id" class="rounded-2xl overflow-hidden border border-gray-100 shadow-sm"><img v-if="member.image_url" :src="mediaUrl(member.image_url)" :alt="member.full_name" class="w-full aspect-square object-cover"/><div v-else class="aspect-square bg-gray-100 flex items-center justify-center"><i class="icofont-user text-5xl text-gray-300"></i></div><div class="p-5"><h2 class="font-bold text-darken">{{ member.full_name }}</h2><p class="text-accent text-sm font-semibold mt-1">{{ member.designation }}</p><p v-if="member.bio" class="text-gray-500 text-sm mt-3">{{ member.bio }}</p></div></article><p v-if="!team.length" class="col-span-full text-center text-gray-400">Team profiles will appear here when published.</p></div></section>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { listTeam } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
const team=ref([]),loading=ref(true)
onMounted(async()=>{try{const r=await listTeam();team.value=r.data?.data||[]}catch{team.value=[]}finally{loading.value=false}})
</script>
