<template>
	<div>

		<div class="page-wrapper compact-wrapper " id="pageWrapper">        
      <div class="error-wrapper">
        <div class="container">
          <div class="error-page1">
            <div class="svg-wrraper mb-0">
             <errorPage2/>
            </div>
            <div class="col-md-8 offset-md-2">
              <h3>page not found</h3>
              <p class="sub-content">The page you are attempting to reach is currently not available. This may be because the page does not exist or has been moved.</p><router-link class="btn btn-primary btn-lg" :to="{ path: '/' }">BACK TO HOME PAGE</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
	</div>
</template>
<script>
import errorPage2 from "./svg/errorPage2.vue"
export default {
  components:{
    errorPage2
  }
}
</script>
