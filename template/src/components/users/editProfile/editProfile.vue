<template>
    <Breadcrumbs main="Users" title="Users Edit"/>
    <div class="container-fluid">
            <div class="edit-profile ">
              <div class="row">
                  <myProfile/>
                  <editProfile1/>
              </div>
            </div>
    </div>
</template>
<script>
import myProfile from "./myProfile.vue"
import editProfile1 from "./editProfile1.vue"
export default {
    components:{
        myProfile,
        editProfile1
    }
}
</script>