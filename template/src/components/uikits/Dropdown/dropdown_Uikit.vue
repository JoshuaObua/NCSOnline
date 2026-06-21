<template>
  <Breadcrumbs main="UI Kits" title="Dropdown"/>
    <div class="container-fluid">
      <div class="row">
      <basicDropdown/>
        <colorDropdown/>
        <splitDropdown/>
        <roundedDropdown/>
        <dropdown/>
      </div>
    </div>
</template>
<script>
import colorDropdown from "./colorDropdown.vue"
import basicDropdown from "./basicDropdown.vue"
import splitDropdown from "./splitDropdown.vue"
import roundedDropdown from "./roundedDropdown.vue"
import dropdown from "./dropdown.vue"
export default {
    components:{
        colorDropdown,
        basicDropdown,
        splitDropdown,
        roundedDropdown,
        dropdown
    }
};
</script>