<template>
    <div class="col-md-12 col-sm-6">
                            <div class="card">
                              <div class="card-header">
                                <h5 class="p-0">
                                  <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon12" :aria-expanded="isActive?'true':'false'"  aria-controls="collapseicon12"  :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Activity Feed</button>
                                </h5>
                              </div>
                              <div  :class="[isActive ? 'block' : ' show']"  v-show="isActive">
                                <div class="card-body social-status filter-cards-view">
                                  <div class="media"><img class="img-50 rounded-circle m-r-15" src="../../../assets/images/user/10.jpg" alt="">
                                    <div class="media-body"><router-link to="/users/profile"><a href="javascript:void(0)"><span class="f-w-600 d-block">Andew Jon
                                    </span></a></router-link>
                                      <p>Commented on Shaun Park's <a href="javascript:void(0)">Photo</a></p><span class="light-span">20 min Ago</span>
                                    </div>
                                  </div>
                                  <div class="media"><img class="img-50 rounded-circle m-r-15" src="../../../assets/images/user/3.jpg" alt="">
                                    <div class="media-body"><router-link to="/users/profile"><a href="javascript:void(0)"><span class="f-w-600 d-block">Johny Waston
                                    </span></a></router-link>
                                      <p>Commented on Shaun Park's <a href="javascript:void(0)">Photo</a></p><span class="light-span">1 hour Ago</span>
                                    </div>
                                  </div>
                                  <div class="media"><img class="img-50 rounded-circle m-r-15" src="../../../assets/images/user/5.jpg" alt="">
                                    <div class="media-body"><router-link to="/users/profile"><a href="javascript:void(0)"><span class="f-w-600 d-block">Comeren Diaz
                                    </span></a></router-link>
                                      <p>Commented on Shaun Park's <a href="javascript:void(0)">Photo</a></p><span class="light-span">1 days Ago</span>
                                    </div>
                                  </div>
                                  <div class="media"><img class="img-50 rounded-circle m-r-15" src="../../../assets/images/user/4.jpg" alt="">
                                    <div class="media-body"><router-link to="/users/profile"><a href="javascript:void(0)"><span class="f-w-600 d-block">Sarah Loren
                                    </span></a></router-link>
                                      <p>Commented on Shaun Park's <a href="javascript:void(0)">Photo</a></p><span class="light-span">2 days Ago</span>
                                    </div>
                                  </div>
                                  <div class="media"><img class="img-50 rounded-circle m-r-15" src="../../../assets/images/user/3.jpg" alt="">
                                    <div class="media-body"><router-link to="/users/profile"><a href="javascript:void(0)"><span class="f-w-600 d-block">Johny Waston
                                    </span></a></router-link>
                                      <p>Commented on Shaun Park's <a href="javascript:void(0)">Photo</a></p><span class="light-span">5 days Ago</span>
                                    </div>
                                  </div>
                                  <div class="media"><img class="img-50 rounded-circle m-r-15" src="../../../assets/images/user/5.jpg" alt="">
                                    <div class="media-body"><router-link to="/users/profile"><a href="javascript:void(0)"><span class="f-w-600 d-block">Comeren Diaz
                                    </span></a></router-link>
                                    <div class="media-body"><router-link to="/users/profile"><a href=""><span class="f-w-600 d-block">Comeren Diaz</span></a></router-link>
                                      <p>Commented on Shaun Park's <a href="javascript:void(0)">Photo</a></p><span class="light-span">6 days Ago</span>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            </div>
                          </div>
    </div>                     
  
</template>
<script>
export default {
    data(){
        return{
            show:true,
            isActive: true,
        }
    },
    methods:{
       toggle() {
      this.isActive = !this.isActive
    },
    }
}
</script>