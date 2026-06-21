<template>
  <div :class="filtered ? 'container-fluid product-wrapper sidebaron' : 'container-fluid product-wrapper'">
    <div class="product-grid">
      <div class="feature-products">
        <div class="row m-b-10">
          <div class="col-md-3 col-sm-2 products-total">
            <div class="square-product-setting d-inline-block">
              <a class="icon-grid grid-layout-view" @click="gridView()">
                <vue-feather type="grid"></vue-feather>
              </a>
            </div>
            <div class="square-product-setting d-inline-block">
              <a class="icon-grid m-0 list-layout-view" @click="listView()">
                <vue-feather type="list"></vue-feather>
              </a>
            </div>
            <div class="grid-options d-inline-block">
              <ul>
                <li @click="grid2()">
                  <a class="product-2-layout-view">
                    <span class="line-grid line-grid-1 bg-primary"></span>
                    <span class="line-grid line-grid-2 bg-primary"></span>
                  </a>
                </li>
                <li @click="grid3()">
                  <a class="product-3-layout-view">
                    <span class="line-grid line-grid-3 bg-primary"></span>
                    <span class="line-grid line-grid-4 bg-primary"></span>
                    <span class="line-grid line-grid-5 bg-primary"></span>
                  </a>
                </li>
                <li @click="grid4()">
                  <a class="product-4-layout-view">
                    <span class="line-grid line-grid-6 bg-primary"></span>
                    <span class="line-grid line-grid-7 bg-primary"></span>
                    <span class="line-grid line-grid-8 bg-primary"></span>
                    <span class="line-grid line-grid-9 bg-primary"></span>
                  </a>
                </li>
                <li @click="grid6()">
                  <a class="product-6-layout-view">
                    <span class="line-grid line-grid-10 bg-primary"></span>
                    <span class="line-grid line-grid-11 bg-primary"></span>
                    <span class="line-grid line-grid-12 bg-primary"></span>
                    <span class="line-grid line-grid-13 bg-primary"></span>
                    <span class="line-grid line-grid-14 bg-primary"></span>
                    <span class="line-grid line-grid-15 bg-primary"></span>
                  </a>
                </li>
              </ul>
            </div>
            <span class="d-none-productlist filter-toggle" @click="collapseFilter()">
              <h6>Filters
                <span class="ms-2">
                  <vue-feather type="chevron-down" class="middle"></vue-feather>
                </span>
              </h6>
            </span>
          </div>
          <productLength />
        </div>
        <div class="row">
          <div class="col-md-12">
            <div class="pro-filter-sec">
              <div class="product-sidebar" :class="filtered ? 'open' : ''">
                <div class="filter-section">
                  <div class="card">
                    <div class="card-header">
                      <h6 class="mb-0 f-w-600 "> Filters
                        <span class="pull-right" @click="collapseFilter()">
                          <i class="fa fa-chevron-down toggle-data"></i>
                        </span>
                      </h6>
                    </div>
                    <Slider @allFilters="allfilter" @priceVal="pricefilterArray" />
                  </div>
                </div>
              </div>
              <div class="product-search">
                <form>
                  <div class="form-group m-0">
                    <input class="form-control" type="text" placeholder="Search.." />
                    <i class="fa fa-search"></i>
                  </div>
                </form>
                <div class="col-sm-12">
                  <div v-if="filterProduct.length == 0">
                    <div class="search-not-found text-center">
                      <p>Sorry, We didn't find any results matching this search</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <productView />
    </div>

  </div>
</template>
<script>


import { mapGetters } from 'vuex'
import productLength from "./productLength.vue"
import produdectFilters from "./produdectFilters.vue"
import productView from "./productView.vue"
import productModel from "./productModel.vue"
import Slider from "./filterBar.vue"
export default {
  components: {

    productLength,
    produdectFilters,
    productView,
    productModel,
    Slider
  },
  data() {
    return {
      modalShow: false,
      quickViewProduct: [],
      counter: 1,
      priceArray: [],
      allfilters: [],
      items: [],
      filtered: false,
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
    addToCart: function (product, qty) {
      product.quantity = qty ? qty : 1;
      this.$store.dispatch('products/addToCart', product);
    },
    quickView: function (product) {
      this.modalShow = true;
      return (this.quickViewProduct = product);
    },
    pricefilterArray(item) {
      this.$store.dispatch('products/priceFilter', item);
    },
    grid2() {
      this.$store.commit('products/grid2', true);
    },
    grid3() {
      this.$store.commit('products/grid3', true);
    },
    grid4() {
      this.$store.commit('products/grid4');
    },
    grid6() {
      this.$store.commit('products/grid6');
    },
    listView() {
      this.$store.commit('products/listView');
    },
    gridView() {
      this.$store.commit('products/grid4');
    },
  },
}
</script>