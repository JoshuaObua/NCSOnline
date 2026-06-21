<template>
          <div class="col-lg-6">
                <div class="card">
                  <div class="card-body">
                    <div class="product-info">
                          <h4>Description</h4>
                      <form>
                        <div class="product-group">
                          <div class="row"> 
                            <div class="col-sm-12">
                              <div class="mb-3">
                                <label class="form-label">Product Name</label>
                                <input class="form-control" placeholder="Enter Product Name" type="text"><span class="text-danger"></span>
                              </div>
                            </div>
                          </div>
                          <div class="row"> 
                            <div class="col-sm-12">
                              <div class="mb-3">
                                <label class="form-label">Product Description</label>
                                <input class="form-control" placeholder="Enter Product Description" type="text"><span class="text-danger"></span>
                              </div>
                            </div>
                          </div>
                        </div>
                      </form>
                      <productCategories/>
                    </div>
                  </div>
                </div>
          </div>
</template>
<script>
import productCategories from "./productCategories.vue"
export default {
    components:{
        productCategories
    }
}
</script>