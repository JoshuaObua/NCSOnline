<template>
  <div class="col-xl-6 box-col-6 top-sell-sec">
    <div class="card">
      <div class="card-header pb-0">
        <div
          class="
            header-top
            d-sm-flex
            justify-content-between
            align-items-center
          "
        >
          <h5>Top Selling Product</h5>
          <div class="center-content">
            <ul class="week-date">
              <li class="font-primary">Today</li>
              <li>Month</li>
            </ul>
          </div>
        </div>
      </div>
      <div class="card-body">
        <div class="media">
          <img
            class="img-fluid"
            src="../../../assets/images/dashboard-2/9.png"
            alt=""
          />
          <div class="media-body">
            <a href="#"><router-link to="/ecommerce/details/1"><h6>Trending Nike shoes</h6></router-link></a>
            <p>New Offer Only $126.00</p>
            <ul class="rating-star">
              <li></li>
              <li><i class="fa fa-star"></i></li>
              <li><i class="fa fa-star"></i></li>
              <li><i class="fa fa-star"></i></li>
              <li><i class="fa fa-star"></i></li>
              <li><i class="fa fa-star"></i></li>
            </ul>
          </div>
          <router-link to="/ecommerce/cartView"><a class="btn btn-iconsolid" href="#"><i class="icon-bag"></i></a></router-link>
        </div>
        <div class="code-box-copy"></div>
      </div>
    </div>
  </div>
  
</template>
<script>

export default {
  components: {
   
  },
};
</script>