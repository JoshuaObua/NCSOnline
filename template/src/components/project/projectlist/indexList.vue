<template>
    <Breadcrumbs main="Project" title="project list"/>
    <div class="container-fluid">
            <div class="row project-cards">
                <projectCards/>
            </div>
    </div>
</template>
<script>
import projectCards from "./projectCards.vue";
import projectList from "./projectList.vue";
import projectDoing from "./projectDoing.vue"
export default {
  components: {
    projectCards,
    projectDoing,
    projectList
  },
   data() {
    return {
      showtab:'all',
    };
  },
  
  methods: {
   changetab(tab) {
        this.showtab = tab;
      }
  },
};
</script>