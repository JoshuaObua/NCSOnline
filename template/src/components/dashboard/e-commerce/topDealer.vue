<template>
  <div class="col-xl-6 box-col-12 des-xl-100 top-dealer-sec">
    <div class="card">
      <div class="card-header pb-0">
        <div class=" header-top d-sm-flex justify-content-between align-items-center">
          <h5>Top Dealer</h5>
          <div class="center-content">
            <p class="d-sm-flex align-items-center">
              <span class="m-r-10">845 Dealer</span><i
                class="toprightarrow-primary fa fa-arrow-up m-r-10 d-grid"></i>86%
              More Than Last Year
            </p>
          </div>
        </div>
      </div>
      <div class="card-body pb-3">
        <div class="col-12">
          <div class="row">
            <swiper :space-between="20" :breakpoints="swiperOptions.breakpoints" @swiper="setSecondSwiper"
              :controller="{ control: firstSwiper }" :modules="modules" :loop="true">
              <swiper-slide v-for="card in cards" :key="card">
                <div class="item">
                  <div class="card">
                    <div class="top-dealerbox text-center">
                      <img class="card-img-top" :src="getImageUrl(card.src1)" alt="..." />
                      <h6>{{ card.name }}</h6>
                      <p>{{ card.name1 }}</p>
                      <router-link to="/pages/social-app" class="btn btn-rounded">View More</router-link>
                    </div>
                  </div>
                </div>

              </swiper-slide>
            </swiper>
            <swiper :space-between="20" :pagination="{
              clickable: true,
            }" :modules="modules" @swiper="setFirstSwiper" :breakpoints="swiperOptions.breakpoints"
              :controller="{ control: secondSwiper }" :loop="true">
              <swiper-slide v-for="card in cards1" :key="card">

                <div class="item">
                  <div class="card">
                    <div class="top-dealerbox text-center">
                      <img class="card-img-top" :src="getImageUrl(card.src1)" alt="..." />
                      <h6>{{ card.name }}</h6>
                      <p>{{ card.name1 }}</p>
                      <router-link to="/pages/social-app" class="btn btn-rounded">View More</router-link>
                    </div>
                  </div>
                </div>

              </swiper-slide>

            </swiper>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { ref } from 'vue';
import { Swiper, SwiperSlide } from "swiper/vue";
import { Controller } from 'swiper';
import "swiper/css";

import "swiper/css/pagination";

import { Pagination } from "swiper";
export default {
  components: {
    Swiper,
    SwiperSlide,
  },
  data() {
    return {
      allSlides: [1, 2, 3, 4, 5],
      cards: [
        { id: "1", src1: "1.png", name: "Thompson lee", name1: "Malasiya" },
        { id: "1", src1: "8.png", name: "Johnson allon", name1: "Bangladesh" },
        { id: "1", src1: "3.png", name: "williams reed", name1: "Belgium" },
        { id: "1", src1: "4.png", name: "Joneson king", name1: "Canada" },
      ],
      cards1: [
        { id: "1", src1: "5.png", name: "Brown davis", name1: "China" },
        { id: "1", src1: "6.png", name: "Wilsons Hill", name1: "Denmark" },
        { id: "1", src1: "7.png", name: "Anderson ban", name1: "Japan" },
        { id: "1", src1: "8.png", name: "Thompson lee", name1: "Malasiya" },
      ],
      swiperOptions: {
        breakpoints: {
          1669: {
            slidesPerView: 4,

          },
          480: {
            slidesPerView: 3,

          },
          0: {
            slidesPerView: 2,

          },
        }
      },
    };
  },
  methods: {
    getImageUrl(path) {
      return require("../../../assets/images/dashboard-2/" + path);

    }
  },
  setup() {
    const firstSwiper = ref(null);
    const secondSwiper = ref(null);
    const setFirstSwiper = (swiper) => {
      firstSwiper.value = swiper;
    };
    const setSecondSwiper = (swiper) => {
      secondSwiper.value = swiper;
    };
    return {
      modules: [Pagination, Controller],
      firstSwiper,
      secondSwiper,
      setFirstSwiper,
      setSecondSwiper,

    };
  },
};
</script>
