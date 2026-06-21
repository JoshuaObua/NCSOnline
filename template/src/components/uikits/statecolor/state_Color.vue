<template>
    <Breadcrumbs main="Base" title="State Color"/>
    <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                  <defaultColor/>
              </div>
            </div>
    </div>
</template>
<script>
import defaultColor from "./defaultColor.vue"
export default {
    components:{
        defaultColor
    }
}
</script>