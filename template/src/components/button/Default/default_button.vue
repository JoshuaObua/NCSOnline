<template>
    <Breadcrumbs main="Buttons" title="Default Style"/>
          <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                  <Default/>
                  <Default2/>
                  <Default3/>
              </div>
            </div>
          </div>
</template>
<script>
import Default from "./Default.vue"
import Default2 from "./Default2.vue"
import Default3 from "./Default3.vue"
export default {
    components:{
        Default,
        Default2,
        Default3
    }
}
</script>