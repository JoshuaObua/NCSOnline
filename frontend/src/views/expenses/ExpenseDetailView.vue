<template>
  <ExpensePage :title="expense ? expense.reference : 'Expense Details'">
    <template #actions><router-link to="/expenses" class="expense-button expense-secondary">Back to Expenses</router-link></template>
    <p v-if="route.query.saved && expense" role="status" class="expense-success">Expense recorded successfully.</p>
    <p v-if="loading" role="status">Loading expense…</p><p v-if="error" role="alert" class="expense-error">{{ error }}</p>
    <template v-if="expense">
      <div class="expense-card"><h2>{{ expense.title }}</h2><strong>Recorded by {{ expense.recorded_by_name }}</strong><p>Department: {{ expense.department_name }}</p><p class="expense-muted">Recorded {{ new Date(expense.recorded_at).toLocaleString() }}</p></div>
      <dl class="expense-card expense-grid expense-details"><div><dt>Expense date</dt><dd>{{ expense.expense_date }}</dd></div><div><dt>Category</dt><dd>{{ expense.category_name }}</dd></div><div><dt>Amount</dt><dd>UGX {{ expenseMoney(expense.amount) }}</dd></div><div><dt>Payee / supplier</dt><dd>{{ expense.payee }}</dd></div><div><dt>Payment method</dt><dd>{{ paymentMethods[expense.payment_method] }}</dd></div><div><dt>Payment / receipt reference</dt><dd>{{ expense.payment_reference || 'Not supplied' }}</dd></div><div class="expense-wide"><dt>Description / reason</dt><dd>{{ expense.description }}</dd></div></dl>
      <div class="expense-card"><h2>Supporting receipt</h2><button v-if="expense.attachment" @click="download" :disabled="downloading" class="expense-button expense-secondary">{{ downloading ? 'Downloading…' : 'Download '+expense.attachment.name }}</button><p v-else>No receipt attached.</p></div>
    </template>
  </ExpensePage>
</template>
<script setup>
import {ref,watch} from 'vue'
import {useRoute} from 'vue-router'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import {expenseData,expenseError,expenseMoney,paymentMethods} from '@/api/expenses'
const route=useRoute(),expense=ref(null),loading=ref(false),downloading=ref(false),error=ref('')
let loadVersion = 0
watch(()=>route.params.id,async id=>{const version=++loadVersion;loading.value=true;expense.value=null;error.value='';try{const data=expenseData(await client.get('/api/v1/expenses/'+encodeURIComponent(id)));if(version===loadVersion)expense.value=data}catch(e){if(version===loadVersion)error.value=expenseError(e)}finally{if(version===loadVersion)loading.value=false}},{immediate:true})
async function download(){downloading.value=true;try{const res=await client.get('/api/v1/expenses/'+expense.value.id+'/attachment',{responseType:'blob'});const url=URL.createObjectURL(res.data),a=document.createElement('a');a.href=url;a.download=expense.value.attachment.name;a.click();setTimeout(()=>URL.revokeObjectURL(url),1000)}catch(e){error.value='Could not download the receipt. Please try again.'}finally{downloading.value=false}}
</script>
