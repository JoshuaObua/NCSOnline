<template>
     <Breadcrumbs main="Ecommerce" title="Add Product" />
      <div class="container-fluid add-product">
            <div class="row"> 
                <productDescription/>
                <productImages/>
            </div>
      </div>
</template>
<script>
import productDescription from "./productDescription.vue"
import productImages from "./productImages.vue"
export default {
    components:{
        productDescription,
        productImages
    }
}
</script>