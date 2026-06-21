<template>
     <div class="row">
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-bookmark ribbon-primary">Ribbon</div>
                    <p>{{ Lorem}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-bookmark ribbon-secondary">Ribbon</div>
                    <p>{{ Lorem}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-bookmark ribbon-success">Ribbon</div>
                    <p>{{ Lorem}}</p>
                  </div>
                </div>
              </div>
            </div>
            <div class="row">
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="card ribbon-wrapper-right">
                  <div class="card-body">
                    <div class="ribbon ribbon-bookmark ribbon-right ribbon-primary">Ribbon</div>
                    <p>{{ Lorem}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="card ribbon-wrapper-right">
                  <div class="card-body">
                    <div class="ribbon ribbon-bookmark ribbon-right ribbon-secondary">Ribbon</div>
                    <p>{{ Lorem}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="card ribbon-wrapper-right">
                  <div class="card-body">
                    <div class="ribbon ribbon-bookmark ribbon-right ribbon-success">Ribbon</div>
                    <p>{{ Lorem}}</p>
                  </div>
                </div>
              </div>
            </div>
</template>
<script>
export default {
  data(){
    return{
      Lorem:"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text."
    }
  }
}
</script>