<template>
    <div class="col-xl-12 col-sm-6 box-col-12">
                        <div class="card">
                          <div class="card-header">
                            <h5 class="p-0">
                              <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon8" :aria-expanded="isActive?'true':'false'"  aria-controls="collapseicon8" :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Followers</button>
                            </h5>
                          </div>
                          <div  :class="[isActive ? 'block' : ' show']"  v-show="isActive">
                            <div class="card-body social-list filter-cards-view">
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/2.png">
                                <div class="media-body"><span class="d-block">Bucky Barnes</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/3.png">
                                <div class="media-body"><span class="d-block">Sarah Loren</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/3.jpg">
                                <div class="media-body"><span class="d-block">Jason Borne</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/10.jpg">
                                <div class="media-body"><span class="d-block">Comeren Diaz</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/11.png">
                                <div class="media-body"><span class="d-block">Andew Jon</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="col-xl-12 col-sm-6 box-col-12">
                        <div class="card">
                          <div class="card-header">
                            <h5 class="p-0">
                              <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon11" aria-expanded="true" aria-controls="collapseicon11" :class="[Active ? 'active' : '', 'collapsible']" v-on:click="toggle1">Followings</button>
                            </h5>
                          </div>
                          <div  :class="[isActive ? 'block' : ' show']"  v-show="Active">
                            <div class="card-body social-list filter-cards-view">
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/3.png">
                                <div class="media-body"><span class="d-block">Sarah Loren</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/2.png">
                                <div class="media-body"><span class="d-block">Bucky Barnes</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/10.jpg">
                                <div class="media-body"><span class="d-block">Comeren Diaz</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/3.jpg">
                                <div class="media-body"><span class="d-block">Jason Borne</span><a href="javascript:void(0)">Add Friend</a></div>
                              </div>
                              <div class="media"><img class="img-50 img-fluid m-r-20 rounded-circle" alt="" src="../../../assets/images/user/11.png">
                                <div class="media-body"><span class="d-block">Andew Jon</span><a href="javascript:void(0)">Add Friend</a></div>
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