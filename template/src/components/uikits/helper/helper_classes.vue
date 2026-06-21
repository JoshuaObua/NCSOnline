<template>
  <Breadcrumbs main="UI Kits" title="Helper Classes"/>
    <div class="container-fluid">
            <div class="row">
              <paddingMargin/>
              <paddingLeftRight/>
              <paddingTopButtom/>
              <marginLeftRight/>
              <marginTopBottom/>
              <verticalAlign/>
              <textColor/>
              <backgroundColor/>
              <borderRadius/>
              <fontSize/>
              <bordersColor/>
              <borderWidth/>
            </div>
    </div>
</template>
<script>
import paddingMargin from "./paddingMargin.vue"
import paddingLeftRight from "./paddingLeftRight.vue"
import paddingTopButtom from "./paddingTopButtom.vue"
import marginLeftRight from "./marginLeftRight.vue"
import marginTopBottom from "./marginTopBottom.vue"
import verticalAlign from "./verticalAlign.vue"
import textColor from "./textColor.vue"
import backgroundColor from "./backgroundColor.vue"
import borderRadius from "./borderRadius.vue"
import fontSize from "./fontSize.vue"
import bordersColor from "./bordersColor.vue"
import borderWidth from "./borderWidth.vue"
export default {
  components:{
    paddingMargin,
    paddingLeftRight,
    paddingTopButtom,
    marginLeftRight,
    marginTopBottom,
    verticalAlign,
    textColor,
    backgroundColor,
    borderRadius,
    fontSize,
    bordersColor,
    borderWidth
  }
}
</script>