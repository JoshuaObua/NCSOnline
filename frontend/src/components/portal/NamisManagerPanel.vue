<template>
  <div class="namis-container">
    <!-- Main Workspace -->
    <main class="namis-main">
      <!-- Analytics Tab -->
      <section v-if="activeTab === 'analytics'" class="namis-section">
        <header class="section-header">
          <div>
            <p>NCS National Statistics</p>
            <h2>Athlete Analytics Dashboard</h2>
            <span>Demographics, geographical spread, and federation representations.</span>
          </div>
          <button type="button" class="action-btn-secondary" @click="loadAnalytics">
            <i class="icofont-refresh"></i> Refresh
          </button>
        </header>

        <div v-if="loading" class="loading-state">Loading dashboard analytics...</div>
        <div v-else-if="analyticsError" class="error-state">{{ analyticsError }}</div>
        <div v-else class="dashboard-grid">
          <!-- KPI Metrics Row -->
          <div class="kpi-row">
            <div class="kpi-card">
              <i class="icofont-users-alt-5 text-primary"></i>
              <div>
                <h4>{{ analytics.total_registered_athletes || 0 }}</h4>
                <span>Registered Athletes</span>
              </div>
            </div>
            <div class="kpi-card">
              <i class="icofont-female-user text-warning"></i>
              <div>
                <h4>{{ analytics.gender?.FEMALE || 0 }}</h4>
                <span>Female Athletes</span>
              </div>
            </div>
            <div class="kpi-card">
              <i class="icofont-male-user text-success"></i>
              <div>
                <h4>{{ analytics.gender?.MALE || 0 }}</h4>
                <span>Male Athletes</span>
              </div>
            </div>
          </div>

          <!-- Charts Grid -->
          <div class="charts-row">
            <!-- Region stats -->
            <div class="chart-card">
              <h3>Regional Development Statistics</h3>
              <div class="chart-list">
                <div v-for="item in analytics.by_region" :key="item.region" class="chart-item">
                  <div class="item-label">
                    <span>{{ item.region }}</span>
                    <strong>{{ item.count }} athletes</strong>
                  </div>
                  <div class="progress-bar">
                    <div class="progress-fill bg-success" :style="{ width: percentOfTotal(item.count) }"></div>
                  </div>
                </div>
                <p v-if="!analytics.by_region?.length" class="empty-text">No region records compiled.</p>
              </div>
            </div>

            <!-- Federation Stats -->
            <div class="chart-card">
              <h3>Athletes by Federation</h3>
              <div class="chart-list">
                <div v-for="item in analytics.by_federation?.slice(0, 5)" :key="item.federation_id" class="chart-item">
                  <div class="item-label">
                    <span>{{ item.federation_name }}</span>
                    <strong>{{ item.count }} athletes</strong>
                  </div>
                  <div class="progress-bar">
                    <div class="progress-fill bg-primary" :style="{ width: percentOfTotal(item.count) }"></div>
                  </div>
                </div>
                <p v-if="!analytics.by_federation?.length" class="empty-text">No federation records compiled.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Generic Domain Module Tabs -->
      <section v-else class="namis-section">
        <header class="section-header">
          <div>
            <p>NAMIS Registry Manager</p>
            <h2>{{ currentTabLabel }}</h2>
            <span>Create, edit, search, and manage {{ currentTabLabel.toLowerCase() }} entries.</span>
          </div>
          <button type="button" class="action-btn-primary" @click="openCreate">
            <i class="icofont-plus"></i> Add Entry
          </button>
        </header>

        <!-- Filters and Searching -->
        <div class="filter-controls">
          <div class="search-input-wrap">
            <i class="icofont-search-1"></i>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search records..."
              @input="onSearchInput"
            />
          </div>
          <button type="button" class="refresh-btn" @click="loadData">
            <i class="icofont-refresh"></i>
          </button>
        </div>

        <!-- Data table -->
        <div class="table-container">
          <table class="namis-table">
            <thead>
              <tr>
                <th v-for="col in currentColumns" :key="col.key">{{ col.label }}</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading"><td :colspan="currentColumns.length + 1" class="text-center py-4 text-muted">Loading data...</td></tr>
              <tr v-else-if="!items.length"><td :colspan="currentColumns.length + 1" class="text-center py-4 text-muted">No records found matching filters.</td></tr>
              <tr v-for="item in items" v-else :key="item.id">
                <td v-for="col in currentColumns" :key="col.key">
                  <span v-if="col.key === 'athlete_id'">{{ resolveAthleteName(item.athlete_id) }}</span>
                  <span v-else-if="col.key === 'federation_id'">{{ resolveFederationName(item.federation_id) }}</span>
                  <span v-else-if="col.key === 'competition_id'">{{ resolveCompetitionName(item.competition_id) }}</span>
                  <span v-else-if="col.type === 'boolean'">
                    <span :class="item[col.key] ? 'badge bg-success-light text-success' : 'badge bg-danger-light text-danger'">{{ item[col.key] ? 'Yes' : 'No' }}</span>
                  </span>
                  <span v-else-if="col.type === 'date'">{{ formatDate(item[col.key]) }}</span>
                  <span v-else-if="col.type === 'json'">
                    <span class="json-summary" :title="JSON.stringify(item[col.key])">{{ formatJsonField(item[col.key]) }}</span>
                  </span>
                  <span v-else>{{ item[col.key] || '-' }}</span>
                </td>
                <td class="table-actions">
                  <button type="button" class="action-btn edit-btn" title="Edit Entry" @click="openEdit(item)"><i class="icofont-ui-edit"></i></button>
                  <button type="button" class="action-btn delete-btn" title="Delete Entry" @click="deleteItem(item)"><i class="icofont-ui-delete"></i></button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <footer v-if="totalPages > 1" class="namis-pagination">
          <button type="button" :disabled="page === 1" @click="changePage(page - 1)"><i class="icofont-rounded-left"></i> Previous</button>
          <span>Page {{ page }} of {{ totalPages }} ({{ totalCount }} items)</span>
          <button type="button" :disabled="page === totalPages" @click="changePage(page + 1)">Next <i class="icofont-rounded-right"></i></button>
        </footer>
      </section>
    </main>

    <!-- Slide-out Modal / Drawer Form -->
    <div v-if="showModal" class="namis-overlay" @click.self="showModal = false">
      <div class="namis-drawer">
        <header class="drawer-header">
          <h2>{{ editId ? 'Modify Record' : 'Create New Record' }}</h2>
          <button type="button" class="close-btn" @click="showModal = false">&times;</button>
        </header>
        <form class="drawer-form" @submit.prevent="submitForm">
          <div v-for="field in currentFields" :key="field.key" class="form-group">
            <label>{{ field.label }}</label>
            
            <select v-if="field.type === 'select'" v-model="formPayload[field.key]" class="form-control" :required="field.required">
              <option value="">Select option</option>
              <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
            </select>

            <select v-else-if="field.type === 'federation_select'" v-model="formPayload[field.key]" class="form-control" :required="field.required">
              <option value="">Select Sports Federation</option>
              <option v-for="fed in federationsList" :key="fed.id" :value="fed.id">{{ fed.name }} ({{ fed.acronym }})</option>
            </select>

            <select v-else-if="field.type === 'athlete_select'" v-model="formPayload[field.key]" class="form-control" :required="field.required">
              <option value="">Select Registered Athlete</option>
              <option v-for="ath in athletesList" :key="ath.id" :value="ath.id">{{ ath.full_name }} ({{ ath.athlete_number }})</option>
            </select>

            <select v-else-if="field.type === 'competition_select'" v-model="formPayload[field.key]" class="form-control" :required="field.required">
              <option value="">Select Competition</option>
              <option v-for="comp in competitionsList" :key="comp.id" :value="comp.id">{{ comp.name }} - {{ comp.venue }}</option>
            </select>

            <div v-else-if="field.type === 'boolean'" class="checkbox-container">
              <input type="checkbox" v-model="formPayload[field.key]" />
              <span>Enable / Active Protection</span>
            </div>

            <input v-else-if="field.type === 'date'" type="date" v-model="formPayload[field.key]" class="form-control" :required="field.required" />
            <input v-else-if="field.type === 'number'" type="number" v-model.number="formPayload[field.key]" class="form-control" :required="field.required" />
            
            <textarea v-else-if="field.type === 'textarea'" v-model="formPayload[field.key]" class="form-control text-area" :required="field.required"></textarea>
            
            <textarea v-else-if="field.type === 'json'" :value="stringifyJson(formPayload[field.key])" @input="updateJsonField(field.key, $event.target.value)" class="form-control text-area json-input" placeholder='{"key": "value"}' :required="field.required"></textarea>

            <input v-else type="text" v-model="formPayload[field.key]" class="form-control" :required="field.required" />
          </div>
          
          <p v-if="formError" class="form-error">{{ formError }}</p>
          
          <footer class="drawer-actions">
            <button type="button" class="btn-cancel" @click="showModal = false">Cancel</button>
            <button type="submit" class="btn-submit" :disabled="formSaving">
              {{ formSaving ? 'Saving...' : 'Save Record' }}
            </button>
          </footer>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import {
  listNsmisDomain,
  createNsmisDomain,
  updateNsmisDomain,
  deleteNsmisDomain,
  getAthleteDashboard
} from '@/api/nsmis.js'

const props = defineProps({
  userScope: { type: Array, default: () => [] },
  tab: { type: String, default: 'analytics' }
})

const tabs = [
  { id: 'analytics', label: 'Dashboard & KPIs', icon: 'icofont-chart-bar-graph' },
  { id: 'athletes', label: 'Athletes Registry', icon: 'icofont-users-alt-2' },
  { id: 'clubs', label: 'Clubs & Academies', icon: 'icofont-home' },
  { id: 'coaches', label: 'Coaches Registry', icon: 'icofont-businessman' },
  { id: 'competitions', label: 'Competitions Logs', icon: 'icofont-runner-alt-1' },
  { id: 'competition-results', label: 'Athlete Results', icon: 'icofont-listing-number' },
  { id: 'medals', label: 'Medals Standings', icon: 'icofont-medal' },
  { id: 'talent', label: 'Talent Pathways', icon: 'icofont-bulb-alt' },
  { id: 'national-team', label: 'National Squads', icon: 'icofont-flag' },
  { id: 'technical-officials', label: 'Technical Officials', icon: 'icofont-referee' },
  { id: 'medical-records', label: 'Medical Files', icon: 'icofont-first-aid' },
  { id: 'safeguarding-records', label: 'Safeguarding Records', icon: 'icofont-shield' },
  { id: 'anti-doping', label: 'Anti-Doping Compliance', icon: 'icofont-test-bulb' }
]

const activeTab = ref(props.tab || 'analytics')

watch(() => props.tab, (newVal) => {
  if (newVal) {
    activeTab.value = newVal
  }
})
const loading = ref(false)
const items = ref([])
const page = ref(1)
const perPage = ref(20)
const totalCount = ref(0)
const searchQuery = ref('')
const searchTimeout = ref(null)

const showModal = ref(false)
const editId = ref('')
const formPayload = ref({})
const formSaving = ref(false)
const formError = ref('')

// Analytics States
const analytics = ref({})
const analyticsError = ref('')

// Reference Select Options Lists
const federationsList = ref([])
const athletesList = ref([])
const competitionsList = ref([])

const currentTabLabel = computed(() => tabs.find(t => t.id === activeTab.value)?.label || '')
const totalPages = computed(() => Math.ceil(totalCount.value / perPage.value))

// Columns mappings
const columnsConfig = {
  athletes: [
    { key: 'athlete_number', label: 'Athlete ID' },
    { key: 'full_name', label: 'Full Name' },
    { key: 'gender', label: 'Gender' },
    { key: 'date_of_birth', label: 'DOB', type: 'date' },
    { key: 'discipline', label: 'Discipline' },
    { key: 'status', label: 'Status' }
  ],
  clubs: [
    { key: 'name', label: 'Name' },
    { key: 'acronym', label: 'Acronym' },
    { key: 'federation_id', label: 'Federation' },
    { key: 'district', label: 'District' },
    { key: 'status', label: 'Status' }
  ],
  coaches: [
    { key: 'full_name', label: 'Full Name' },
    { key: 'certification_level', label: 'Certification' },
    { key: 'license_number', label: 'License' },
    { key: 'expiry_date', label: 'Expires', type: 'date' },
    { key: 'status', label: 'Status' }
  ],
  competitions: [
    { key: 'name', label: 'Name' },
    { key: 'level', label: 'Level' },
    { key: 'venue', label: 'Venue' },
    { key: 'starts_on', label: 'Starts', type: 'date' },
    { key: 'ends_on', label: 'Ends', type: 'date' }
  ],
  'competition-results': [
    { key: 'competition_id', label: 'Competition' },
    { key: 'athlete_name', label: 'Athlete' },
    { key: 'event', label: 'Event' },
    { key: 'position', label: 'Finish' },
    { key: 'is_national_record', label: 'NR', type: 'boolean' }
  ],
  medals: [
    { key: 'athlete_name', label: 'Athlete' },
    { key: 'event', label: 'Event' },
    { key: 'medal_type', label: 'Medal' },
    { key: 'won_on', label: 'Date Won', type: 'date' },
    { key: 'status', label: 'Status' }
  ],
  talent: [
    { key: 'athlete_name', label: 'Athlete' },
    { key: 'age_at_identification', label: 'Age Identified' },
    { key: 'district', label: 'District' },
    { key: 'identified_on', label: 'Scouted On', type: 'date' },
    { key: 'status', label: 'Status' }
  ],
  'national-team': [
    { key: 'athlete_id', label: 'Athlete' },
    { key: 'team_name', label: 'Squad' },
    { key: 'category', label: 'Tier' },
    { key: 'appearances_count', label: 'Caps' }
  ],
  'technical-officials': [
    { key: 'full_name', label: 'Full Name' },
    { key: 'official_type', label: 'Type' },
    { key: 'level', label: 'Tier' },
    { key: 'status', label: 'Status' }
  ],
  'medical-records': [
    { key: 'athlete_id', label: 'Athlete' },
    { key: 'blood_group', label: 'Blood Type' },
    { key: 'current_injury_status', label: 'Injury Status' }
  ],
  'safeguarding-records': [
    { key: 'athlete_id', label: 'Athlete' },
    { key: 'guardian_details', label: 'Guardian Details', type: 'json' },
    { key: 'anti_doping_education_completed', label: 'Doping Ed.', type: 'boolean' }
  ],
  'anti-doping': [
    { key: 'athlete_id', label: 'Athlete' },
    { key: 'testing_status', label: 'Pool Status' },
    { key: 'wada_education_completed', label: 'WADA Complete', type: 'boolean' }
  ]
}

const currentColumns = computed(() => columnsConfig[activeTab.value] || [])

// Forms mapping layout
const fieldsConfig = {
  athletes: [
    { key: 'full_name', label: 'Full Name', required: true },
    { key: 'gender', label: 'Gender', type: 'select', options: ['MALE', 'FEMALE', 'OTHER', 'NOT_STATED'], required: true },
    { key: 'date_of_birth', label: 'Date of Birth', type: 'date', required: true },
    { key: 'national_id_passport', label: 'National ID/Passport' },
    { key: 'district', label: 'District of Origin' },
    { key: 'region', label: 'Region', type: 'select', options: ['North', 'East', 'Central', 'West'] },
    { key: 'club', label: 'Club Name' },
    { key: 'discipline', label: 'Discipline/Sport', required: true },
    { key: 'national_team_status', label: 'National Team Status', type: 'select', options: ['NO', 'DEVELOPMENT', 'SENIOR', 'FORMER'] },
    { key: 'status', label: 'Status', type: 'select', options: ['ACTIVE', 'INACTIVE', 'SUSPENDED', 'RETIRED'] },
    { key: 'consent_basis', label: 'Consent Basis (e.g. FEDERATION_MANDATE)' },
    { key: 'age_category', label: 'Age Category', type: 'select', options: ['U10', 'U12', 'U15', 'U17', 'U20', 'Senior'] },
    { key: 'phone_contact', label: 'Phone Number' },
    { key: 'email_address', label: 'Email Address' },
    { key: 'next_of_kin', label: 'Next of Kin' },
    { key: 'emergency_contact', label: 'Emergency Contact' },
    { key: 'education_institution', label: 'School/Institution' },
    { key: 'highest_education_level', label: 'Highest Education Level' },
    { key: 'sports_scholarship_status', label: 'Sports Scholarship Status', type: 'boolean' },
    { key: 'current_occupation', label: 'Current Occupation' }
  ],
  clubs: [
    { key: 'federation_id', label: 'Federation', type: 'federation_select', required: true },
    { key: 'name', label: 'Club Name', required: true },
    { key: 'acronym', label: 'Acronym (Abbrev.)', required: true },
    { key: 'contact_person', label: 'Contact Person' },
    { key: 'email', label: 'Contact Email' },
    { key: 'phone', label: 'Contact Phone' },
    { key: 'district', label: 'District' },
    { key: 'region', label: 'Region', type: 'select', options: ['North', 'East', 'Central', 'West'] },
    { key: 'date_founded', label: 'Date Founded', type: 'date' },
    { key: 'status', label: 'Status', type: 'select', options: ['ACTIVE', 'INACTIVE', 'SUSPENDED'] }
  ],
  coaches: [
    { key: 'federation_id', label: 'Federation', type: 'federation_select', required: true },
    { key: 'full_name', label: 'Full Name', required: true },
    { key: 'certification_level', label: 'Certification Level', required: true },
    { key: 'license_number', label: 'License Number', required: true },
    { key: 'expiry_date', label: 'Expiry Date', type: 'date' },
    { key: 'status', label: 'Status', type: 'select', options: ['ACTIVE', 'INACTIVE', 'SUSPENDED'] },
    { key: 'phone', label: 'Phone' },
    { key: 'email', label: 'Email' }
  ],
  competitions: [
    { key: 'federation_id', label: 'Federation', type: 'federation_select', required: true },
    { key: 'name', label: 'Competition Name', required: true },
    { key: 'venue', label: 'Venue Location', required: true },
    { key: 'host_country', label: 'Host Country' },
    { key: 'level', label: 'Competition Level', type: 'select', options: ['DISTRICT', 'REGIONAL', 'NATIONAL', 'EAST_AFRICAN', 'AFRICAN', 'COMMONWEALTH', 'OLYMPIC', 'WORLD_CHAMPIONSHIP', 'CONTINENTAL', 'INTERNATIONAL'], required: true },
    { key: 'starts_on', label: 'Starts On', type: 'date', required: true },
    { key: 'ends_on', label: 'Ends On', type: 'date', required: true },
    { key: 'status', label: 'Status', type: 'select', options: ['DRAFT', 'SUBMITTED', 'APPROVED', 'CANCELLED'] }
  ],
  'competition-results': [
    { key: 'competition_id', label: 'Competition', type: 'competition_select', required: true },
    { key: 'athlete_id', label: 'Registered Athlete (Optional)', type: 'athlete_select' },
    { key: 'athlete_name', label: 'Athlete Name', required: true },
    { key: 'event', label: 'Event Name', required: true },
    { key: 'position', label: 'Finished Position', type: 'number' },
    { key: 'time_result', label: 'Time Result (e.g. 9.58s)' },
    { key: 'distance_result', label: 'Distance Result (e.g. 8.21m)' },
    { key: 'weight_result', label: 'Weight Result (e.g. 105kg)' },
    { key: 'score_result', label: 'Score Result (e.g. 98.4)' },
    { key: 'ranking_result', label: 'Ranking Result' },
    { key: 'is_national_record', label: 'Is National Record', type: 'boolean' },
    { key: 'is_personal_best', label: 'Is Personal Best', type: 'boolean' },
    { key: 'is_seasonal_best', label: 'Is Seasonal Best', type: 'boolean' }
  ],
  medals: [
    { key: 'competition_id', label: 'Competition', type: 'competition_select', required: true },
    { key: 'federation_id', label: 'Federation', type: 'federation_select', required: true },
    { key: 'athlete_id', label: 'Registered Athlete (Optional)', type: 'athlete_select' },
    { key: 'athlete_name', label: 'Athlete Name', required: true },
    { key: 'event', label: 'Event Name (e.g. 100m Dash)', required: true },
    { key: 'country', label: 'Host Country', required: true },
    { key: 'won_on', label: 'Date Won', type: 'date', required: true },
    { key: 'medal_type', label: 'Medal Type', type: 'select', options: ['GOLD', 'SILVER', 'BRONZE'], required: true },
    { key: 'coach_responsible', label: 'Coach Responsible' },
    { key: 'team_manager', label: 'Team Manager' },
    { key: 'funding_source', label: 'Funding Source' },
    { key: 'status', label: 'Status', type: 'select', options: ['DRAFT', 'SUBMITTED', 'APPROVED', 'VOIDED'] },
    { key: 'level', label: 'Medal Tier/Level', type: 'select', options: ['DISTRICT', 'REGIONAL', 'NATIONAL', 'EAST_AFRICAN', 'AFRICAN', 'COMMONWEALTH', 'OLYMPIC', 'WORLD_CHAMPIONSHIP', 'CONTINENTAL', 'INTERNATIONAL'] },
    { key: 'coach_at_win_id', label: 'Coach Entourage ID (Optional)' },
    { key: 'is_team_event', label: 'Is Team Event', type: 'boolean' },
    { key: 'prize_money', label: 'Prize Money (UGX)', type: 'number' },
    { key: 'ncs_recognition_status', label: 'NCS Recognition', type: 'select', options: ['PENDING', 'APPROVED', 'REJECTED'] }
  ],
  talent: [
    { key: 'federation_id', label: 'Federation', type: 'federation_select', required: true },
    { key: 'athlete_id', label: 'Registered Athlete (Optional)', type: 'athlete_select' },
    { key: 'athlete_name', label: 'Athlete Name', required: true },
    { key: 'age_at_identification', label: 'Age at Identification', type: 'number', required: true },
    { key: 'school', label: 'School/Institution' },
    { key: 'district', label: 'Scouting District', required: true },
    { key: 'region', label: 'Region', type: 'select', options: ['North', 'East', 'Central', 'West'] },
    { key: 'identified_by', label: 'Scout Name', required: true },
    { key: 'identified_on', label: 'Scouting Date', type: 'date', required: true },
    { key: 'status', label: 'Status', type: 'select', options: ['IDENTIFIED', 'ASSESSED', 'SELECTED', 'NATIONAL_TEAM', 'CLOSED'] },
    { key: 'talent_centre', label: 'Talent Centre' },
    { key: 'talent_category', label: 'Talent Category', type: 'select', options: ['AMATEUR', 'EMERGING', 'ELITE'] },
    { key: 'recommended_pathway', label: 'Recommended Pathway' },
    { key: 'scholarship_status', label: 'Scholarship Status', type: 'select', options: ['NONE', 'APPLIED', 'ACTIVE', 'EXPIRED'] }
  ],
  'national-team': [
    { key: 'athlete_id', label: 'Registered Athlete', type: 'athlete_select', required: true },
    { key: 'team_name', label: 'Squad/Team Name', required: true },
    { key: 'category', label: 'Tier', type: 'select', options: ['SENIOR', 'DEVELOPMENT', 'JUNIOR'], required: true },
    { key: 'first_call_up_on', label: 'First Call Up Date', type: 'date' },
    { key: 'last_appearance_on', label: 'Last Appearance Date', type: 'date' },
    { key: 'appearances_count', label: 'Total Caps Count', type: 'number' },
    { key: 'notes', label: 'Technical Notes', type: 'textarea' }
  ],
  'technical-officials': [
    { key: 'federation_id', label: 'Federation', type: 'federation_select', required: true },
    { key: 'full_name', label: 'Full Name', required: true },
    { key: 'official_type', label: 'Official Type', type: 'select', options: ['REFEREE', 'UMPIRE', 'JUDGE', 'ASSESSOR', 'OTHER'], required: true },
    { key: 'level', label: 'Certification Tier', required: true },
    { key: 'certification', label: 'Credential Name', required: true },
    { key: 'valid_until', label: 'Valid Until', type: 'date' },
    { key: 'status', label: 'Status', type: 'select', options: ['ACTIVE', 'INACTIVE', 'SUSPENDED'] }
  ],
  'medical-records': [
    { key: 'athlete_id', label: 'Athlete', type: 'athlete_select', required: true },
    { key: 'blood_group', label: 'Blood Group' },
    { key: 'allergies', label: 'Allergies Description' },
    { key: 'injury_history', label: 'Injury History (comma-separated)', type: 'textarea' },
    { key: 'current_injury_status', label: 'Injury Status', type: 'select', options: ['FIT', 'INJURED', 'RECOVERING'] },
    { key: 'medical_insurance', label: 'Medical Insurance Details' }
  ],
  'safeguarding-records': [
    { key: 'athlete_id', label: 'Athlete', type: 'athlete_select', required: true },
    { key: 'guardian_details', label: 'Guardian Details (JSON)', type: 'json' },
    { key: 'manager_details', label: 'Manager Details (JSON)', type: 'json' },
    { key: 'consent_forms_url', label: 'Consent Forms Link' },
    { key: 'anti_doping_education_completed', label: 'Anti Doping Education Completed', type: 'boolean' }
  ],
  'anti-doping': [
    { key: 'athlete_id', label: 'Athlete', type: 'athlete_select', required: true },
    { key: 'testing_status', label: 'Testing Pool Status', type: 'select', options: ['NOT_TESTED', 'IN_POOL', 'TESTED'] },
    { key: 'last_tested_on', label: 'Last Tested Date', type: 'date' },
    { key: 'last_test_result', label: 'Last Test Result' },
    { key: 'wada_education_completed', label: 'WADA Course Completed', type: 'boolean' },
    { key: 'suspension_history', label: 'Suspension History Logs', type: 'textarea' }
  ]
}

const currentFields = computed(() => fieldsConfig[activeTab.value] || [])

watch(activeTab, () => {
  page.value = 1
  searchQuery.value = ''
  items.value = []
  if (activeTab.value === 'analytics') {
    loadAnalytics()
  } else {
    loadData()
  }
})

onMounted(() => {
  loadAnalytics()
  loadSelectionReferences()
})

async function loadAnalytics() {
  loading.value = true
  analyticsError.value = ''
  try {
    const res = await getAthleteDashboard()
    analytics.value = res.data?.data || res.data || {}
  } catch (err) {
    analyticsError.value = 'Failed to load athlete stats'
  } finally {
    loading.value = false
  }
}

async function loadSelectionReferences() {
  try {
    const [feds, aths, comps] = await Promise.all([
      listNsmisDomain('federations', { per_page: 100 }),
      listNsmisDomain('athletes', { per_page: 200 }),
      listNsmisDomain('competitions', { per_page: 100 })
    ])
    federationsList.value = feds.data?.data?.items || feds.data?.items || []
    athletesList.value = aths.data?.data?.items || aths.data?.items || []
    competitionsList.value = comps.data?.data?.items || comps.data?.items || []
  } catch (e) {
    console.error('Could not load drop-down selects reference arrays', e)
  }
}

async function loadData() {
  loading.value = true
  try {
    const res = await listNsmisDomain(activeTab.value, {
      page: page.value,
      per_page: perPage.value,
      search: searchQuery.value.trim()
    })
    const body = res.data?.data || res.data || {}
    items.value = body.items || []
    totalCount.value = body.pagination?.total || items.value.length
  } catch (err) {
    console.error('Failed to load nsmis domain rows', err)
  } finally {
    loading.value = false
  }
}

function selectTab(tabId) {
  activeTab.value = tabId
}

function onSearchInput() {
  clearTimeout(searchTimeout.value)
  searchTimeout.value = setTimeout(() => {
    page.value = 1
    loadData()
  }, 350)
}

function changePage(p) {
  page.value = p
  loadData()
}

function percentOfTotal(count) {
  const total = analytics.value.total_registered_athletes || 1
  return `${Math.min(100, Math.round((count / total) * 100))}%`
}

function formatDate(val) {
  if (!val) return '-'
  const d = new Date(val)
  return isNaN(d.getTime()) ? val : d.toLocaleDateString('en-GB')
}

// Helpers for reference name mappings
function resolveAthleteName(id) {
  return athletesList.value.find(a => a.id === id)?.full_name || id
}

function resolveFederationName(id) {
  const f = federationsList.value.find(x => x.id === id)
  return f ? `${f.name} (${f.acronym})` : id
}

function resolveCompetitionName(id) {
  return competitionsList.value.find(c => c.id === id)?.name || id
}

// JSON formatting
function formatJsonField(val) {
  if (!val) return '-'
  try {
    const obj = typeof val === 'string' ? JSON.parse(val) : val
    return Object.entries(obj).map(([k, v]) => `${k}: ${v}`).join(', ')
  } catch {
    return String(val)
  }
}

function stringifyJson(val) {
  if (!val) return '{}'
  return typeof val === 'object' ? JSON.stringify(val, null, 2) : String(val)
}

function updateJsonField(key, text) {
  try {
    formPayload.value[key] = JSON.parse(text)
  } catch {
    formPayload.value[key] = text
  }
}

function openCreate() {
  formPayload.value = {}
  currentFields.value.forEach(f => {
    formPayload.value[f.key] = f.type === 'boolean' ? false : f.type === 'json' ? {} : ''
  })
  editId.value = ''
  formError.value = ''
  showModal.value = true
}

function openEdit(item) {
  formPayload.value = { ...item }
  // format dates to YYYY-MM-DD for date inputs
  currentFields.value.forEach(f => {
    if (f.type === 'date' && formPayload.value[f.key]) {
      formPayload.value[f.key] = String(formPayload.value[f.key]).slice(0, 10)
    }
  })
  editId.value = item.id
  formError.value = ''
  showModal.value = true
}

async function submitForm() {
  formSaving.value = true
  formError.value = ''
  try {
    if (editId.value) {
      await updateNsmisDomain(activeTab.value, editId.value, formPayload.value)
    } else {
      await createNsmisDomain(activeTab.value, formPayload.value)
    }
    showModal.value = false
    loadData()
    loadSelectionReferences() // Reload lists to refresh options
  } catch (err) {
    formError.value = err.response?.data?.error?.message || 'Could not save record'
  } finally {
    formSaving.value = false
  }
}

async function deleteItem(item) {
  if (!confirm('Are you sure you want to delete this record?')) return
  try {
    await deleteNsmisDomain(activeTab.value, item.id)
    loadData()
    loadSelectionReferences()
  } catch (err) {
    alert(err.response?.data?.error?.message || 'Could not delete item')
  }
}
</script>

<style scoped>
.namis-container {
  display: flex;
  min-height: 580px;
  background: var(--bg-card, #fff);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
  border: 1px solid rgba(0,0,0,0.06);
  overflow: hidden;
}

.namis-sidebar {
  width: 250px;
  background: #f8f9fa;
  border-right: 1px solid rgba(0,0,0,0.06);
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(0,0,0,0.06);
}

.sidebar-header i {
  font-size: 24px;
  color: #6777ef;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 700;
  color: #34395e;
}

.sidebar-header span {
  font-size: 11px;
  color: #98a6ad;
}

.sidebar-nav {
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
  flex: 1;
}

.nav-tab-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 8px;
  background: transparent;
  border: none;
  text-align: left;
  font-size: 13px;
  font-weight: 600;
  color: #555;
  cursor: pointer;
  transition: all 0.2s ease;
}

.nav-tab-btn:hover {
  background: rgba(103,119,239,0.05);
  color: #6777ef;
}

.nav-tab-btn.active {
  background: #6777ef;
  color: #fff;
}

.namis-main {
  flex: 1;
  padding: 24px;
  background: #fff;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.section-header p {
  margin: 0 0 4px;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 800;
  color: #98a6ad;
}

.section-header h2 {
  margin: 0 0 4px;
  font-size: 20px;
  font-weight: 700;
  color: #34395e;
}

.section-header span {
  font-size: 13px;
  color: #6c757d;
}

.action-btn-primary {
  background: #6777ef;
  color: #fff;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: background 0.2s;
}

.action-btn-primary:hover {
  background: #394eea;
}

.action-btn-secondary {
  background: #f8f9fa;
  color: #555;
  border: 1px solid rgba(0,0,0,0.08);
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.2s;
}

.action-btn-secondary:hover {
  background: #e9ecef;
}

/* Filters */
.filter-controls {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.search-input-wrap {
  position: relative;
  flex: 1;
}

.search-input-wrap i {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #98a6ad;
}

.search-input-wrap input {
  width: 100%;
  padding: 8px 12px 8px 36px;
  border-radius: 6px;
  border: 1px solid rgba(0,0,0,0.1);
  font-size: 13px;
  outline: none;
}

.refresh-btn {
  background: #f8f9fa;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 6px;
  width: 38px;
  height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #6c757d;
}

/* Tables */
.table-container {
  border: 1px solid rgba(0,0,0,0.06);
  border-radius: 8px;
  overflow-x: auto;
  flex: 1;
}

.namis-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 13px;
}

.namis-table th {
  background: #f8f9fa;
  padding: 12px 16px;
  font-weight: 700;
  color: #34395e;
  border-bottom: 1px solid rgba(0,0,0,0.06);
}

.namis-table td {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(0,0,0,0.04);
  color: #555;
  vertical-align: middle;
}

.badge {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 700;
}

.bg-success-light { background: rgba(40,167,69,0.1); }
.bg-danger-light { background: rgba(220,53,69,0.1); }

.table-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  border: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 12px;
  transition: opacity 0.2s;
}

.action-btn:hover {
  opacity: 0.8;
}

.edit-btn { background: #3abaf4; color: #fff; }
.delete-btn { background: #fc544b; color: #fff; }

/* Pagination */
.namis-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
  font-size: 13px;
  color: #6c757d;
}

.namis-pagination button {
  background: #f8f9fa;
  border: 1px solid rgba(0,0,0,0.08);
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
}

.namis-pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Modals & Drawer Forms */
.namis-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0,0,0,0.4);
  display: flex;
  justify-content: flex-end;
  z-index: 1050;
}

.namis-drawer {
  width: 480px;
  background: #fff;
  height: 100%;
  box-shadow: -4px 0 20px rgba(0,0,0,0.15);
  display: flex;
  flex-direction: column;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

.drawer-header {
  padding: 20px 24px;
  border-bottom: 1px solid rgba(0,0,0,0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.drawer-header h2 {
  margin: 0;
  font-size: 18px;
  color: #34395e;
}

.close-btn {
  background: transparent;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #98a6ad;
}

.drawer-form {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 12px;
  font-weight: 700;
  color: #34395e;
}

.form-control {
  padding: 8px 12px;
  border: 1px solid rgba(0,0,0,0.12);
  border-radius: 6px;
  font-size: 13px;
  outline: none;
  background: #fff;
}

.checkbox-container {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.text-area {
  min-height: 80px;
  resize: vertical;
}

.json-input {
  font-family: monospace;
  font-size: 11px;
}

.form-error {
  color: #fc544b;
  font-size: 12px;
  margin: 0;
}

.drawer-actions {
  border-top: 1px solid rgba(0,0,0,0.08);
  padding-top: 16px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: auto;
}

.btn-cancel {
  background: #f8f9fa;
  border: 1px solid rgba(0,0,0,0.08);
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
}

.btn-submit {
  background: #6777ef;
  color: #fff;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
}

/* Analytics Dashboard style */
.dashboard-grid {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.kpi-row {
  display: flex;
  gap: 16px;
}

.kpi-card {
  flex: 1;
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 16px;
  border: 1px solid rgba(0,0,0,0.04);
}

.kpi-card i {
  font-size: 32px;
}

.kpi-card h4 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #34395e;
}

.kpi-card span {
  font-size: 12px;
  color: #98a6ad;
}

.charts-row {
  display: flex;
  gap: 16px;
}

.chart-card {
  flex: 1;
  border: 1px solid rgba(0,0,0,0.06);
  border-radius: 8px;
  padding: 20px;
}

.chart-card h3 {
  margin: 0 0 16px;
  font-size: 14px;
  font-weight: 700;
  color: #34395e;
  border-bottom: 1px solid rgba(0,0,0,0.04);
  padding-bottom: 8px;
}

.chart-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.chart-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.item-label {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.progress-bar {
  background: #e9ecef;
  height: 8px;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
}

.empty-text {
  font-size: 13px;
  color: #98a6ad;
  text-align: center;
  margin: 20px 0;
}
</style>
