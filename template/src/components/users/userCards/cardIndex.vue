<template>
    <Breadcrumbs main="Users" title="Users Cards"/>
    <div class="container-fluid user-card">
            <div class="row">
                <card1/>
                <card2/>
            </div>
    </div>
</template>
<script>
import card1 from "./card1.vue"
import card2 from "./card2.vue" 
export default {
    components:{
        card1,
        card2
    }
}
</script>