<template>
    <Breadcrumbs main="Bonus UI" title="ribbons"/>
     <div class="container-fluid">
            <div class="row">
                <coloredBreadcrumb/>
            </div>
            <ribbons/>
            <ribbons2/>
            <ribbons3/>
            <ribbons4/>
            <ribbons5/>
            <ribbons6/>
            </div>
</template>
<script>
import coloredBreadcrumb from "./coloredBreadcrumb.vue"
import ribbons from "./ribbons.vue"
import ribbons2 from "./ribbons2.vue"
import ribbons3 from "./ribbons3.vue"
import ribbons4 from "./ribbons4.vue"
import ribbons5 from "./ribbons5.vue"
import ribbons6 from "./ribbons6.vue"
export default {
    components:{
        coloredBreadcrumb,
        ribbons,
        ribbons2,
        ribbons3,
        ribbons4,
        ribbons5,
        ribbons6
    }
}
</script>