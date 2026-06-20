<template>
  <LayoutDefault title="Federation Reports">
    <div class="space-y-6">
      <section class="admin-card">
        <div class="admin-card-body flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div><p class="text-xs font-semibold uppercase tracking-wider text-primary-700">NSMIS reporting workspace</p><h1 class="mt-1 text-xl font-bold text-gray-900">Reporting obligations</h1><p class="mt-1 text-sm text-gray-500">Complete, submit, approve and monitor reports within your assigned federation scope.</p></div>
          <label class="block min-w-64 text-sm font-medium text-gray-700">Reporting period<select v-model="periodID" class="mt-1 w-full rounded-lg border-gray-300" @change="loadObligations"><option value="">All periods</option><option v-for="p in periods" :key="p.id" :value="p.id">{{ p.name }}</option></select></label>
        </div>
      </section>

      <div v-if="message" role="status" class="rounded-lg border border-green-200 bg-green-50 p-4 text-sm text-green-700">{{ message }}</div>
      <div v-if="error" role="alert" class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700">{{ error }}</div>

      <section class="admin-card overflow-hidden">
        <div class="admin-card-header"><h2 class="text-sm font-semibold text-gray-900">Assigned reports</h2><button class="text-xs font-semibold text-primary-700" @click="loadObligations">Refresh</button></div>
        <div v-if="loading" class="p-10 text-center text-sm text-gray-500">Loading reporting obligations…</div>
        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead><tr><th class="table-th">Federation</th><th class="table-th">Report</th><th class="table-th">Period</th><th class="table-th">Due</th><th class="table-th">Status</th><th class="table-th"><span class="sr-only">Actions</span></th></tr></thead>
            <tbody v-if="obligations.length" class="divide-y divide-gray-100"><tr v-for="item in obligations" :key="item.id" class="hover:bg-gray-50"><td class="table-td font-medium text-gray-900">{{ item.federation_name }}</td><td class="table-td">{{ title(item.report_type) }}</td><td class="table-td">{{ item.period_name }}</td><td class="table-td"><span :class="item.days_overdue ? 'font-semibold text-red-700' : ''">{{ date(item.due_on) }}</span><span v-if="item.days_overdue" class="block text-xs text-red-600">{{ item.days_overdue }} days overdue</span></td><td class="table-td"><span :class="badge(item.report_status || item.status)">{{ title(item.report_status || item.status) }}</span></td><td class="table-td text-right"><button v-if="item.report_type === 'GOVERNANCE' && canEdit" class="font-semibold text-primary-700 hover:underline" @click="edit(item)">{{ item.report_id ? 'Continue' : 'Start' }}</button><button v-else-if="item.report_id && canAct(item)" class="font-semibold text-primary-700 hover:underline" @click="transitionItem(item)">{{ actionLabel(item) }}</button></td></tr></tbody>
            <tbody v-else><tr><td colspan="6" class="p-12 text-center text-sm text-gray-500">No reporting obligations are assigned for this selection.</td></tr></tbody>
          </table>
        </div>
      </section>

      <section v-if="selected" class="admin-card" aria-labelledby="governance-form-title">
        <div class="admin-card-header"><div><h2 id="governance-form-title" class="font-semibold text-gray-900">Governance report · {{ selected.federation_name }}</h2><p class="text-xs text-gray-500">{{ selected.period_name }}</p></div><button class="text-gray-500 hover:text-gray-900" aria-label="Close report form" @click="selected=null">✕</button></div>
        <form class="admin-card-body space-y-5" @submit.prevent="save">
          <div v-for="field in yesNoFields" :key="field.key" class="grid gap-3 border-b border-gray-100 pb-4 md:grid-cols-2 md:items-end">
            <label class="text-sm font-medium text-gray-800">{{ field.label }}<select v-model="form[field.key]" class="mt-1 block w-full rounded-lg border-gray-300"><option :value="false">No</option><option :value="true">Yes</option></select></label>
            <label v-if="form[field.key]" class="text-sm font-medium text-gray-800">Date<input v-model="form[field.date]" type="date" required class="mt-1 block w-full rounded-lg border-gray-300" /></label>
          </div>
          <label class="block text-sm font-medium text-gray-800">Disciplinary cases handled<input v-model.number="form.disciplinary_cases_handled" type="number" min="0" required class="mt-1 block w-full rounded-lg border-gray-300 md:w-64" /></label>
          <div class="flex flex-wrap justify-end gap-3"><button type="button" class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold" @click="selected=null">Cancel</button><button type="submit" :disabled="saving" class="rounded-lg bg-primary-700 px-4 py-2 text-sm font-semibold text-white disabled:opacity-50">{{ saving ? 'Saving…' : 'Save draft' }}</button><button v-if="selected.report_id && (selected.report_status === 'DRAFT' || !selected.report_status)" type="button" :disabled="saving" class="rounded-lg bg-green-700 px-4 py-2 text-sm font-semibold text-white disabled:opacity-50" @click="submitReport">Submit report</button></div>
        </form>
      </section>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import apiClient from '@/api/client.js'
import { useAuthStore } from '@/stores/auth.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'

const auth=useAuthStore(), breadcrumb=useBreadcrumbStore(), loading=ref(true), saving=ref(false), error=ref(''), message=ref(''), periods=ref([]), obligations=ref([]), periodID=ref(''), selected=ref(null)
const yesNoFields=[{key:'executive_meeting_held',date:'executive_meeting_date',label:'Executive meeting held'},{key:'agm_conducted',date:'agm_date',label:'AGM conducted'},{key:'board_meeting_held',date:'board_meeting_date',label:'Board meeting held'},{key:'elections_conducted',date:'election_date',label:'Elections conducted'}]
const form=reactive({version:0,executive_meeting_held:false,executive_meeting_date:'',agm_conducted:false,agm_date:'',board_meeting_held:false,board_meeting_date:'',elections_conducted:false,election_date:'',disciplinary_cases_handled:0})
const canEdit=computed(()=>auth.hasAnyRole('super_admin','admin','federation_general_secretary'))
function title(v=''){return v.toLowerCase().replaceAll('_',' ').replace(/\b\w/g,c=>c.toUpperCase())}
function date(v){return new Intl.DateTimeFormat('en-UG',{dateStyle:'medium',timeZone:'Africa/Kampala'}).format(new Date(v))}
function badge(v){const s=(v||'').toUpperCase();return ['inline-flex rounded-full px-2.5 py-1 text-xs font-semibold',s==='LOCKED'||s==='APPROVED'||s==='ACCEPTED'?'bg-green-100 text-green-700':s==='OVERDUE'||s==='NEEDS_CORRECTION'?'bg-red-100 text-red-700':s==='DRAFT'||s==='NOT_STARTED'?'bg-gray-100 text-gray-700':'bg-blue-100 text-blue-700']}
function edit(item){selected.value=item;Object.assign(form,{version:item.report_version||0,executive_meeting_held:false,executive_meeting_date:'',agm_conducted:false,agm_date:'',board_meeting_held:false,board_meeting_date:'',elections_conducted:false,election_date:'',disciplinary_cases_handled:0});message.value='';error.value=''}
function canAct(item){return (item.report_status==='SUBMITTED'&&auth.hasAnyRole('federation_president'))||(['PRESIDENT_APPROVED','NCS_UNDER_REVIEW','APPROVED'].includes(item.report_status)&&auth.hasAnyRole('super_admin','admin','technical_department'))}
function actionLabel(item){return item.report_status==='SUBMITTED'?'Approve':item.report_status==='APPROVED'?'Lock':'Review'}
async function transitionItem(item){const to=item.report_status==='SUBMITTED'?'PRESIDENT_APPROVED':item.report_status==='PRESIDENT_APPROVED'?'NCS_UNDER_REVIEW':item.report_status==='NCS_UNDER_REVIEW'?'APPROVED':'LOCKED';await transition(item.report_id,to)}
async function transition(id,to){saving.value=true;error.value='';try{await apiClient.post(`/api/v1/nsmis/reports/${id}/transitions`,{to_status:to,reason:''});message.value=`Report moved to ${title(to)}.`;selected.value=null;await loadObligations()}catch(e){error.value=e.response?.data?.error?.message||'The report status could not be changed.'}finally{saving.value=false}}
async function save(){saving.value=true;error.value='';try{const res=await apiClient.put(`/api/v1/nsmis/report-obligations/${selected.value.id}/governance-draft`,{...form});const data=res.data?.data;selected.value={...selected.value,report_id:data.id,report_version:data.version,report_status:'DRAFT'};form.version=data.version;message.value='Governance draft saved.';await loadObligations()}catch(e){error.value=e.response?.data?.error?.message||'The draft could not be saved.'}finally{saving.value=false}}
async function submitReport(){if(!selected.value?.report_id)return;await transition(selected.value.report_id,'SUBMITTED')}
async function loadObligations(){loading.value=true;error.value='';try{const params=periodID.value?{period_id:periodID.value}:{};const res=await apiClient.get('/api/v1/nsmis/report-obligations',{params});obligations.value=res.data?.data||[]}catch(e){error.value=e.response?.data?.error?.message||'Reporting obligations could not be loaded.'}finally{loading.value=false}}
onMounted(async()=>{breadcrumb.set('Federation Reports');try{const res=await apiClient.get('/api/v1/nsmis/reporting-periods',{params:{open:false}});periods.value=res.data?.data||[]}catch{}await loadObligations()})
</script>
