<template>
  <div class="container-fluid ">
    <div class="row product-page-main p-0">
      <div class="col-xl-5 col-md-6 box-col-12 xl-50">
        <div class="card">
          <div class="card-body">
            <div class="row">
              <div class="col-xl-9 product-main">
                <div class="product-slider" id="sync1">
                  <div class="pro-slide-single">
                    <swiper :slidesPerView="1" :loop="true" :thumbs="{ swiper: thumbsSwiper }"
                      :autoplay="{ delay: 3500, disableOnInteraction: false, }" :spaceBetween="30"
                      :centeredSlides="true" :modules="modules" class="mySwiper">
                      <swiper-slide v-for="(product, index) in products.images" :key="index">
                        <img :src="getImgUrl(product)" :id="index" class="img-fluid bg-img" alt="index" />
                      </swiper-slide>
                    </swiper>
                  </div>
                </div>
              </div>
              <div class="col-xl-3 product-thumbnail ">
                <div class="pro-slide-right">
                  <div>
                    <swiper :breakpoints="swiperOptions.breakpoints" @swiper="setThumbsSwiper" :loop="true"
                      :autoHeight="true" :slidesPerView="3" :spaceBetween="3" :watchSlidesProgress="true" :pagination="{
                        clickable: true,
                      }" :modules="modules" class="Swiper">
                      <swiper-slide v-for="(product, index) in products.images" :key="index">
                        <img :src="getImgUrl(product)" :id="index" class="img-fluid bg-img" alt="index"
                          style="height:auto; " />
                      </swiper-slide>
                    </swiper>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <pageDetails />
      <filterBox />
    </div>
    <tablistView />
  </div>
</template>
<script>
import { mapState } from 'vuex';
import router from '@/router';
import { Swiper, SwiperSlide } from "swiper/vue";
import pageDetails from "./pageDetails.vue"
import filterBox from "./filterBox.vue"
import tablistView from "./tablistView.vue"
import { FreeMode, Autoplay, Thumbs } from 'swiper';
export default {
  name: 'ProductDetails',
  components: {
    Swiper,
    SwiperSlide,
    pageDetails,
    filterBox,
    tablistView
  },
  data() {
    return {
      counter: 1,
      modalShow: false,
      swiperOptions: {
        breakpoints: {
          1215: {
            direction: 'vertical'
          }
        }
      },
      swiperOption1: {
      },
      activeTab: 'fabric',
    };
  },
  computed: {
    ...mapState({
      products: (state) =>
        state.products.products.find((item) => {
          if (parseInt(router.currentRoute._rawValue.params.id) === item.sku) return item;
        }),
    }),
    swiper() {
      return this.$refs.mySwiper.swiper;
    },
  },
  methods: {
    getImgUrl(path) {
      return require('../../../assets/images/' + path);
    },
    slideTo(id) {
      this.swiper.slideTo(id, 1000, false);
    },
    addToCart: function (product, qty) {
      product.quantity = qty ? qty : 1;
      this.$store.dispatch('products/addToCart', product);
    },
    increment() {
      if (this.counter < this.quickViewProduct.stock) this.counter++;
    },
    decrement() {
      if (this.counter > 1) this.counter--;
    },
    tabContent(val) {
      this.activeTab = val;
    },
  },
  setup() {
    return {
      modules: [FreeMode, Autoplay, Thumbs]
    };
  }
};
</script>