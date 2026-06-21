<template>
    <Breadcrumbs main="Learning" title="Learning List"/>
         <div class="container-fluid">
            <div class="row learning-block">
                <div class="col-xl-9 xl-60 ">
					<div class="row">
						<div class="col-xl-12 col-sm-6" v-for="(course,index) in specialcourses" :key="'special'+index">
							<div class="card">
								<div class="blog-box blog-list row">
									<div class="col-xl-5 col-12"><img class="img-fluid sm-100-w" :src="getImgUrl(course.image)" alt=""></div>
									<div class="col-xl-7 col-12">
										<div class="blog-details">
											<div class="blog-date digits"><span v-text="course.dd"></span> {{ course.mm }} {{ course.yy }}</div>
											<router-link :to="{ name: 'coursedetailed', params:{ id: course.id } }" >
											    <div class="bottom-details"> 
                                                    <h6 v-text="course.title"></h6>
                                                </div>
										    </router-link>
											<div class="blog-bottom-content">
												<ul class="blog-social">
													<li v-text="' by: '+course.author"></li>
													<li class="digits" v-text="course.hits+' Hits'"></li>
												</ul>
												<hr>
												<p class="mt-0" v-text="course.descr"></p>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>

						<div class="col-xl-4 xl-50 col-sm-6 box-col-6" v-for="(course,index) in normalcourses" :key="'normal'+index">
							<div class="card">
								<div class="product-box learning-box">
									<div class="product-img"><img class="img-fluid top-radius-blog" :src="getImgUrl(course.image)" alt="">
										<div class="product-hover">
											<ul>
												<li>
												<router-link :to="{ name: 'coursedetailed', params:{ id: course.id } }" ><i class="icon-link"></i> 
												</router-link>											
												</li>
												<li>
												<router-link :to="{ name: 'coursedetailed', params:{ id: course.id } }" ><i class="icon-import"></i> 
												</router-link>											
												</li>
											</ul>
										</div>
									</div>
									<div class="details-main">
										<router-link :to="{ name: 'coursedetailed', params:{ id: course.id } }" >
											<div class="bottom-details"> 
                                                <h6 v-text="course.title"></h6>
                                            </div>
										</router-link>                                     
                          <p v-text="course.descr"></p>
                        </div>
								</div>
							</div>
						</div>

					</div>
				</div>
				<div class="col-xl-3 xl-40">
					<div class="job-sidebar"><a class="btn btn-primary job-toggle" href="javascript:void(0)" @click="collapseFilter()">learning filter</a>
                  <div class="job-left-aside custom-scrollbar " :class="filtered?'open':''">
                    <div class="default-according style-1 faq-accordion job-accordion" id="accordionoc">
                      <div class="row">
						  <learningtab/>
					  </div>
					</div>
				  </div>
				  </div>
				</div>
            </div>
         </div>
</template>
<script>
import learningtab from "./learningtab.vue"
 import { mapGetters } from 'vuex';
export default {
	components:{
		learningtab
	},
	data(){
      return {
        
        filtered: false,
      };
    },
   
     computed: {
      ...mapGetters({
        specialcourses:'courses/specialcourses',
        normalcourses:'courses/normalcourses'
      })
    },
    methods:{
      getImgUrl(filename) {
        var images = require.context('../../../assets/images/faq/', false, /\.jpg$/);
        return images('./' + filename);
      },
	  collapseFilter() {
        this.filtered = !this.filtered;
      },
    }
}
</script>