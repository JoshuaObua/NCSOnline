<template>
     <div class="col-md-9 col-sm-10 text-end">
             
                <span class="f-w-600 m-r-5">Showing Products 1 - {{filterProduct.length}} Results</span>
                <div class="select2-drpdwn-product select-options d-inline-block">
                  <select
                    class="form-control btn-square"
                    name="select"
                    @change="onChangeSort($event)"> 
                      <option>Order Product By</option>
                      <option value="a-z">Alphabetically, A-Z</option>
                      <option value="z-a">Alphabetically, Z-A</option>
                      <option value="low">price, low to high</option>
                      <option value="high">price, high to low</option>
                  </select>
                </div>
             
            </div>
</template>
<script>
import { mapGetters } from 'vuex';
export default {
     data() {
      return {
        modalShow: false,
        quickViewProduct: [],
        counter: 1,
        priceArray: [],
        allfilters: [],
        items: [],
        filtered: false,
        col2: false,
        col3: false,
        col4: true,
        col6: false,
        listViewEnable: false,
        list: false
      };
    },
    computed: {
      
      ...mapGetters({
        filterProduct: 'products/filterProducts',
        tags: 'products/setTags'
      }),
    },
    methods: {
      getImgUrl(path) {
        return require('../../../assets/images/ecommerce' + path);
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
      pricefilterArray(item) {
        this.$store.dispatch('products/priceFilter', item);
      },
      addToCart: function(product, qty) {
        product.quantity = qty ? qty : 1;
        this.$store.dispatch('products/addToCart', product);
      },
      quickView: function(product) {
        this.modalShow = true;
        return (this.quickViewProduct = product);
      },
      quickViewClose: function() {
        this.modalShow = false;
      },
    }
}
</script>