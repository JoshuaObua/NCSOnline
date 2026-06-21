<template>
    <Breadcrumbs main="Chat" title="Video Chat"/>
    <div class="container-fluid">
      <div class="row">
         <searchChat/>
         <chatBody/>
      </div>
    </div>
</template>
<script>
import chatBody from "./chatBody.vue"
import searchChat from "../chatList/searchChat.vue"
export default {
    components:{
      searchChat,
      chatBody
    }
}
</script>