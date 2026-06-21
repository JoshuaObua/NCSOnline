<template>
    <Breadcrumbs sub="Tables" main="Bootstrap tables" title="Styling Tables"/>
         <div class="container-fluid">
            <div class="row">
                <defultStyling/>
                <tableHead/>
                <primaryBackground/>
                <tableFooter/>
                <customTable/>
                <hoverAndStripped/>
            </div>
         </div>
</template>
<script>
import defultStyling from "./defultStyling.vue"
import tableHead from "./tableHead.vue"
import primaryBackground from "./primaryBackground.vue"
import tableFooter from "./tableFooter.vue"
import customTable from "./customTable.vue"
import hoverAndStripped from "./hoverAndStripped.vue"
export default {
    components:{
        defultStyling,
        tableHead,
        primaryBackground,
        tableFooter,
        customTable,
        hoverAndStripped
    }
}
</script>