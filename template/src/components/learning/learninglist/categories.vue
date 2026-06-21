<template>
    <div class="col-xl-12 ">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link" data-bs-toggle="collapse" data-bs-target="#collapseicon1" :aria-expanded="isActive?'true':'false'" aria-controls="collapseicon1" :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Categories</button>
                              </h5>
                            </div>
                            <div  :class="[isActive ? 'block' : ' show']"  v-show="isActive">
                              <div class="categories">
                                <div class="learning-header"><span class="f-w-600">Design</span></div>
                                <ul>
                                  <li><a href="#">UI Design </a><span class="badge badge-primary pull-right">28</span></li>
                                  <li><a href="#">UX Design </a><span class="badge badge-primary pull-right">35</span></li>
                                  <li><a href="#">Interface Design </a><span class="badge badge-primary pull-right">17</span></li>
                                  <li><a href="#">User Experience </a><span class="badge badge-primary pull-right">26</span></li>
                                </ul>
                              </div>
                              <div class="categories pt-0">
                                <div class="learning-header"><span class="f-w-600">Development</span></div>
                                <ul>
                                  <li><a href="#">Frontend Development</a><span class="badge badge-primary pull-right">48</span></li>
                                  <li><a href="#">Backend Development</a><span class="badge badge-primary pull-right">19</span></li>
                                </ul>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="col-xl-12">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link" data-bs-toggle="collapse" data-bs-target="#collapseicon2" :aria-expanded="isActive?'true':'false'" aria-controls="collapseicon2" :class="[Active ? 'active' : '', 'collapsible']" v-on:click="toggle1">Upcoming Courses</button>
                              </h5>
                            </div>
                            <div  :class="[isActive ? 'block' : ' show']"  v-show="Active">
                              <div class="upcoming-course card-body">
                                <div class="media">
                                  <div class="media-body"><span class="f-w-600">Upcoming course</span>
                                    <h6>Course By <a href="#">Development Team</a></h6>
                                    <ul class="rating">
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star-half-o font-warning"></i></li>
                                    </ul>
                                  </div>
                                  <div>
                                    <h5 class="mb-0 p-0 font-primary">18</h5><span class="d-block">Dec</span>
                                  </div>
                                </div>
                                <div class="media">
                                  <div class="media-body"><span class="f-w-600">Business Analyst</span>
                                    <h6>Course By <a href="#">Analyst Team</a></h6>
                                    <ul class="rating">
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star-half-o font-warning"></i></li>
                                    </ul>
                                  </div>
                                  <div>
                                    <h5 class="mb-0 p-0 font-primary">28</h5><span class="d-block">Dec</span>
                                  </div>
                                </div>
                                <div class="media">
                                  <div class="media-body"><span class="f-w-600">Web Development</span>
                                    <h6>Course By <a href="#">Designer</a></h6>
                                    <ul class="rating">
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star font-warning"></i></li>
                                      <li><i class="fa fa-star-half-o font-warning"></i></li>
                                    </ul>
                                  </div>
                                  <div>
                                    <h5 class="mb-0 p-0 font-primary">5</h5><span class="d-block">Jan</span>
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