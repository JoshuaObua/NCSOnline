<template>
  <div class="col-sm-6 col-xl-3 col-lg-6">
    <div class="card o-hidden border-0">
      <div class="bg-primary b-r-4 card-body">
        <div class="media static-top-widget">
          <div class="align-self-center text-center">
            <i data-feather="shopping-bag"
              ><vue-feather type="message-circle"></vue-feather
            ></i>
          </div>
          <div class="media-body">
            <span class="m-0">Messages</span>
            <h4 class="mb-0 counter">
              <number
                class="bold transition"
                ref="number1"
                :from="0"
                :to="893"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              />
            </h4>
            <vue-feather class="icon-bg" type="message-circle"></vue-feather>
          </div>
        </div>
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