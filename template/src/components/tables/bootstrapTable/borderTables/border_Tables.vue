<template>
     <Breadcrumbs sub="Tables" main="Bootstrap Tables" title="Border Tables"/>
          <div class="container-fluid">
            <div class="row">
                <horizontalBorders/>
                <verticalBorders/>
                <bothBordeds/>
                <borderlessTable/>
                <defaultTable/>
                <solidBorder/>
                <borderBottomColor/>
            </div>
          </div>
</template>
<script>
import horizontalBorders from "./horizontalBorders.vue"
import verticalBorders from "./verticalBorders.vue"
import borderlessTable from "./borderlessTable.vue"
import bothBordeds from "./bothBordeds.vue"
import defaultTable from "./defaultTable.vue"
import solidBorder from "./solidBorder.vue"
import borderBottomColor from "./borderBottomColor.vue"
export default {
    components:{
        horizontalBorders,
        verticalBorders,
        borderlessTable,
        bothBordeds,
        defaultTable,
        solidBorder,
        borderBottomColor
    }
}
</script>