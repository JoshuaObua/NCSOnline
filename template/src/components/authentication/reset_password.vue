<template>
	
      <div class="container-fluid p-0">
        <div class="row">
          <div class="col-12 ">     
            <div class="login-card">
               
                  <form class="theme-form login-form">
                    <h4>Create your account</h4>
                    <div class="form-group">
                  <label>New Password</label>
                  <div class="input-group"><span class="input-group-text"><i class="icon-lock"></i></span>
                    <input class="form-control" type="password" name="login[password]" required="" placeholder="*********">
                    <div class="show-hide"><span class="show"></span></div>
                  </div>
                </div>
                <div class="form-group">
                  <label>Rtype Password</label>
                  <div class="input-group"><span class="input-group-text"><i class="icon-lock"></i></span>
                    <input class="form-control" type="password" name="login[password]" required="" placeholder="*********">
                    <div class="show-hide"><span class="show"></span></div>
                  </div>
                </div>
                <div class="form-group">
                  <div class="checkbox">
                    <input id="checkbox1" type="checkbox">
                    <label for="checkbox1">Remember password</label>
                  </div>
                </div>
                  <div class="form-group"><a class="btn btn-primary btn-block" href="/" type="submit">Sign in</a></div>
                    <p class="mt-4 mb-0">Don't have account?
						      <router-link
                    class="ms-2"
                    tag="a"
                    to="/auth/register"
                  >
                    Create Account
                  </router-link>
					</p>
                  </form>
              
            </div>
          </div>
        </div>
      </div>
   
</template>
