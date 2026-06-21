<template>
    <Breadcrumbs title="Social App"/>
    <div class="container-fluid">
            <div class="user-profile social-app-profile">
              <div class="row"> 
                 <socialHeader/> 
              </div>
            </div>
    </div>
</template>
<script>

import socialHeader from "./socialHeader.vue"
export default {
    components:{
        socialHeader
    }
}
</script>