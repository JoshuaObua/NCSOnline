<template>
  <div class="tab-pane fade" :class="{ 'active show': customizer == 'layouts' }" id="c-pills-home">
    <ul class="sidebar-type layout-grid layout-types">
      <li data-attr="compact-sidebar">
        <div class="layout-img">
          <a @click="handlePageLayputs('compact-wrapper', 'Defaul')"><img src="../../assets/images/demo/1.jpg"
              class="img-fluid" alt="oneone" /></a>

          <h6>Defaul layout</h6>
        </div>
      </li>
      <li class="only-body" data-attr="default-body">
        <div class="layout-img">
          <a @click="
            handlePageLayputs('compact-wrapper compact-sidebar', 'Compact')
          "><img src="../../assets/images/demo/2.jpg" class="img-fluid" alt="one" /></a>

          <h6>Compact layout</h6>
        </div>
      </li>
      <li data-attr="modern-layout">
        <div class="layout-img">
          <a @click="
            handlePageLayputs('compact-wrapper modern-sidebar ', 'Modern')
          ">
            <img src="../../assets/images/demo/3.jpg" class="img-fluid" alt="one" /></a>
          <h6>Modern layout</h6>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { mapState } from "vuex";
import { layoutClasses } from "../../constants/layout";
export default {
  name: "Customlayout",
  computed: {
    ...mapState({
      customizer: (state) => state.menu.customizer,
      layout: (state) => state.layout.layout,
      menuItems: (state) => state.menu.data,
      togglesidebar: (state) => state.menu.togglesidebar,
    }),
  },
  methods: {
    handlePageLayputs(val, layout) {
      const layoutobj = layoutClasses.find(
        (item) => Object.keys(item).pop() === layout
      );
      this.$store.dispatch("layout/setCustomizeSidebarType", layoutobj[layout]);
      this.layout.settings.layout = layout;
      this.$router.push({ query: { layout: layout } }).catch((err) => err);
      if (layoutobj[layout].split(" ").includes("compact-sidebar")) {
        this.$store.state.menu.activeoverlay = true;
        document.getElementById("mainnav");
      } else {
        this.$store.state.menu.activeoverlay = false;
        document.getElementById("mainnav");
      }
      if (layout === "Defaul") {
        this.$store.state.layout.boxlayout = true;
      } else {
        this.$store.state.layout.boxlayout = false;
      }

    },
  },
};
</script>
