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
          <small>{{ departmentLabel(form.department_name) }} · {{ stepCount(form) }} step{{ stepCount(form) === 1 ? '' : 's' }}</small>
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
import { buildSectionSteps } from '@/utils/formBuilder.js'

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
const stepCount = form => buildSectionSteps(form).length
const formatMoney = value => new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0))
const departmentLabel = value => /(^|_)(super_?admin|admin|user)($|_)/i.test(String(value || '')) ? 'NCS Service' : (value || 'NCS Service')
</script>

<style scoped>
.open-forms-section{margin-top:8px}.open-forms-section>header{display:flex;align-items:end;justify-content:space-between;margin-bottom:12px}.open-forms-section>header p{margin:0;color:#7e8598;font-size:10px;font-weight:800;text-transform:uppercase}.open-forms-section>header h2{margin:3px 0;color:#303345;font-size:18px}.open-forms-section>header button{border:0;background:transparent;color:#5266d8;font-size:12px;font-weight:800}.forms-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:14px}.forms-grid article{display:grid;grid-template-columns:78px 1fr;min-height:170px;overflow:hidden;border:1px solid #e3e6ed;border-radius:6px;background:#fff}.forms-grid article>img,.form-mark{width:78px;height:100%;object-fit:cover}.form-mark{display:grid;place-items:center;background:#10233f;color:#f5a623;font-size:29px}.forms-grid article>div{display:flex;flex-direction:column;min-width:0;padding:15px}.forms-grid small{color:#d17b00;font-size:9px;font-weight:800;text-transform:uppercase}.forms-grid h3{margin:4px 0;color:#303345;font-size:14px;line-height:1.35}.forms-grid p{display:-webkit-box;overflow:hidden;margin:0;color:#7b8295;font-size:11px;line-height:1.5;-webkit-line-clamp:2;-webkit-box-orient:vertical}.forms-grid footer{display:flex;align-items:center;justify-content:space-between;gap:8px;margin-top:auto;padding-top:10px}.forms-grid footer strong{color:#303345;font-size:11px}.forms-grid footer button{border:0;background:transparent;color:#5266d8;font-size:11px;font-weight:800}.forms-empty{padding:48px;border:1px dashed #ccd2de;border-radius:6px;background:#fff;text-align:center;color:#7c8497}@media(max-width:1120px){.forms-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:600px){.forms-grid{grid-template-columns:1fr}.forms-grid article{grid-template-columns:68px 1fr}.forms-grid article>img,.form-mark{width:68px}}
.open-forms-section{margin-top:8px!important}.open-forms-section>header p{color:#6777ef!important}.open-forms-section>header h2{color:#34395e!important;font-weight:700!important}.open-forms-section>header button{color:#6777ef!important}.forms-grid{gap:18px!important}.forms-grid article{min-height:176px!important;border:0!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.forms-grid article>img,.form-mark{width:78px!important}.form-mark{background:#6777ef!important;color:#fff!important}.forms-grid small{color:#98a6ad!important}.forms-grid h3{color:#34395e!important;font-weight:700!important}.forms-grid p{color:#6c757d!important}.forms-grid footer strong{color:#34395e!important}.forms-grid footer button{color:#6777ef!important}.forms-empty{border-color:#e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#98a6ad!important}:global(.dark) .forms-grid article{background:#1f2937!important;box-shadow:0 4px 25px rgba(0,0,0,.28)!important}:global(.dark) .open-forms-section>header h2,:global(.dark) .forms-grid h3,:global(.dark) .forms-grid footer strong{color:#f8fafc!important}:global(.dark) .forms-grid p,:global(.dark) .forms-empty{color:#cbd5e1!important}:global(.dark) .forms-empty{background:#111827!important;border-color:#334155!important}
.open-forms-section,.forms-grid,.forms-grid article,.forms-grid article>div,.forms-grid footer{min-width:0!important}.open-forms-section>header{align-items:flex-start!important;gap:12px!important;flex-wrap:wrap!important}.open-forms-section>header>div{min-width:0!important}.open-forms-section>header button{flex:0 0 auto!important;white-space:nowrap!important}.forms-grid{grid-template-columns:repeat(2,minmax(0,1fr))!important;align-items:stretch!important}.forms-grid article{grid-template-columns:minmax(76px,96px) minmax(0,1fr)!important;min-height:190px!important;height:100%!important}.forms-grid article>img,.form-mark{width:100%!important;min-width:0!important}.forms-grid small,.forms-grid h3,.forms-grid p,.forms-grid footer strong,.forms-grid footer button{overflow-wrap:anywhere!important}.forms-grid footer{align-items:flex-start!important;flex-wrap:wrap!important}.forms-grid footer strong{flex:1 1 110px!important;line-height:1.35!important}.forms-grid footer button{flex:0 1 auto!important;line-height:1.35!important;text-align:left!important}@media(max-width:860px){.forms-grid{grid-template-columns:1fr!important}.forms-grid article{grid-template-columns:minmax(76px,96px) minmax(0,1fr)!important}}@media(max-width:520px){.open-forms-section>header{display:grid!important}.open-forms-section>header button{justify-self:start!important}.forms-grid article{grid-template-columns:1fr!important}.forms-grid article>img,.form-mark{height:120px!important}.forms-grid footer button{width:100%!important}}
.forms-grid{gap:12px!important}.forms-grid article{grid-template-columns:64px minmax(0,1fr)!important;min-height:132px!important;height:auto!important}.forms-grid article>img,.form-mark{width:64px!important;height:100%!important}.form-mark{font-size:22px!important}.forms-grid article>div{padding:10px 12px!important;gap:3px!important}.forms-grid small{display:-webkit-box!important;overflow:hidden!important;font-size:8.5px!important;line-height:1.25!important;-webkit-line-clamp:1!important;-webkit-box-orient:vertical!important}.forms-grid h3{display:-webkit-box!important;overflow:hidden!important;margin:2px 0!important;font-size:12.5px!important;line-height:1.25!important;-webkit-line-clamp:2!important;-webkit-box-orient:vertical!important}.forms-grid p{font-size:10.5px!important;line-height:1.35!important;-webkit-line-clamp:2!important}.forms-grid footer{align-items:center!important;gap:6px!important;padding-top:6px!important}.forms-grid footer strong{flex:1 1 84px!important;font-size:10.5px!important}.forms-grid footer button{font-size:10.5px!important;white-space:normal!important}@media(max-width:760px){.forms-grid{grid-template-columns:1fr!important}}@media(max-width:520px){.forms-grid article{grid-template-columns:58px minmax(0,1fr)!important}.forms-grid article>img,.form-mark{width:58px!important;height:100%!important}.forms-grid footer button{width:auto!important}}
.open-forms-section{width:85%!important;max-width:1224px!important;margin-left:auto!important;margin-right:auto!important}.forms-grid article{max-width:100%!important}@media(max-width:991px){.open-forms-section{width:100%!important;max-width:none!important}}
.open-forms-section{margin-left:0!important;margin-right:auto!important}
.open-forms-section{width:100%!important;max-width:none!important}

/* Final responsive layout layer for the applicant dashboard. */
.open-forms-section{min-width:0!important;width:100%!important;max-width:100%!important;margin:0!important;overflow:hidden!important}
.open-forms-section>header{display:flex!important;align-items:flex-end!important;justify-content:space-between!important;gap:12px!important;min-width:0!important;margin-bottom:12px!important}
.open-forms-section>header>div{min-width:0!important}
.open-forms-section>header h2{font-size:16px!important;line-height:1.25!important;overflow-wrap:anywhere!important}
.open-forms-section>header button{flex:0 0 auto!important;max-width:100%!important;white-space:normal!important}
.forms-grid{display:grid!important;grid-template-columns:repeat(2,minmax(0,1fr))!important;gap:12px!important;width:100%!important;min-width:0!important}
.forms-grid article{display:grid!important;grid-template-columns:64px minmax(0,1fr)!important;min-width:0!important;min-height:134px!important;overflow:hidden!important}
.forms-grid article>img,.form-mark{width:64px!important;min-width:64px!important;height:100%!important;object-fit:cover!important}
.forms-grid article>div{min-width:0!important;padding:11px 12px!important}
.forms-grid small,.forms-grid h3,.forms-grid p,.forms-grid footer strong,.forms-grid footer button{min-width:0!important;overflow-wrap:anywhere!important;word-break:normal!important}
.forms-grid h3{font-size:12.5px!important;line-height:1.25!important}
.forms-grid footer{display:flex!important;align-items:center!important;justify-content:space-between!important;gap:6px!important;min-width:0!important;flex-wrap:wrap!important}
.forms-grid footer strong{flex:1 1 78px!important}
.forms-grid footer button{flex:0 1 auto!important;max-width:100%!important;text-align:left!important}
.forms-empty{width:100%!important;max-width:100%!important;padding:28px 14px!important;overflow-wrap:anywhere!important}
@media(max-width:760px){.forms-grid{grid-template-columns:1fr!important}.forms-grid article{grid-template-columns:58px minmax(0,1fr)!important}.forms-grid article>img,.form-mark{width:58px!important;min-width:58px!important}}
@media(max-width:420px){.open-forms-section>header{align-items:flex-start!important;flex-direction:column!important}.forms-grid article{grid-template-columns:52px minmax(0,1fr)!important}.forms-grid article>img,.form-mark{width:52px!important;min-width:52px!important}.forms-grid footer{align-items:flex-start!important;flex-direction:column!important}.forms-grid footer button{width:100%!important}}
</style>
