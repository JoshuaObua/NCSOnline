<template>
  <LayoutDefault :title="`${title} Dashboard`">
    <div class="space-y-6">
      <section class="rounded-xl bg-gradient-to-r from-primary-700 to-primary-500 p-6 text-white shadow-lg"><p class="text-xs font-semibold uppercase tracking-widest text-primary-100">NSMIS official analytics</p><h1 class="mt-1 text-2xl font-bold">{{ title }} dashboard</h1><p class="mt-1 text-sm text-primary-100">Approved records only · Updated {{ asOf }}</p></section>
      <div v-if="error" role="alert" class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700">{{ error }}</div>
      <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-3"><div v-for="i in 3" :key="i" class="h-28 animate-pulse rounded-xl bg-gray-100"></div></div>
      <template v-else>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3"><article v-for="card in cards" :key="card.label" class="stat-card"><div class="flex h-11 w-11 items-center justify-center rounded-xl bg-primary-50"><i :class="[card.icon,'text-xl text-primary-700']"></i></div><div><p class="text-2xl font-bold text-gray-900">{{ format(card.value,card.money) }}</p><p class="text-xs text-gray-500">{{ card.label }}</p></div></article></div>
        <section v-if="dashboard==='athletes'" class="grid gap-5 xl:grid-cols-2"><DataPanel title="Gender distribution" :rows="objectRows(data.gender,'Category')"/><DataPanel title="Athletes by region" :rows="namedRows(data.by_region,'region')"/></section>
        <section v-if="dashboard==='performance'" class="grid gap-5 xl:grid-cols-2"><DataPanel title="Medal table by federation" :rows="medalRows"/><DataPanel title="Medals by host country" :rows="namedRows(data.by_country,'country')"/></section>
        <section v-if="dashboard==='finance'" class="admin-card"><div class="admin-card-body text-sm text-gray-600">Amounts are sourced from the NCS disbursement ledger and accepted financial accountability records. Outstanding amounts include released or partially accounted funds past their due date.</div></section>
        <section v-if="dashboard==='talent'" class="admin-card"><div class="admin-card-body text-sm text-gray-600">Talent progression is counted from governed identification records, national-team progression states, and currently active scholarship records.</div></section>
      </template>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { computed, defineComponent, h, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import apiClient from '@/api/client.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'

const route=useRoute(),breadcrumb=useBreadcrumbStore(),loading=ref(true),error=ref(''),data=ref({})
const dashboard=computed(()=>route.params.dashboard),title=computed(()=>({athletes:'Athlete',performance:'Performance',finance:'Financial',talent:'Talent'}[dashboard.value]||'NSMIS'))
const asOf=computed(()=>data.value.as_of?new Intl.DateTimeFormat('en-UG',{dateStyle:'medium',timeStyle:'short',timeZone:'Africa/Kampala'}).format(new Date(data.value.as_of)):'—')
const cards=computed(()=>({athletes:[{label:'Registered athletes',value:data.value.total_registered_athletes,icon:'icofont-runner-alt-1'},{label:'Female athletes',value:data.value.gender?.FEMALE||0,icon:'icofont-user-female'},{label:'National scope',value:(data.value.by_federation||[]).length,icon:'icofont-building-alt'}],performance:[{label:'Total medals',value:data.value.total_medals,icon:'icofont-medal'},{label:'Gold medals',value:data.value.medals?.GOLD||0,icon:'icofont-trophy'},{label:'International representation',value:data.value.international_representation,icon:'icofont-world'}],finance:[{label:'Funds released',value:data.value.funds_released,icon:'icofont-money',money:true},{label:'Accepted accountabilities',value:data.value.accountability_submitted,icon:'icofont-check-circled'},{label:'Outstanding',value:data.value.outstanding_accountabilities,icon:'icofont-warning-alt',money:true}],talent:[{label:'Athletes identified',value:data.value.athletes_identified,icon:'icofont-search-user'},{label:'Progressed to national teams',value:data.value.progressed_to_national_teams,icon:'icofont-flag-alt-2'},{label:'Active scholarships',value:data.value.receiving_scholarships,icon:'icofont-graduate-alt'}]}[dashboard.value]||[]))
const medalRows=computed(()=>(data.value.by_federation||[]).map(x=>({Name:x.federation_name,Gold:x.gold,Silver:x.silver,Bronze:x.bronze})))
function format(v,money){if(v==null)return '0';return money?new Intl.NumberFormat('en-UG',{style:'currency',currency:'UGX',maximumFractionDigits:0}).format(v):new Intl.NumberFormat('en-UG').format(v)}
function objectRows(obj,label){return Object.entries(obj||{}).map(([k,v])=>({[label]:k.replaceAll('_',' '),Count:v}))}
function namedRows(items,key){return (items||[]).map(x=>({Name:x[key]?.replaceAll?.('_',' ')||x[key],Count:x.count}))}
async function load(){loading.value=true;error.value='';breadcrumb.set(`${title.value} Dashboard`);try{const res=await apiClient.get(`/api/v1/nsmis/dashboards/${dashboard.value}`);data.value=res.data?.data||{}}catch(e){error.value=e.response?.data?.error?.message||'This dashboard could not be loaded.'}finally{loading.value=false}}
watch(()=>route.params.dashboard,load);onMounted(load)

const DataPanel=defineComponent({props:{title:String,rows:Array},setup(props){return()=>h('article',{class:'admin-card overflow-hidden'},[h('div',{class:'admin-card-header'},h('h2',{class:'text-sm font-semibold text-gray-900'},props.title)),props.rows?.length?h('div',{class:'overflow-x-auto'},h('table',{class:'w-full'},[h('thead',h('tr',Object.keys(props.rows[0]).map(k=>h('th',{class:'table-th'},k)))),h('tbody',{class:'divide-y divide-gray-100'},props.rows.map((row,i)=>h('tr',{key:i},Object.values(row).map(v=>h('td',{class:'table-td'},String(v))))))])):h('p',{class:'p-8 text-center text-sm text-gray-500'},'No approved data available.')])}})
</script>
