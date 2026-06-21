<template>
     <Breadcrumbs main="UI Kits" title="Tag & Pills"/>
          <div class="container-fluid">
            <div class="row">
                <contextualVariations/>
                <tagPillsNumber/>
                <tagPillsIcon/>
                <badgesExample/>
                <badgesButtons/>
            </div>
          </div>
</template>
<script>
import contextualVariations from "./contextualVariations.vue"
import tagPillsNumber from "./tagPillsNumber.vue"
import tagPillsIcon from "./tagPillsIcon.vue"
import badgesExample from "./badgesExample.vue"
import badgesButtons from "./badgesButtons.vue"
export default {
    components:{
        contextualVariations,
        tagPillsNumber,
        tagPillsIcon,
        badgesExample,
        badgesButtons
    }
}
</script>