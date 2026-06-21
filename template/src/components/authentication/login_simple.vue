<template>
     <div class="container-fluid p-0">
        <div class="row">
          <div class="col-12 ">
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
                    <div class="show-hide"><span class="show"></span></div>
                  </div>
                </div>
                <div class="form-group">
                  <div class="checkbox">
                    <input id="checkbox1" type="checkbox">
                    <label for="checkbox1">Remember password</label>
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
                <p>Don't have account?<a class="ms-2" href="#">Create Account</a></p>
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
             password: '',
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