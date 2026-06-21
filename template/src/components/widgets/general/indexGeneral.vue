<template>
 <Breadcrumbs main="Widgets" title="General" />
 <div class="container-fluid general-widget">
      <div class="row">
        <earningsGeneral />
        <productsGeneral/>
        <messagesGeneral/>
        <newUse/>
        <calenderGeneral/>
        <clockGeneral/>
        <weatherGeneral/>
        <sale/>
        <pendingGeneral/>
        <poioklotGeneral/>
        <recentActivity/>
        <managerProfile/>
        <facebookGeneral/>
        <twitterGeneral/>
        <linkdinGeneral/>
        <googleGeneral/>
        <chromeGeneral/>
        <firefoxGeneral/>
        <safariGeneral/>
        <productsCart/>
        <employeeStatus/>
        <explainGeneral/>
        <contactUs/>
      </div>
 </div>
</template>
<script>
import earningsGeneral from "./earningsGeneral.vue";
import productsGeneral from "./productsGeneral.vue";
import messagesGeneral from './messagesGeneral.vue';
import newUse from './newUse.vue';
import calenderGeneral from './calenderGeneral.vue';
import weatherGeneral from './weatherGeneral.vue';
import clockGeneral from "./clockGeneral.vue";
import sale from './sale/saleGeneral.vue';
import pendingGeneral from "./pendingGeneral.vue";
import poioklotGeneral from "./poioklotGeneral.vue";
import recentActivity from './recentActivity.vue';
import managerProfile from './managerProfile.vue';
import facebookGeneral from './facebookGeneral.vue';
import twitterGeneral from './twitterGeneral.vue';
import linkdinGeneral from './linkdinGeneral.vue';
import googleGeneral from './googleGeneral.vue';
import chromeGeneral from './chromeGeneral.vue';
import firefoxGeneral from "./firefoxGeneral.vue";
import safariGeneral from './safariGeneral.vue';
import productsCart from "./productsCart.vue";
import employeeStatus from './employeeStatus.vue';
import explainGeneral from './explainGeneral.vue';
import contactUs from './contactUs.vue';
export default {
  components: {
    earningsGeneral,
    productsGeneral,
    messagesGeneral,
    newUse,
    calenderGeneral,
   weatherGeneral,
   clockGeneral,
   sale,
   pendingGeneral,
   poioklotGeneral,
   recentActivity,
   managerProfile,
   facebookGeneral,
   twitterGeneral,
   linkdinGeneral,
   googleGeneral,
   chromeGeneral,
   firefoxGeneral,
   safariGeneral,
   productsCart,
   employeeStatus,
   explainGeneral,
   contactUs
  },
};
</script>