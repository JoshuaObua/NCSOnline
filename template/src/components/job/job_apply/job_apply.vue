<template>
    <Breadcrumbs main="Job Search" title="Apply"/>
        <div class="container-fluid">
				<div class="row">
                    <jobtab/>
                    <div class="col-xl-9 xl-60 box-col-7">
						<div class="card">
							<div class="job-search">
								<div class="card-body pb-0">
									<div class="media"><img class="img-40 img-fluid m-r-20" :src="getImgUrl(jobs.image)" alt="">
										<div class="media-body">
											<h5 class="f-w-700 mb-0">
												<router-link :to="{ name: 'jobdetails', params: { id: jobs.id }}">{{ jobs.title }}</router-link>
											</h5>
											<p>{{ jobs.company }} <span>{{ jobs.city }}, {{ jobs.country }} </span><span class="rating" v-html="stars(jobs.stars)"></span></p>
                      <button class="btn btn-primary job-apply-btn" type="button">
                        <span><i class="fa fa-check text-white"></i></span>
                        Save this job</button>
										</div>
									</div>
                                    <div class="job-description">
                                        <personal/>
                                        <education/>
                                        <experience/>
                                        <upload/>
                                    </div>
                                </div>
                                	<div class="card-footer">
									<button class="btn btn-primary m-r-10" type="submit">Submit</button>
									<router-link class="btn btn-light" :to="{ name: 'jobdetails', params: { id: jobs.id }}">Cancel</router-link>
								</div>
                            </div>
                        </div>
                    </div>
                </div>
        </div>
</template>
<script>
 import { mapState } from 'vuex';
  import router from '@/router';
import jobtab from "../job-list/jobtab.vue"
import Personal from "./personal/personal.vue"
import education from "./education.vue"
import experience from "./experience.vue"
import upload from "./upload.vue"
export default {
    components:{
      jobtab,
      Personal,
      education,
      experience,
      upload
    },
 computed:{
      ...mapState({
        jobs: state => state.jobs.jobs.find(job => {
          if(parseInt(router.currentRoute._rawValue.params.id) === job.id)
            return job;
        })
      })
    }, 
    methods:{
      getImgUrl(filename) {
        var images = require.context('../../../assets/images/job-search/', false, /\.jpg$/);
        return images('./' + filename);
      },
      stars(count) {
        var stars = '';

        for (var i = 0; i < 5; i++) {
          if(count > i) {
            stars = stars + '<i class="fa fa-star font-warning"></i>';
          } else {
            stars = stars + '<i class="fa fa-star font-warning-o"></i>';
          }
        }

        return stars;
      }
    }
}
</script>