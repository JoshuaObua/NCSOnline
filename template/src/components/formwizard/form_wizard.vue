<template>
     <Breadcrumbs main="Forms" title="Form Wizard"/>
          <div class="container-fluid ">
            <div class="row">
                <SimpleWizard/>
            </div>
          </div>
</template>
<script>
import SimpleWizard from "./SimpleWizard.vue"
export default {
    components:{
        SimpleWizard
    }
}
</script>