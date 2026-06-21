<template>
  <Breadcrumbs main="Gallery" title="Masonry Gallery" />
    <div class="container-fluid">
      <div class="row">
         <masonry/>
      </div>
    </div>
</template>
<script>
import masonry from "./masonry.vue"
export default {
  components:{
    masonry
  },
    
};
</script>