<template>
  <div class="left-filter ">
    <div class="card-body filter-cards-view animate-chk">
      <div class="product-filter">
        <h6 class="f-w-600">Category</h6>
        <div class="checkbox-animated mt-0">
          <label class="d-block" for="edo-ani5" v-for="(product, index) in getCategory" :key="index">
            <input @change="appliedFilter(product.category)" :value="product.category" v-model="applyFilter"
              class="checkbox_animated" :id="product.category" type="checkbox">{{ product.category }}
          </label>
        </div>
      </div>
      <div class="product-filter">
        <h6 class="f-w-600">Brand</h6>
        <div class="checkbox-animated mt-0">
          <label class="d-block" v-for="(brand, index) in getBrand" :key="index">
            <input @change="appliedFilter(brand)" class="checkbox_animated" :value="brand" v-model="applyFilter"
              :id="brand" type="checkbox"> {{ brand }}
          </label>
        </div>
      </div>
      <div class="product-filter slider-product">
        <h6 class="f-w-600">Colors</h6>
        <div class="color-selector">
          <ul>
            <li id="white"></li>
            <li class="gray"></li>
            <li class="black"></li>
            <li class="orange"></li>
            <li class="green"></li>
            <li class="pink"></li>
            <li class="yellow
        list: false"></li>
            <li class="blue"></li>
            <li class="red"></li>
          </ul>
        </div>
      </div>
      <div class="product-filter pb-0">
        <h6 class="f-w-500">Price</h6>
        <VueSlider class="my-5" v-model="one.value" :data="one.data" :marks="true" :tooltip="'always'"
          :tooltip-placement="'top'" />

      </div>
      <newProduct />
    </div>
  </div>
</template>
<script>
import { mapGetters, mapState } from 'vuex';
import { Swiper, SwiperSlide } from "swiper/vue";
import newProduct from "./newProduct.vue"
import VueSlider from 'vue-slider-component'
import 'vue-slider-component/theme/antd.css'

import "swiper/css";
export default {
  components: {
    Swiper,
    SwiperSlide,

    newProduct,
    VueSlider
  },
  data() {
    return {
      filter: false,
      value: [0, 450],
      applyFilter: [],
      swiperOption: {
        loop: true,
        navigation: {
          nextEl: '.swiper-button-next',
          prevEl: '.swiper-button-prev'
        }
      },
      one: {
        value: [200, 800],
        data: [10, 100, 200, 250, 500, 750, 800, 1000],
      },
    };
  },
  computed: {
    ...mapState({
      products: state => state.products.products
    }),
    ...mapGetters({
      getCategory: 'products/getCategory',
      getBrand: 'products/getBrands',
      getColors: 'products/getColors'
    })
  },
  mounted() {
    this.$emit('priceVal', this.value);
  },
  methods: {
    isActive(filterItem) {
      return this.applyFilter.indexOf(filterItem) > -1;
    },

    appliedFilter() {
      this.$emit('allFilters', this.applyFilter);
    },
    sliderChange(event) {
      this.$emit('priceVal', event);
    },
    getImgUrl(path) {
      return require('../../../assets/images/' + path);
    },
  }
}
</script>