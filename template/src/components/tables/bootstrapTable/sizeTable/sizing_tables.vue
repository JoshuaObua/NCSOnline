<template>
     <Breadcrumbs sub="Tables" main="Bootstrap Tables" title="Sizing Tables"/>
         <div class="container-fluid">
            <div class="row">
                <extraLarge/>
                <largeTable/>
                <defaultTable/>
                <smallTable/>
                <extraSmall/>
            </div>
         </div>
</template>
<script>
import extraLarge from "./extraLarge.vue"
import largeTable from "./largeTable.vue"
import defaultTable from "./defaultTable.vue"
import smallTable from "./smallTable.vue"
import extraSmall from "./extraSmall.vue"
export default {
    components:{
        extraLarge,
        largeTable,
        defaultTable,
        smallTable,
        extraSmall
    }
}
</script>