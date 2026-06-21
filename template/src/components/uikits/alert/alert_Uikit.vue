<template>
    <Breadcrumbs main="UI Kits" title="alert"/>
       <div class="container-fluid">
            <div class="row">
                <themeColor/>
                <linkColor/>
                <dismissing/>
                <alertIcon/>
                <lightAlerts/>
                <linkColor2/>
                <dismissing1/>
                <outlineAlerts/>
                <outlineDark/>
                <iconOutline/>
                <iconInverse/>
                <textAction/>
                <additionalContent/>
            </div>
       </div>
</template>
<script>
import themeColor from "./themeColor.vue"
import linkColor from "./linkColor.vue"
import dismissing from "./dismissing.vue"
import alertIcon from "./alertIcon.vue"
import lightAlerts from "./lightAlerts.vue"
import linkColor2 from "./linkColor2.vue"
import dismissing1 from "./dismissing1.vue"
import outlineAlerts from "./outlineAlerts.vue"
import outlineDark from "./outlineDark.vue"
import iconOutline from "./iconOutline.vue"
import iconInverse from "./iconInverse.vue"
import textAction from "./textAction.vue"
import additionalContent from "./additionalContent.vue"
export default {
    components:{
        themeColor,
        linkColor,
        dismissing,
        alertIcon,
        lightAlerts,
        linkColor2,
        dismissing1,
        outlineAlerts,
        outlineDark,
        iconOutline,
        iconInverse,
        textAction,
        additionalContent
    }
}
</script>