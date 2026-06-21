<template>
    <Breadcrumbs main="UI Kits" title="grid"/>
         <div class="container-fluid grid-wrrapper">
            <div class="row">
                <gridOptions/>
                <gridColumn/>
                <settingOne/>
                <verticalAlignment/>
                <horizontalAlignment/>
                <nestingColumn/>
                <order/>
                <offset/>
            </div>
         </div>
</template>
<script>
import gridOptions from "./gridOptions.vue"
import gridColumn from "./gridColumn.vue"
import settingOne from "./settingOne.vue"
import verticalAlignment from "./verticalAlignment.vue"
import horizontalAlignment from "./horizontalAlignment.vue"
import nestingColumn from "./nestingColumn.vue"
import order from "./order.vue"
import offset from "./offset.vue"
export default {
    components:{
        gridOptions,
        gridColumn,
        settingOne,
        verticalAlignment,
        horizontalAlignment,
        nestingColumn,
        order,
        offset
    }
}
</script>