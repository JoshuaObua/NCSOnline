<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-xl-6 order-1 b-center bg-size" :style="{
        backgroundImage:
          'url(' + require('../../assets/images/login/1.jpg') + ')',
      }">
        <img class="bg-img-cover bg-center" alt="looginpage" style="display: none" />
      </div>
      <div class="col-xl-6 p-0 ">
        <div class="login-card">
          <form class="theme-form login-form" novalidate="" @submit.prevent="validate">
            <h4>Login</h4>
            <h6>Welcome back! Log in to your account.</h6>
            <div class="form-group">
              <label>Email Address</label>
              <div class="input-group">
                <span class="input-group-text"><i class="icon-email"></i></span>
                <input type="email" class="form-control"
                  v-bind:class="formSubmitted ? emailError ? 'is-invalid' : 'is-valid' : ''" id="email"
                  placeholder="Your email" v-model="email" />
              </div>
            </div>
            <div class="form-group">
              <label>Password</label>
              <div class="input-group">
                <span class="input-group-text"><i class="icon-lock"></i></span>
                <input class="form-control" type="password" name="login[password]"
                  v-bind:class=" formSubmitted ? passwordError ? 'is-invalid'  : 'is-valid' : '' " required=""
                  placeholder="*********" v-model="password" />
                <div class="show-hide"><span class="show"></span></div>
              </div>
            </div>
            <div class="form-group">
              <div class="checkbox">
                <input id="checkbox1" type="checkbox" />
                <label class="text-muted" for="checkbox1">Remember password</label>
              </div>
              <router-link class="link" to="/authentication/forgetpassword">Forgot password?</router-link>
            </div>
            <div class="form-group">
              <button class="btn btn-primary btn-block" type="submit">
                Sign in
              </button>
            </div>
            <div class="login-social-title">
              <h5>Sign in with</h5>
            </div>
            <div class="form-group">
              <ul class="login-social">
                <li> <a href="https://www.linkedin.com/login" target="_blank"><i data-feather="linkedin"><vue-feather
                        type="linkedin"></vue-feather></i></a> </li>
                <li> <a href="https://www.linkedin.com/login" target="_blank"><i data-feather="twitter"><vue-feather
                        type="twitter"></vue-feather></i></a></li>
                <li> <a href="https://www.linkedin.com/login" target="_blank"><i data-feather="facebook"><vue-feather
                        type="facebook"></vue-feather></i></a></li>
                <li> <a href="https://www.instagram.com/login" target="_blank"><i data-feather="instagram"><vue-feather
                        type="instagram"></vue-feather></i></a> </li>
              </ul>
            </div>
            <p> Don't have account?<router-link class="ms-2" to="/auth/register">Create Account</router-link></p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "home",
  data() {
    return {
      email: "",
      formSubmitted: false,
      emailError: false,
      emailFieldClasses: "",
      errors: [],
      password: "",
      passwordError: false,
    };
  },
  methods: {
    validate() {
      this.formSubmitted = true;
      this.errors = [];
      if (this.email.length < 10 || this.email.search("@") == -1) {
        this.emailError = true;
        this.errors.push({});
      } else {
        this.emailError = false;
        this.errors.push({});

      }
      if (this.password.length < 7) {
        this.passwordError = true;

      } else {
        this.passwordError = false;

      }
    },
  },
};
</script>