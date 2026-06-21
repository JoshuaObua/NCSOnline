<template>
     <Breadcrumbs main="Project" title="Create New"/>
     <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-body">
                    <createProject1/>
                  </div>
                </div>
              </div>
            </div>
     </div>
</template>
<script>
import createProject1 from './createProject1.vue'
export default {
    components:{
        createProject1
    }
}
</script>