<template>
    <div>
        <Breadcrumbs title="Email" />
        <div class="container-fluid ">
            <div class="email-wrap">
                <div class="row">
                    <navMenu />
                    <emailInbox />
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import navMenu from "./navMenu.vue"
import emailInbox from "./emailInbox.vue"

export default {
    components: {
        navMenu,
        emailInbox,

    },

}
</script>