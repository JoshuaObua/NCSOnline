<template>
     <Breadcrumbs main="UI Kits" title="lists"/>
          <div class="container-fluid">
            <div class="row">
                <defaultList/>
                <withBadges/>
                <customContent/>
                <javaScriptBehavior/>
                <withIcon/>
            </div>
          </div>
</template>
<script>
import defaultList from "./defaultList.vue"
import withBadges from "./withBadges.vue"
import customContent from "./customContent.vue"
import javaScriptBehavior from "./javaScriptBehavior.vue"
import withIcon from "./withIcon.vue"
export default {
    components:{
        defaultList,
        withBadges,
        customContent,
        javaScriptBehavior,
        withIcon
    }
}
</script>