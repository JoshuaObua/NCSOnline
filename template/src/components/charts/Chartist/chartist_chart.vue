<template>
<Breadcrumbs main="Charts" title="Chartist Chart"/>
     <div class="container-fluid">
            <div class="row">
                <chartist/>
                <pathAnimation/>
                <donut/>
                <bipolar/>
                <linechart/>
                <barchart/>
                <stacked/>
                <horizontal/>
                <extreme/>
                <simplechart/>
                <holes/>
                <filledholes/>
            </div>
     </div>
</template>
<script>
import chartist from "./chartist.vue" 
import pathAnimation from "./pathAnimation.vue"
import donut from "./donut.vue"
import bipolar from "./bipolar.vue"
import linechart from "./linechart.vue"
import barchart from "./barchart.vue"
import stacked from "./stacked.vue"
import horizontal from "./horizontal.vue"
import extreme from "./extreme.vue"
import simplechart from "./simplechart.vue"
import holes from "./holes.vue"
import filledholes from "./filledholes.vue"
export default {
   components:{
     chartist,
     pathAnimation,
     donut,
     bipolar,
     linechart,
     barchart,
     stacked,
     horizontal,
     extreme,
     simplechart,
     holes,
     filledholes
   }
}
</script>