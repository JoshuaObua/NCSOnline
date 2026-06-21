<template>
    <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>Basic Examples</h5>
                  </div>
                  <div class="card-body btn-showcase">
                    <button class="btn btn-primary sweet-1" type="button" v-on:click="basic" >Basic</button>
                    <button class="btn btn-primary sweet-2" type="button" v-on:click="basic_with_title" >With Title alert</button>
                    <button class="btn btn-success sweet-3" type="button" v-on:click="basic_success_alert" >Success alert</button>
                    <button class="btn btn-info sweet-4" type="button" v-on:click="basic_info_alert">Info alert</button>
                    <button class="btn btn-warning sweet-5" type="button" v-on:click="basic_warning_alert">Warning alert</button>
                  </div>
                </div>
              </div>
</template>
<script>
export default {
    methods:{
         basic:function(){
        this.$swal({
          title:'Hello world!',
            confirmButtonColor: '#24695c ',
        });
      },
      basic_warning_alert:function(){
        this.$swal({
            icon: 'warning',
            title:"Are you sure?",
          text:'Once deleted, you will not be able to recover this imaginary file!',
          showCancelButton: true,
           confirmButtonText: 'Ok',
          confirmButtonColor: '#24695c ',
          cancelButtonText: 'Cancel',
        }).then((result)=>{
          if(result.value){
            this.$swal({
             icon: 'success',
              text:'Poof! Your imaginary file has been deleted!',
              type:'success',
               confirmButtonColor: '#24695c ',
            });
          }else{
            this.$swal({
              text:'Your imaginary file is safe!',
               confirmButtonColor: '#24695c ',
            });
          }
        });
      },
      basic_with_title:function(){
        this.$swal({
          title:'Here\'s the title',
          text:'...and here\'s the text!',
           confirmButtonColor: '#24695c ',
        });
      },
      basic_success_alert:function(){
        this.$swal({
          icon: 'warning',
          title:'Good job!',
          text:'You clicked the button!',
          type:'info',
           confirmButtonColor: '#24695c ',
        });
      },
      basic_info_alert:function(){
        this.$swal({
          text:' Click on either the button or outside the modal.',
          confirmButtonColor: '#24695c ',
        }).then((result)=>{
          if(result.value){
            this.$swal({
              text:'The returned value is:'+result.value,
               confirmButtonColor: '#24695c ',
            });
          }else{
            this.$swal({
              text:'The returned value is:'+result.value
            });
          }
        });
      }
    }
}
</script>