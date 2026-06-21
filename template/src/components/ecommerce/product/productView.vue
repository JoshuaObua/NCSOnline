<template>
    <div class="product-wrapper-grid " :class="listViewEnable ? 'list-view' : ''">
    <div class="row">
      <div
        :class="[ col2 ? 'col-md-6' : col3 ? 'col-xl-4 col-sm-4' : col4 ? 'col-xl-3 col-md-6' : col6 ? 'col-xl-2' : list ? 'col-xl-12' : 'col-sm-6 col-xl-3 col-sm-3 col-xl-12',]"
        v-for="(product, index) in filterProduct"
        :key="index">
        <div class="card"> 
          <div class="product-box">
            <div class="product-img">
              <div class="ribbon ribbon-danger" v-if="product.sale">Sale</div>
              <img class="img-fluid" :src="getImgUrl(product.images[0])" alt />
              <div class="product-hover">
                <ul>
                  <router-link :to="'/ecommerce/cartView'">
                    <li @click="addToCart(product)">
                      <i class="icon-shopping-cart"></i>
                    </li>
                  </router-link>
                  <li
                    @click="quickView(product)"
                    data-bs-toggle="modal"
                    data-bs-target="#exampleModal">
                    <i class="icon-eye"></i>
                  </li>
                </ul>
              </div>
            </div> 
            <div class="modal fade"  id="exampleModal" aria-hidden="true">
         <div class="modal-dialog modal-lg modal-dialog-centered" >
         <div class="modal-content">
        <div class="modal-header">
          <div class="product-box row">
            <div class="product-img col-lg-6">
               <img class="img-fluid" v-if="quickViewProduct.images" :src="getImgUrl(quickViewProduct.images[0])" alt />
            </div>
            <div class="product-details col-lg-6">
              <h4 @click.prevent="this.$router.push('/ecommerce/details/'+quickViewProduct.sku)" data-bs-dismiss="modal">{{ quickViewProduct.name }}</h4>
              <div class="product-price">
                ${{ quickViewProduct.price }}
                <del>${{ quickViewProduct.salePrice }} </del>
              </div>
              <div class="product-view">
                <h6 class="f-w-600">Product Details</h6>
                <p class="mb-0">Cotton Solid Men's Fashion Full Sleeves Latest Jacket for Men With Button Closure Long Sleeve Casual Torn Lycra Denim Jacket.</p>
              </div>
              <div class="product-size">
                 <ul>
                   <li> <button class="btn btn-outline-light" type="button" data-bs-original-title="" title="">M</button></li>
                   <li> 
                    <button class="btn btn-outline-light" type="button" data-bs-original-title="" title="">L</button>
                     </li>
                      <li> <button class="btn btn-outline-light" type="button" data-bs-original-title="" title="">Xl</button> </li>
                    </ul>
                </div>
              <div class="product-qnty">
                <h6 class="f-w-600">Quantity</h6>
               <fieldset class="qty-box">
                <div class="input-group bootstrap-touchspin">
                <button class="btn btn-primary btn-square bootstrap-touchspin-down" v-on:click="decrement(item)" type="button"><i class="fa fa-minus"></i></button>
                 <input class="touchspin text-center form-control" v-model="counter" name="item.quantity" type="text">
                 <button class="btn btn-primary btn-square bootstrap-touchspin-up" v-on:click="increment(item)" type="button">
                  <i class="fa fa-plus"></i></button>
                                  </div>
                                </fieldset>
                <div class="addcart-btn mt-3">
                  <router-link to="/ecommerce/cartView" >
                    <button class="btn btn-primary m-r-10" type="button" data-original-title="btn btn-info-gradien" title  @click="addToCart(quickViewProduct, counter)" data-bs-dismiss="modal">
                      Add To Cart
                    </button>
                  </router-link>
                  <router-link :to="'/ecommerce/checkout'" class="btn btn-primary cart-btn-transform" >Buy Now</router-link>
                </div>
              </div>
            </div>
          </div>
          <button class="btn-close" type="button" data-bs-dismiss="modal" ></button>
        </div>
     </div>
   </div>
        </div>
            <div class="product-details">
              <router-link :to="'/ecommerce/details/' + product.sku">
                <h4>{{ product.name }}</h4>
              </router-link>
              <p>{{ product.shortDescription }}</p>
              <div class="product-price">
                ${{ product.price }}
                <del>${{ product.salePrice }}</del>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
</div>
</template>
<script>

import { mapState, mapGetters } from "vuex";

export default {
  components: {
   
  },
  data() {
    return {
      modalShow: false,
      quickViewProduct: [],
      counter: 1,
      priceArray: [],
      items: [],
      filtered: false,  
    };
  },
  computed: {
    ...mapState({
      products: (state) => state.products.products,
      productslist: (state) => state.products.productslist,
    }),
    ...mapGetters({
      filterProduct: "products/filterProducts",
      tags: "products/setTags",
      col2: 'products/col2',
      col3:"products/col3",
      col4:"products/col4",
      col6:"products/col6",
      listViewEnable:"products/listViewEnable"
    }),
  },
  methods: {
    grid(){
       this.$emit("grid2")
    },
    getImgUrl(path) {
      return require("../../../assets/images/" + path);
    },
    onChangeSort(event) {
      this.$store.dispatch("products/sortProducts", event.target.value);
    },
    collapseFilter() {
      this.filtered = !this.filtered;
    },
    addToCart: function (product, qty) {
      product.quantity = qty ? qty : 1;
      this.$store.dispatch("products/addToCart", product);
    },
    quickView: function (product) {
      this.modalShow = true;
      return (this.quickViewProduct = product);
    },
    quickViewClose: function () {
      this.modalShow = false;
    },
    increment() {
      if (this.counter < this.quickViewProduct.stock) this.counter++;
    },
    decrement() {
      if (this.counter > 1) this.counter--;
    },
  },
};
</script>