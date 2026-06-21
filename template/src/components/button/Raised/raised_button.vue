<template>
    <Breadcrumbs main="Buttons" title="Raised Style"/>
         <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                  <raised/>
                  <raised2/>
                  <raised3/>
              </div>
            </div>
         </div>
</template>
<script>
import raised from "./raised.vue"
import raised2 from "./raised2.vue"
import raised3 from "./raised3.vue"
export default {
    components:{
        raised,
        raised2,
        raised3
    }
}
</script>