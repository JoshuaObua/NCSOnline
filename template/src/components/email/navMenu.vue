<template>
    <div class="col-xl-3 col-md-6"> <div class="email-sidebar"><a class="btn btn-primary email-aside-toggle" href="javascript:void(0)" @click="collapseFilter()">email filter</a>
                    <div class="email-left-aside" :class="filtered?'open':''">
                        <div class="card">
                            <div class="card-body">
                                <div class="email-app-sidebar">
                                    <div class="media">
                                        <div class="media-size-email"><img class="me-3 rounded-circle img-50 me-2"  src="../../assets/images/user/1.jpg" alt=""></div>
                                        <div class="media-body">
                                            <h6 class="f-w-600">MARKJENCO</h6>
                                            <p>Markjecno@gmail.com</p>
                                        </div>
                                    </div>
                                    <ul class="nav main-menu" role="tablist">
                                        <li @click="setactiveemails(0)" class="nav-item"><a class="btn-primary btn-block btn-mail" id="pills-darkhome-tab" data-toggle="pill" href="#pills-darkhome" role="tab" aria-controls="pills-darkhome" aria-selected="true"><i class="icofont icofont-envelope me-2"></i> NEW MAIL</a></li>
                                        <li @click='type = "inbox"' class="nav-item"><a class="show" id="pills-darkprofile-tab" data-toggle="pill" href="#pills-darkprofile" role="tab" aria-controls="pills-darkprofile" aria-selected="false"><span class="title"><i class="icon-import"></i> Inbox</span><span class="badge pull-right digits">({{index}})</span></a></li>
                                        <li @click='type = "all mail"' ><a href="#"><span class="title"><i class="icon-folder"></i> All mail</span></a></li>
                                        <li @click='type = "sent"'><a href="#"><span class="title"><i class="icon-new-window"></i> Sent</span><span class="badge pull-right digits">({{send}})</span></a></li>
                                        <li @click='type = "draft"'><a href="#"><span class="title"><i class="icon-pencil-alt"></i> DRAFT</span><span class="badge pull-right digits">({{draft}})</span></a></li>
                                        <li @click='type = "trash"'><a href="#"><span class="title"><i class="icon-trash"></i> TRASH</span><span class="badge pull-right digits">({{trash}})</span></a></li>
                                        <li><a href="#"><span class="title"><i class="icon-info-alt"></i> IMPORTANT</span></a></li>
                                        <li @click='type = "starred"'><a href="#"><span class="title"><i class="icon-star"></i> Starred</span><span class="badge pull-right digits">({{starred}})</span></a></li>
                                        <li>
                                            <hr>
                                        </li>
                                        <li><a href="#"><span class="title"><i class="icon-email"></i> UNREAD</span></a></li>
                                        <li><a href="#"><span class="title"><i class="icon-export"></i> SPAM</span></a></li>
                                        <li><a href="#"><span class="title"><i class="icon-share"></i> OUTBOX</span></a></li>
                                        <li><a href="#"><span class="title"><i class="icon-file"></i> UPDATE</span></a></li>
                                        <li><a href="#"><span class="title"><i class="icon-bell"></i> ALERT</span></a></li>
                                        <li><a href="#"><span class="title"><i class="icon-notepad"></i> NOTES</span><span class="badge pull-right digits">(20)</span></a></li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                    </div>
                </div>
                <div class="col-xl-3">
                    <div class="email-right-aside">
                        <div class="card email-body">
                            <div class="pe-0 b-r-light">
                                <div class="email-top">
                                    <div class="row">
                                        <div class="col">
                                            <h5 class="text-capitalize">{{type}}</h5>
                                        </div>
                                        <div class="col text-end d-flex justify-content-end">
                                            <div class="dropdown">
                                                <button class="btn bg-transparent dropdown-toggle p-0 text-muted" id="dropdownMenuButton" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" v-on:click="menutoogle = !menutoogle">More</button>
                                                <div class="dropdown-menu" aria-labelledby="dropdownMenuButton" :class="{ show: menutoogle }">
                                                    <a class="dropdown-item" href="#">Action</a>
                                                    <a class="dropdown-item" href="#">Another action</a>
                                                    <a class="dropdown-item" href="#">Something else here</a>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="inbox">
                                    <div class="media active" v-for="(list,index) in emaillist" :key="index" @click="setactiveemails(list.id)">
                                        <div class="media-size-email">
                                            <img class="me-3 rounded-circle img-fluid img-50 me-2" :src='getImgUrl(list.image)' alt=""></div>
                                        <div class="media-body">
                                            <h6>{{list.name}}<small><span class="digits"></span><span class="digits"></span></small></h6>
                                            <p>{{list.subject}}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
</template>
<script>
import { mapState,mapGetters } from 'vuex';
export default {
  components:{
  },
    data(){
     return{
          type: 'inbox',
           menutoogle:false,
           filtered: false,
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
      collapseFilter() {
        this.filtered = !this.filtered;
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