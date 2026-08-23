<template>
  <section class="open-forms-section">
    <!-- Header & Search Toolbar -->
    <header class="section-toolbar-header mb-4">
      <div class="toolbar-titles">
        <span class="toolbar-tag"><i class="icofont-cube me-1"></i> Available Online Services</span>
        <h2 class="toolbar-heading">{{ title || 'Services & Application Catalog' }}</h2>
        <p class="toolbar-subtext">Select an official National Council of Sports service below to begin or continue your application.</p>
      </div>

      <div v-if="!expanded && forms.length > 6" class="toolbar-actions">
        <button type="button" class="btn btn-outline-primary btn-sm view-all-btn" @click="$emit('view-all')">
          View All ({{ forms.length }}) <i class="icofont-arrow-right ms-1"></i>
        </button>
      </div>
    </header>

    <!-- Search & Filter Controls (When in Expanded Apply View) -->
    <div v-if="expanded" class="catalog-filters-bar mb-4">
      <div class="search-input-wrap">
        <i class="icofont-search search-icon"></i>
        <input
          v-model="searchQuery"
          type="text"
          class="form-control catalog-search-input"
          placeholder="Search by service title, department or keyword..."
        />
        <button v-if="searchQuery" type="button" class="btn-clear-search" @click="searchQuery = ''">
          <i class="icofont-close-line"></i>
        </button>
      </div>

      <div class="category-filter-chips">
        <button
          type="button"
          class="category-chip"
          :class="{ active: selectedCategory === 'ALL' }"
          @click="selectedCategory = 'ALL'"
        >
          All Services ({{ forms.length }})
        </button>
        <button
          v-for="cat in availableCategories"
          :key="cat.name"
          type="button"
          class="category-chip"
          :class="{ active: selectedCategory === cat.name }"
          @click="selectedCategory = cat.name"
        >
          {{ cat.name }} ({{ cat.count }})
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="forms-loading-grid">
      <div v-for="n in 6" :key="n" class="service-card skeleton-card">
        <div class="skeleton-thumbnail"></div>
        <div class="skeleton-body">
          <div class="skeleton-line short"></div>
          <div class="skeleton-line title"></div>
          <div class="skeleton-line desc"></div>
          <div class="skeleton-btn"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!filteredForms.length" class="empty-state card shadow-sm text-center border-0 p-5">
      <div class="empty-icon-wrap mb-3">
        <i class="icofont-search-folder text-muted" style="font-size: 54px;"></i>
      </div>
      <h4 class="fw-bold empty-title">No matching services found</h4>
      <p class="text-muted max-w-md mx-auto">
        We couldn't find any available application forms matching "{{ searchQuery }}". Try checking your spelling or clearing the category filters.
      </p>
      <button v-if="searchQuery || selectedCategory !== 'ALL'" type="button" class="btn btn-outline-secondary btn-sm mt-2 mx-auto" @click="resetFilters">
        Reset Search & Filters
      </button>
    </div>

    <!-- Modern Professional Cards Grid -->
    <div v-else class="services-cards-grid">
      <article
        v-for="form in filteredForms"
        :key="form.id"
        class="service-card"
        :class="{ 'has-draft': draftFor(form.id), 'is-pending': pendingFor(form.id) }"
      >
        <!-- Thumbnail / Banner Area -->
        <div class="service-card-banner">
          <img
            v-if="form.banner_image_url"
            :src="mediaUrl(form.banner_image_url)"
            :alt="form.title"
            class="banner-img"
            loading="lazy"
          />
          <div v-else class="banner-fallback" :class="getGradientClass(form.id)">
            <i class="icofont-file-document fallback-icon"></i>
            <span class="fallback-pattern"></span>
          </div>

          <!-- Overlay Badges -->
          <div class="banner-overlay-top">
            <span class="department-pill">
              <i class="icofont-building-alt me-1"></i> {{ departmentLabel(form.department_name) }}
            </span>
            <span class="fee-pill" :class="Number(form.price_ugx) ? 'paid-fee' : 'free-fee'">
              {{ Number(form.price_ugx) ? `UGX ${formatMoney(form.price_ugx)}` : 'FREE' }}
            </span>
          </div>
        </div>

        <!-- Service Content Body -->
        <div class="service-card-body">
          <div class="service-meta-line">
            <span class="steps-indicator">
              <i class="icofont-listine-dots me-1"></i> {{ stepCount(form) }} Step{{ stepCount(form) === 1 ? '' : 's' }}
            </span>
            <span v-if="draftFor(form.id)" class="badge bg-warning-light text-warning fw-bold">
              <i class="icofont-edit-alt me-1"></i> Draft in Progress
            </span>
            <span v-else-if="pendingFor(form.id)" class="badge bg-info-light text-info fw-bold">
              <i class="icofont-clock-time me-1"></i> Under Review
            </span>
            <span v-else class="badge bg-light text-muted">
              <i class="icofont-check-circled text-success me-1"></i> Open for Submission
            </span>
          </div>

          <h3 class="service-title" :title="form.title">{{ form.title }}</h3>
          
          <p class="service-description">
            {{ form.description || 'Complete and submit this official application form online through the NCS portal.' }}
          </p>

          <!-- Card Footer & Action Button -->
          <div class="service-card-footer">
            <button
              type="button"
              class="service-action-btn"
              :class="getActionBtnClass(form.id)"
              :disabled="!!pendingFor(form.id)"
              @click="$emit('start', form)"
            >
              <span>{{ actionLabel(form.id) }}</span>
              <i :class="getActionBtnIcon(form.id)"></i>
            </button>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { ref, computed } from 'vue'
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

const searchQuery = ref('')
const selectedCategory = ref('ALL')

const availableCategories = computed(() => {
  const map = {}
  props.forms.forEach(f => {
    const dept = departmentLabel(f.department_name)
    map[dept] = (map[dept] || 0) + 1
  })
  return Object.keys(map).sort().map(name => ({ name, count: map[name] }))
})

const visibleForms = computed(() => {
  if (!props.expanded) {
    return props.forms.slice(0, 6)
  }
  return props.forms
})

const filteredForms = computed(() => {
  let list = visibleForms.value
  if (selectedCategory.value !== 'ALL') {
    list = list.filter(f => departmentLabel(f.department_name) === selectedCategory.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(f => 
      String(f.title || '').toLowerCase().includes(q) ||
      String(f.description || '').toLowerCase().includes(q) ||
      String(f.department_name || '').toLowerCase().includes(q)
    )
  }
  return list
})

function resetFilters() {
  searchQuery.value = ''
  selectedCategory.value = 'ALL'
}

const draftFor = id => props.submissions.find(item => item.template_id === id && ['DRAFT', 'NEEDS_INFORMATION'].includes(item.status))
const pendingFor = id => props.submissions.find(item => item.template_id === id && isPendingSubmission(item))

function actionLabel(id) {
  if (pendingFor(id)) return 'Application Submitted'
  const draft = draftFor(id)
  if (draft?.status === 'NEEDS_INFORMATION') return 'Edit Queried Application'
  if (draft) return 'Continue Application'
  return 'Start Application'
}

function getActionBtnClass(id) {
  if (pendingFor(id)) return 'btn-action-pending'
  const draft = draftFor(id)
  if (draft?.status === 'NEEDS_INFORMATION') return 'btn-action-queried'
  if (draft) return 'btn-action-continue'
  return 'btn-action-start'
}

function getActionBtnIcon(id) {
  if (pendingFor(id)) return 'icofont-check-circled'
  const draft = draftFor(id)
  if (draft?.status === 'NEEDS_INFORMATION') return 'icofont-warning-alt'
  if (draft) return 'icofont-arrow-right'
  return 'icofont-plus-circle'
}

function getGradientClass(id) {
  const gradients = ['gradient-blue', 'gradient-purple', 'gradient-teal', 'gradient-indigo', 'gradient-orange']
  const index = Math.abs(String(id).split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)) % gradients.length
  return gradients[index]
}

const stepCount = form => buildSectionSteps(form).length
const formatMoney = value => new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0))
function departmentLabel(value) {
  if (!value) return 'NCS General Service'
  if (/(^|_)(super_?admin|admin|user)($|_)/i.test(String(value))) return 'NCS Central Secretariat'
  return String(value).trim()
}
const isPendingSubmission = item => item?.status && !['DRAFT', 'NEEDS_INFORMATION', 'APPROVED', 'REJECTED'].includes(item.status)
</script>

<style scoped>
.open-forms-section {
  width: 100%;
  max-width: 100%;
}

/* ── Section Header Toolbar ───────────────────────────────────────── */
.section-toolbar-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}
.toolbar-tag {
  display: inline-block;
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #6777ef;
  margin-bottom: 4px;
}
.toolbar-heading {
  font-size: 22px;
  font-weight: 800;
  color: #1e293b;
  margin: 0 0 4px 0;
  line-height: 1.25;
}
.toolbar-subtext {
  font-size: 13px;
  color: #64748b;
  margin: 0;
}
.view-all-btn {
  border-radius: 20px;
  font-weight: 700;
  font-size: 12px;
  padding: 8px 18px;
  border-color: #6777ef;
  color: #6777ef;
  transition: all 0.2s ease;
}
.view-all-btn:hover {
  background: #6777ef;
  color: #fff;
}

/* ── Filters & Search Bar ────────────────────────────────────────── */
.catalog-filters-bar {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.search-input-wrap {
  position: relative;
  width: 100%;
}
.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
  font-size: 18px;
  pointer-events: none;
}
.catalog-search-input {
  width: 100%;
  height: 48px;
  padding-left: 48px;
  padding-right: 42px;
  font-size: 14px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  background: #fff;
  color: #1e293b;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.02);
  transition: all 0.2s ease;
}
.catalog-search-input:focus {
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15);
}
.btn-clear-search {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: #94a3b8;
  font-size: 18px;
  cursor: pointer;
  padding: 4px;
}
.btn-clear-search:hover {
  color: #64748b;
}

.category-filter-chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}
.category-chip {
  padding: 7px 16px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: #fff;
  border: 1px solid #e2e8f0;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}
.category-chip:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}
.category-chip.active {
  background: #6777ef;
  border-color: #6777ef;
  color: #fff;
  box-shadow: 0 2px 8px rgba(103, 119, 239, 0.3);
}

/* ── Modern Services Cards Grid ──────────────────────────────────── */
.services-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 22px;
  align-items: stretch;
}

.service-card {
  display: flex;
  flex-direction: column;
  border-radius: 14px;
  overflow: hidden;
  background: #fff;
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}
.service-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.1);
  border-color: rgba(103, 119, 239, 0.4);
}

/* Banner Area */
.service-card-banner {
  position: relative;
  width: 100%;
  height: 160px;
  overflow: hidden;
  background: #0f172a;
}
.banner-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.35s ease;
}
.service-card:hover .banner-img {
  transform: scale(1.05);
}

.banner-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
.fallback-icon {
  font-size: 54px;
  color: rgba(255, 255, 255, 0.35);
  z-index: 1;
}
.gradient-blue { background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%); }
.gradient-purple { background: linear-gradient(135deg, #4c1d95 0%, #8b5cf6 100%); }
.gradient-teal { background: linear-gradient(135deg, #0f766e 0%, #14b8a6 100%); }
.gradient-indigo { background: linear-gradient(135deg, #312e81 0%, #6366f1 100%); }
.gradient-orange { background: linear-gradient(135deg, #7c2d12 0%, #ea580c 100%); }

.banner-overlay-top {
  position: absolute;
  top: 12px;
  left: 12px;
  right: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  z-index: 2;
}
.department-pill {
  font-size: 11px;
  font-weight: 700;
  padding: 5px 10px;
  border-radius: 20px;
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(4px);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.15);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 65%;
}
.fee-pill {
  font-size: 11px;
  font-weight: 800;
  padding: 5px 11px;
  border-radius: 20px;
  letter-spacing: 0.3px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}
.fee-pill.free-fee {
  background: #10b981;
  color: #fff;
}
.fee-pill.paid-fee {
  background: #3b82f6;
  color: #fff;
}

/* Service Card Body */
.service-card-body {
  display: flex;
  flex-direction: column;
  padding: 18px;
  flex-grow: 1;
}
.service-meta-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 10px;
  font-size: 11px;
}
.steps-indicator {
  color: #64748b;
  font-weight: 600;
}
.service-title {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 8px 0;
  line-height: 1.35;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.service-description {
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
  margin: 0 0 16px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex-grow: 1;
}

/* Card Footer & Action Button */
.service-card-footer {
  margin-top: auto;
  padding-top: 14px;
  border-top: 1px solid #f1f5f9;
}
.service-action-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 11px 18px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 700;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}
.btn-action-start {
  background: #6777ef;
  color: #fff;
  box-shadow: 0 4px 12px rgba(103, 119, 239, 0.25);
}
.btn-action-start:hover {
  background: #5563db;
  box-shadow: 0 6px 16px rgba(103, 119, 239, 0.35);
}
.btn-action-continue {
  background: #f59e0b;
  color: #fff;
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.25);
}
.btn-action-continue:hover {
  background: #d97706;
}
.btn-action-queried {
  background: #ef4444;
  color: #fff;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.25);
}
.btn-action-queried:hover {
  background: #dc2626;
}
.btn-action-pending {
  background: #e2e8f0;
  color: #94a3b8;
  cursor: not-allowed;
}

/* Skeleton Loading Cards */
.forms-loading-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 22px;
}
.skeleton-card {
  height: 330px;
}
.skeleton-thumbnail {
  height: 160px;
  background: #e2e8f0;
  animation: skeletonPulse 1.5s infinite;
}
.skeleton-body {
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.skeleton-line {
  background: #e2e8f0;
  border-radius: 4px;
  animation: skeletonPulse 1.5s infinite;
}
.skeleton-line.short { width: 30%; height: 14px; }
.skeleton-line.title { width: 85%; height: 20px; }
.skeleton-line.desc { width: 100%; height: 36px; }
.skeleton-btn { height: 42px; background: #e2e8f0; border-radius: 10px; margin-top: auto; animation: skeletonPulse 1.5s infinite; }

@keyframes skeletonPulse {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}

.empty-icon-wrap {
  font-size: 54px;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .services-cards-grid,
  .forms-loading-grid {
    grid-template-columns: 1fr;
  }
}

/* ── Dark Mode Support ─────────────────────────────────────────────── */
:global(.dark) .toolbar-heading {
  color: #f8fafc;
}
:global(.dark) .toolbar-subtext {
  color: #94a3b8;
}
:global(.dark) .catalog-search-input {
  background: #1e293b;
  border-color: #334155;
  color: #f8fafc;
}
:global(.dark) .category-chip {
  background: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}
:global(.dark) .category-chip:hover {
  background: #334155;
  color: #f8fafc;
}
:global(.dark) .category-chip.active {
  background: #6777ef;
  border-color: #6777ef;
  color: #fff;
}
:global(.dark) .service-card {
  background: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25);
}
:global(.dark) .service-card:hover {
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4);
  border-color: #6777ef;
}
:global(.dark) .service-title {
  color: #f8fafc;
}
:global(.dark) .service-description {
  color: #94a3b8;
}
:global(.dark) .service-card-footer {
  border-top-color: #334155;
}
:global(.dark) .btn-action-pending {
  background: #334155;
  color: #64748b;
}
:global(.dark) .skeleton-thumbnail,
:global(.dark) .skeleton-line,
:global(.dark) .skeleton-btn {
  background: #334155;
}
</style>

