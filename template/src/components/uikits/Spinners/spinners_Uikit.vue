<template>
     <Breadcrumbs main="UI Kits" title="Spinners"/>
              <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Loader Style</h5>
                  </div>
                  <div class="card-body row">
                      <loader/>
                      <loader2/>
                      <loader3/>
                  </div>
                </div>
              </div>
            </div>
              </div>
</template>
<script>
import loader from "./loader.vue"
import loader2 from "./loader2.vue"
import loader3 from "./loader3.vue"
export default {
    components:{
        loader,
        loader2,
        loader3
    }
}
</script>