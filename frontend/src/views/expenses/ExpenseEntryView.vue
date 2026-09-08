<template>
  <ExpensePage title="Record Expense" description="Record an expense and attach its supporting receipt.">
    <template #actions><router-link to="/expenses" class="expense-button expense-secondary">Back to Expenses</router-link></template>
    <p v-if="loading" role="status">Loading expense options…</p>
    <p v-if="error" role="alert" class="expense-error">{{ error }}</p>
    <template v-if="!loading && options">
      <div class="expense-card"><strong>Recorded by {{ options.recorded_by_name }}</strong><p class="expense-muted">Your name and the save time are recorded automatically.</p></div>
      <p v-if="!categories.length" class="expense-error">No active expense categories are available. An administrator must create or activate a category before you can record an expense.</p>
      <form @submit.prevent="save" class="expense-card" :aria-busy="saving">
        <fieldset :disabled="saving || !categories.length">
          <div class="expense-grid">
            <label>Expense date<input required v-model="form.expense_date" type="date" :max="options.today" class="expense-input" /></label>
            <label>Expense category<select aria-label="Expense category" required v-model="form.category_id" class="expense-input"><option value="" disabled>Select a category</option><option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option></select></label>
            <label class="expense-wide">Title<input required v-model.trim="form.title" maxlength="200" class="expense-input" placeholder="e.g. Office stationery purchase" /></label>
            <label class="expense-wide">Description / reason<textarea required v-model.trim="form.description" maxlength="5000" rows="4" class="expense-input" /></label>
            <label>Amount (UGX)<input required v-model="form.amount" type="number" min="0.01" max="999999999999.99" step="0.01" class="expense-input" /></label>
            <label>Expense department<select aria-label="Expense department" required v-model="form.department_id" class="expense-input"><option value="" disabled>Select the department charged</option><option v-for="department in options.departments" :key="department.id" :value="department.id">{{ department.name }}</option></select></label>
            <label>Payee / supplier<input required v-model.trim="form.payee" maxlength="200" class="expense-input" /></label>
            <label>Payment method<select aria-label="Payment method" required v-model="form.payment_method" class="expense-input"><option v-for="(label,value) in paymentMethods" :key="value" :value="value">{{ label }}</option></select></label>
            <label class="expense-wide">Payment / receipt reference (optional)<input v-model.trim="form.payment_reference" maxlength="200" class="expense-input" /></label>
            <label class="expense-wide">Receipt / supporting document (optional)<input type="file" accept="application/pdf,image/jpeg,image/png,image/webp" @change="chooseFile" class="expense-input" /><span class="expense-muted">One PDF or image (JPEG, PNG, WebP), up to 5 MB.</span></label>
          </div>
          <div class="expense-actions mt-6"><button type="submit" class="expense-button">{{ saving ? 'Saving…' : 'Record Expense' }}</button><router-link to="/expenses" class="expense-link">Cancel</router-link></div>
        </fieldset>
      </form>
    </template>
  </ExpensePage>
</template>
<script setup>
import {ref,onMounted} from 'vue'
import {useRouter} from 'vue-router'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import {expenseData,expenseError,paymentMethods} from '@/api/expenses'
const router=useRouter(),options=ref(null),categories=ref([]),error=ref(''),loading=ref(true),saving=ref(false),file=ref(null)
const form=ref({expense_date:'',category_id:'',title:'',description:'',amount:'',payee:'',payment_method:'CASH',payment_reference:'',department_id:''})
onMounted(async()=>{try{const [o,c]=await Promise.all([client.get('/api/v1/expenses/options'),client.get('/api/v1/expenses/categories')]);options.value=expenseData(o);categories.value=(expenseData(c)||[]).filter(c=>c.is_active);form.value.expense_date=options.value.today}catch(e){error.value=expenseError(e)}finally{loading.value=false}})
function chooseFile(event){error.value='';const selected=event.target.files[0];if(selected&&(selected.size>5*1024*1024||!['application/pdf','image/jpeg','image/png','image/webp'].includes(selected.type))){error.value='Choose a PDF, JPEG, PNG or WebP receipt up to 5 MB.';event.target.value='';file.value=null;return}file.value=selected||null}
async function save(){if(saving.value)return;saving.value=true;error.value='';try{const data=new FormData();for(const [key,value] of Object.entries(form.value))data.append(key,value);if(file.value)data.append('file',file.value);const res=await client.post('/api/v1/expenses',data,{headers:{'Content-Type':'multipart/form-data'}});await router.push({name:'ExpenseDetail',params:{id:expenseData(res).id},query:{saved:'1'}})}catch(e){error.value=expenseError(e)}finally{saving.value=false}}
</script>
