<template>
     <div class="row">
            <div class="col-sm-3">
              <div class="product-sidebar" :class="filtered?'open':''">
                <div class="filter-section">
                  <div class="card">
                    <div class="card-header">
                      <h6 class="mb-0 f-w-600">
                        Filters
                        <span class="pull-right" @click="collapseFilter()">
                          <i class="fa fa-chevron-down toggle-data"><vue-feather type="chevron-down" class="middle"></vue-feather></i>
                        </span>
                      </h6>
                    </div>
                    <Slider @allFilters="allfilter" @priceVal="pricefilterArray" />
                  </div>
                </div>
              </div>
            </div>
            <div class="col-xl-9 col-sm-12">
              <form>
                <div class="form-group m-0">
                  <input
                    class="form-control"
                    type="text"
                    placeholder="Search.."
                  />
                  <i class="fa fa-search"></i>
                </div>
              </form>
              <div class="col-sm-12">
                <div v-if="filterProduct.length == 0" >
                  <div class="search-not-found text-center">
                   
                    <p>Sorry, We didn't find any results matching this search</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
</template>
<script>

 import { mapGetters } from 'vuex';
 import Slider from "./filterBar.vue"
export default {
    components:{
     
        Slider
    },
    data(){
        return{
        modalShow: false,
        quickViewProduct: [],
        counter: 1,
        priceArray: [],
        allfilters: [],
        items: [],
        filtered: false,
        listViewEnable: false,
        list: false
        }
    },
 computed: {
      ...mapGetters({
        filterProduct: 'products/filterProducts',
        tags: 'products/setTags'
      }),
    },
    methods: {
      getImgUrl(path) {
        return require('../../../assets/images/' + path);
      },
      onChangeSort(event) {
        this.$store.dispatch('products/sortProducts', event.target.value);
      },
      allfilter(selectedVal) {
        this.allfilters = selectedVal;
        this.$store.dispatch('products/setTags', selectedVal);
      },
     
      collapseFilter() {
        this.filtered = !this.filtered;
      },
      addToCart: function(product, qty) {
        product.quantity = qty ? qty : 1;
        this.$store.dispatch('products/addToCart', product);
      },
      quickView: function(product) {
        this.modalShow = true;
        return (this.quickViewProduct = product);
      },
    },
}
</script>