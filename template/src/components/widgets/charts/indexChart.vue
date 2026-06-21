<template>
    <Breadcrumbs main="Widgets" title="chart" />
   <div class="container-fluid chart-widget">
            <div class="row">
                <saleChart/>
                <projectsChart/>
                <productChart/>
                <marketingExpenses/>
                <totalEarning/>
                <liveProducts/>
                <turnoverChart/>
                <monthlySales/>
                <usesChart/>
                <browserUses/>
            </div>
   </div>
</template>
<script>
import saleChart from './saleChart.vue';
import projectsChart from './projectsChart.vue';
import productChart from './productChart.vue';
import marketingExpenses from './marketingExpenses.vue';
import totalEarning from './totalEarning.vue';
import liveProducts from './liveProducts.vue';
import turnoverChart from './turnoverChart.vue';
import monthlySales from './monthlySales.vue';
import usesChart from './usesChart.vue';
import browserUses from './browserUses.vue'
export default {
    components:{
        saleChart,
        projectsChart,
        productChart,
        marketingExpenses,
        totalEarning,
        liveProducts,
        turnoverChart,
        monthlySales,
        usesChart,
        browserUses
    }
}
</script>