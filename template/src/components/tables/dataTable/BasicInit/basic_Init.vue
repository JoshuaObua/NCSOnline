<template>
    <Breadcrumbs main="Tables" title="Basic Init"/>
        <div class="container-fluid">
            <div class="row">
                <zeroConfiguration/>   
            </div>
        </div>
</template>
<script>
import zeroConfiguration from "./zeroConfiguration.vue"

export default {
    components:{
     zeroConfiguration,   
    },
    data(){
        return{
        }
    },
}
</script>