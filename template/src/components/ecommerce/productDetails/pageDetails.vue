<template>
  <div class="col-xl-5 box-col-6 proorder-xl-3 xl-100">
    <div class="card">
      <div class="card-body">
        <div class="product-page-details">
          <h3>{{ products.name }}</h3>
          <star-rating :star-size="20" :show-rating="false"></star-rating>
        </div>
        <div class="product-price">
          ${{ products.price }}
          <del>${{ products.salePrice }} </del>
        </div>
        <ul class="product-color ">
          <li class="bg-primary"></li>
          <li class="bg-secondary"></li>
          <li class="bg-success"></li>
          <li class="bg-danger"></li>
          <li class="bg-info"></li>
          <li class="bg-warning"></li>
        </ul>
        <div class="pro-group">
          <p>{{ products.description }}</p>
        </div>
        <div class="pro-group">
          <div class="row">
            <div class="col-md-6">
              <table>
                <tbody>
                  <tr>
                    <td> <b>Availability &nbsp;: &nbsp;</b></td>
                    <td class="txt-success">In stock</td>
                  </tr>
                  <tr>
                    <td> <b>Brand &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;:</b></td>
                    <td>Pixelstrap</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="col-md-6">
              <table>
                <tbody>
                  <tr>
                    <td><b>Seller &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;: &nbsp;&nbsp;&nbsp;</b></td>
                    <td>ABC</td>
                  </tr>
                  <tr>
                    <td> <b>Fabric &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;: &nbsp;&nbsp;&nbsp;</b></td>
                    <td>Cotton </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="pro-group">
              <div class="row">
                <div class="col-md-4 xl-50">
                  <h6 class="product-title">share it</h6>
                </div>
                <div class="col-md-7 xl-50">
                  <div class="product-icon">
                    <ul class="product-social">
                      <li class="d-inline-block"><a href="https://www.facebook.com"><i class="fa fa-facebook"></i></a>
                      </li>
                      <li class="d-inline-block"><a href="https://accounts.google.com"><i
                            class="fa fa-google-plus"></i></a></li>
                      <li class="d-inline-block"><a href="https://twitter.com"><i class="fa fa-twitter"></i></a></li>
                      <li class="d-inline-block"><a href="https://www.instagram.com"><i class="fa fa-instagram"></i></a>
                      </li>
                      <li class="d-inline-block"><a href="https://dashboard.rss.com/auth/sign-in"><i
                            class="fa fa-rss"></i></a></li>
                    </ul>
                    <form class="d-inline-block f-right"></form>
                  </div>
                </div>
              </div>

            </div>
            <div class="m-t-15">
              <router-link :to="'/ecommerce/cartView'">
                <button class="btn btn-primary m-r-10" type="button" @click="addToCart(products)">
                  <i class="fa fa-shopping-basket me-2"></i>Add To Cart
                </button>
              </router-link>
              <router-link :to="'/ecommerce/wishList'">
                <button class="btn btn-danger m-r-10" type="button">
                  <i class="fa fa-heart me-2"></i>Add To Wishlist
                </button>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>

import { mapState } from 'vuex';
import router from '@/router';
import StarRating from 'vue-star-rating';
export default {
  components: {
   
    StarRating
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
}
</script>