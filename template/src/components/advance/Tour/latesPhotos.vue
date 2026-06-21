<template>
    <div class="col-xl-12 col-sm-6 box-col-12">
                        <div class="card">
                          <div class="card-header">
                            <h5 class="p-0">
                              <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon4" :aria-expanded="isActive?'true':'false'"  aria-controls="collapseicon4" :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Latest Photos</button>
                            </h5>
                          </div>
                          <div  :class="[isActive ? 'block' : ' show']"  v-show="isActive">
                            <div class="card-body photos filter-cards-view">
                              <ul>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-1.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-2.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-3.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-4.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-5.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-6.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-7.png"></div>
                                </li>
                                <li>
                                  <div class="latest-post"><img class="img-fluid" alt="" src="../../../assets/images/social-app/post-8.png"></div>
                                </li>
                              </ul>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="col-xl-12 col-sm-6 box-col-12">
                        <div class="card">
                          <div class="card-header">
                            <h5 class="p-0">
                              <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon13" aria-expanded="true" aria-controls="collapseicon13" :class="[Active ? 'active' : '', 'collapsible']" v-on:click="toggle1">Friends</button>
                            </h5>
                          </div>
                          <div  :class="[isActive ? 'block' : ' show']"  v-show="Active">
                            <div class="card-body avatar-showcase filter-cards-view">
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/3.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/5.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/1.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/2.png" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/3.png" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/6.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/10.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/14.png" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/1.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/4.jpg" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/11.png" alt="#"></div>
                              <div class="d-inline-block friend-pic"><img class="img-50 rounded-circle" src="../../../assets/images/user/8.jpg" alt="#"></div>
                            </div>
                          </div>
                        </div>
                      </div>
</template>

<script>
export default {
  data(){
        return{
            isActive: true,
           Active:true
        }
    },
     methods: {
toggle() {
      this.isActive = !this.isActive
    },
    toggle1() {
      this.Active = !this.Active
    },
}
}
</script>