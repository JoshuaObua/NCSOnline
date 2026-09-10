<template>
  <ExpensePage :title="editing ? 'Edit Expense' : 'Record Expense'" :description="editing ? 'Update your own recorded expense while it is still pending.' : 'Record an expense and attach its supporting receipt.'">
    <template #actions>
      <router-link to="/expenses" class="btn btn-sm btn-icon icon-left btn-outline-secondary">
        <i class="icofont-rounded-left"></i> Back to Expenses
      </router-link>
    </template>

    <div v-if="loading" class="card">
      <div class="card-body text-center py-5 text-muted">
        <i class="icofont-spinner-alt-4 spin mr-2"></i> Loading expense options...
      </div>
    </div>

    <div v-if="error" role="alert" class="alert alert-danger expense-otika-alert">
      <i class="icofont-warning-alt"></i>
      <span>{{ error }}</span>
    </div>

    <template v-if="!loading && options">
      <div class="card expense-otika-summary">
        <div class="card-body expense-summary-body">
          <div class="expense-summary-user">
            <div class="expense-otika-avatar"><i class="icofont-user-alt-5"></i></div>
            <div>
              <h6 class="mb-1">Recorded by {{ options.recorded_by_name }}</h6>
              <p class="mb-0 text-muted small">Your name and save time are recorded automatically.</p>
            </div>
          </div>
          <span class="badge badge-primary badge-shadow">Expense Register</span>
        </div>
      </div>

      <div v-if="!categories.length" class="alert alert-danger expense-otika-alert">
        <i class="icofont-warning-alt"></i>
        <span>No active expense categories are available. An administrator must create or activate a category before you can record an expense.</span>
      </div>

      <form @submit.prevent="save" class="card expense-otika-form" :aria-busy="saving">
        <div class="card-header">
          <h4><i class="icofont-ui-note mr-2"></i>Expense Details</h4>
        </div>
        <fieldset :disabled="saving || !categories.length" class="card-body">
          <div class="expense-form-row">
            <div class="form-group expense-form-col">
              <label for="expense-date">Expense Date</label>
              <input id="expense-date" required v-model="form.expense_date" type="date" :max="options.today" class="form-control" />
            </div>
            <div class="form-group expense-form-col">
              <label for="expense-category">Expense Category</label>
              <select id="expense-category" aria-label="Expense category" required v-model="form.category_id" class="form-control custom-select">
                <option value="" disabled>Select a category</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label for="expense-title">Title</label>
            <input id="expense-title" required v-model.trim="form.title" maxlength="200" class="form-control" placeholder="e.g. Office stationery purchase" />
          </div>

          <div class="form-group">
            <label for="expense-description">Description / Reason</label>
            <textarea id="expense-description" required v-model.trim="form.description" maxlength="5000" rows="4" class="form-control" placeholder="Briefly describe what this expense covers"></textarea>
          </div>

          <div class="expense-form-row expense-form-row-three">
            <div class="form-group expense-form-col">
              <label for="expense-amount">Amount (UGX)</label>
              <div class="expense-input-group">
                <span class="input-group-text">UGX</span>
                <input id="expense-amount" required v-model="form.amount" type="number" min="0.01" max="999999999999.99" step="0.01" class="form-control text-right" placeholder="0.00" />
              </div>
            </div>
            <div class="form-group expense-form-col">
              <label for="expense-department">Expense Department</label>
              <select id="expense-department" aria-label="Expense department" required v-model="form.department_id" class="form-control custom-select">
                <option value="" disabled>Select the department charged</option>
                <option v-for="department in options.departments" :key="department.id" :value="department.id">{{ department.name }}</option>
              </select>
            </div>
            <div class="form-group expense-form-col">
              <label for="payment-method">Payment Method</label>
              <select id="payment-method" aria-label="Payment method" required v-model="form.payment_method" class="form-control custom-select">
                <option v-for="(label,value) in paymentMethods" :key="value" :value="value">{{ label }}</option>
              </select>
            </div>
          </div>

          <div class="expense-form-row">
            <div class="form-group expense-form-col">
              <label for="expense-payee">Payee / Supplier</label>
              <input id="expense-payee" required v-model.trim="form.payee" maxlength="200" class="form-control" placeholder="Supplier or staff member paid" />
            </div>
            <div class="form-group expense-form-col">
              <label for="payment-reference">Payment / Receipt Reference</label>
              <input id="payment-reference" v-model.trim="form.payment_reference" maxlength="200" class="form-control" placeholder="Optional receipt or voucher number" />
            </div>
          </div>

          <div class="section-title mt-2">Receipt / Supporting Document</div>
          <label
            class="dropzone expense-dropzone"
            :class="{ 'expense-dropzone-active': dragOver, 'expense-dropzone-loaded': file }"
            @dragover.prevent="dragOver = true"
            @dragleave.prevent="dragOver = false"
            @drop.prevent="dropFile"
          >
            <input ref="fileInput" type="file" accept="application/pdf,image/jpeg,image/png,image/webp" class="expense-upload-input" @change="chooseFile" />
            <div class="dz-message">
              <i class="icofont-cloud-upload"></i>
              <h6>{{ uploadTitle }}</h6>
              <p>Drop a PDF or image here, or click to browse. Maximum size 5 MB.</p>
              <button v-if="file" type="button" class="btn btn-sm btn-outline-danger mt-2" @click.prevent="clearFile">
                <i class="icofont-close"></i> Remove File
              </button>
            </div>
          </label>

          <div class="card-footer expense-form-actions">
            <router-link to="/expenses" class="btn btn-icon icon-left btn-outline-secondary">
              <i class="icofont-close"></i> Cancel
            </router-link>
            <button type="submit" class="btn btn-icon icon-left btn-primary" :disabled="saving || !categories.length">
              <i v-if="saving" class="icofont-spinner-alt-4 spin"></i>
              <i v-else class="icofont-check-circled"></i>
              {{ saving ? 'Saving...' : 'Record Expense' }}
            </button>
          </div>
        </fieldset>
      </form>
    </template>
  </ExpensePage>
</template>
<script setup>
import {computed,ref,onMounted} from 'vue'
import {useRouter,useRoute} from 'vue-router'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import {expenseData,expenseError,paymentMethods} from '@/api/expenses'
const router=useRouter(),route=useRoute(),options=ref(null),categories=ref([]),error=ref(''),loading=ref(true),saving=ref(false),file=ref(null),dragOver=ref(false),fileInput=ref(null)
const editing=ref(Boolean(route.params.id))
const form=ref({expense_date:'',category_id:'',title:'',description:'',amount:'',payee:'',payment_method:'CASH',payment_reference:'',department_id:''})
const uploadTitle=computed(()=>file.value?file.value.name:'Drop files here to upload')
onMounted(async()=>{try{const requests=[client.get('/api/v1/expenses/options'),client.get('/api/v1/expenses/categories')];if(editing.value)requests.push(client.get('/api/v1/expenses/'+encodeURIComponent(route.params.id)));const [o,c,e]=await Promise.all(requests);options.value=expenseData(o);categories.value=(expenseData(c)||[]).filter(c=>c.is_active);if(editing.value){const current=expenseData(e);form.value={expense_date:current.expense_date,category_id:current.category_id,title:current.title,description:current.description,amount:current.amount,payee:current.payee,payment_method:current.payment_method,payment_reference:current.payment_reference,department_id:current.department_id}}else form.value.expense_date=options.value.today}catch(e){error.value=expenseError(e)}finally{loading.value=false}})
function setFile(selected){error.value='';dragOver.value=false;if(selected&&(selected.size>5*1024*1024||!['application/pdf','image/jpeg','image/png','image/webp'].includes(selected.type))){error.value='Choose a PDF, JPEG, PNG or WebP receipt up to 5 MB.';if(fileInput.value)fileInput.value.value='';file.value=null;return}file.value=selected||null}
function chooseFile(event){setFile(event.target.files[0])}
function dropFile(event){setFile(event.dataTransfer?.files?.[0])}
function clearFile(){file.value=null;if(fileInput.value)fileInput.value.value=''}
async function save(){if(saving.value)return;saving.value=true;error.value='';try{let res;if(editing.value)res=await client.put('/api/v1/expenses/'+encodeURIComponent(route.params.id),form.value);else{const data=new FormData();for(const [key,value] of Object.entries(form.value))data.append(key,value);if(file.value)data.append('file',file.value);res=await client.post('/api/v1/expenses',data,{headers:{'Content-Type':'multipart/form-data'}})}await router.push({name:'ExpenseDetail',params:{id:editing.value?route.params.id:expenseData(res).id},query:editing.value?{updated:'1'}:{saved:'1'}})}catch(e){error.value=expenseError(e)}finally{saving.value=false}}
</script>
<style scoped>
.expense-otika-alert {
  align-items: center;
  border-radius: 4px;
  display: flex;
  gap: 10px;
}
.expense-otika-summary,
.expense-otika-form {
  border: 1px solid #e4e6fc;
  border-radius: 4px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, .06);
}
.expense-otika-form .card-body {
  padding: 22px 24px 24px;
}
.expense-summary-body {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: space-between;
}
.expense-summary-user {
  align-items: center;
  display: flex;
  gap: 12px;
  min-width: 0;
}
.expense-otika-form .card-header {
  align-items: center;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  min-height: 56px;
}
.expense-form-row {
  display: grid;
  gap: 18px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}
.expense-form-row-three {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}
.expense-form-col {
  min-width: 0;
}
.expense-form-actions {
  align-items: center;
  background: transparent;
  border-top: 1px solid #f1f5f9;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 8px;
  padding: 18px 0 0;
}
.expense-otika-form .card-header h4 {
  color: #34395e;
  font-size: 16px;
  font-weight: 700;
  margin: 0;
}
.expense-otika-form .form-group label {
  color: #34395e;
  display: block;
  font-size: 12px;
  font-weight: 700;
  margin-bottom: 7px;
}
.expense-otika-form .form-control,
.expense-otika-form .custom-select {
  appearance: auto;
  box-sizing: border-box;
  border-color: #e4e6fc;
  border-radius: 4px;
  color: #495057;
  display: block;
  font-size: 13px;
  height: 42px;
  line-height: 1.5;
  padding: 9px 12px;
  width: 100%;
}
.expense-otika-form textarea.form-control {
  height: 118px;
  line-height: 1.5;
  min-height: 108px;
  resize: vertical;
}
.expense-otika-form .form-control:focus,
.expense-otika-form .custom-select:focus {
  border-color: #6777ef;
  box-shadow: 0 0 0 0.2rem rgba(103, 119, 239, .15);
}
.expense-otika-form .input-group-text {
  align-items: center;
  background: #f4f6f9;
  border-color: #e4e6fc;
  border-radius: 4px 0 0 4px;
  border-style: solid;
  border-width: 1px 0 1px 1px;
  color: #6777ef;
  display: inline-flex;
  font-size: 12px;
  font-weight: 800;
  height: 42px;
  padding: 8px 12px;
  white-space: nowrap;
}
.expense-input-group {
  display: flex;
  height: 42px;
  width: 100%;
}
.expense-input-group .form-control {
  border-radius: 0 4px 4px 0;
  height: 42px;
  min-width: 0;
}
.expense-otika-avatar {
  align-items: center;
  background: #6777ef;
  border-radius: 4px;
  color: #fff;
  display: inline-flex;
  font-size: 20px;
  height: 44px;
  justify-content: center;
  width: 44px;
}
.expense-dropzone {
  align-items: center;
  background: #fdfdff;
  border: 2px dashed #cdd3ff;
  border-radius: 4px;
  box-sizing: border-box;
  color: #6c757d;
  cursor: pointer;
  display: flex;
  justify-content: center;
  margin-bottom: 0;
  min-height: 150px;
  padding: 22px;
  position: relative;
  text-align: center;
  transition: background .2s, border-color .2s, box-shadow .2s;
  width: 100%;
}
.expense-upload-input {
  height: 1px;
  left: -9999px;
  opacity: 0;
  overflow: hidden;
  pointer-events: none;
  position: absolute;
  width: 1px;
}
.expense-dropzone:hover,
.expense-dropzone-active {
  background: #f4f6ff;
  border-color: #6777ef;
  box-shadow: 0 0 0 4px rgba(103, 119, 239, .08);
}
.expense-dropzone-loaded {
  background: #f0fdf4;
  border-color: #47c363;
}
.expense-dropzone .dz-message {
  margin: 0;
}
.expense-dropzone i {
  color: #6777ef;
  display: block;
  font-size: 34px;
  margin-bottom: 6px;
}
.expense-dropzone h6 {
  color: #34395e;
  font-size: 14px;
  font-weight: 800;
  margin-bottom: 5px;
  overflow-wrap: anywhere;
}
.expense-dropzone p {
  color: #6c757d;
  font-size: 12px;
  line-height: 1.45;
  margin: 0;
}
.spin {
  animation: spin 1s linear infinite;
}
.btn-outline-secondary {
  background: #fff;
  border-color: #e4e6fc;
  color: #6c757d;
}
.btn-outline-secondary:hover {
  background: #f4f6f9;
  border-color: #6777ef;
  color: #6777ef;
}
.btn-outline-danger {
  background: #fff;
  border-color: #fc544b;
  color: #fc544b;
}
.btn-outline-danger:hover {
  background: #fc544b;
  color: #fff;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
html.dark .expense-otika-summary,
html.dark .expense-otika-form {
  background: #0f172a;
  border-color: #1e293b;
}
html.dark .expense-otika-form .card-header {
  border-color: #1e293b;
}
html.dark .expense-form-actions {
  border-color: #1e293b;
}
html.dark .expense-otika-form .card-header h4,
html.dark .expense-otika-form .form-group label,
html.dark .expense-dropzone h6 {
  color: #e2e8f0;
}
html.dark .expense-otika-form .form-control,
html.dark .expense-otika-form .custom-select {
  background: #111827;
  border-color: #334155;
  color: #e2e8f0;
}
html.dark .expense-dropzone {
  background: #111827;
  border-color: #334155;
}
@media (max-width: 992px) {
  .expense-form-row,
  .expense-form-row-three {
    grid-template-columns: 1fr;
  }
}
@media (max-width: 640px) {
  .expense-summary-body,
  .expense-summary-user,
  .expense-form-actions {
    align-items: stretch;
    flex-direction: column;
  }
  .expense-form-actions .btn {
    width: 100%;
  }
}
</style>
