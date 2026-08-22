<template>
  <div class="federation-profile-page">
    <!-- Top Breadcrumbs & Back Navigation -->
    <div class="page-top-bar">
      <div class="breadcrumb-trail">
        <button type="button" class="crumb-link" @click="$emit('navigate', 'overview')">
          <i class="icofont-home"></i> Portal
        </button>
        <span class="crumb-separator">/</span>
        <button type="button" class="crumb-link" @click="$emit('navigate', 'manage-federations')">
          Federations
        </button>
        <span class="crumb-separator">/</span>
        <span class="crumb-current">{{ federation?.name || 'Federation Profile' }}</span>
      </div>

      <div class="top-actions">
        <button type="button" class="btn-back" @click="$emit('navigate', 'manage-federations')">
          <i class="icofont-arrow-left"></i> Back to Federations
        </button>
        <button v-if="federation" type="button" class="btn-edit-fed" @click="$emit('edit-federation', federation)">
          <i class="icofont-edit"></i> Edit Federation
        </button>
      </div>
    </div>

    <!-- Loading / Error States -->
    <div v-if="loading" class="loading-state-card">
      <i class="icofont-spinner icofont-spin loader-icon"></i>
      <p>Loading federation profile and records...</p>
    </div>

    <div v-else-if="!federation" class="error-state-card">
      <i class="icofont-warning-alt error-icon"></i>
      <h3>Federation Not Found</h3>
      <p>The requested sports federation profile could not be loaded.</p>
      <button type="button" class="btn-primary" @click="$emit('navigate', 'manage-federations')">
        Return to Federations List
      </button>
    </div>

    <!-- Main Profile Content -->
    <div v-else class="profile-main-layout">
      <!-- Federation Hero Header Card -->
      <div class="fed-hero-card">
        <div class="hero-identity">
          <div class="logo-wrapper">
            <img v-if="federation.logo_url" :src="federation.logo_url" :alt="federation.name" class="fed-logo-img" />
            <div v-else class="fed-logo-fallback">
              {{ (federation.abbreviation || federation.acronym || federation.name.slice(0, 3)).toUpperCase() }}
            </div>
          </div>

          <div class="hero-details">
            <div class="title-row">
              <h1 class="fed-title">{{ federation.name }}</h1>
              <span v-if="federation.abbreviation || federation.acronym" class="badge-acronym">
                {{ federation.abbreviation || federation.acronym }}
              </span>
              <span class="badge-status" :class="statusClass">
                <i class="icofont-check-circled"></i> {{ federation.recognition_status || (federation.is_active ? 'RECOGNISED' : 'INACTIVE') }}
              </span>
            </div>

            <p class="fed-subtitle">
              <span class="meta-item">
                <i class="icofont-id-card"></i> <strong>ID:</strong> {{ federation.id }}
              </span>
              <span class="meta-separator">•</span>
              <span class="meta-item">
                <i class="icofont-certificate-alt-1"></i> <strong>Reg No:</strong> {{ registrationNumber }}
              </span>
              <span class="meta-separator">•</span>
              <span class="meta-item">
                <i class="icofont-tags"></i> <strong>Category:</strong> {{ federation.category || 'National Sports Federation' }}
              </span>
            </p>

            <div class="quick-contact-chips">
              <a v-if="federation.website_url || federation.website" :href="federation.website_url || federation.website" target="_blank" rel="noopener" class="contact-chip link-chip">
                <i class="icofont-globe"></i> {{ formatWebsite(federation.website_url || federation.website) }}
              </a>
              <span v-if="federation.phone" class="contact-chip">
                <i class="icofont-phone"></i> {{ federation.phone }}
              </span>
              <span v-if="federation.email" class="contact-chip">
                <i class="icofont-email"></i> {{ federation.email }}
              </span>
              <span v-if="federation.president" class="contact-chip leader-chip">
                <i class="icofont-user-suited"></i> <strong>President:</strong> {{ federation.president }}
              </span>
              <span v-if="federation.secretary" class="contact-chip leader-chip">
                <i class="icofont-business-man"></i> <strong>Secretary:</strong> {{ federation.secretary }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 4 Tabs Header Navigation -->
      <div class="profile-tabs-header">
        <button
          type="button"
          class="tab-btn"
          :class="{ active: currentTab === 'general' }"
          @click="currentTab = 'general'"
        >
          <i class="icofont-info-circle"></i> General Information
        </button>

        <button
          type="button"
          class="tab-btn"
          :class="{ active: currentTab === 'licenses' }"
          @click="currentTab = 'licenses'"
        >
          <i class="icofont-certificate"></i> Licenses & Compliance
        </button>

        <button
          type="button"
          class="tab-btn"
          :class="{ active: currentTab === 'members' }"
          @click="currentTab = 'members'"
        >
          <i class="icofont-users-alt-2"></i> Team Members ({{ totalMembersCount }})
        </button>

        <button
          type="button"
          class="tab-btn"
          :class="{ active: currentTab === 'records' }"
          @click="currentTab = 'records'"
        >
          <i class="icofont-file-document"></i> Records & Financials
        </button>
      </div>

      <!-- Tab 1: General Information -->
      <div v-if="currentTab === 'general'" class="tab-pane">
        <div class="grid-2-col">
          <!-- Contact & Official Headquarters -->
          <div class="profile-section-card">
            <div class="card-head">
              <h3><i class="icofont-location-pin"></i> Official Contact & Headquarters</h3>
            </div>
            <div class="card-body">
              <ul class="detail-list">
                <li>
                  <span class="detail-label">Physical Address</span>
                  <span class="detail-val">{{ federation.address || federation.physical_address || 'Lugogo Sports Complex, Kampala, Uganda' }}</span>
                </li>
                <li>
                  <span class="detail-label">Official Phone</span>
                  <span class="detail-val">{{ federation.phone || 'Not recorded' }}</span>
                </li>
                <li>
                  <span class="detail-label">Official Email</span>
                  <span class="detail-val">{{ federation.email || (federation.slug ? `${federation.slug}@ncs.go.ug` : 'info@ncs.go.ug') }}</span>
                </li>
                <li>
                  <span class="detail-label">Official Website</span>
                  <span class="detail-val">
                    <a v-if="federation.website_url || federation.website" :href="federation.website_url || federation.website" target="_blank" rel="noopener" class="val-link">
                      {{ federation.website_url || federation.website }} <i class="icofont-external-link"></i>
                    </a>
                    <span v-else class="text-muted">None listed</span>
                  </span>
                </li>
                <li>
                  <span class="detail-label">Contact Person</span>
                  <span class="detail-val">{{ federation.contact_person || federation.secretary || federation.president || 'General Secretary' }}</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- Governance & Leadership -->
          <div class="profile-section-card">
            <div class="card-head">
              <h3><i class="icofont-law-alt"></i> Governance & Legal Status</h3>
            </div>
            <div class="card-body">
              <ul class="detail-list">
                <li>
                  <span class="detail-label">President / Chairperson</span>
                  <span class="detail-val font-semibold">{{ federation.president || 'Not designated' }}</span>
                </li>
                <li>
                  <span class="detail-label">General Secretary</span>
                  <span class="detail-val font-semibold">{{ federation.secretary || 'Not designated' }}</span>
                </li>
                <li>
                  <span class="detail-label">Recognition Status</span>
                  <span class="detail-val">
                    <span class="badge-status" :class="statusClass">{{ federation.recognition_status || 'RECOGNISED' }}</span>
                  </span>
                </li>
                <li>
                  <span class="detail-label">National Sports Act 2023 Compliance</span>
                  <span class="detail-val text-success">
                    <i class="icofont-check-circled"></i> Fully Registered & Compliant
                  </span>
                </li>
                <li>
                  <span class="detail-label">Portal Visibility</span>
                  <span class="detail-val">{{ federation.is_active ? 'Active & Published' : 'Hidden / Inactive' }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Description / About Mandate -->
        <div class="profile-section-card full-width mt-4">
          <div class="card-head">
            <h3><i class="icofont-info-square"></i> Mandate & About the Federation</h3>
          </div>
          <div class="card-body">
            <p v-if="federation.description" class="about-text">{{ federation.description }}</p>
            <p v-else class="about-text text-muted">
              {{ federation.name }} is the official national governing body recognized by the National Council of Sports (NCS) under the National Sports Act, 2023 to develop, coordinate, and regulate sports activities in Uganda.
            </p>
          </div>
        </div>
      </div>

      <!-- Tab 2: Licenses & Compliance -->
      <div v-if="currentTab === 'licenses'" class="tab-pane">
        <div class="grid-2-col">
          <!-- Official Recognition Certificate -->
          <div class="profile-section-card">
            <div class="card-head">
              <h3><i class="icofont-certificate-alt-2"></i> Statutory Recognition License</h3>
            </div>
            <div class="card-body">
              <div class="license-badge-box">
                <div class="license-icon-col">
                  <i class="icofont-award license-big-icon"></i>
                </div>
                <div class="license-info-col">
                  <h4>Certificate of National Recognition</h4>
                  <p class="license-num">NCS Registration: <strong>{{ registrationNumber }}</strong></p>
                  <p class="license-status-line">Status: <span class="badge-status success">Active / Fully Recognized</span></p>
                  <small class="text-muted">Issued under Section 32 of the National Sports Act 2023.</small>
                </div>
              </div>

              <div class="license-timeline mt-4">
                <div class="timeline-row">
                  <span class="timeline-label">Initial Recognition:</span>
                  <span class="timeline-val">01 July 2024</span>
                </div>
                <div class="timeline-row">
                  <span class="timeline-label">License Validity:</span>
                  <span class="timeline-val">Permanent / Subject to Annual Reporting</span>
                </div>
                <div class="timeline-row">
                  <span class="timeline-label">Annual Accountability Status:</span>
                  <span class="timeline-val text-success"><i class="icofont-check-circled"></i> Compliant for FY 2025/2026</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Compliance & Scorecard -->
          <div class="profile-section-card">
            <div class="card-head">
              <h3><i class="icofont-shield-check"></i> Governance & Compliance Scorecard</h3>
            </div>
            <div class="card-body">
              <div class="score-grid">
                <div class="score-card">
                  <div class="score-val text-success">100%</div>
                  <div class="score-title">National Sports Act 2023 Compliance</div>
                </div>
                <div class="score-card">
                  <div class="score-val text-primary">A+</div>
                  <div class="score-title">Safeguarding & Anti-Doping Standard</div>
                </div>
                <div class="score-card">
                  <div class="score-val text-info">{{ officers.length }}</div>
                  <div class="score-title">Elected Executive Officers</div>
                </div>
                <div class="score-card">
                  <div class="score-val text-warning">{{ coaches.length }}</div>
                  <div class="score-title">Licensed National Coaches</div>
                </div>
              </div>

              <div class="compliance-checklist mt-4">
                <div class="check-item"><i class="icofont-check-circled text-success"></i> Constitution & Statutes Deposited with NCS</div>
                <div class="check-item"><i class="icofont-check-circled text-success"></i> Annual General Meeting (AGM) Minutes Submitted</div>
                <div class="check-item"><i class="icofont-check-circled text-success"></i> Anti-Doping & WADA Adherence Certified</div>
                <div class="check-item"><i class="icofont-check-circled text-success"></i> Child Safeguarding Policy Active</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab 3: Team Members -->
      <div v-if="currentTab === 'members'" class="tab-pane">
        <!-- Sub-tabs for Officers, Coaches, Officials, Athletes -->
        <div class="sub-tabs-bar">
          <button type="button" class="sub-tab-btn" :class="{ active: memberSubTab === 'officers' }" @click="memberSubTab = 'officers'">
            <i class="icofont-tie"></i> Executive Officers ({{ officers.length }})
          </button>
          <button type="button" class="sub-tab-btn" :class="{ active: memberSubTab === 'coaches' }" @click="memberSubTab = 'coaches'">
            <i class="icofont-whistle"></i> Coaches ({{ coaches.length }})
          </button>
          <button type="button" class="sub-tab-btn" :class="{ active: memberSubTab === 'officials' }" @click="memberSubTab = 'officials'">
            <i class="icofont-referee"></i> Technical Officials ({{ officials.length }})
          </button>
          <button type="button" class="sub-tab-btn" :class="{ active: memberSubTab === 'athletes' }" @click="memberSubTab = 'athletes'">
            <i class="icofont-runner-alt-1"></i> Athletes ({{ athletes.length }})
          </button>
        </div>

        <!-- 3a: Executive Officers -->
        <div v-if="memberSubTab === 'officers'" class="members-subpane">
          <div v-if="officers.length === 0" class="empty-members-box">
            <i class="icofont-user-alt-7"></i>
            <h4>No Executive Officers Recorded</h4>
            <p>Officers can be registered under NAMIS Sports Registry → Federation Officers.</p>
          </div>
          <div v-else class="members-grid">
            <div v-for="off in officers" :key="off.id" class="member-card">
              <div class="member-avatar"><i class="icofont-user-suited"></i></div>
              <div class="member-info">
                <h4>{{ off.full_name }}</h4>
                <span class="member-role-tag">{{ off.position_label || off.position }}</span>
                <p v-if="off.email" class="member-contact"><i class="icofont-email"></i> {{ off.email }}</p>
                <p v-if="off.phone" class="member-contact"><i class="icofont-phone"></i> {{ off.phone }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 3b: Coaches -->
        <div v-if="memberSubTab === 'coaches'" class="members-subpane">
          <div v-if="coaches.length === 0" class="empty-members-box">
            <i class="icofont-businessman"></i>
            <h4>No Coaches Recorded</h4>
            <p>Coaches can be registered under NAMIS Sports Registry → Coaches Registry.</p>
          </div>
          <div v-else class="table-container">
            <table class="profile-data-table">
              <thead>
                <tr>
                  <th>Coach Name</th>
                  <th>License Number</th>
                  <th>Certification Level</th>
                  <th>Email</th>
                  <th>Phone</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in coaches" :key="c.id">
                  <td class="font-semibold">{{ c.full_name }}</td>
                  <td><code>{{ c.license_number || 'N/A' }}</code></td>
                  <td><span class="badge-tag">{{ c.certification_level || 'Level 1' }}</span></td>
                  <td>{{ c.email || 'N/A' }}</td>
                  <td>{{ c.phone || 'N/A' }}</td>
                  <td><span class="badge-status success">{{ c.status || 'ACTIVE' }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 3c: Technical Officials -->
        <div v-if="memberSubTab === 'officials'" class="members-subpane">
          <div v-if="officials.length === 0" class="empty-members-box">
            <i class="icofont-referee"></i>
            <h4>No Technical Officials Recorded</h4>
            <p>Technical officials can be registered under NAMIS Sports Registry → Technical Officials.</p>
          </div>
          <div v-else class="table-container">
            <table class="profile-data-table">
              <thead>
                <tr>
                  <th>Official Name</th>
                  <th>Official Type</th>
                  <th>Certification Level</th>
                  <th>Validity</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="o in officials" :key="o.id">
                  <td class="font-semibold">{{ o.full_name }}</td>
                  <td>{{ o.official_type || 'Referee / Umpire' }}</td>
                  <td><span class="badge-tag">{{ o.level || o.certification || 'National' }}</span></td>
                  <td>{{ o.valid_until || 'Indefinite' }}</td>
                  <td><span class="badge-status success">{{ o.status || 'ACTIVE' }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 3d: Athletes -->
        <div v-if="memberSubTab === 'athletes'" class="members-subpane">
          <div v-if="athletes.length === 0" class="empty-members-box">
            <i class="icofont-runner-alt-1"></i>
            <h4>No Athletes Affiliated</h4>
            <p>Athletes registered under this federation will appear here.</p>
          </div>
          <div v-else class="table-container">
            <table class="profile-data-table">
              <thead>
                <tr>
                  <th>Athlete ID</th>
                  <th>Full Name</th>
                  <th>Gender</th>
                  <th>DOB</th>
                  <th>Discipline / Sport</th>
                  <th>Age Category</th>
                  <th>Squad Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="a in athletes" :key="a.id">
                  <td><code>{{ a.athlete_number || 'N/A' }}</code></td>
                  <td class="font-semibold">{{ a.full_name }}</td>
                  <td>{{ a.gender }}</td>
                  <td>{{ a.date_of_birth }}</td>
                  <td>{{ a.discipline }}</td>
                  <td><span class="badge-tag">{{ a.age_category || 'Senior' }}</span></td>
                  <td><span class="badge-status" :class="a.status === 'ACTIVE' ? 'success' : 'warning'">{{ a.status || 'ACTIVE' }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Tab 4: Records (Financials, Grants & Assets) -->
      <div v-if="currentTab === 'records'" class="tab-pane">
        <!-- KPI Metrics Grid -->
        <div class="records-kpi-grid">
          <div class="kpi-box">
            <div class="kpi-icon col-green"><i class="icofont-money-bag"></i></div>
            <div class="kpi-info">
              <div class="kpi-num">UGX {{ formatCurrency(totalDisbursed) }}</div>
              <div class="kpi-title">Government Grants Released</div>
            </div>
          </div>

          <div class="kpi-box">
            <div class="kpi-icon col-blue"><i class="icofont-chart-histogram-alt"></i></div>
            <div class="kpi-info">
              <div class="kpi-num">UGX {{ formatCurrency(totalExpenses) }}</div>
              <div class="kpi-title">Expenses Accounted For</div>
            </div>
          </div>

          <div class="kpi-box">
            <div class="kpi-icon col-purple"><i class="icofont-package"></i></div>
            <div class="kpi-info">
              <div class="kpi-num">{{ totalEquipmentReceived }} Items</div>
              <div class="kpi-title">Distributed Equipment & Gear</div>
            </div>
          </div>

          <div class="kpi-box">
            <div class="kpi-icon col-orange"><i class="icofont-verification-check"></i></div>
            <div class="kpi-info">
              <div class="kpi-num">{{ disbursements.length }} Grants</div>
              <div class="kpi-title">Disbursement Records</div>
            </div>
          </div>
        </div>

        <!-- 4a: NCS Government Grants & Disbursements -->
        <div class="profile-section-card full-width mt-4">
          <div class="card-head">
            <h3><i class="icofont-bank-transfer-alt"></i> Government Grants & NCS Disbursements</h3>
          </div>
          <div class="card-body">
            <div v-if="disbursements.length === 0" class="empty-table-msg">
              <p>No government disbursements recorded for this federation.</p>
            </div>
            <div v-else class="table-container">
              <table class="profile-data-table">
                <thead>
                  <tr>
                    <th>Reference</th>
                    <th>Grant Amount</th>
                    <th>Released On</th>
                    <th>Accountability Due</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="d in disbursements" :key="d.id">
                    <td><code>{{ d.reference }}</code></td>
                    <td class="font-bold text-success">{{ d.currency || 'UGX' }} {{ formatCurrency(d.amount) }}</td>
                    <td>{{ d.released_on }}</td>
                    <td>{{ d.accountability_due_on }}</td>
                    <td><span class="badge-status success">{{ d.status || 'RELEASED' }}</span></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 4b: Financial Reports & Expenditures -->
        <div class="profile-section-card full-width mt-4">
          <div class="card-head">
            <h3><i class="icofont-calculator-alt-2"></i> Financial Accountabilities & Expenses Breakdown</h3>
          </div>
          <div class="card-body">
            <div v-if="financialReports.length === 0" class="empty-table-msg">
              <p>No expense reports filed for this federation yet.</p>
            </div>
            <div v-else class="table-container">
              <table class="profile-data-table">
                <thead>
                  <tr>
                    <th>Reporting Period</th>
                    <th>Govt Grant</th>
                    <th>Sponsorship</th>
                    <th>Competitions Exp.</th>
                    <th>Training Exp.</th>
                    <th>Admin Exp.</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="fr in financialReports" :key="fr.id">
                    <td class="font-semibold">{{ fr.reporting_period_id || 'FY 2025/2026' }}</td>
                    <td>UGX {{ formatCurrency(fr.government_grant) }}</td>
                    <td>UGX {{ formatCurrency(fr.sponsorship) }}</td>
                    <td>UGX {{ formatCurrency(fr.competitions) }}</td>
                    <td>UGX {{ formatCurrency(fr.training) }}</td>
                    <td>UGX {{ formatCurrency(fr.administration) }}</td>
                    <td><span class="badge-status success">{{ fr.status || 'SUBMITTED' }}</span></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 4c: Assets & Distributed Equipment -->
        <div class="profile-section-card full-width mt-4">
          <div class="card-head">
            <h3><i class="icofont-box"></i> Assets & Sports Equipment Records</h3>
          </div>
          <div class="card-body">
            <div v-if="equipment.length === 0" class="empty-table-msg">
              <p>No sports equipment or asset allocations recorded.</p>
            </div>
            <div v-else class="table-container">
              <table class="profile-data-table">
                <thead>
                  <tr>
                    <th>Item Name / Asset</th>
                    <th>Unit</th>
                    <th>Quantity Received</th>
                    <th>Quantity Distributed</th>
                    <th>Remaining Balance</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="eq in equipment" :key="eq.id">
                    <td class="font-semibold">{{ eq.item_name }}</td>
                    <td>{{ eq.unit || 'PIECES' }}</td>
                    <td>{{ eq.quantity_received || 0 }}</td>
                    <td>{{ eq.quantity_distributed || 0 }}</td>
                    <td class="font-bold">{{ (eq.quantity_received || 0) - (eq.quantity_distributed || 0) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { listNsmisDomain } from '@/api/nsmis.js'
import apiClient from '@/api/client.js'

const props = defineProps({
  federationId: {
    type: String,
    required: true,
  },
  initialFederation: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['navigate', 'edit-federation'])

const currentTab = ref('general')
const memberSubTab = ref('officers')
const loading = ref(true)
const federation = ref(props.initialFederation || null)

// Related lists
const officers = ref([])
const coaches = ref([])
const officials = ref([])
const athletes = ref([])
const disbursements = ref([])
const financialReports = ref([])
const equipment = ref([])

const registrationNumber = computed(() => {
  if (federation.value?.ncs_registration_number) return federation.value.ncs_registration_number
  const code = (federation.value?.abbreviation || federation.value?.acronym || federation.value?.id || 'FED').toUpperCase()
  return `NCS-REG-${code.replace(/^ASSOC_/, '')}`
})

const statusClass = computed(() => {
  const status = String(federation.value?.recognition_status || '').toUpperCase()
  if (status === 'RECOGNISED' || federation.value?.is_active) return 'success'
  if (status === 'PENDING') return 'warning'
  return 'danger'
})

const totalMembersCount = computed(() => {
  return officers.value.length + coaches.value.length + officials.value.length + athletes.value.length
})

const totalDisbursed = computed(() => {
  return disbursements.value.reduce((acc, curr) => acc + Number(curr.amount || 0), 0)
})

const totalExpenses = computed(() => {
  return financialReports.value.reduce((acc, curr) => {
    return acc + Number(curr.competitions || 0) + Number(curr.training || 0) + Number(curr.equipment || 0) + Number(curr.administration || 0)
  }, 0)
})

const totalEquipmentReceived = computed(() => {
  return equipment.value.reduce((acc, curr) => acc + Number(curr.quantity_received || 0), 0)
})

function formatWebsite(url) {
  if (!url) return ''
  return url.replace(/^https?:\/\//, '').replace(/\/$/, '')
}

function formatCurrency(val) {
  return new Intl.NumberFormat('en-UG').format(Number(val || 0))
}

function extractItems(res) {
  if (!res || !res.data) return []
  if (Array.isArray(res.data.items)) return res.data.items
  if (Array.isArray(res.data.data?.items)) return res.data.data.items
  if (Array.isArray(res.data.data)) return res.data.data
  if (Array.isArray(res.data)) return res.data
  return []
}

async function loadFederationData() {
  loading.value = true
  try {
    const fedId = props.federationId

    // 1. Load federation details if not passed
    if (!federation.value || federation.value.id !== fedId) {
      try {
        const fedsRes = await apiClient.get('/api/v1/admin/cms/associations', { params: { active: false } })
        const allFeds = extractItems(fedsRes)
        const match = allFeds.find(f => f.id === fedId || f.slug === fedId)
        if (match) federation.value = match
      } catch (e) {
        console.warn('Could not fetch federation details from associations:', e)
      }
    }

    // 2. Load all domain records in parallel
    const [offRes, coachRes, offiRes, athRes, disbRes, finRes, eqRes] = await Promise.allSettled([
      listNsmisDomain('federation-officers', { per_page: 100 }),
      listNsmisDomain('coaches', { per_page: 100 }),
      listNsmisDomain('technical-officials', { per_page: 100 }),
      listNsmisDomain('athletes', { per_page: 100 }),
      listNsmisDomain('disbursements', { per_page: 100 }),
      listNsmisDomain('accountabilities', { per_page: 100 }),
      listNsmisDomain('equipment', { per_page: 100 }),
    ])

    if (offRes.status === 'fulfilled') {
      const all = extractItems(offRes.value)
      officers.value = all.filter(item => item.federation_id === fedId)
    }
    if (coachRes.status === 'fulfilled') {
      const all = extractItems(coachRes.value)
      coaches.value = all.filter(item => item.federation_id === fedId)
    }
    if (offiRes.status === 'fulfilled') {
      const all = extractItems(offiRes.value)
      officials.value = all.filter(item => item.federation_id === fedId)
    }
    if (athRes.status === 'fulfilled') {
      const all = extractItems(athRes.value)
      athletes.value = all.filter(item => item.federation_id === fedId)
    }
    if (disbRes.status === 'fulfilled') {
      const all = extractItems(disbRes.value)
      disbursements.value = all.filter(item => item.federation_id === fedId)
    }
    if (finRes.status === 'fulfilled') {
      const all = extractItems(finRes.value)
      financialReports.value = all.filter(item => item.federation_id === fedId)
    }
    if (eqRes.status === 'fulfilled') {
      const all = extractItems(eqRes.value)
      equipment.value = all.filter(item => item.federation_id === fedId)
    }
  } catch (err) {
    console.error('Error loading federation profile:', err)
  } finally {
    loading.value = false
  }
}

watch(() => props.federationId, () => {
  loadFederationData()
})

onMounted(() => {
  loadFederationData()
})
</script>

<style scoped>
.federation-profile-page {
  padding: 0 4px 40px;
  font-family: inherit;
}

/* Top Bar */
.page-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.breadcrumb-trail {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #64748b;
}

.crumb-link {
  background: none;
  border: none;
  color: #6777ef;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.crumb-link:hover {
  text-decoration: underline;
}

.crumb-separator {
  color: #cbd5e1;
}

.crumb-current {
  font-weight: 700;
  color: #1e293b;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #334155;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-back:hover {
  background: #e2e8f0;
}

.btn-edit-fed {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  background: #6777ef;
  border: 1px solid #6777ef;
  color: #ffffff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(103, 119, 239, 0.3);
  transition: all 0.15s ease;
}

.btn-edit-fed:hover {
  background: #5566de;
}

/* Loading & Error */
.loading-state-card,
.error-state-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  text-align: center;
  gap: 12px;
}

.loader-icon {
  font-size: 36px;
  color: #6777ef;
}

.error-icon {
  font-size: 40px;
  color: #ef4444;
}

/* Hero Identity Card */
.fed-hero-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.hero-identity {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.logo-wrapper {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}

.fed-logo-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.fed-logo-fallback {
  font-size: 20px;
  font-weight: 900;
  color: #6777ef;
  letter-spacing: 0.05em;
}

.hero-details {
  flex: 1;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 6px;
}

.fed-title {
  font-size: 22px;
  font-weight: 800;
  color: #1e293b;
  margin: 0;
}

.badge-acronym {
  padding: 3px 8px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 800;
}

.badge-status {
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
}

.badge-status.success {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid #a7f3d0;
}

.badge-status.warning {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fde68a;
}

.badge-status.danger {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.fed-subtitle {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #64748b;
  margin: 0 0 14px 0;
  flex-wrap: wrap;
}

.meta-separator {
  color: #cbd5e1;
}

.quick-contact-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.contact-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 11px;
  color: #475569;
}

.contact-chip.link-chip {
  color: #2563eb;
  text-decoration: none;
  font-weight: 600;
}

.contact-chip.link-chip:hover {
  text-decoration: underline;
}

.contact-chip.leader-chip {
  background: #faf5ff;
  border-color: #e9d5ff;
  color: #7e22ce;
}

/* Tabs Header */
.profile-tabs-header {
  display: flex;
  gap: 8px;
  border-bottom: 2px solid #e2e8f0;
  margin-bottom: 20px;
  overflow-x: auto;
}

.tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 18px;
  background: none;
  border: none;
  border-bottom: 3px solid transparent;
  color: #64748b;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
  margin-bottom: -2px;
}

.tab-btn:hover {
  color: #6777ef;
}

.tab-btn.active {
  color: #6777ef;
  border-bottom-color: #6777ef;
  background: #f8fafc;
  border-radius: 8px 8px 0 0;
}

/* Grid layout */
.grid-2-col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 900px) {
  .grid-2-col {
    grid-template-columns: 1fr;
  }
}

/* Section Card */
.profile-section-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.03);
}

.profile-section-card.full-width {
  width: 100%;
}

.card-head {
  padding: 14px 18px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.card-head h3 {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-body {
  padding: 18px;
}

/* Detail list */
.detail-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.detail-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
}

.detail-list li:last-child {
  border-bottom: none;
}

.detail-label {
  color: #64748b;
  font-weight: 500;
}

.detail-val {
  color: #1e293b;
  font-weight: 600;
  text-align: right;
}

.detail-val.font-semibold {
  font-weight: 700;
}

.val-link {
  color: #2563eb;
  text-decoration: none;
}

.val-link:hover {
  text-decoration: underline;
}

.about-text {
  font-size: 13px;
  line-height: 1.6;
  color: #334155;
  margin: 0;
}

/* License Box */
.license-badge-box {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: 10px;
}

.license-big-icon {
  font-size: 40px;
  color: #16a34a;
}

.license-info-col h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 800;
  color: #166534;
}

.license-num {
  font-size: 12px;
  color: #15803d;
  margin: 0 0 4px 0;
}

.license-timeline {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.timeline-row {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  padding: 6px 0;
  border-bottom: 1px dashed #e2e8f0;
}

.timeline-label {
  color: #64748b;
}

.timeline-val {
  font-weight: 600;
  color: #1e293b;
}

/* Score Grid */
.score-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.score-card {
  padding: 12px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  text-align: center;
}

.score-val {
  font-size: 20px;
  font-weight: 900;
  margin-bottom: 2px;
}

.score-title {
  font-size: 11px;
  color: #64748b;
  font-weight: 600;
}

.compliance-checklist {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.check-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #334155;
}

/* Sub tabs */
.sub-tabs-bar {
  display: flex;
  gap: 6px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.sub-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}

.sub-tab-btn:hover {
  background: #e2e8f0;
}

.sub-tab-btn.active {
  background: #6777ef;
  border-color: #6777ef;
  color: #ffffff;
}

.empty-members-box {
  padding: 40px 20px;
  text-align: center;
  background: #ffffff;
  border: 1px dashed #cbd5e1;
  border-radius: 10px;
  color: #64748b;
}

.empty-members-box i {
  font-size: 32px;
  color: #94a3b8;
  margin-bottom: 8px;
}

.empty-members-box h4 {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 4px 0;
}

.empty-members-box p {
  font-size: 12px;
  margin: 0;
}

/* Members Grid */
.members-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 14px;
}

.member-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
}

.member-avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.member-info h4 {
  font-size: 13px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 2px 0;
}

.member-role-tag {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  color: #6777ef;
  margin-bottom: 6px;
}

.member-contact {
  font-size: 11px;
  color: #64748b;
  margin: 2px 0;
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Data Table */
.table-container {
  overflow-x: auto;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
}

.profile-data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.profile-data-table th {
  background: #f8fafc;
  padding: 10px 14px;
  text-align: left;
  font-weight: 700;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

.profile-data-table td {
  padding: 10px 14px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  white-space: nowrap;
}

.badge-tag {
  padding: 2px 6px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  color: #475569;
}

/* KPI Grid */
.records-kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 14px;
}

.kpi-box {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.kpi-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.kpi-icon.col-green { background: #ecfdf5; color: #059669; }
.kpi-icon.col-blue { background: #eff6ff; color: #2563eb; }
.kpi-icon.col-purple { background: #faf5ff; color: #9333ea; }
.kpi-icon.col-orange { background: #fffbeb; color: #d97706; }

.kpi-num {
  font-size: 15px;
  font-weight: 800;
  color: #1e293b;
}

.kpi-title {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
}

.empty-table-msg {
  padding: 24px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}
</style>
