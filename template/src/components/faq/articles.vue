<template>
     <div class="col-lg-12 faq-articles ">
                  <div class="header-faq">
                    <h5 class="mb-0">Latest articles and videos</h5>
                  </div>
                  <div class="row">
                    <div class="col-xl-4 col-md-6">
                      <div class="row">
                        <div class="col-sm-12">
                          <div class="card">
                            <div class="card-body">
                              <div class="media"><i class="m-r-30" data-feather="codepen"><vue-feather type="codepen"></vue-feather></i>
                                <div class="media-body">
                                  <h6 class="f-w-600">Article Base Video</h6>
                                  <p>The web is a very big place, and if you are the typical internet user.</p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="col-sm-12">
                          <div class="card">
                            <div class="card-body">
                              <div class="media"><i class="m-r-30" data-feather="codepen"><vue-feather type="codepen"></vue-feather></i>
                                <div class="media-body">
                                  <h6 class="f-w-600">Validate website</h6>
                                  <p>website is the process of ensuring that the pages on the website conform.</p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="col-sm-12">
                          <div class="card">
                            <div class="card-body">
                              <div class="media"><i class="m-r-30" data-feather="youtube"><vue-feather type="youtube"></vue-feather></i>
                                <div class="media-body">
                                  <h6 class="f-w-600">sources Demos</h6>
                                  <p>simple demos of frequently asked questions about using the Libraries and information resources</p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="col-xl-4 col-md-6">
                      <div class="row">
                        <div class="col-sm-12">
                          <div class="card">
                            <div class="card-body">
                              <div class="media"><i class="m-r-30" data-feather="file-text"><vue-feather type="file-text"></vue-feather></i>
                                <div class="media-body">
                                  <h6 class="f-w-600">knows your sources</h6>
                                  <p>A book giving information on many subjects or on many aspects of one subject.</p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="col-sm-12">
                          <div class="card">
                            <div class="card-body">
                              <div class="media"><i class="m-r-30" data-feather="file-text"><vue-feather type="file-text"></vue-feather></i>
                                <div class="media-body">
                                  <h6 class="f-w-600">Tailwind Design</h6>
                                  <p>Tailwind is so low-level, it never encourages you to design the same site twice.</p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="col-sm-12">
                          <div class="card">
                            <div class="card-body">
                              <div class="media"><i class="m-r-30" data-feather="codepen"><vue-feather type="codepen"></vue-feather></i>
                                <div class="media-body">
                                  <h6 class="f-w-600">Validate Html</h6>
                                  <p>website is the process of ensuring that the pages on the website conform.</p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                   <articles2/>
                  </div>
                </div>
</template>
<script>
import articles2 from "./articles2.vue"
export default {
    components:{articles2}
}
</script>