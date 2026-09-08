<template>
  <ExpensePage title="Expense Categories" description="Admin-managed categories. Deactivated categories remain on historical expense records.">
    <template #actions><router-link to="/expenses" class="expense-button expense-secondary">Back to Expenses</router-link></template>
    <p v-if="error" role="alert" class="expense-error">{{ error }}</p><p v-if="notice" role="status" class="expense-success">{{ notice }}</p>
    <form @submit.prevent="save" class="expense-card"><h2>{{ form.id ? 'Edit Category' : 'New Category' }}</h2><fieldset :disabled="saving"><div class="expense-grid"><label>Name<input required v-model.trim="form.name" maxlength="120" class="expense-input" /></label><label>Description<input v-model.trim="form.description" maxlength="1000" class="expense-input" /></label><label><input type="checkbox" v-model="form.is_active" /> Active (available for new expenses)</label></div><div class="expense-actions mt-4"><button class="expense-button">{{ saving ? 'Saving…' : 'Save Category' }}</button><button v-if="form.id" type="button" @click="form=empty()" class="expense-button expense-secondary">Cancel Edit</button></div></fieldset></form>
    <div class="expense-card expense-table-wrap"><p v-if="loading" role="status">Loading categories…</p><table v-else><thead><tr><th>Name</th><th>Description</th><th>Status</th><th>Action</th></tr></thead><tbody><tr v-for="category in categories" :key="category.id"><td>{{ category.name }}</td><td>{{ category.description }}</td><td>{{ category.is_active ? 'Active' : 'Inactive' }}</td><td><button :disabled="saving" @click="edit(category)" class="expense-button expense-secondary" :aria-label="'Edit '+category.name">Edit</button></td></tr><tr v-if="!categories.length"><td colspan="4">No categories yet. Create the first category above.</td></tr></tbody></table></div>
  </ExpensePage>
</template>
<script setup>
import {ref,onMounted} from 'vue'
import ExpensePage from '@/components/expenses/ExpensePage.vue'
import client from '@/api/client'
import {expenseData,expenseError} from '@/api/expenses'
const empty=()=>({id:'',name:'',description:'',is_active:true}),form=ref(empty()),categories=ref([]),error=ref(''),notice=ref(''),loading=ref(false),saving=ref(false)
async function load(){loading.value=true;try{categories.value=expenseData(await client.get('/api/v1/expenses/categories'))}catch(e){error.value=expenseError(e)}finally{loading.value=false}}
function edit(category){form.value={...category};notice.value='';window.scrollTo({top:0,behavior:'smooth'})}
async function save(){if(saving.value)return;saving.value=true;error.value='';notice.value='';try{const {id,...payload}=form.value;await client[id?'put':'post']('/api/v1/expenses/categories'+(id?'/'+id:''),payload);form.value=empty();notice.value='Category saved.';await load()}catch(e){error.value=expenseError(e)}finally{saving.value=false}}
onMounted(load)
</script>
