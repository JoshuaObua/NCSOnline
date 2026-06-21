<template>
    <Breadcrumbs sub="Forms" main="Form Controls" title="Input Groups"/>
         <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12">
                <input1/>
                <input2/>
              </div>
            </div>
         </div>
</template>
<script>
import input1 from "./input1.vue"
import input2 from "./input2.vue"
export default {
    components:{
        input1,
        input2
    }
}
</script>