<template>
  
                <readEmail/>
</template>
<script>
import readEmail from "./readEmail.vue"
import { mapState,mapGetters } from 'vuex';
export default {
    components:{
          readEmail
    },
    data(){
        return{
             type: 'inbox',
           menutoogle:false,
        }
    },
    computed:{
         newmail () {
        if(this.$store.state.email.activeemail === 0)
          return true;
        else
          return false;
      },
      ...mapState({
        emaillist (state) {
          return state.email.emaillist.filter(list => {
            if (this.type === 'starred' && list.favourite === true)
              return list;
            else if (list.type === this.type)
              return list;
            else if (this.type === 'all mail')
              return list;
          });
        }
      }),
      ...mapGetters({
        activeemail: 'email/activeemail',
        index: 'email/index',
        send: 'email/send',
        draft: 'email/draft',
        trash: 'email/trash',
        starred: 'email/starred',
      })
    },
     methods: {
      getImgUrl(path) {
        return require('../../assets/images/'+path);
      },
      setactiveemails(id) {
        this.$store.dispatch('email/setactiveemails',id);
      },
      addfavorites(id) {
        this.$store.dispatch('email/addfavorites',id);
      }
    }
}
</script>