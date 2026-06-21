<template>
    <div class="col-sm-12 col-xl-12 ">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>Colored breadcrumb</h5><span>use class <code>.breadcrumb-colored .bg-primary</code></span>
                  </div>
                  <div class="card-body">
                    <div class="row">
                      <div class="col-sm-12 col-md-6 col-lg-4">
                        <div class="card ribbon-wrapper">
                          <div class="card-body">
                            <div class="ribbon ribbon-primary ribbon-space-bottom">Ribbon</div>
                            <p>{{lorn}}</p>
                          </div>
                        </div>
                      </div>
                      <div class="col-sm-12 col-md-6 col-lg-4">
                        <div class="card ribbon-wrapper">
                          <div class="card-body">
                            <div class="ribbon ribbon-primary ribbon-space-bottom">Ribbon</div>
                            <p>{{lorn}}</p>
                          </div>
                        </div>
                      </div>
                      <div class="col-sm-12 col-md-6 col-lg-4">
                        <div class="card ribbon-wrapper">
                          <div class="card-body">
                            <div class="ribbon ribbon-primary ribbon-space-bottom">Ribbon</div>
                            <p>{{lorn}}</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-primary">Ribbon</div>
                    <p>{{lorn}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-secondary">Ribbon</div>
                    <p>{{lorn}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-success">Ribbon</div>
                    <p>{{lorn}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-info">Ribbon</div>
                    <p>{{lorn}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-warning">Ribbon</div>
                    <p>{{lorn}}</p>
                  </div>
                </div>
              </div>
              <div class="col-sm-14 col-md-6 col-lg-4 ">
                <div class="ribbon-wrapper card">
                  <div class="card-body">
                    <div class="ribbon ribbon-danger">Ribbon</div>
                    <p>{{lorn}}</p>
                  </div>
                </div>
              </div>
</template>
<script>
export default {
  data(){
    return{
      lorn:"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text."
    }
  }
}
</script>