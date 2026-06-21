<template>
     <div class="col-sm-6 col-xl-3 box-col-6 des-xl-50">
                <div class="card">
                  <div class="blog-box blog-grid">
                    <div class="blog-wrraper"><a href="#"><router-link to="/blog/single"><img class="img-fluid top-radius-blog" src="../../../assets/images/blog/blog-5.jpg" alt=""></router-link></a></div>
                    <div class="blog-details-second">
                      <div class="blog-post-date"><span class="blg-month">April</span><span class="blg-date">09</span></div><a href="#"><router-link to="/blog/single">
                        <h6 class="blog-bottom-details">Encounter every day.</h6></router-link></a>
                      <p>A huge part of it is the incomparable beauty you can encounter every day.</p>
                      <div class="detail-footer"> 
                        <ul class="sociyal-list">
                          <li><i class="fa fa-user-o"></i>admin</li>
                          <li><i class="fa fa-comments-o"></i>5</li>
                          <li><i class="fa fa-thumbs-o-up"></i>2 like</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-sm-6 col-xl-3 box-col-6 des-xl-50">
                <div class="card">
                  <div class="blog-box blog-grid">
                    <div class="blog-wrraper"><a href="#"><router-link to="/blog/single"><img class="img-fluid top-radius-blog" src="../../../assets/images/blog/blog-6.jpg" alt=""></router-link></a></div>
                    <div class="blog-details-second">
                      <div class="blog-post-date"><span class="blg-month">April</span><span class="blg-date">30</span></div><a href="#"><router-link to="/blog/single">
                        <h6 class="blog-bottom-details">Black color is important.</h6></router-link></a>
                      <p>A garden where plants that flower are grown and displayed in floral garden.</p>
                      <div class="detail-footer"> 
                        <ul class="sociyal-list">
                          <li><i class="fa fa-user-o"></i>admin</li>
                          <li><i class="fa fa-comments-o"></i>5</li>
                          <li><i class="fa fa-thumbs-o-up"></i>0 like</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-sm-6 col-xl-3 box-col-6 des-xl-50">
                <div class="card">
                  <div class="blog-box blog-grid">
                    <div class="blog-wrraper"><a href="#"><router-link to="/blog/single"><img class="img-fluid top-radius-blog" src="../../../assets/images/blog/blog-5.jpg" alt=""></router-link></a></div>
                    <div class="blog-details-second">
                      <div class="blog-post-date"><span class="blg-month">April</span><span class="blg-date">06</span></div><a href="#"><router-link to="/blog/single">
                        <h6 class="blog-bottom-details">Encounter every day.</h6></router-link></a>
                      <p>A huge part of it is the incomparable beauty you can encounter every day.</p>
                      <div class="detail-footer"> 
                        <ul class="sociyal-list">
                          <li><i class="fa fa-user-o"></i>admin</li>
                          <li><i class="fa fa-comments-o"></i>5</li>
                          <li><i class="fa fa-thumbs-o-up"></i>0 like</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-sm-6 col-xl-3 box-col-6 des-xl-50">
                <div class="card">
                  <div class="blog-box blog-grid">
                    <div class="blog-wrraper"><a href="#"><router-link to="/blog/single"><img class="img-fluid top-radius-blog" src="../../../assets/images/blog/blog-6.jpg" alt=""></router-link></a></div>
                    <div class="blog-details-second">
                      <div class="blog-post-date"><span class="blg-month">April</span><span class="blg-date">05</span></div><a href="#"><router-link to="/blog/single">
                        <h6 class="blog-bottom-details">Black color is important.</h6></router-link></a>
                      <p>A garden where plants that flower are grown and displayed in floral garden.</p>
                      <div class="detail-footer"> 
                        <ul class="sociyal-list">
                          <li><i class="fa fa-user-o"></i>admin</li>
                          <li><i class="fa fa-comments-o"></i>5</li>
                          <li><i class="fa fa-thumbs-o-up"></i>0 like</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
</template>
<script>
export default {
  data(){
    return{
      items:[{id:1, title:"Perspiciatis unde omnis.", desc:"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt.", image:"blog-5.jpg"},
      {id:1, title:"Perspiciatis unde omnis.", desc:"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt.", image:"blog-5.jpg"},
      {id:1, title:"Perspiciatis unde omnis.", desc:"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt.", image:"blog-5.jpg"},
      {id:1, title:"Perspiciatis unde omnis.", desc:"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt.", image:"blog-5.jpg"},
      ]
    }
  },
  methods: {
    getImageUrl(path) {
      return require("../../../assets/images/blog/" + path);
    },
  },
}
</script>