<template>
    <div class="container-fluid">
        <div class="row">
              <Breadcrumbs main="UI Kits" title="Progress"/>
                  <basicProgress/>
                  <largeProgress/>
                  <smallProgress/>
                  <customHeight/>
                  <progressStriped/>
                  <progressAnimated/>
                  <progressStates/>
                  <multipleBars/>
            </div>
         </div>
</template>
<script>
import basicProgress from "./basicProgress.vue"
import smallProgress from "./smallProgress.vue"
import largeProgress from "./largeProgress.vue"
import customHeight from "./customHeight.vue"
import progressStates from "./progressStates.vue"
import multipleBars from "./multipleBars.vue"
import progressStriped from "./progressStriped.vue"
import progressAnimated from "./progressAnimated.vue"
export default {
    components:{
        basicProgress,
        smallProgress,
        largeProgress,
        customHeight,
        progressStates,
        multipleBars,
        progressStriped,
        progressAnimated
    }
}
</script>