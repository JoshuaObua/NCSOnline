<template>
     <Breadcrumbs main="Buttons" title="Edge Style"/>
          <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                  <edge/>
                  <edge2/>
                  <edge3/>
              </div>
            </div>
          </div>
</template>
<script>
import edge from "./edge.vue"
import edge2 from "./edge2.vue"
import edge3 from "./edge3.vue"
export default {
    components:{
        edge,
        edge2,
        edge3
    }
}
</script>