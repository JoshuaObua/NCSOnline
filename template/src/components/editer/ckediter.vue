<template>
  <Breadcrumbs main="Editors" title="CK Editors" />
  <div class="container-fluid">
    <div class="row">
      <div class="col-md-12 ">
        <div class="card">
          <div class="card-header pb-0">
            <h5 class="card-title">Ck Editor</h5>
          </div>
          <div class="card-body">
            <ckeditor :editor="editor" v-model="editorData">

            </ckeditor>
          </div>
        </div>
      </div>

    </div>
    <div class="row">
      <div class="col-sm-12 ">
        <div class="card">
          <div class="card-header pb-0">
            <h5>Inline Editor</h5>
          </div>
          <div class="card-body">
            <div id="area1" contenteditable="true">
              <h1 class="f-w-700">Your title</h1>
              <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec at vulputate urna, sed dignissim arcu.
                Aliquam at ligula imperdiet, faucibus ante a, interdum enim. Sed in mauris a lectus lobortis
                condimentum. Sed in nunc magna. Quisque massa urna, cursus vitae commodo eget, rhoncus nec erat. Sed
                sapien turpis, elementum sit amet elit vitae, elementum gravida eros. In ornare tempus nibh ut lobortis.
                Nam venenatis justo ex, vitae vulputate neque laoreet a.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import CKEditor from '@ckeditor/ckeditor5-vue';
import ClassicEditor from '@ckeditor/ckeditor5-build-classic';
export default {
  components: {
    ckeditor: CKEditor.component
  },
  data() {
    return {
      editor: ClassicEditor,
      editorData: '<p>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p>',

    };
  }

}
</script>