<template>
    <Breadcrumbs sub="Forms" main="Form Widgets" title="Select2"/>
   
         <div class="container-fluid ">
            <div class="select2-drpdwn">
              <div class="row">
                 <selectView/>
                 <Outline/>
                 <fullcolored/>
              </div>
            </div>
         </div>
</template>
<script>
import Outline from "./Outline.vue"
import fullcolored from "./fullcolored.vue"
import selectView from "./select.vue"
export default {
   components:{
     Outline,
     fullcolored,
     selectView
   }
}
</script>