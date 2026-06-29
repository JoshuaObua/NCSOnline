use this exact css /* HTML: <div class="loader"></div> */
.loader {
  height: 80px;
  aspect-ratio: 1;
  padding: 10px;
  border-radius: 20px;
  box-sizing: border-box;
  position: relative;
  mask: conic-gradient(#000 0 0) content-box exclude,conic-gradient(#000 0 0);
  filter: blur(12px);
}
.loader:before {
  content: "";
  position: absolute;
  inset: 0;
  background: repeating-conic-gradient(#0000 0 5%,#C02942,#0000 20% 50%);
  animation: l3 1.5s linear infinite;
}
@keyframes l3 {
  to {rotate: 1turn}
}
add preloader to all public pages
add the main logo in between the preloader and ake it fit and center everthing 
use this as the html 
<div class="loader"></div>
