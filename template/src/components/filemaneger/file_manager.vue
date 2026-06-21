<template>
  <Breadcrumbs title="File Manager"/>
 <div class="container-fluid">
      <div class="row">
        <fileManeger1 />
        <fileManeger2 />
      </div>
 </div>
</template>
<script>
import fileManeger1 from "./fileSidebar.vue";
import fileManeger2 from "./fileManeger2.vue";
export default {
  components: {
    fileManeger1,
    fileManeger2,
  },
};
</script>