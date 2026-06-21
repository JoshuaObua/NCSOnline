<template>
  <div>
    <Breadcrumbs main="Cards" title="Basic Card" />
    <div class="container-fluid">
      <div class="row">
       <card1/>
       <card2/>
      </div>
    </div>
  </div>
</template>
<script>
import card1 from "./card1.vue"
import card2 from "./card2.vue"
export default {
  components: {
    card1,
    card2
  },
};
</script>