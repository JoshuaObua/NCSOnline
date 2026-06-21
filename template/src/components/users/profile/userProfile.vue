<template>
     <Breadcrumbs main="Users" title="Users Profile"/>
     <div class="container-fluid">
            <div class="user-profile">
              <div class="row">
                  <profileHeader/>
              </div>
            </div>
     </div>
</template>
<script>
import profileHeader from "./profileHeader.vue"
export default {
    components:{
        profileHeader,
    }
}
</script>