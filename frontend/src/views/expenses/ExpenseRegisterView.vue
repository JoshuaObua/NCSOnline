<template>
  <ExpensePage title="Expenses" description="Expense register with the responsible department and recording officer.">
    <template #actions><div class="expense-actions"><router-link v-if="auth.isAdmin" to="/expenses/categories" class="expense-button expense-secondary">Manage Categories</router-link><router-link to="/expenses/new" class="expense-button">Record Expense</router-link></div></template>
    <p v-if="error" role="alert" class="expense-error">{{ error }}</p>
    <form @submit.prevent="offset=0;load()" class="expense-card expense-grid">
      <label>Search<input v-model="filters.search" class="expense-input" placeholder="Title, reference, payee or recorded by" /></label>
      <label>Category<select aria-label="Category" v-model="filters.category_id" class="expense-input"><option value="">All categories</option><option v-for="c in categories" :value="c.id" :key="c.id">{{ c.name }}</option></select></label>
      <label>Department<select aria-label="Department" v-model="filters.department_id" class="expense-input"><option value="">All departments</option><option v-for="d in departments" :value="d.id" :key="d.id">{{ d.name }}</option></select></label>
      <label>From date<input v-model="filters.from" type="date" class="expense-input" /></label>
      <label>To date<input v-model="filters.to" type="date" class="expense-input" /></label>
      <div class="expense-actions"><button :disabled="loading" class="expense-button">Apply Filters</button><button type="button" @click="reset" class="expense-button expense-secondary">Reset</button></div>
    </form>
    <div class="expense-card" aria-live="polite"><strong>{{ total }} expenses · UGX {{ expenseMoney(totalAmount) }}</strong><p class="expense-muted">Total for the applied filters.</p></div>
    <div class="expense-card expense-table-wrap">
      <p v-if="loading" role="status">Loading expenses…</p>
      <table v-else><thead><tr><th>Reference / date</th><th>Title / category</th><th>Amount (UGX)</th><th>Department</th><th>Recorded by</th><th>Receipt</th></tr></thead><tbody>
        <tr v-for="expense in expenses" :key="expense.id"><td><router-link :to="'/expenses/'+expense.id" class="expense-link">{{ expense.reference }}</router-link><div>{{ expense.expense_date }}</div></td><td>{{ expense.title }}<div class="expense-muted">{{ expense.category_name }}</div></td><td>{{ expenseMoney(expense.amount) }}</td><td>{{ expense.department_name }}</td><td><strong>{{ expense.recorded_by_name }}</strong><div class="expense-muted">{{ new Date(expense.recorded_at).toLocaleString() }}</div></td><td>{{ expense.attachment ? 'Attached' : 'None' }}</td></tr>
        <tr v-if="!expenses.length"><td colspan="6">No expenses match these filters.</td></tr>
      </tbody></table>
    </div>
    <div class="expense-actions"><button :disabled="loading||offset===0" @click="offset-=50;load()" class="expense-button expense-secondary">Previous</button><span>{{ total ? offset+1 : 0 }}–{{ Math.min(offset+expenses.length,total) }} of {{ total }}</span><button :disabled="loading||offset+50>=total" @click="offset+=50;load()" class="expense-button expense-secondary">Next</button></div>
  </ExpensePage>
</template>
<script setup>
import {ref,onMounted} from 'vue'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import {useAuthStore} from '@/stores/auth'
import {expenseData,expenseError,expenseMoney} from '@/api/expenses'
const auth=useAuthStore(),expenses=ref([]),categories=ref([]),departments=ref([]),total=ref(0),totalAmount=ref('0'),offset=ref(0),loading=ref(false),error=ref('')
const blank=()=>({search:'',category_id:'',department_id:'',from:'',to:''}),filters=ref(blank())
let version=0
async function load(){const current=++version;loading.value=true;error.value='';try{const res=expenseData(await client.get('/api/v1/expenses',{params:{...filters.value,offset:offset.value}}));if(current!==version)return;expenses.value=res.expenses;total.value=res.total;totalAmount.value=res.total_amount}catch(e){if(current===version){error.value=expenseError(e);expenses.value=[];total.value=0;totalAmount.value='0'}}finally{if(current===version)loading.value=false}}
function reset(){filters.value=blank();offset.value=0;load()}
onMounted(async()=>{await load();try{const [c,o]=await Promise.all([client.get('/api/v1/expenses/categories'),client.get('/api/v1/expenses/options')]);categories.value=expenseData(c);departments.value=expenseData(o).departments}catch(e){error.value=expenseError(e)}})
</script>
