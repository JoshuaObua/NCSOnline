<template>
     <div class="col-xl-12 ">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon2" :aria-expanded="isActive?'true':'false'" aria-controls="collapseicon2" :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Job Title</button>
                              </h5>
                            </div>
                            <div :class="[isActive ? 'block' : 'none', 'content']" v-show="isActive">
                              <div class="card-body animate-chk">
                                <label class="d-block" for="chk-ani9">
                                  <input class="checkbox_animated" id="chk-ani9" type="checkbox">UI/Ux designer(25)
                                </label>
                                <label class="d-block" for="chk-ani10">
                                  <input class="checkbox_animated" id="chk-ani10" type="checkbox">Graphic designer(10)
                                </label>
                                <label class="d-block" for="chk-ani11">
                                  <input class="checkbox_animated" id="chk-ani11" type="checkbox">Front end designer(15)
                                </label>
                                <label class="d-block" for="chk-ani12">
                                  <input class="checkbox_animated" id="chk-ani12" type="checkbox">PHP developer(42)
                                </label>
                                <label class="d-block mb-0" for="chk-ani13">
                                  <input class="checkbox_animated" id="chk-ani13" type="checkbox">React Developer(5)
                                </label>
                              </div>
                              <button class="btn btn-block btn-primary text-center" type="button">All Job Title</button>
                            </div>
                          </div>
                        </div>
                        <div class="col-xl-12 ">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon3" :aria-expanded="isActive1?'true':'false'" aria-controls="collapseicon3" :class="[isActive1 ? 'active' : '', 'collapsible']" v-on:click="open1()">Industry</button>
                              </h5>
                            </div>
                            <div :class="[isActive1 ? 'block' : 'none', 'content']" v-show="isActive1">
                              <div class="card-body animate-chk">
                                <label class="d-block" for="chk-ani14">
                                  <input class="checkbox_animated" id="chk-ani14" type="checkbox">Computer Software(14)
                                </label>
                                <label class="d-block" for="chk-ani15">
                                  <input class="checkbox_animated" id="chk-ani15" type="checkbox">IT Engineer(10)
                                </label>
                                <label class="d-block" for="chk-ani16">
                                  <input class="checkbox_animated" id="chk-ani16" type="checkbox">Service industry(20)
                                </label>
                                <label class="d-block" for="chk-ani17">
                                  <input class="checkbox_animated" id="chk-ani17" type="checkbox">Accounting (34)
                                </label>
                                <label class="d-block mb-0" for="chk-ani18">
                                  <input class="checkbox_animated" id="chk-ani18" type="checkbox">Financial Services(5)
                                </label>
                              </div>
                              <button class="btn btn-block btn-primary text-center" type="button">All Industries</button>
                            </div>
                          </div>
                        </div>
                        <div class="col-xl-12 ">
                          <div class="card">
                            <div class="card-header">
                              <h5 class="mb-0 p-0">
                                <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon4" :aria-expanded="isActive2?'true':'false'" aria-controls="collapseicon4" :class="[isActive2 ? 'active' : '', 'collapsible']" v-on:click="open2()">Specific skills</button>
                              </h5>
                            </div>
                            <div :class="[isActive2 ? 'block' : 'none', 'content']" v-show="isActive2">
                              <div class="card-body animate-chk">
                                <label class="d-block" for="chk-ani19">
                                  <input class="checkbox_animated" id="chk-ani19" type="checkbox">HTML,scss & sass
                                </label>
                                <label class="d-block" for="chk-ani20">
                                  <input class="checkbox_animated" id="chk-ani20" type="checkbox">Javascript
                                </label>
                                <label class="d-block" for="chk-ani21">
                                  <input class="checkbox_animated" id="chk-ani21" type="checkbox">Node.js
                                </label>
                                <label class="d-block" for="chk-ani22">
                                  <input class="checkbox_animated" id="chk-ani22" type="checkbox">Gulp & Pug
                                </label>
                                <label class="d-block mb-0" for="chk-ani23">
                                  <input class="checkbox_animated" id="chk-ani23" type="checkbox">Angular.js
                                </label>
                              </div>
                              <button class="btn btn-block btn-primary text-center" type="button">All Skills</button>
                            </div>
                          </div>
                        </div>
</template>
<script>
export default {
  data(){
    return{
        isActive: true,
         isActive1: true,
          isActive2: true,
    }
  },
  methods: {
    toggle(){
       this.isActive = !this.isActive
    },
      open1(){
       this.isActive1 = !this.isActive1
    },
      open2(){
       this.isActive2 = !this.isActive2
    }
}
}
</script>