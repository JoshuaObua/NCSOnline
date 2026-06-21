<template>
    <Breadcrumbs main="Buttons" title="Flat Style"/>
         <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
              <flat/>
              <flat2/>
              <flat3/>
              </div>
            </div>
         </div>
</template>
<script>
import flat from "./flat.vue"
import flat2 from "./flat2.vue"
import flat3 from "./flat3.vue"
export default {
    components:{
        flat,
        flat2,
        flat3
    }
}
</script>