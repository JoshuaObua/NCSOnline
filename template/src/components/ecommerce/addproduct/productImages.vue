<template>
     <div class="col-lg-6">
                <div class="card"> 
                  <div class="card-body">
                    <div class="product-info">
                      <h4>Product Image</h4>
                     <DropZone
                      :maxFileSize="Number(60000000)"
                      :uploadOnDrop="true"
                      :dropzoneMessageClassName="customClass"
                     > 
                     </DropZone>  
                      <h4 class="mt-4">Select Size</h4>
                      <form>
                        <div class="product-group">
                          <div class="row"> 
                            <div class="col-sm-6"> 
                              <div class="mb-3">  
                                <label class="form-label">Size</label>
                                <select class="form-select">
                                  <option>S </option>
                                  <option>M </option>
                                  <option>L </option>
                                  <option>XL </option>
                                  <option>XXL</option>
                                </select>
                              </div>
                            </div>
                            <div class="col-sm-6"> 
                              <div class="mb-3">  
                                <label class="form-label">Date</label>
                                <input class="form-control" type="date">
                              </div>
                            </div>
                          </div>
                          <div class="row"> 
                            <div class="col-sm-12">
                              <div class="mb-3">  
                                <label class="form-label">Price</label>
                                <input class="form-control" type="number">
                              </div>
                            </div>
                          </div>
                        </div>
                      </form>
                    </div>
                    <div class="mt-4">
                      <div class="row"> 
                        <div class="col-sm-12 text-end"><router-link class="btn btn-primary me-3" :to="'/ecommerce/indexProduct'"> ADD </router-link>
                        <a class="btn btn-secondary">Cancel</a></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
</template>
<script>
import DropZone from "dropzone-vue";
export default {
  components:{
    DropZone
  },
}
</script>
<style scoped>
@import 'dropzone-vue/dist/dropzone-vue.common.css';
</style>