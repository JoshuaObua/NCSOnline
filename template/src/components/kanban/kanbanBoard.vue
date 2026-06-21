<template>
     <Breadcrumbs title="Kanban Board"/>
        <div class="container-fluid">
            <div class="row">
                <defultDemo/>
                <customBoard/>
                <apiDemo/>
            </div>
        </div>
</template>
<script>
import defultDemo from './defultDemo.vue'
import customBoard from "./customBoard.vue"
import apiDemo from "./apiDemo.vue"
export default {
    components:{
        defultDemo,
        customBoard,
        apiDemo
    }
}
</script>