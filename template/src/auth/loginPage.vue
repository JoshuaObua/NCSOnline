<template>
    <div class="container-fluid p-0 ">
      <div class="row">
        <div class="col-12">
          <div class="login-card">
            <div>
            </div>
            <form class="theme-form login-form" @submit.prevent="">
              <h4>Login</h4>
              <h6>Welcome back! Log in to your account.</h6>
              <div class="form-group">
                <label>Email Address</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="icon-email"></i></span>
                  <input v-model="email" class="form-control" type="email" placeholder="Test@gmail.com">
                 
                </div>
              </div>
              <div class="form-group">
                <label>Password</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="icon-lock"></i></span>
                  <input v-model="password" :type="type" class="form-control" name="login[password]"
                      placeholder="*********">
                    <div class="show-hide"><span class="show" @click="showPassword"> </span></div>
                </div>
              </div>
              <div class="form-group">
                <div class="checkbox p-0 mb-3 w-100">
                  <input id="checkbox1" type="checkbox" />
                  <label class="text-muted" for="checkbox1">Remember password</label>
                  <a class="link" href="Forget Password">Forgot password?</a>
                </div>
                <button class="btn btn-primary btn-block" type="button" @click="doLogin">
                  Login
                </button>
              </div>
              <div class="login-social-title">
                <h5>Sign in with</h5>
              </div>
              <div class="form-group">
                <ul class="login-social">
                    <li><a href="https://www.linkedin.com" target="_blank"><vue-feather type="linkedin"></vue-feather></a>
                    </li>
                    <li><a href="https://twitter.com" target="_blank"><vue-feather type="twitter"></vue-feather></a></li>
                    <li><a href="https://www.facebook.com" target="_blank"><vue-feather type="facebook"></vue-feather></a>
                    </li>
                    <li><a href="https://www.instagram.com" target="_blank"><vue-feather
                          type="instagram"></vue-feather></a></li>
                  </ul>
              </div>
              <p class="mt-4 mb-0">
                Don't have account?<router-link class="ms-2"  to="/auth/register">
                  Create Account</router-link>
              </p>
            </form>
          </div>
        </div>
      </div>
    </div>
  </template>
  <script>

export default {
  data() {
    return {
      type: 'password',
      email: "test@admin.com",
      password: "test@123456",
    }
  },

  methods: {
    showPassword: function () {
      if (this.type === 'password') {
        this.type = 'text';
      } else {
        this.type = 'password';
      }
    },
    doLogin() {
      if (this.email === "test@admin.com" && this.password === "test@123456") {

        this.$router.replace('/');
        localStorage.setItem('user', this.email)
      } else {
        this.$toast.show('invalid email and password ', { theme: 'outline', position: 'top-right', icon: 'times', type: 'error', duration: 2000 });
      }
    },

  }
}
</script>