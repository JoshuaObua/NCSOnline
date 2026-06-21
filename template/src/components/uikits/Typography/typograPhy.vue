<template>
       <Breadcrumbs main="UI Kits" title="typography"/>
            <div class="container-fluid">
            <div class="row">
                <headings/>
                <headings2/>
                <headings3/>
            </div>
            </div>
</template>
<script>
import headings from "./headings.vue"
import headings2 from "./headings2.vue"
import headings3 from "./headings3.vue"
export default {
    components:{
        headings,
       headings2,
       headings3
    }
}
</script>