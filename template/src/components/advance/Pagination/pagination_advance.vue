<template>
    <Breadcrumbs main="Bonus UI" title="pagination"/>
  
         <div class="container-fluid">
            <div class="row">
                <pagination/>
                <paginationIcons/>
                <paginationAlignment/>
                <activeAndDisabled/>
                <paginationColor/>
                <paginationSizing/>
            </div>
         </div>
   
</template>
<script>
import pagination from "./pagination.vue"
import paginationIcons from "./paginationIcons.vue"
import paginationAlignment from "./paginationAlignment.vue"
import activeAndDisabled from "./activeAndDisabled.vue"
import paginationColor from "./paginationColor.vue"
import paginationSizing from "./paginationSizing.vue"
export default {
    components:{
        pagination,
        paginationIcons,
        paginationAlignment,
        activeAndDisabled,
        paginationColor,
        paginationSizing
    }
}
</script>