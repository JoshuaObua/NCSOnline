<template>
  <div class="container-fluid dashboard-default-sec">
    <div class="row">
      <div class="col-xl-5 box-col-12 des-xl-100">
        <div class="row">
          <div class="col-xl-12 col-md-6 box-col-6 des-xl-50">
            <div class="card profile-greeting">
              <div class="card-header">
                <div class="header-top">
                  <div class="setting-list bg-primary position-unset">
                    <ul class="list-unstyled setting-option">
                      <li>
                        <div class="setting-white">
                          <i class="icon-settings"></i>
                        </div>
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
              <div class="card-body text-center p-t-0">
                <h3 class="font-light">Welcome Back, John!!</h3>
                <p>
                  Welcome to the viho Family! we are glad that you are visite
                  this dashboard. we will be happy to help you grow your
                  business.
                </p>
                <button class="btn btn-light">Update</button>
              </div>
              <div class="confetti">
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
                <div class="confetti-piece"></div>
              </div>
            </div>
          </div>
         
          <annualIncome />
          <annualLosses/>
        </div>
      </div>
       <salesOverview/>
      <growthOverview/>
      <userActivations/>
    </div>
  </div>
</template>
<script>
import annualIncome from "../default/annualIncome.vue";
import annualLosses from "../default/annual_Losses.vue";
import growthOverview from "../default/growthOverview.vue";
import salesOverview from "../default/salesOverview.vue";
import userActivations from '../default/userActivations.vue';
export default {
  components: {
    annualIncome,
    annualLosses,
    growthOverview,
    salesOverview,
    userActivations
  },
};
</script>