<template>
     <div class="card">
                  <div class="card-header pb-0">
                    <h5>Custom state buttons</h5>
                   <span>The <code>.btn</code> class used with <code>&lt;button&gt;</code>, <code>&lt;a&gt;</code> and <code>&lt;input&gt;</code> elements.</span>
                  </div>
                  <div class="card-body btn-showcase"><a class="btn btn-square btn-primary" href="javascript:void(0)" data-bs-toggle="tooltip" title="btn btn-primary" role="button">Link</a>
                    <input class="btn btn-square btn-secondary" type="button" value="Input" data-bs-toggle="tooltip" title="btn btn-secondary">
                    <input class="btn btn-square btn-success" type="submit" value="Submit" data-bs-toggle="tooltip" title="btn btn-success">
                    <button class="btn btn-square btn-info" type="submit">Button</button>    
                  </div>
                </div>
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>outline buttons</h5>
                   <span>Add <code>.btn-square</code> and <code>.btn-outline-*</code> class for button with outline</span>
                  </div>
                  <div class="card-body btn-showcase">
                    <button class="btn btn-square btn-outline-primary" type="button">Primary Button</button>
                    <button class="btn btn-square btn-outline-secondary" type="button">Secondary Button</button>
                    <button class="btn btn-square btn-outline-success" type="button">Success Button</button>
                    <button class="btn btn-square btn-outline-info" type="button">Info Button</button>
                    <button class="btn btn-square btn-outline-warning" type="button">Warning Button</button>
                    <button class="btn btn-square btn-outline-danger" type="button">Danger Button</button>
                    <button class="btn btn-square btn-outline-light txt-dark" type="button">Light Button</button> 
                  </div>
                </div>
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>bold Border outline buttons</h5>
                    <span>Add <code>.btn-square</code> and <code>.btn-outline-*-2x</code> class for button with bold outline</span>
                  </div>
                  <div class="card-body btn-showcase">
                    <button class="btn btn-square btn-outline-primary-2x" type="button">Primary Button</button>
                    <button class="btn btn-square btn-outline-secondary-2x" type="button">Secondary Button</button>
                    <button class="btn btn-square btn-outline-success-2x" type="button">Success Button</button>
                    <button class="btn btn-square btn-outline-info-2x" type="button">Info Button</button>
                    <button class="btn btn-square btn-outline-warning-2x" type="button">Warning Button</button>
                    <button class="btn btn-square btn-outline-danger-2x" type="button">Danger Button</button>
                    <button class="btn btn-square btn-outline-light-2x txt-dark" type="button">Light Button</button>  
                  </div>
                </div>
                 <div class="card">
                  <div class="card-header pb-0">
                    <h5>outline Large buttons</h5>
                    <span>Add <code>.btn-square</code>,<code>.btn-outline-*</code> and <code>.btn-lg</code> class for large button</span>
                  </div>
                  <div class="card-body btn-showcase">
                    <button class="btn btn-square btn-outline-primary btn-lg" type="button">Primary Button</button>
                    <button class="btn btn-square btn-outline-secondary btn-lg" type="button">Secondary Button</button>
                    <button class="btn btn-square btn-outline-success btn-lg" type="button">Success Button</button>
                    <button class="btn btn-square btn-outline-info btn-lg" type="button">Info Button</button>
                    <button class="btn btn-square btn-outline-warning btn-lg" type="button">Warning Button</button>
                    <button class="btn btn-square btn-outline-danger btn-lg" type="button">Danger Button</button>
                    <button class="btn btn-square btn-outline-light btn-lg txt-dark" type="button">Light Button</button> 
                  </div>
                </div>
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>outline small buttons</h5>
                    <span>Add <code>.btn-square</code>,<code>.btn-outline-*</code> and <code>.btn-sm</code> class for small button</span>
                  </div>
                  <div class="card-body btn-showcase">
                    <button class="btn btn-square btn-outline-primary btn-sm" type="button">Primary Button</button>
                    <button class="btn btn-square btn-outline-secondary btn-sm" type="button">Secondary Button</button>
                    <button class="btn btn-square btn-outline-success btn-sm" type="button">Success Button</button>
                    <button class="btn btn-square btn-outline-info btn-sm" type="button">Info Button</button>
                    <button class="btn btn-square btn-outline-warning btn-sm" type="button">Warning Button</button>
                    <button class="btn btn-square btn-outline-danger btn-sm" type="button">Danger Button</button>
                    <button class="btn btn-square btn-outline-light btn-sm txt-dark" type="button">Light Button</button>  
                  </div>
                </div>
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>outline Extra small buttons</h5>
                    <span>Add <code>.btn-square</code>,<code>.btn-outline-*</code> and <code>.btn-xs</code> class for extra small button</span>
                  </div>
                  <div class="card-body btn-showcase">
                    <button class="btn btn-square btn-outline-primary btn-xs" type="button">Primary Button</button>
                    <button class="btn btn-square btn-outline-secondary btn-xs" type="button">Secondary Button</button>
                    <button class="btn btn-square btn-outline-success btn-xs" type="button">Success Button</button>
                    <button class="btn btn-square btn-outline-info btn-xs" type="button">Info Button</button>
                    <button class="btn btn-square btn-outline-warning btn-xs" type="button">Warning Button</button>
                    <button class="btn btn-square btn-outline-danger btn-xs" type="button">Danger Button</button>
                    <button class="btn btn-square btn-outline-light btn-xs txt-dark" type="button">Light Button</button>
                    </div>
                </div>
</template>
<script>
import { Tooltip } from 'bootstrap/dist/js/bootstrap.esm.min.js'
export default {
   mounted() {
     new Tooltip(document.body, {
      selector: "[data-bs-toggle='tooltip']",
    })
}
}
</script>