<template>
     <div class="col-xl-6 col-md-12">
                    <div class="email-right-aside">
                        <div class="card email-body radius-left">
                            <div class="ps-0">
                                <div class="tab-content">
                                    <div class="tab-pane fade active show" v-if="newmail != 0 "   id="pills-darkhome" role="tabpanel" aria-labelledby="pills-darkhome-tab">
                                        <div class="email-compose">
                                            <div class="email-top compose-border">
                                                <div class="row">
                                                    <div class="col-sm-8 xl-50">
                                                        <h4 class="mb-0">New Message</h4>
                                                    </div>
                                                    <div class="col-sm-4 btn-middle xl-50">
                                                        <button class="btn btn-primary btn-block btn-mail text-center mb-0 mt-0" type="button"><i class="fa fa-paper-plane me-2"></i> SEND</button>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="email-wrapper">
                                                <form class="theme-form">
                                                    <div class="form-group">
                                                        <label class="col-form-label pt-0" for="exampleInputEmail1">To</label>
                                                        <input class="form-control" id="exampleInputEmail1" type="email">
                                                    </div>
                                                    <div class="form-group">
                                                        <label for="exampleInputPassword1">Subject</label>
                                                        <input class="form-control" id="exampleInputPassword1" type="text">
                                                    </div>
                                                    <div class="form-group mb-0">
                                                        <label class="text-muted">Message</label>
                                                        <ckeditor :editor="editor"  v-model="editorData" :config="editorConfig"></ckeditor>
                                                    </div>
                                                </form>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="tab-pane fade active show " v-if="newmail == 0 " id="pills-darkprofile" role="tabpanel" aria-labelledby="pills-darkprofile-tab">
                                        <div class="email-content">
                                            <div class="email-top">
                                                <div class="row">
                                                    <div class="col-md-6 xl-100 col-sm-12">
                                                        <div class="media"><img class="me-3 rounded-circle img-fluid img-50 me-2" :src='getImgUrl(activeemail.image)' alt="">
                                                            <div class="media-body">
                                                                <h6 class="d-block">{{activeemail.name}} <small> <span class="digits">({{activeemail.date}})</span><span class="digits"> </span></small></h6>
                                                                <p>{{activeemail.subject}}</p>
                                                            </div>
                                                        </div>
                                                    </div>
                                                    <div class="col-md-6 col-sm-12">
                                                        <div class="d-flex user-details">
                                                            <p class="user-emailid">{{activeemail.email}}</p>
                                                            <span><samp @click="addfavorites(activeemail.id)" :class='{"text-warning":activeemail.favourite}'><i class="fa fa-star-o f-18 mt-1"></i></samp></span>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="email-wrapper">
                                                 <div class="emailread-group">
                                            <div class="read-group">
                                                <div v-html="activeemail.text"></div>
                                            </div>
                                                 </div>
                                                <div class="emailread-group">
                                                    <h6 class="text-muted"><i class="icofont icofont-clip"></i> ATTACHMENTS</h6><a class="text-muted text-end right-download font-primary f-w-600" href="#"><i class="fa fa-long-arrow-down me-2"></i>Download All</a>
                                                    <div class="clearfix"></div>
                                               <div class="attachment">
                                                    <ul>
                                                        <li ><img class="img-fluid" src="../../assets/images/email/1.jpg" alt=""></li>
                                                        <li ><img class="img-fluid" src="../../assets/images/email/2.jpg" alt=""></li>
                                                        <li ><img class="img-fluid" src="../../assets/images/email/3.jpg" alt=""></li>
                                                    </ul>
                                                </div> 
                                                 </div>
                                                <div class="action-wrapper">
                                                    <ul class="actions">
                                                       <li><a class="text-muted" href="#"><i class="fa fa-reply me-2"></i>Reply</a></li>
                                                        <li><a class="text-muted" href="#"><i class="fa fa-reply-all me-2"></i>Reply All</a></li>
                                                        <li><a class="text-muted" href="#"><i class="fa fa-share me-2"></i></a>Forward</li>
                                                    </ul>
                                                </div>
                                            </div>
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
 import CKEditor from '@ckeditor/ckeditor5-vue';
  import ClassicEditor from '@ckeditor/ckeditor5-build-classic';
export default {
    components: { ckeditor: CKEditor.component },
    data(){
      return {
        type: 'inbox',
        editor: ClassicEditor,
        editorData: '<p>Content of the editor.</p>',
        editorConfig: {
          config:{ height: '300px' }
        },
        menutoogle:false,
      };
    },
    computed: {
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