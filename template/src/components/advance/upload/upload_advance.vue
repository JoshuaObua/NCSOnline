<template>
  <Breadcrumbs main="Bonus UI" title="Upload"/>
    <div class="container-fluid ">
      <div class="row">
        <singleFile/>
         <multiFile/>
       <fileTypeValidation/>
      </div>
    </div>
</template>
<script>
import singleFile from "./singleFile.vue"
import multiFile from "./multiFile.vue"
import fileTypeValidation from "./fileTypeValidation.vue"
export default {
    components:{
         multiFile,
         singleFile,
        fileTypeValidation
    },
  data() {
    return {
       
    };
  },
};
</script>