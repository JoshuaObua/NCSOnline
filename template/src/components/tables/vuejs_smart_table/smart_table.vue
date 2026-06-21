<template>
    <Breadcrumbs main="Tables" title="VueJs Smart Table"/>
    <div class="container-fluid ">
			<div class="row">
             <basic/>
             <filterView/>
             <agefilter/>
             <rowSelect/>
            </div>
    </div>
</template>
<script>
import basic from "./basic.vue"
import filterView from "./filterView.vue"
import agefilter from "./agefilter.vue"
import rowSelect from "./rowSelect.vue"
export default {
    components:{
        basic,
        filterView,
        agefilter,
        rowSelect
    }
}
</script>