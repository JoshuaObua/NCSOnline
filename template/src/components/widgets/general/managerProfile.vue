<template>
     <div class="col-xl-4 col-md-6 box-col-4">
                <div class="card custom-card">
                  <div class="card-header pb-0"><img class="img-fluid" src="../../../assets/images/user-card/3.jpg" alt=""></div>
                  <div class="card-profile"><img class="rounded-circle" src="../../../assets/images/avtar/3.jpg" alt=""></div>
                  <ul class="card-social">
                    <li><a href="https://www.facebook.com"><i class="fa fa-facebook"></i></a></li>
                    <li><a href="https://accounts.google.com"><i class="fa fa-google-plus"></i></a></li>
                    <li><a href="https://twitter.com"><i class="fa fa-twitter"></i></a></li>
                    <li><a href="https://www.instagram.com"><i class="fa fa-instagram"></i></a></li>
                    <li><a href="https://dashboard.rss.com/auth/sign-in"><i class="fa fa-rss"></i></a></li>
                  </ul>
                  <div class="text-center profile-details">
                    <h4><router-link to="/pages/social-app" alt="">Scarlett john</router-link></h4>
                    <h6>Manager</h6>
                  </div>
                  <div class="card-footer row">
                    <div class="col-4 col-sm-4 pb-0">
                      <h6>Follower</h6>
                      <h3 class="counter"><number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="9564"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut"/></h3>
                    </div>
                    <div class="col-4 col-sm-4 pb-0">
                      <h6>Following</h6>
                      <h3><span class="counter"><number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="49"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut" /></span>K</h3>
                    </div>
                    <div class="col-4 col-sm-4 pb-0">
                      <h6>Total Post</h6>
                      <h3><span class="counter"><number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="96"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut" /></span>M</h3>
                    </div>
                  </div>
                </div>
              </div>
</template>
<script>

export default {
     components: {
   
  },
}
</script>