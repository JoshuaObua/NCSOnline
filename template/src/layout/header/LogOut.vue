<template>
     <button class="btn btn-primary-light" type="button"><a @click="logout" class="d-flex"><vue-feather type="log-out"></vue-feather>Log out</a></button>
</template>
<script>

export default {
    components:{
       
    },
     methods: {
      logout: function() {
        this.$router.replace('/auth/login');
            localStorage.clear('user')
      },
    },
}
</script>