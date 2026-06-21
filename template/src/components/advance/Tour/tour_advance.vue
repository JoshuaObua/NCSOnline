<template>
  <Breadcrumbs main="Bonus UI" title="tour"/>
    <div class="container-fluid ">
      <div class="user-profile">
        <div class="row">
          <tourProfile />
          <tourDetail />
          <div class="col-xl-9 col-lg-12 xl-65">
            <div class="row">
              <tourPost />
              <profilePost1 />
            </div>
          </div>
          <v-tour name="myTour" :steps="steps" :options="tourOptions"></v-tour>
        </div>
      </div>
    </div>
</template>
<script>
import tourProfile from "./tourProfile.vue";
import tourDetail from "./tourDetail.vue";
import tourPost from "./tourPost.vue";
import profilePost1 from "../../users/profile/profilePost1.vue";
export default {
  components: {
    profilePost1,
    tourProfile,
    tourDetail,
    tourPost,
  },
  data() {
    return {
      tourOptions: {
        useKeyboardNavigation: true,
        labels: {
          buttonSkip: "Skip",
          buttonPrevious: "Back",
          buttonNext: "Next",
          buttonStop: "Done",
        },
      },
      toogle: true,
      steps: [
        {
          target: "#profile-tour",
          content: "This is Profile image",
        },
        {
          target: "#update-profile-tour",
          content: "Change Profile image here",
        },
        {
          target: "#info-bar-tour",
          content: "This is the your details",
        },
        {
          target: "#social-bar-tour",
          content: "This is your Social details",
        },
        {
          target: "#first-post-tour",
          content: "This is the your first Post",
          params: {
            placement: "top",
          },
        },
        {
          target: "#last-post-tour",
          content: "This is the your last Post",
          params: {
            placement: "top",
          },
        },
      ],
    };
  },
  mounted: function () {
    this.$tours["myTour"].start();
  },
};
</script>