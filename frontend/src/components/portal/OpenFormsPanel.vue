<template>
  <section class="open-forms-section">
    <header>
      <div><p>Available Services</p><h2>{{ title || 'Applications you can start' }}</h2></div>
      <button v-if="!expanded && forms.length > 6" type="button" @click="$emit('view-all')">View all <i class="icofont-rounded-right"></i></button>
    </header>
    <div v-if="loading" class="forms-empty">Loading open forms...</div>
    <div v-else-if="!visibleForms.length" class="forms-empty">There are no open application forms right now.</div>
    <div v-else class="forms-grid">
      <article v-for="form in visibleForms" :key="form.id">
        <img v-if="form.banner_image_url" :src="mediaUrl(form.banner_image_url)" alt="" />
        <span v-else class="form-mark"><i class="icofont-file-document"></i></span>
        <div>
          <small>{{ departmentLabel(form.department_name) }}</small>
          <h3>{{ form.title }}</h3>
          <p>{{ form.description || 'Complete this application online.' }}</p>
          <footer>
            <strong>{{ Number(form.price_ugx) ? `UGX ${formatMoney(form.price_ugx)}` : 'Free' }}</strong>
            <button type="button" @click="$emit('start', form)">{{ draftFor(form.id) ? 'Continue' : 'Start application' }} <i class="icofont-rounded-right"></i></button>
          </footer>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  title: { type: String, default: '' },
  forms: { type: Array, default: () => [] },
  submissions: { type: Array, default: () => [] },
  loading: Boolean,
  expanded: Boolean,
})
defineEmits(['start', 'view-all'])

const visibleForms = computed(() => props.expanded ? props.forms : props.forms.slice(0, 6))
const draftFor = id => props.submissions.find(item => item.template_id === id && item.status === 'DRAFT')
const formatMoney = value => new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0))
const departmentLabel = value => /(^|_)(super_?admin|admin|user)($|_)/i.test(String(value || '')) ? 'NCS Service' : (value || 'NCS Service')
</script>

<style scoped>
.open-forms-section{margin-top:8px}.open-forms-section>header{display:flex;align-items:end;justify-content:space-between;margin-bottom:12px}.open-forms-section>header p{margin:0;color:#7e8598;font-size:10px;font-weight:800;text-transform:uppercase}.open-forms-section>header h2{margin:3px 0;color:#303345;font-size:18px}.open-forms-section>header button{border:0;background:transparent;color:#5266d8;font-size:12px;font-weight:800}.forms-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:14px}.forms-grid article{display:grid;grid-template-columns:78px 1fr;min-height:170px;overflow:hidden;border:1px solid #e3e6ed;border-radius:6px;background:#fff}.forms-grid article>img,.form-mark{width:78px;height:100%;object-fit:cover}.form-mark{display:grid;place-items:center;background:#10233f;color:#f5a623;font-size:29px}.forms-grid article>div{display:flex;flex-direction:column;min-width:0;padding:15px}.forms-grid small{color:#d17b00;font-size:9px;font-weight:800;text-transform:uppercase}.forms-grid h3{margin:4px 0;color:#303345;font-size:14px;line-height:1.35}.forms-grid p{display:-webkit-box;overflow:hidden;margin:0;color:#7b8295;font-size:11px;line-height:1.5;-webkit-line-clamp:2;-webkit-box-orient:vertical}.forms-grid footer{display:flex;align-items:center;justify-content:space-between;gap:8px;margin-top:auto;padding-top:10px}.forms-grid footer strong{color:#303345;font-size:11px}.forms-grid footer button{border:0;background:transparent;color:#5266d8;font-size:11px;font-weight:800}.forms-empty{padding:48px;border:1px dashed #ccd2de;border-radius:6px;background:#fff;text-align:center;color:#7c8497}@media(max-width:1120px){.forms-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:600px){.forms-grid{grid-template-columns:1fr}.forms-grid article{grid-template-columns:68px 1fr}.forms-grid article>img,.form-mark{width:68px}}
.open-forms-section{margin-top:8px!important}.open-forms-section>header p{color:#6777ef!important}.open-forms-section>header h2{color:#34395e!important;font-weight:700!important}.open-forms-section>header button{color:#6777ef!important}.forms-grid{gap:18px!important}.forms-grid article{min-height:176px!important;border:0!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.forms-grid article>img,.form-mark{width:78px!important}.form-mark{background:#6777ef!important;color:#fff!important}.forms-grid small{color:#98a6ad!important}.forms-grid h3{color:#34395e!important;font-weight:700!important}.forms-grid p{color:#6c757d!important}.forms-grid footer strong{color:#34395e!important}.forms-grid footer button{color:#6777ef!important}.forms-empty{border-color:#e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#98a6ad!important}:global(.dark) .forms-grid article{background:#1f2937!important;box-shadow:0 4px 25px rgba(0,0,0,.28)!important}:global(.dark) .open-forms-section>header h2,:global(.dark) .forms-grid h3,:global(.dark) .forms-grid footer strong{color:#f8fafc!important}:global(.dark) .forms-grid p,:global(.dark) .forms-empty{color:#cbd5e1!important}:global(.dark) .forms-empty{background:#111827!important;border-color:#334155!important}
</style>
