<template>
    <Breadcrumbs main="UI Kits" title="Modal"/>
 <div class="container-fluid">
            <div class="row">
               <staticExample/>
               <basicModal/>
               <varyingModal/>
               <sizesModal/>
            </div>
 </div>
</template>
<script>
import basicModal from "./basicModal.vue"
import varyingModal from "./varyingModal.vue"
import sizesModal from "./sizesModal.vue"
import staticExample from "./staticExample.vue"
export default {
   components:{
       basicModal,
       varyingModal,
       sizesModal,
       staticExample
   }
}
</script>