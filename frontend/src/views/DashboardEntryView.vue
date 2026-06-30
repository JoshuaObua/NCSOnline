<template><OrganisationDashboardView v-if="organisation"/><DashboardView v-else/></template>
<script setup>
import DashboardView from './DashboardView.vue'; import OrganisationDashboardView from './OrganisationDashboardView.vue'
let organisation=null; try{organisation=JSON.parse(localStorage.getItem('ncsms_active_organisation')||'null')}catch{organisation=null}
</script>
