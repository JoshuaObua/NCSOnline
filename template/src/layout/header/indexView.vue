<template>

  <div class="main-header-right row m-0">

    <div class="main-header-left">
      <SideBar />
    </div>
    <div class="left-menu-header col">
      <ul>
        <li>
          <searchBar />
        </li>
      </ul>
    </div>
    <div class="nav-right col pull-right right-menu p-0">
      <ul class="nav-menus" :class="filtered ? 'open' : ''">
        <li>
          <maxiMaize />
        </li>
        <li class="language-nav">
          <languageView/>
          </li>
        <li class="onhover-dropdown">
          <bookMark />
        </li>
        <li class="onhover-dropdown">
          <notifiCations />
        </li>
        <li class="onhover-dropdown">
          <modeView />
        </li>
        <li class="onhover-dropdown">
          <ChatView />
        </li>
        <li class="onhover-dropdown p-0">
          <LogOut />
        </li>
      </ul>
    </div>
    <div class="d-lg-none mobile-toggle pull-right w-auto"><vue-feather
          type="more-horizontal" @click="collapseFilter()"></vue-feather></div>
  </div>

</template>
<script>
import { mapState } from 'vuex';
import SideBar from "./sidebar.vue"
import maxiMaize from "./maxiMaize.vue"
import bookMark from "./BookMark.vue"
import notifiCations from "./notifiCations.vue"
import modeView from "./modeView.vue"
import ChatView from "./ChatView.vue"
import searchBar from "./searchBar.vue"
import LogOut from "./LogOut.vue"
import languageView from "./languageView.vue"

export default {
  components: {
    SideBar,
    maxiMaize,
    bookMark,
    notifiCations,
    modeView,
    ChatView,
    searchBar,
    LogOut,
    languageView

  },
  data() {
    return {
      bookmark: false,
      toogle: false,
      filtered: false,
    };
  },
  computed: {
    ...mapState({
      menuItems: (state) => state.menu.searchData,
      megamenuItems: (state) => state.menu.megamenu,
      searchOpen: (state) => state.menu.searchOpen,
    }),
  },
  methods: {
    search_open() {
      this.$store.state.menu.searchOpen = true;
    },
    collapseFilter() {
      this.filtered = !this.filtered;
    },
    bookmark_open() {
      this.bookmark = !this.bookmark;
    },
  },
}
</script>