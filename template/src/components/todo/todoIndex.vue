<template>
  <Breadcrumbs title="todo" />
  <div class="container-fluid email-wrap bookmark-wrap todo-wrap">
    <div class="row">
      <div class="col-xl-3">
        <todoSidebar />
      </div>
      <div class="col-xl-9">
        <addTask />
      </div>
    </div>
  </div>
</template>
<script>
import addTask from "./addTask.vue"
import todoSidebar from "./todoSidebar.vue"

export default {
  components: {
    addTask,
    todoSidebar
  },

}
</script>