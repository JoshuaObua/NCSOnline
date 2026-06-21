<template>
    <Breadcrumbs main="Buttons" title="Button Group"/>
          <div class="container-fluid">
            <div class="row">
                <group/>
                <edges/>
                <flat/>
                <large/>
                <largeEdges/>
                <outline/>
                <outline2/>
                <outline3/>
                <outlineEdges/>
                <outlineFlat/>
                <radio/>
                <checkBox/>
                <nesting/>
            </div>
          </div>
</template>
<script>
import group from "./group.vue"
import edges from "./edges.vue"
import flat from "./flat.vue"
import large from "./large.vue"
import largeEdges from "./largeEdges.vue"
import outline from "./outline.vue"
import outline2 from "./outline2.vue"
import outline3 from "./outline3.vue"
import outlineEdges from "./outlineEdges.vue"
import outlineFlat from "./outlineFlat.vue"
import radio from "./radio.vue"
import checkBox from "./checkBox.vue"
import nesting from "./nesting.vue"
export default {
    components:{
        group,
        edges,
        flat,
        large,
        largeEdges,
        outline,
        outline2,
        outline3,
        outlineEdges,
        outlineFlat,
        radio,
        checkBox,
        nesting
    }
}
</script>