<template>
    <Breadcrumbs main="UI Kits" title="avatars"/>
        <div class="container-fluid">
            <div class="user-profile">
              <div class="row">
                  <sizing/>
                  <shape/>
              </div>
            </div>
        </div>
</template>
<script>
import sizing from "./sizing.vue"
import shape from "./shape.vue"
export default {
    components:{
        sizing,
        shape
    }
}
</script>