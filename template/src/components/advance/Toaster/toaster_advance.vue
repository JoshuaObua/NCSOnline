<template>
    <Breadcrumbs main="Bonus UI" title="toaster"/>
         <div class="container-fluid">
            <div class="row">
            <toastr/>
            <toastrIcon/>
            <toastrPosition/>
            </div>
         </div>
</template>
<script>
import toastr from "./toastr.vue"
import toastrIcon from "./toastrIcon.vue"
import toastrPosition from "./toastrPosition.vue"
export default {
   components:{
       toastr,
       toastrIcon,
       toastrPosition
   }
}
</script>