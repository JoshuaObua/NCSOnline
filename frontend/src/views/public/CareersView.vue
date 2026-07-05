<template>
  <div class="careers-page bg-white">
    <section class="careers-hero">
      <div class="careers-hero__media" aria-hidden="true">
        <img src="/main-logo.png" alt="" />
      </div>
      <div class="careers-shell careers-hero__grid">
        <div class="careers-hero__copy">
          <span class="careers-kicker">{{ page.eyebrow }}</span>
          <h1>{{ page.title }}</h1>
          <p>{{ page.intro }}</p>
          <div class="careers-hero__actions">
            <a href="#open-roles" class="careers-button careers-button--primary">View Open Roles</a>
            <router-link to="/contact-us" class="careers-button careers-button--ghost">Contact HR</router-link>
          </div>
        </div>
        <aside class="careers-hero__panel" aria-label="Careers summary">
          <div>
            <span>Open roles</span>
            <strong>{{ careers.length }}</strong>
          </div>
          <div>
            <span>Departments hiring</span>
            <strong>{{ departmentCount }}</strong>
          </div>
          <div>
            <span>Work location</span>
            <strong>{{ primaryLocation }}</strong>
          </div>
        </aside>
      </div>
    </section>

    <section class="careers-stats">
      <div class="careers-shell careers-stats__grid">
        <article v-for="item in page.stats" :key="item.label">
          <span>{{ item.value }}</span>
          <strong>{{ item.label }}</strong>
          <p>{{ item.text }}</p>
        </article>
      </div>
    </section>

    <section id="open-roles" class="careers-shell careers-openings" aria-labelledby="open-roles-title">
      <div class="careers-section-head">
        <span class="careers-kicker">{{ page.jobs_eyebrow }}</span>
        <h2 id="open-roles-title">{{ page.jobs_title }}</h2>
        <p>{{ page.jobs_intro }}</p>
      </div>

      <div class="careers-toolbar">
        <label>
          <span>Search jobs</span>
          <input v-model.trim="search" type="search" placeholder="Search by role, department, or location" />
        </label>
        <label>
          <span>Department</span>
          <select v-model="selectedDepartment">
            <option value="">All departments</option>
            <option v-for="dept in departments" :key="dept" :value="dept">{{ dept }}</option>
          </select>
        </label>
        <label>
          <span>Job type</span>
          <select v-model="selectedType">
            <option value="">All types</option>
            <option v-for="type in jobTypes" :key="type" :value="type">{{ labelJobType(type) }}</option>
          </select>
        </label>
      </div>

      <div v-if="loading" class="careers-list">
        <article v-for="i in 4" :key="i" class="career-card career-card--loading">
          <span></span>
          <strong></strong>
          <p></p>
        </article>
      </div>

      <div v-else-if="filteredCareers.length" class="careers-list">
        <router-link v-for="career in filteredCareers" :key="career.id" :to="`/careers/${career.id}`" class="career-card">
          <div class="career-card__top">
            <span class="career-card__badge">{{ career.department || career.department_name || 'NCS Uganda' }}</span>
            <span v-if="career.deadline_at" class="career-card__deadline">Closes {{ formatDate(career.deadline_at) }}</span>
          </div>
          <h3>{{ career.title }}</h3>
          <p>{{ excerpt(career.description) }}</p>
          <div class="career-card__meta">
            <span>
              <svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 21s7-4.7 7-11a7 7 0 1 0-14 0c0 6.3 7 11 7 11Z"/><circle cx="12" cy="10" r="2.5"/></svg>
              {{ career.location || 'Kampala, Uganda' }}
            </span>
            <span>
              <svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><path d="M3 8h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8Z"/><path d="M3 13h18"/></svg>
              {{ labelJobType(career.job_type) }}
            </span>
            <span v-if="career.salary_range">
              <svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 7h16v10H4z"/><circle cx="12" cy="12" r="2.5"/><path d="M7 7v10M17 7v10"/></svg>
              {{ career.salary_range }}
            </span>
          </div>
          <span class="career-card__cta">View Details</span>
        </router-link>
      </div>

      <div v-else class="careers-empty">
        <div aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none"><path d="M8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><path d="M3 8h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8Z"/><path d="M3 13h18"/></svg>
        </div>
        <h3>{{ hasFilters ? 'No matching roles found' : page.empty_title }}</h3>
        <p>{{ hasFilters ? 'Try a different search, department, or job type.' : page.empty_text }}</p>
      </div>
    </section>

    <section class="careers-process">
      <div class="careers-shell">
        <div class="careers-section-head">
          <span class="careers-kicker">{{ page.process_eyebrow }}</span>
          <h2>{{ page.process_title }}</h2>
        </div>
        <div class="careers-process__grid">
          <article v-for="(step, index) in page.process" :key="step.title">
            <span>{{ String(index + 1).padStart(2, '0') }}</span>
            <h3>{{ step.title }}</h3>
            <p>{{ step.text }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="careers-shell careers-cta">
      <div>
        <span class="careers-kicker">{{ page.cta_eyebrow }}</span>
        <h2>{{ page.cta_title }}</h2>
        <p>{{ page.cta_text }}</p>
      </div>
      <router-link to="/contact-us" class="careers-button careers-button--primary">Reach Recruitment</router-link>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getSettings, listCareers } from '@/api/cms.js'
import { sanitizePlainText } from '@/utils/sanitize.js'

const defaults = {
  eyebrow: 'Careers at NCS',
  title: 'Join the team shaping Uganda sport',
  intro: 'Work with the National Council of Sports to strengthen federations, support athletes, and build a more active Uganda.',
  jobs_eyebrow: 'Open opportunities',
  jobs_title: 'Current vacancies',
  jobs_intro: 'Explore published career opportunities and review the role details before contacting the recruitment team.',
  empty_title: 'No open positions at this time',
  empty_text: 'Please check back later for new opportunities with the National Council of Sports.',
  process_eyebrow: 'Recruitment process',
  process_title: 'What to expect',
  cta_eyebrow: 'Need support?',
  cta_title: 'Have a question about a vacancy?',
  cta_text: 'Our team can help with application guidance, deadlines, and role-specific enquiries.',
  stats: [
    { value: '1964', label: 'Established', text: 'Serving Uganda sport through a statutory national mandate.' },
    { value: '50+', label: 'Sports bodies', text: 'Working alongside recognised national associations.' },
    { value: '1', label: 'National mission', text: 'Maximising opportunities for all Ugandans in sport.' },
  ],
  process: [
    { title: 'Review the role', text: 'Read the job description, requirements, deadline, and department details.' },
    { title: 'Prepare documents', text: 'Match your CV, references, and supporting documents to the published requirements.' },
    { title: 'Contact recruitment', text: 'Use the listed role details and official NCS contacts for application guidance.' },
  ],
}

const page = reactive(JSON.parse(JSON.stringify(defaults)))
const careers = ref([])
const loading = ref(true)
const search = ref('')
const selectedDepartment = ref('')
const selectedType = ref('')

const departments = computed(() => unique(careers.value.map(c => c.department || c.department_name).filter(Boolean)))
const jobTypes = computed(() => unique(careers.value.map(c => c.job_type).filter(Boolean)))
const departmentCount = computed(() => departments.value.length || 1)
const primaryLocation = computed(() => careers.value.find(c => c.location)?.location || 'Uganda')
const hasFilters = computed(() => Boolean(search.value || selectedDepartment.value || selectedType.value))

const filteredCareers = computed(() => {
  const q = search.value.toLowerCase()
  return careers.value.filter(career => {
    const haystack = [career.title, career.department, career.department_name, career.location, career.job_type].filter(Boolean).join(' ').toLowerCase()
    return (!q || haystack.includes(q))
      && (!selectedDepartment.value || [career.department, career.department_name].includes(selectedDepartment.value))
      && (!selectedType.value || career.job_type === selectedType.value)
  })
})

function unique(items) {
  return [...new Set(items)]
}

function labelJobType(type = '') {
  return String(type || 'full_time').replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}

function excerpt(value = '') {
  const plain = sanitizePlainText(String(value).replace(/<[^>]+>/g, ''))
  return plain.length > 150 ? `${plain.slice(0, 150).trim()}...` : plain
}

function mergePageSettings(value) {
  if (!value || typeof value !== 'object') return
  Object.assign(page, value)
  page.stats = Array.isArray(value.stats) && value.stats.length ? value.stats : defaults.stats
  page.process = Array.isArray(value.process) && value.process.length ? value.process : defaults.process
}

onMounted(async () => {
  try {
    const [settingsRes, careersRes] = await Promise.allSettled([
      getSettings('careers_page'),
      listCareers({ status: 'published', per_page: 100 }),
    ])
    if (settingsRes.status === 'fulfilled') mergePageSettings(settingsRes.value.data?.data?.value)
    if (careersRes.status === 'fulfilled') careers.value = careersRes.value.data.data?.items || []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.careers-page{color:#1f2937}.careers-shell{width:min(1120px,calc(100% - 2rem));margin:0 auto}.careers-hero{position:relative;overflow:hidden;background:#fff2e1;padding:5.5rem 0 4.25rem}.careers-hero__media{position:absolute;inset:0;display:flex;justify-content:flex-end;align-items:center;opacity:.08;pointer-events:none}.careers-hero__media img{width:min(38rem,70vw);height:auto;transform:translateX(12%)}.careers-hero__grid{position:relative;display:grid;grid-template-columns:minmax(0,1fr) minmax(18rem,.42fr);gap:3rem;align-items:end}.careers-kicker{display:inline-flex;color:#a84b00;font-size:.75rem;font-weight:800;letter-spacing:.12em;text-transform:uppercase}.careers-hero h1,.careers-section-head h2,.careers-cta h2{margin:.8rem 0 1rem;color:#2f327d;font-size:clamp(2.4rem,5vw,4.8rem);font-weight:800;line-height:1.05;letter-spacing:0}.careers-hero p{max-width:43rem;color:#696984;font-size:1.08rem;line-height:1.85}.careers-hero__actions{display:flex;flex-wrap:wrap;gap:.8rem;margin-top:2rem}.careers-button{display:inline-flex;align-items:center;justify-content:center;min-height:44px;border-radius:6px;padding:.8rem 1.15rem;font-size:.9rem;font-weight:800;transition:background .2s,color .2s,border-color .2s}.careers-button--primary{background:#a84b00;color:#fff}.careers-button--primary:hover{background:#873d00}.careers-button--ghost{border:1px solid rgb(47 50 125 / .24);color:#2f327d;background:#fff}.careers-button--ghost:hover{border-color:#a84b00;color:#a84b00}.careers-hero__panel{display:grid;gap:1px;border-radius:8px;overflow:hidden;background:rgb(47 50 125 / .14);box-shadow:0 24px 60px rgb(47 50 125 / .14)}.careers-hero__panel div{background:#fff;padding:1.3rem}.careers-hero__panel span{display:block;color:#696984;font-size:.78rem;font-weight:700;text-transform:uppercase;letter-spacing:.08em}.careers-hero__panel strong{display:block;margin-top:.35rem;color:#2f327d;font-size:1.7rem;font-weight:900}.careers-stats{background:#fff;padding:2.5rem 0;border-bottom:1px solid #eef0f6}.careers-stats__grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:1rem}.careers-stats article{padding:1.5rem;border:1px solid #eef0f6;border-radius:8px;background:#fff}.careers-stats span{display:block;color:#a84b00;font-size:2rem;font-weight:900}.careers-stats strong{display:block;margin:.3rem 0;color:#2f327d}.careers-stats p{color:#696984;font-size:.9rem;line-height:1.6}.careers-openings{padding:4.5rem 0}.careers-section-head{max-width:42rem;margin-bottom:2rem}.careers-section-head h2,.careers-cta h2{font-size:clamp(2rem,3vw,3rem)}.careers-section-head p,.careers-cta p{color:#696984;line-height:1.75}.careers-toolbar{display:grid;grid-template-columns:minmax(0,1.5fr) minmax(12rem,.75fr) minmax(12rem,.65fr);gap:.8rem;margin-bottom:1.4rem}.careers-toolbar label{display:grid;gap:.35rem;color:#2f327d;font-size:.78rem;font-weight:800}.careers-toolbar input,.careers-toolbar select{width:100%;min-height:44px;border:1px solid #e3e6ef;border-radius:6px;background:#fff;padding:.7rem .85rem;color:#111827}.careers-list{display:grid;gap:1rem}.career-card{display:grid;gap:1rem;border:1px solid #edf0f6;border-radius:8px;background:#fff;padding:1.35rem;color:inherit;box-shadow:0 12px 28px rgb(47 50 125 / .06);transition:transform .2s,border-color .2s,box-shadow .2s}.career-card:hover{transform:translateY(-2px);border-color:rgb(168 75 0 / .35);box-shadow:0 22px 48px rgb(47 50 125 / .12)}.career-card__top,.career-card__meta{display:flex;flex-wrap:wrap;align-items:center;gap:.6rem}.career-card__badge{border-radius:999px;background:#fff2e1;color:#a84b00;padding:.3rem .65rem;font-size:.74rem;font-weight:800}.career-card__deadline{margin-left:auto;color:#b91c1c;font-size:.78rem;font-weight:800}.career-card h3{color:#2f327d;font-size:1.25rem;font-weight:800;letter-spacing:0}.career-card p{max-width:52rem;color:#696984;line-height:1.7}.career-card__meta span{display:inline-flex;align-items:center;gap:.35rem;color:#4b5563;font-size:.84rem;font-weight:700}.career-card__meta svg{width:1rem;height:1rem;stroke:currentColor;stroke-width:1.8}.career-card__cta{justify-self:start;color:#a84b00;font-weight:900;font-size:.9rem}.career-card--loading span,.career-card--loading strong,.career-card--loading p{display:block;border-radius:999px;background:#eef0f6;min-height:1rem;animation:pulse 1.3s infinite}.career-card--loading strong{width:55%;height:1.35rem}.career-card--loading p{width:80%}@keyframes pulse{50%{opacity:.45}}.careers-empty{display:grid;justify-items:center;text-align:center;border:1px dashed #d9deea;border-radius:8px;padding:3rem 1rem}.careers-empty div{display:grid;place-items:center;width:4rem;height:4rem;border-radius:50%;background:#fff2e1;color:#a84b00}.careers-empty svg{width:2rem;height:2rem;stroke:currentColor;stroke-width:1.7}.careers-empty h3{margin:1rem 0 .4rem;color:#2f327d;font-size:1.2rem;font-weight:900}.careers-empty p{max-width:34rem;color:#696984}.careers-process{background:#f8fafc;padding:4rem 0}.careers-process__grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:1rem}.careers-process article{border-radius:8px;background:#fff;border:1px solid #edf0f6;padding:1.4rem}.careers-process span{color:#a84b00;font-size:.82rem;font-weight:900}.careers-process h3{margin:.7rem 0 .45rem;color:#2f327d;font-weight:900}.careers-process p{color:#696984;line-height:1.7}.careers-cta{display:grid;grid-template-columns:minmax(0,1fr) auto;gap:2rem;align-items:center;margin-top:4rem;margin-bottom:4rem;border-radius:8px;background:#2f327d;padding:2rem}.careers-cta .careers-kicker,.careers-cta h2,.careers-cta p{color:#fff}.careers-cta p{opacity:.78;max-width:42rem}.careers-cta .careers-button{white-space:nowrap}
@media(max-width:900px){.careers-hero__grid,.careers-toolbar,.careers-cta{grid-template-columns:1fr}.careers-hero__panel{max-width:28rem}.careers-stats__grid,.careers-process__grid{grid-template-columns:1fr}.career-card__deadline{margin-left:0}.careers-cta{padding:1.5rem}}
</style>
