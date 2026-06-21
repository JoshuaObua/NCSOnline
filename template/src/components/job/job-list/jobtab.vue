<template>
      <div class="col-xl-3 xl-40 box-col-4"> 
        <div class=" job-sidebar">
        <a class="btn btn-primary job-toggle" href="javascript:void(0)" @click="collapseFilter()">job filter</a>
          <div class="job-left-aside custom-scrollbar" :class="filtered?'open':''">
                    <div class="default-according style-1 faq-accordion job-accordion" id="accordionoc">
                      <div class="row">
                           <div class="col-xl-12 ">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon" :aria-expanded="isActive?'true':'false'" aria-controls="collapseicon" :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Filter</button>
                              </h5>
                            </div>
                            <div :class="[isActive ? 'block' : 'none', 'content']" v-show="isActive">
                              <div class="card-body filter-cards-view animate-chk">
                                <div class="job-filter mb-3">
                                  <div class="faq-form">
                                    <input class="form-control" type="text" placeholder="Search.."><i data-feather="search"><vue-feather class="search-icon" type="search"></vue-feather></i>
                                  </div>
                                </div>
                                <div class="job-filter">
                                  <div class="faq-form">
                                    <input class="form-control" type="text" placeholder="location.."><i data-feather="map-pin"><vue-feather class="search-icon" type="map-pin"></vue-feather></i>
                                  </div>
                                </div>
                                <div class="checkbox-animated">
                                  <label class="d-block" for="chk-ani">
                                    <input class="checkbox_animated" id="chk-ani" type="checkbox"> Full-time (866)
                                  </label>
                                  <label class="d-block" for="chk-ani1">
                                    <input class="checkbox_animated" id="chk-ani1" type="checkbox"> Contract (503)
                                  </label>
                                  <label class="d-block" for="chk-ani2">
                                    <input class="checkbox_animated" id="chk-ani2" type="checkbox"> Part-time (288)
                                  </label>
                                  <label class="d-block" for="chk-ani3">
                                    <input class="checkbox_animated" id="chk-ani3" type="checkbox"> Internship (146)
                                  </label>
                                  <label class="d-block" for="chk-ani4">
                                    <input class="checkbox_animated" id="chk-ani4" type="checkbox">  Temporary (25)
                                  </label>
                                  <label class="d-block" for="chk-ani5">
                                    <input class="checkbox_animated" id="chk-ani5" type="checkbox"> Commission (868)
                                  </label>
                                </div>
                                <button class="btn btn-primary text-center" type="button">Find jobs</button>
                              </div>
                            </div>
                          </div>
                        </div>
                         <div class="col-xl-12 ">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon1" :aria-expanded="isActive2?'true':'false'" aria-controls="collapseicon1" :class="[isActive2 ? 'active' : '', 'collapsible']" v-on:click="open2()" >Location</button>
                              </h5>
                            </div>
                            <div :class="[isActive2 ? 'block' : 'none', 'content']" v-show="isActive2">
                              <div class="card-body animate-chk">
                                <div class="location-checkbox">
                                  <label class="d-block" for="chk-ani6">
                                    <input class="checkbox_animated" id="chk-ani6" type="checkbox">London<span class="d-block">London (69)</span>
                                  </label>
                                  <label class="d-block" for="chk-ani7">
                                    <input class="checkbox_animated" id="chk-ani7" type="checkbox">New York<span class="d-block">NY (39)</span>
                                  </label>
                                  <label class="d-block" for="chk-ani7">
                                    <input class="checkbox_animated" id="chk-ani7" type="checkbox">Contrary<span class="d-block">Con (45)</span>
                                  </label>
                                  <label class="d-block" for="chk-ani7">
                                    <input class="checkbox_animated" id="chk-ani7" type="checkbox">Chicago<span class="d-block">IL (63)</span>
                                  </label>
                                  <label class="d-block mb-0" for="chk-ani8">
                                    <input class="checkbox_animated" id="chk-ani8" type="checkbox">Seattle<span class="d-block">WA (249)</span>
                                  </label>
                                </div>
                              </div>
                              <button class="btn btn-block btn-primary text-center" type="button">All Locations</button>
                            </div>
                          </div>
                        </div>
                   <jobtab2/>
                      </div>
                    </div>
          </div>
        </div>
      </div>    
</template>
<script>
import jobtab2 from "./jobtab2.vue"
export default {
    components:{
        jobtab2
    },
    data(){
      return {
        isActive: true,
        filtered: false,
        isActive2: true,
      };
    },
    methods:{
      collapseFilter() {
        this.filtered = !this.filtered;
      },
      toggle() {
      this.isActive = !this.isActive
    },
    open2(){
       this.isActive2 = !this.isActive2
    }
    }
}
</script>