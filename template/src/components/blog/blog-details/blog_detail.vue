<template>
    <Breadcrumbs main="Blog" title="Blog Details"/>
          <div class="container-fluid blog-page ">
            <div class="row">
                 <div class="col-xxl-6 set-col-12 box-col-5 xl-40">
                <div class="card">
                  <div class="blog-box blog-shadow bg-size">
                   
                    <div class="blog-details">
                      <p>25 July 2022</p>
                      <h4>A huge part of it is the incomparable beauty you can encounter every day.</h4>
                      <ul class="blog-social">
                        <li><i class="icofont icofont-user"></i>Mark Jecno</li>
                        <li><i class="icofont icofont-thumbs-up"></i>02 Hits</li>
                        <li><i class="icofont icofont-ui-chat"></i>598 Comments</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
              <detail/>
              <detail1/>
            </div>
          </div>
</template>
<script>
import detail from "./detail.vue"
import detail1 from "./detail1.vue"
export default {
    components:{
      detail,
      detail1
    }
}
</script>