<template>
  <ExpensePage title="Expenses Register" description="Comprehensive expense tracking, verification, approval, and management workspace.">
    <template #actions>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <router-link to="/expenses/reports" class="btn btn-sm btn-outline-info" title="Dedicated Expense Reports Page">
          <i class="fas fa-chart-line mr-1"></i> Reports
        </router-link>
        <button @click="exportExcel" class="btn btn-sm btn-outline-success" title="Export to Excel">
          <i class="fas fa-file-excel mr-1"></i> Export Excel
        </button>
        <button @click="exportPDF" class="btn btn-sm btn-outline-danger" title="Export / Print PDF">
          <i class="fas fa-file-pdf mr-1"></i> Export PDF
        </button>
        <router-link v-if="auth.isAdmin" to="/expenses/categories" class="btn btn-sm btn-outline-primary">
          <i class="fas fa-tags mr-1"></i> Categories
        </router-link>
        <router-link to="/expenses/new" class="btn btn-sm btn-emerald text-white">
          <i class="fas fa-plus mr-1"></i> Record Expense
        </router-link>
      </div>
    </template>

    <div v-if="error" role="alert" class="alert alert-danger alert-dismissible fade show">
      <i class="fas fa-exclamation-triangle mr-2"></i> {{ error }}
      <button type="button" class="close" @click="error=''" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
    </div>

    <div v-if="successMsg" role="alert" class="alert alert-success alert-dismissible fade show">
      <i class="fas fa-check-circle mr-2"></i> {{ successMsg }}
      <button type="button" class="close" @click="successMsg=''" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
    </div>

    <!-- Structured Filter Card -->
    <div class="card mb-4 shadow-sm border-0">
      <div class="card-header bg-white border-bottom py-3 d-flex justify-content-between align-items-center">
        <h4 class="card-title m-0 text-primary font-weight-bold text-sm">
          <i class="fas fa-filter mr-2"></i> Filter Expenses
        </h4>
        <div class="card-header-action">
          <button type="button" @click="reset" class="btn btn-sm btn-light text-muted mr-2" title="Reset Filters">
            <i class="fas fa-redo mr-1"></i> Reset
          </button>
          <button type="submit" form="expense-filter-form" :disabled="loading" class="btn btn-sm btn-primary" title="Apply Filters">
            <i class="fas fa-search mr-1"></i> Filter
          </button>
        </div>
      </div>
      <div class="card-body p-3 bg-light">
        <form id="expense-filter-form" @submit.prevent="offset=0;load()">
          <div class="row">
            <div class="col-md-4 col-sm-6 mb-3">
              <label class="font-weight-bold text-xs text-uppercase text-muted mb-1">Search Query</label>
              <div class="input-group input-group-sm">
                <div class="input-group-prepend">
                  <span class="input-group-text bg-white border-right-0"><i class="fas fa-search text-muted"></i></span>
                </div>
                <input v-model="filters.search" class="form-control form-control-sm border-left-0" placeholder="Title, reference, payee..." />
              </div>
            </div>
            <div class="col-md-4 col-sm-6 mb-3">
              <label class="font-weight-bold text-xs text-uppercase text-muted mb-1">Category</label>
              <select aria-label="Category" v-model="filters.category_id" class="form-control form-control-sm">
                <option value="">All Categories</option>
                <option v-for="c in categories" :value="c.id" :key="c.id">{{ c.name }}</option>
              </select>
            </div>
            <div class="col-md-4 col-sm-6 mb-3">
              <label class="font-weight-bold text-xs text-uppercase text-muted mb-1">Department</label>
              <select aria-label="Department" v-model="filters.department_id" class="form-control form-control-sm">
                <option value="">All Departments</option>
                <option v-for="d in departments" :value="d.id" :key="d.id">{{ d.name }}</option>
              </select>
            </div>
            <div class="col-md-4 col-sm-6 mb-2">
              <label class="font-weight-bold text-xs text-uppercase text-muted mb-1">From Date</label>
              <input v-model="filters.from" type="date" class="form-control form-control-sm" />
            </div>
            <div class="col-md-4 col-sm-6 mb-2">
              <label class="font-weight-bold text-xs text-uppercase text-muted mb-1">To Date</label>
              <input v-model="filters.to" type="date" class="form-control form-control-sm" />
            </div>
            <div class="col-md-4 col-sm-12 mb-2 d-flex align-items-end">
              <div class="btn-group btn-group-sm w-100">
                <button type="submit" :disabled="loading" class="btn btn-primary font-weight-bold">
                  <i class="fas fa-filter mr-1"></i> Apply Filter
                </button>
                <button type="button" @click="reset" class="btn btn-outline-secondary">
                  <i class="fas fa-undo mr-1"></i> Clear
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- Summary Statistics Bar -->
    <div class="row mb-4">
      <div class="col-md-6 col-sm-12 mb-3 mb-md-0">
        <div class="card card-statistic-1 mb-0 border shadow-sm">
          <div class="card-icon bg-primary text-white d-flex align-items-center justify-content-center rounded-left" style="width: 60px;">
            <i class="fas fa-receipt fa-2x"></i>
          </div>
          <div class="card-wrap p-3 bg-white rounded-right">
            <div class="card-header p-0 bg-transparent border-0">
              <h4 class="text-xs font-weight-bold text-uppercase text-muted m-0">Total Expenses Found</h4>
            </div>
            <div class="card-body p-0 mt-1">
              <span class="h4 font-weight-bold text-dark mb-0">{{ total }}</span>
              <span class="text-muted text-xs ml-2">records matching criteria</span>
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-6 col-sm-12">
        <div class="card card-statistic-1 mb-0 border shadow-sm">
          <div class="card-icon bg-success text-white d-flex align-items-center justify-content-center rounded-left" style="width: 60px;">
            <i class="fas fa-coins fa-2x"></i>
          </div>
          <div class="card-wrap p-3 bg-white rounded-right">
            <div class="card-header p-0 bg-transparent border-0">
              <h4 class="text-xs font-weight-bold text-uppercase text-muted m-0">Aggregate Total Value</h4>
            </div>
            <div class="card-body p-0 mt-1">
              <span class="h4 font-weight-bold text-success mb-0">UGX {{ expenseMoney(totalAmount) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Otika State Save DataTable Container -->
    <div class="row">
      <div class="col-12">
        <div class="card shadow-sm border-0">
          <div class="card-header bg-white border-bottom d-flex justify-content-between align-items-center py-3">
            <h4 class="card-title m-0 text-primary font-weight-bold">
              <i class="fas fa-list-alt mr-2"></i> Expenses Register
            </h4>
            <div class="card-header-action">
              <span class="badge badge-light text-muted">Page {{ currentPage }} of {{ totalPages || 1 }}</span>
            </div>
          </div>
          <div class="card-body p-3">
            <div class="table-responsive">
              <div id="save-stage_wrapper" class="dataTables_wrapper container-fluid dt-bootstrap4 no-footer p-3">
                
                <!-- Length Selector & Search Row -->
                <div class="row mb-3 align-items-center">
                  <div class="col-sm-12 col-md-6 mb-2 mb-md-0">
                    <div class="dataTables_length" id="save-stage_length">
                      <label class="d-inline-flex align-items-center text-xs font-weight-bold text-muted mb-0">
                        Show 
                        <select name="save-stage_length" v-model="limit" @change="offset=0;load()" aria-controls="save-stage" class="form-control form-control-sm mx-2" style="width: auto;">
                          <option :value="10">10</option>
                          <option :value="25">25</option>
                          <option :value="50">50</option>
                          <option :value="100">100</option>
                        </select>
                        entries
                      </label>
                    </div>
                  </div>
                  <div class="col-sm-12 col-md-6 text-md-right">
                    <div id="save-stage_filter" class="dataTables_filter">
                      <label class="d-inline-flex align-items-center text-xs font-weight-bold text-muted mb-0">
                        Search:
                        <input type="search" v-model="filters.search" @input="debouncedSearch" class="form-control form-control-sm ml-2" placeholder="Search reference, title, payee..." aria-controls="save-stage" style="width: 200px;">
                      </label>
                    </div>
                  </div>
                </div>

                <!-- Main Data Table -->
                <div class="row">
                  <div class="col-sm-12">
                    <table class="table table-striped table-hover dataTable no-footer w-100" id="save-stage" role="grid" aria-describedby="save-stage_info">
                      <thead>
                        <tr role="row">
                          <th style="width: 4%;" class="text-center align-middle">#</th>
                          <th style="width: 15%;" class="align-middle sorting" @click="sortBy('reference')">Reference / Date</th>
                          <th style="width: 22%;" class="align-middle sorting" @click="sortBy('title')">Title / Category</th>
                          <th style="width: 14%;" class="text-right align-middle sorting" @click="sortBy('amount')">Amount (UGX)</th>
                          <th style="width: 13%;" class="align-middle sorting" @click="sortBy('department_name')">Department</th>
                          <th style="width: 10%;" class="text-center align-middle">Status</th>
                          <th style="width: 12%;" class="align-middle">Recorded By</th>
                          <th style="width: 5%;" class="text-center align-middle">Receipt</th>
                          <th style="width: 5%;" class="text-center align-middle">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-if="loading">
                          <td colspan="9" class="text-center py-5 text-muted">
                            <i class="fas fa-spinner fa-spin fa-2x mb-2 d-block text-primary"></i>
                            Loading expenses register...
                          </td>
                        </tr>
                        <tr v-else-if="!expenses.length">
                          <td colspan="9" class="text-center py-5 text-muted">
                            <i class="fas fa-inbox fa-2x mb-2 d-block text-muted"></i>
                            No expenses match the specified filter criteria.
                          </td>
                        </tr>
                        <template v-else>
                          <tr v-for="(expense, idx) in expenses" :key="expense.id" class="odd">
                            <td class="text-center align-middle font-weight-bold text-muted">{{ offset + idx + 1 }}</td>
                            <td class="align-middle">
                              <router-link :to="'/expenses/'+expense.id" class="font-weight-bold text-primary text-decoration-none">
                                {{ expense.reference }}
                              </router-link>
                              <div class="text-muted text-xs mt-1"><i class="far fa-calendar-alt mr-1"></i> {{ expense.expense_date }}</div>
                            </td>
                            <td class="align-middle">
                              <div class="font-weight-bold text-dark">{{ expense.title }}</div>
                              <span class="badge badge-light text-muted border font-weight-normal mt-1">{{ expense.category_name || 'General' }}</span>
                            </td>
                            <td class="text-right align-middle font-weight-bold text-dark">
                              {{ expenseMoney(expense.amount) }}
                            </td>
                            <td class="align-middle">
                              <span class="text-dark font-weight-600">{{ expense.department_name || '-' }}</span>
                            </td>
                            <td class="text-center align-middle">
                              <span :class="getStatusBadgeClass(expense.status)">
                                {{ expense.status || 'RECORDED' }}
                              </span>
                            </td>
                            <td class="align-middle">
                              <strong class="text-dark d-block text-xs">{{ expense.recorded_by_name || 'System' }}</strong>
                              <span class="text-muted text-xs">{{ formatDate(expense.recorded_at) }}</span>
                            </td>
                            <td class="text-center align-middle">
                              <a v-if="expense.attachment" :href="'/api/v1/expenses/'+expense.id+'/attachment'" target="_blank" class="badge badge-success text-decoration-none px-2 py-1" title="Download Receipt">
                                <i class="fas fa-paperclip mr-1"></i> Receipt
                              </a>
                              <span v-else class="badge badge-light text-muted border">None</span>
                            </td>
                            <td class="text-center align-middle">
                              <div class="btn-group btn-group-sm" role="group" aria-label="Expense Actions">
                                <!-- Verify Action -->
                                <button 
                                  v-if="canVerify(expense)" 
                                  @click="verifyExpense(expense)" 
                                  class="btn btn-outline-info btn-action" 
                                  title="Verify Expense"
                                  :disabled="actionLoading === expense.id"
                                >
                                  <i class="fas fa-check"></i>
                                </button>

                                <!-- Approve Action -->
                                <button 
                                  v-if="canApprove(expense)" 
                                  @click="approveExpense(expense)" 
                                  class="btn btn-outline-success btn-action" 
                                  title="Approve Expense"
                                  :disabled="actionLoading === expense.id"
                                >
                                  <i class="fas fa-check-circle"></i>
                                </button>

                                <!-- Detail / View Action -->
                                <router-link 
                                  :to="'/expenses/'+expense.id" 
                                  class="btn btn-outline-primary btn-action" 
                                  title="View Details"
                                >
                                  <i class="fas fa-eye"></i>
                                </router-link>

                                <!-- Delete Action -->
                                <button 
                                  v-if="canDelete" 
                                  @click="confirmDelete(expense)" 
                                  class="btn btn-outline-danger btn-action" 
                                  title="Delete Expense"
                                  :disabled="actionLoading === expense.id"
                                >
                                  <i class="fas fa-trash-alt"></i>
                                </button>
                              </div>
                            </td>
                          </tr>
                        </template>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Footer Info & Pagination Row -->
                <div class="row mt-3 align-items-center">
                  <div class="col-sm-12 col-md-5 mb-2 mb-md-0">
                    <div class="dataTables_info text-xs font-weight-bold text-muted" id="save-stage_info" role="status" aria-live="polite">
                      Showing {{ total ? offset + 1 : 0 }} to {{ Math.min(offset + expenses.length, total) }} of {{ total }} entries
                    </div>
                  </div>
                  <div class="col-sm-12 col-md-7 text-md-right">
                    <div class="dataTables_paginate paging_simple_numbers d-inline-block" id="save-stage_paginate">
                      <ul class="pagination pagination-sm m-0 justify-content-end">
                        <li class="paginate_button page-item previous" :class="{ disabled: loading || offset === 0 }">
                          <a href="#" @click.prevent="prevPage" class="page-link"><i class="fas fa-chevron-left mr-1"></i> Previous</a>
                        </li>

                        <li 
                          v-for="p in visiblePages" 
                          :key="p" 
                          class="paginate_button page-item"
                          :class="{ active: p === currentPage }"
                        >
                          <a href="#" @click.prevent="goToPage(p)" class="page-link">{{ p }}</a>
                        </li>

                        <li class="paginate_button page-item next" :class="{ disabled: loading || offset + limit >= total }">
                          <a href="#" @click.prevent="nextPage" class="page-link">Next <i class="fas fa-chevron-right ml-1"></i></a>
                        </li>
                      </ul>
                    </div>
                  </div>
                </div>

              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Print / PDF Report Template (Hidden, rendered when exportPDF triggered) -->
    <div id="print-area" class="d-none">
      <div class="text-center p-4">
        <h2>NATIONAL COUNCIL OF SPORTS - UGANDA</h2>
        <h3>EXPENSES REGISTER REPORT</h3>
        <p>Generated: {{ new Date().toLocaleString() }} | Filter Total: UGX {{ expenseMoney(totalAmount) }}</p>
        <table class="table table-bordered text-xs mt-3 w-100">
          <thead>
            <tr>
              <th>Ref</th>
              <th>Date</th>
              <th>Title</th>
              <th>Category</th>
              <th>Amount (UGX)</th>
              <th>Department</th>
              <th>Status</th>
              <th>Recorded By</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="e in expenses" :key="e.id">
              <td>{{ e.reference }}</td>
              <td>{{ e.expense_date }}</td>
              <td>{{ e.title }}</td>
              <td>{{ e.category_name }}</td>
              <td>{{ expenseMoney(e.amount) }}</td>
              <td>{{ e.department_name }}</td>
              <td>{{ e.status || 'RECORDED' }}</td>
              <td>{{ e.recorded_by_name }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </ExpensePage>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import { expenseData, expenseError, expenseMoney } from '@/api/expenses'

const auth = useAuthStore()
const expenses = ref([])
const categories = ref([])
const departments = ref([])
const total = ref(0)
const totalAmount = ref('0')
const offset = ref(0)
const limit = ref(10)
const loading = ref(false)
const actionLoading = ref(null)
const error = ref('')
const successMsg = ref('')

const blank = () => ({ search: '', category_id: '', department_id: '', from: '', to: '' })
const filters = ref(blank())
let version = 0
let debounceTimer = null

const currentPage = computed(() => Math.floor(offset.value / limit.value) + 1)
const totalPages = computed(() => Math.ceil(total.value / limit.value))

const canDelete = computed(() => auth.isAdmin || (typeof auth.hasRole === 'function' && auth.hasRole('chief_accountant')))

function canVerify(exp) {
  const roleCheck = typeof auth.hasRole === 'function' ? (auth.hasRole('senior_accountant') || auth.hasRole('chief_accountant')) : false
  return (!exp.status || exp.status === 'RECORDED') && (auth.isAdmin || roleCheck)
}

function canApprove(exp) {
  const roleCheck = typeof auth.hasRole === 'function' ? (auth.hasRole('chief_accountant') || auth.hasRole('finance_department')) : false
  return exp.status === 'VERIFIED' && (auth.isAdmin || roleCheck)
}

function getStatusBadgeClass(st) {
  switch (st) {
    case 'VERIFIED':
      return 'badge badge-info text-white'
    case 'APPROVED':
      return 'badge badge-success'
    case 'REJECTED':
      return 'badge badge-danger'
    default:
      return 'badge badge-warning text-dark'
  }
}

function formatDate(dt) {
  if (!dt) return ''
  try {
    return new Date(dt).toLocaleDateString()
  } catch {
    return dt
  }
}

async function load() {
  const current = ++version
  loading.value = true
  error.value = ''
  try {
    const res = expenseData(await client.get('/api/v1/expenses', {
      params: { ...filters.value, offset: offset.value, limit: limit.value }
    }))
    if (current !== version) return
    expenses.value = res.expenses || []
    total.value = res.total || 0
    totalAmount.value = res.total_amount || '0'
  } catch (e) {
    if (current === version) {
      error.value = expenseError(e)
      expenses.value = []
      total.value = 0
      totalAmount.value = '0'
    }
  } finally {
    if (current === version) loading.value = false
  }
}

function debouncedSearch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    offset.value = 0
    load()
  }, 300)
}

function reset() {
  filters.value = blank()
  offset.value = 0
  load()
}

function prevPage() {
  if (offset.value >= limit.value) {
    offset.value -= limit.value
    load()
  }
}

function nextPage() {
  if (offset.value + limit.value < total.value) {
    offset.value += limit.value
    load()
  }
}

function goToPage(p) {
  offset.value = (p - 1) * limit.value
  load()
}

const visiblePages = computed(() => {
  const totalP = totalPages.value
  if (totalP <= 1) return [1]
  const cur = currentPage.value
  const pages = []
  const start = Math.max(1, cur - 2)
  const end = Math.min(totalP, cur + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

async function verifyExpense(exp) {
  if (!confirm(`Are you sure you want to VERIFY expense "${exp.reference}" (${exp.title})?`)) return
  actionLoading.value = exp.id
  try {
    await client.put(`/api/v1/expenses/${exp.id}/verify`)
    successMsg.value = `Expense ${exp.reference} verified successfully.`
    await load()
  } catch (e) {
    error.value = expenseError(e)
  } finally {
    actionLoading.value = null
  }
}

async function approveExpense(exp) {
  if (!confirm(`Are you sure you want to APPROVE expense "${exp.reference}" (${exp.title})?`)) return
  actionLoading.value = exp.id
  try {
    await client.put(`/api/v1/expenses/${exp.id}/approve`)
    successMsg.value = `Expense ${exp.reference} approved successfully.`
    await load()
  } catch (e) {
    error.value = expenseError(e)
  } finally {
    actionLoading.value = null
  }
}

async function confirmDelete(exp) {
  if (!confirm(`CAUTION: Delete expense "${exp.reference}" (${exp.title}) permanently?`)) return
  actionLoading.value = exp.id
  try {
    await client.delete(`/api/v1/expenses/${exp.id}`)
    successMsg.value = `Expense ${exp.reference} deleted successfully.`
    await load()
  } catch (e) {
    error.value = expenseError(e)
  } finally {
    actionLoading.value = null
  }
}

function sortBy(col) {
  expenses.value.sort((a, b) => {
    if (a[col] < b[col]) return -1
    if (a[col] > b[col]) return 1
    return 0
  })
}

function exportExcel() {
  if (!expenses.value.length) return alert('No expenses available to export.')
  const headers = ['Reference', 'Date', 'Title', 'Category', 'Amount (UGX)', 'Department', 'Status', 'Recorded By', 'Payee', 'Payment Method']
  const rows = expenses.value.map(e => [
    `"${e.reference}"`,
    `"${e.expense_date}"`,
    `"${e.title.replace(/"/g, '""')}"`,
    `"${e.category_name}"`,
    `"${e.amount}"`,
    `"${e.department_name}"`,
    `"${e.status || 'RECORDED'}"`,
    `"${e.recorded_by_name}"`,
    `"${e.payee.replace(/"/g, '""')}"`,
    `"${e.payment_method}"`
  ])
  const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', `ncs_expenses_export_${new Date().toISOString().slice(0,10)}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function exportPDF() {
  window.print()
}

onMounted(async () => {
  await load()
  try {
    const [c, o] = await Promise.all([
      client.get('/api/v1/expenses/categories'),
      client.get('/api/v1/expenses/options')
    ])
    categories.value = expenseData(c)
    departments.value = expenseData(o).departments
  } catch (e) {
    error.value = expenseError(e)
  }
})
</script>

<style scoped>
.card-statistic-1 {
  display: flex;
  align-items: stretch;
}
.btn-emerald {
  background-color: #059669;
  border-color: #059669;
}
.btn-emerald:hover {
  background-color: #047857;
  border-color: #047857;
}
.bg-soft-primary {
  background-color: rgba(103, 119, 239, 0.12);
}
.btn-action {
  width: 32px;
  height: 32px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}
.table td, .table th {
  vertical-align: middle !important;
}
@media print {
  body * {
    visibility: hidden;
  }
  #print-area, #print-area * {
    visibility: visible;
  }
  #print-area {
    display: block !important;
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
  }
}
</style>
