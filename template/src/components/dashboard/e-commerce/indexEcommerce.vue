<template>
         <Breadcrumbs main="Dashboard" title="E-Commerce"/>
         <invoiceOverview/>
    
</template>
<script>
import invoiceOverview from '../e-commerce/invoiceOverview.vue'

export default {
    components:{
        invoiceOverview,
    }
}
</script>