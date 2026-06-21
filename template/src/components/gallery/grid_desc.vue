<template>
    <Breadcrumbs main="Gallery" title="Gallery Grid With Description" />
                  <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>IMAGE GALLERY WITH DESCRIPTION</h5>
                  </div>
                  <easyLightbox/>
                </div>
              </div>
            </div>
          </div>  
</template>
<script>
import easyLightbox from "../easyLightbox.vue"
export default {
    components:{
        easyLightbox
    },
}
</script>