<template>
    <Breadcrumbs main="UI Kits" title="Accordion"/>
        <div class="container-fluid">
            <div class="row">
                <basicAccordion/>
                <closeAccordion/>
                <colorAccordion/>
                <accordionIcon/>
                <openCloseIcon/>
            </div>
        </div>
</template>
<script>
import basicAccordion from "./basicAccordion.vue"
import closeAccordion from "./closeAccordion.vue"
import colorAccordion from "./colorAccordion.vue"
import accordionIcon from "./accordionIcon.vue"
import openCloseIcon from "./openCloseIcon.vue"
export default {
    components:{
        basicAccordion,
        closeAccordion,
        colorAccordion,
        accordionIcon,
        openCloseIcon
    }
}
</script>