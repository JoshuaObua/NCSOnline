<template>
 
  <div class="col-xl-4 des-xl-50 box-col-12 activity-sec chart_data_left">
                <div class="card">
                  <div class="card-header">
                    <div class="header-top d-sm-flex justify-content-between align-items-center">
                      <h5 class="m-0">Activity Timeline</h5>
                      <div class="center-content">
                        <p>Yearly User 24.65k</p>
                      </div>
                     
                    </div>
                  </div>
                  <div class="card-body">
                    <div class="chart-main activity-timeline update-line">
                      <div class="media">
                        <div class="activity-line"></div>
                        <div class="activity-dot-primary"></div>    
                        <div class="media-body d-block">
                          <h6> <span class="font-primary">20-04-2021</span>Today </h6>
                          <h5>Updated Product<i class="fa fa-circle circle-dot-primary pull-right"></i></h5>
                          <p>Quisque a consequat ante Sit amet magna at volutapt.</p>
                        </div>
                      </div>
                      <div class="media">
                        <div class="activity-dot-primary"></div>
                        <div class="media-body d-block">
                          <h6> <span class="font-primary">20-04-20121</span>Today<span class="badge pill-badge-primary m-l-10">new</span></h6>
                          <h5>James just like your product     <i class="fa fa-circle circle-dot-primary pull-right"></i></h5>
                          <p>Quisque a consequat ante Sit amet magna at volutapt.</p>
                          <ul class="timeline-pro">
                            <li><img class="img-fluid" src="@/assets/images/dashboard-2/11.png" alt="Product-1"></li>
                            <li><img class="img-fluid" src="@/assets/images/dashboard-2/10.png" alt="Product-2"></li>
                          </ul>
                        </div>
                      </div>
                      <div class="media">
                        <div class="activity-dot-primary"></div>
                        <div class="media-body d-block">
                          <h6> <span class="font-primary">20-04-20121</span>Today</h6>
                          <h5>Jihan Doe just like your product<i class="fa fa-circle circle-dot-primary pull-right"></i></h5>
                          <p>Vestibulum nec mi suscipit, dapibus purus ane.</p>
                        </div>
                      </div>
                      <div class="media">
                        <div class="media-body d-block">
                          <div class="tomorrow-sec">
                            <p>Tomorrow</p>
                          </div>
                        </div>
                      </div>
                      <div class="media">
                        <div class="activity-dot-primary"></div>
                        <div class="media-body d-block">
                          <h6> <span class="font-primary">20-04-20121</span>Tomorrow</h6>
                          <h5>Today Total  Revenue<i class="fa fa-circle circle-dot-primary pull-right"></i></h5>
                          <p>Quisque a consequat ante Sit amet magna at volutapt.</p>
                        </div>
                      </div>
                      <div class="media">
                        <div class="activity-dot-primary"></div>
                        <div class="media-body d-block">                                         
                          <div class="hospital-small-chart">
                            <div id="column-chart">      <activityTimeline2/></div>
                          </div>
                        </div>
                      </div>
                    </div>
                   
                  </div>
                </div>
              </div>
</template>
<script>
import activityTimeline2 from "../e-commerce/activityTimeline2.vue";
export default {
  components: {
    activityTimeline2,
  },
};
</script>
