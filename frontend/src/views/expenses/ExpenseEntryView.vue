<template>
  <ExpensePage title="Record Expense" description="Record an expense and attach its supporting receipt.">
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
        <div class="card-body d-flex flex-wrap align-items-center justify-content-between gap-3">
          <div class="media align-items-center">
            <div class="expense-otika-avatar"><i class="icofont-user-alt-5"></i></div>
            <div class="media-body ml-3">
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
          <div class="form-row">
            <div class="form-group col-md-6">
              <label for="expense-date">Expense Date</label>
              <input id="expense-date" required v-model="form.expense_date" type="date" :max="options.today" class="form-control" />
            </div>
            <div class="form-group col-md-6">
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

          <div class="form-row">
            <div class="form-group col-md-4">
              <label for="expense-amount">Amount (UGX)</label>
              <div class="input-group">
                <div class="input-group-prepend">
                  <span class="input-group-text">UGX</span>
                </div>
                <input id="expense-amount" required v-model="form.amount" type="number" min="0.01" max="999999999999.99" step="0.01" class="form-control text-right" placeholder="0.00" />
              </div>
            </div>
            <div class="form-group col-md-4">
              <label for="expense-department">Expense Department</label>
              <select id="expense-department" aria-label="Expense department" required v-model="form.department_id" class="form-control custom-select">
                <option value="" disabled>Select the department charged</option>
                <option v-for="department in options.departments" :key="department.id" :value="department.id">{{ department.name }}</option>
              </select>
            </div>
            <div class="form-group col-md-4">
              <label for="payment-method">Payment Method</label>
              <select id="payment-method" aria-label="Payment method" required v-model="form.payment_method" class="form-control custom-select">
                <option v-for="(label,value) in paymentMethods" :key="value" :value="value">{{ label }}</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group col-md-6">
              <label for="expense-payee">Payee / Supplier</label>
              <input id="expense-payee" required v-model.trim="form.payee" maxlength="200" class="form-control" placeholder="Supplier or staff member paid" />
            </div>
            <div class="form-group col-md-6">
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
            <input ref="fileInput" type="file" accept="application/pdf,image/jpeg,image/png,image/webp" class="d-none" @change="chooseFile" />
            <div class="dz-message">
              <i class="icofont-cloud-upload"></i>
              <h6>{{ uploadTitle }}</h6>
              <p>Drop a PDF or image here, or click to browse. Maximum size 5 MB.</p>
              <button v-if="file" type="button" class="btn btn-sm btn-outline-danger mt-2" @click.prevent="clearFile">
                <i class="icofont-close"></i> Remove File
              </button>
            </div>
          </label>

          <div class="card-footer px-0 pb-0 d-flex flex-wrap align-items-center justify-content-end gap-2">
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
import {useRouter} from 'vue-router'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import {expenseData,expenseError,paymentMethods} from '@/api/expenses'
const router=useRouter(),options=ref(null),categories=ref([]),error=ref(''),loading=ref(true),saving=ref(false),file=ref(null),dragOver=ref(false),fileInput=ref(null)
const form=ref({expense_date:'',category_id:'',title:'',description:'',amount:'',payee:'',payment_method:'CASH',payment_reference:'',department_id:''})
const uploadTitle=computed(()=>file.value?file.value.name:'Drop files here to upload')
onMounted(async()=>{try{const [o,c]=await Promise.all([client.get('/api/v1/expenses/options'),client.get('/api/v1/expenses/categories')]);options.value=expenseData(o);categories.value=(expenseData(c)||[]).filter(c=>c.is_active);form.value.expense_date=options.value.today}catch(e){error.value=expenseError(e)}finally{loading.value=false}})
function setFile(selected){error.value='';dragOver.value=false;if(selected&&(selected.size>5*1024*1024||!['application/pdf','image/jpeg','image/png','image/webp'].includes(selected.type))){error.value='Choose a PDF, JPEG, PNG or WebP receipt up to 5 MB.';if(fileInput.value)fileInput.value.value='';file.value=null;return}file.value=selected||null}
function chooseFile(event){setFile(event.target.files[0])}
function dropFile(event){setFile(event.dataTransfer?.files?.[0])}
function clearFile(){file.value=null;if(fileInput.value)fileInput.value.value=''}
async function save(){if(saving.value)return;saving.value=true;error.value='';try{const data=new FormData();for(const [key,value] of Object.entries(form.value))data.append(key,value);if(file.value)data.append('file',file.value);const res=await client.post('/api/v1/expenses',data,{headers:{'Content-Type':'multipart/form-data'}});await router.push({name:'ExpenseDetail',params:{id:expenseData(res).id},query:{saved:'1'}})}catch(e){error.value=expenseError(e)}finally{saving.value=false}}
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
.expense-otika-form .card-header {
  align-items: center;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  min-height: 56px;
}
.expense-otika-form .card-header h4 {
  color: #34395e;
  font-size: 16px;
  font-weight: 700;
  margin: 0;
}
.expense-otika-form .form-group label {
  color: #34395e;
  font-size: 12px;
  font-weight: 700;
}
.expense-otika-form .form-control,
.expense-otika-form .custom-select {
  border-color: #e4e6fc;
  border-radius: 4px;
  color: #495057;
  font-size: 13px;
  min-height: 42px;
}
.expense-otika-form textarea.form-control {
  min-height: 108px;
}
.expense-otika-form .form-control:focus,
.expense-otika-form .custom-select:focus {
  border-color: #6777ef;
  box-shadow: 0 0 0 0.2rem rgba(103, 119, 239, .15);
}
.expense-otika-form .input-group-text {
  background: #f4f6f9;
  border-color: #e4e6fc;
  color: #6777ef;
  font-size: 12px;
  font-weight: 800;
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
  color: #6c757d;
  cursor: pointer;
  display: flex;
  justify-content: center;
  margin-bottom: 0;
  min-height: 190px;
  padding: 28px;
  text-align: center;
  transition: background .2s, border-color .2s, box-shadow .2s;
  width: 100%;
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
  font-size: 44px;
  margin-bottom: 8px;
}
.expense-dropzone h6 {
  color: #34395e;
  font-size: 15px;
  font-weight: 800;
  margin-bottom: 6px;
}
.expense-dropzone p {
  color: #6c757d;
  font-size: 12px;
  margin: 0;
}
.spin {
  animation: spin 1s linear infinite;
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
</style>
