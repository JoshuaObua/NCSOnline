<template>
    <div class="container-fluid">
        <div class="row">
           <div class="col-xl-7 b-center bg-size" :style="{ backgroundImage: 'url('+ require('../../assets/images/login/2.jpg') +')' }">
        <img
          class="bg-img-cover bg-center"
          src="../../assets/images/login/2.jpg"
          alt="looginpage"
          style="display: none;"
        />
           </div>
          <div class="col-xl-5 p-0">
            <div class="login-card">
              <form class="theme-form login-form">
                <h4>Login</h4>
                <h6>Welcome back! Log in to your account.</h6>
                <div class="form-group">
                  <label>Email Address</label>
                  <div class="input-group"><span class="input-group-text"><i class="icon-email"></i></span>
                    <input class="form-control" type="email" required="" placeholder="Test@gmail.com">
                  </div>
                </div>
                <div class="form-group">
                  <label>Password</label>
                  <div class="input-group"><span class="input-group-text"><i class="icon-lock"></i></span>
                    <input class="form-control" type="password" name="login[password]" required="" placeholder="*********">
                    <div class="show-hide"><span class="show" @click="showPassword"></span></div>
                  </div>
                </div>
                <div class="form-group">
                  <div class="checkbox">
                    <input id="checkbox1" type="checkbox">
                    <label class="text-muted" for="checkbox1">Remember password</label>
                  </div><router-link class="link" to="/authentication/forgetpassword">Forgot password?</router-link>
                </div>
                <div class="form-group">
                  <button class="btn btn-primary btn-block" type="submit">Sign in</button>
                </div>
                <div class="login-social-title">                
                  <h5>Sign in with</h5>
                </div>
                <div class="form-group">
                  <ul class="login-social">
                    <li><a href="https://www.linkedin.com/login" target="_blank"><i data-feather="linkedin"><vue-feather type="linkedin"></vue-feather></i></a></li>
                    <li><a href="https://www.linkedin.com/login" target="_blank"><i data-feather="twitter"><vue-feather type="twitter"></vue-feather></i></a></li>
                    <li><a href="https://www.linkedin.com/login" target="_blank"><i data-feather="facebook"><vue-feather type="facebook"></vue-feather></i></a></li>
                    <li><a href="https://www.instagram.com/login" target="_blank"><i data-feather="instagram"><vue-feather type="instagram"></vue-feather></i></a></li>
                  </ul>
                </div>
                <p>Don't have account?<router-link class="ms-2" to="/auth/register">Create Account</router-link></p>
              </form>
            </div>
          </div>
        </div>
      </div>
</template>
<script>
export default {
  data(){
    return{
      type:"password",
       submitted: false,
    }
  },
  methods:{
    showPassword: function() {
        if (this.type === 'password') {
          this.type = 'text';
        } else {
          this.type = 'password';
        }
      },
  }
}
</script>