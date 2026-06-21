<template>
      <div class="col-sm-12">
                  <div class="card profile-header">
                    <div class="profile-img-wrrap"></div>
                    <div class="userpro-box">
                      <div class="img-wrraper">                              
                        <div class="avatar" id="profile-tour"><img class="img-fluid" alt="" src="../../../assets/images/user/7.jpg"></div><a class="icon-wrapper" id="update-profile-tour" href=""><router-link to="/users/edit"><i class="icofont icofont-pencil-alt-5"></i></router-link></a>
                      </div>
                      <div class="user-designation">
                        <div class="title"><a target="_blank" href=""> 
                            <h4>Emay Walter</h4>
                            <h6>designer</h6></a></div>
                        <div class="social-media" data-intro="This is your Social details" id="social-bar-tour">
                          <ul class="user-list-social">
                            <li><a href="https://www.facebook.com"><i class="fa fa-facebook"></i></a></li>
                            <li><a href="https://accounts.google.com"><i class="fa fa-google-plus"></i></a></li>
                            <li><a href="https://twitter.com"><i class="fa fa-twitter"></i></a></li>
                            <li><a href="https://www.instagram.com"><i class="fa fa-instagram"></i></a></li>
                            <li><a href="https://dashboard.rss.com/auth/sign-in"><i class="fa fa-rss"></i></a></li>
                          </ul>
                        </div>
                        <div class="follow">
                          <ul class="follow-list">
                            <li>
                              <div class="follow-num"><number
                              class="counter"
                              ref="number"
                              :from="0"
                              :to="325"
                              :duration="5"
                              :delay="0"
                              easing="Power1.easeOut" /></div><span>Follower</span>
                            </li>
                            <li>
                              <div class="follow-num counter"><number
                             class="bold counter"
                             ref="number1"
                             :from="0"
                             :to="450"
                             :duration="5"
                             :delay="0"
                             easing="Power1.easeOut" /></div><span>Following</span>
                            </li>
                            <li>
                              <div class="follow-num counter"><number
                          class="bold counter"
                          ref="number1"
                          :from="0"
                          :to="500"
                          :duration="5"
                          :delay="0"
                          easing="Power1.easeOut" /></div><span>Likes</span>
                            </li>
                          </ul>
                        </div>
                      </div>
                    </div>
                  </div>  
                </div>
               
</template>
<script>
export default {
}
</script>