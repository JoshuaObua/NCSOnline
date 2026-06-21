<template>

 <div id="mainnav" >
    <ul class="nav-menu custom-scrollbar" id="myDIV"
    :style="[
        layoutobject.split(' ').includes('horizontal-wrapper') ? layout.settings.layout_type == 'rtl' ? { '  -right': margin + 'px' } : { 'margin-left': margin + 'px' } : { margin: '0px' },
      ]"
    >
     <li class="back-btn">
        <div class="mobile-back text-end">
          <span>Back</span> <i class="fa fa-angle-right ps-2" aria-hidden="true"></i>
        </div>
      </li>
      <li
        v-for="(menu, index) in menuItems"
        :key="index"
        :class="{
          'sidebar-main-title': menu.type == 'headtitle',
        }"
        class="dropdown d-block"
      >
        <div v-if="menu.type == 'headtitle'">
          <h6 class="lan-1">{{ menu.headTitle1 }}</h6>
        </div>

      
        <a
          href="javascript:void(0)"
           :class="{ active: menu.active}"
          class="nav-link menu-title"
          v-if="menu.type == 'sub'"
          @click="setNavActive(menu, index)"
          ><vue-feather :type="menu.icon" class="home"></vue-feather>
          <span >
            {{ $t(menu.title) }}
          </span>
          <div class="according-menu" v-if="menu.children">
            <i class="pull-right" :class="[menu.active ?'fa fa-angle-down':'fa fa-angle-right' ]" ></i>
          </div></a
        >
      
        <router-link
          :to="menu.path"
          class="nav-link menu-title"
          v-if="menu.type == 'link'"
          router-link-exact-active
          exact
          v-on:click="hidesecondmenu()"
        >
          <vue-feather :type="menu.icon" class="home"></vue-feather>
          <span class="menu1">
            {{$t(menu.title) }}
          </span>
          <vue-feather type="chevron-right" v-if="menu.children"></vue-feather>
        </router-link>
        <a
          :href="menu.path"
          class="sidebar-link sidebar-title"
          v-if="menu.type == 'extLink'"
          @click="setNavActive(menuItem, index)"
        >
          <feather :type="menu.icon" class="top"></feather>
          <span>
            {{ $t(menu.title) }}
          </span>
          <i class="fa fa-angle-right pull-right" v-if="menuItem.children"></i>
        </a>
     
        <a
          :href="menuItem.path"
          target="_blank"
          class="sidebar-link sidebar-title"
          v-if="menu.type == 'extTabLink'"
          @click="setNavActive(menu, index)"
        >
          <feather :type="menuItem.icon" class="top"></feather>
          <span>
            {{ $t(menu.title) }}
          </span>
          <i class="fa fa-angle-right pull-right" v-if="menu.children"></i>
        </a>
      
        <ul
          class="nav-submenu menu-content"
          v-if="menu.children"
          :class="{ 'menu-open d-block': menu.active }"
        >
          <li
            v-for="(childrenItem, index) in menu.children"
            :key="index"
           
          >
         
            <a
              class="submenu-title"
              href="javascript:void(0)"
              v-if="childrenItem.type == 'sub'"
              @click="setNavActive(childrenItem, index)"
            >
              {{ $t(childrenItem.title) }}
              <label
                :class="'badge badge-' + childrenItem.badgeType + ' pull-right'"
                v-if="childrenItem.badgeType"
                >{{ $t(childrenItem.badgeValue) }}</label
              >
             
                <i class=" pull-right" :class="[childrenItem.active ?'fa fa-angle-down':'fa fa-angle-right' ]" v-if="childrenItem.children"></i>
            </a>
          
            <router-link
              class="submenu-title"
              :to="childrenItem.path"
              v-if="childrenItem.type == 'link'"
              router-link-exact-active
              exact
              v-on:click="hidesecondmenu()"
            >
              {{ $t(childrenItem.title) }}
              <label
                :class="'badge badge-' + childrenItem.badgeType + ' pull-right'"
                v-if="childrenItem.badgeType"
                >{{ $t(childrenItem.badgeValue) }}</label
              >
              <i
                class="fa fa-angle-right pull-right mt-1"
                v-if="childrenItem.children"
              ></i>
            </router-link>
          
            <a
              :href="childrenItem.path"
              v-if="childrenItem.type == 'extLink'"
              class="submenu-title"
            >
              {{ $t(childrenItem.title) }}
              <label
                :class="'badge badge-' + childrenItem.badgeType + ' pull-right'"
                v-if="childrenItem.badgeType"
                >{{ $t(childrenItem.badgeValue) }}</label
              >
              <i
                class="fa fa-angle-right pull-right mt-1"
                v-if="childrenItem.children"
              ></i>
            </a>
          
            <a
              class="submenu-title"
              :href="childrenItem.path"
              target="_blank"
              v-if="childrenItem.type == 'extTabLink'"
            >
              {{ $t(childrenItem.title) }}
              <label
                :class="'badge badge-' + childrenItem.badgeType + ' pull-right'"
                v-if="childrenItem.badgeType"
                >{{ $t(childrenItem.badgeValue) }}</label
              >
              <i
                class="fa fa-angle-right pull-right mt-1"
                v-if="childrenItem.children"
              ></i>
            </a>
           
            <ul
              class="nav-sub-childmenu submenu-content"
              v-if="childrenItem.children"
              :class="{ opensubchild: childrenItem.active }"
            >
              <li
                v-for="(childrenSubItem, index) in childrenItem.children"
                :key="index"
              >
              
                <router-link
                  :to="childrenSubItem.path"
                  v-if="childrenSubItem.type == 'link'"
                  router-link-exact-active
                  exact
                  v-on:click="hidesecondmenu()"
                >
                  {{ $t(childrenSubItem.title )}}
                  <label
                    :class="
                      'badge badge-' + childrenSubItem.badgeType + ' pull-right'
                    "
                    v-if="childrenSubItem.badgeType"
                    >{{ childrenSubItem.badgeValue }}</label
                  >
                  <i
                    class="fa fa-angle-right pull-right"
                    v-if="childrenSubItem.children"
                  ></i>
                </router-link>
               
                <router-link
                  :to="childrenSubItem.path"
                  v-if="childrenSubItem.type == 'extLink'"
                  router-link-exact-active
                >
                  {{ childrenSubItem.title }}
                  <label
                    :class="
                      'badge badge-' + childrenSubItem.badgeType + ' pull-right'
                    "
                    v-if="childrenSubItem.badgeType"
                    >{{ childrenSubItem.badgeValue }}</label
                  >
                  <i
                    class="fa fa-angle-right pull-right"
                    v-if="childrenSubItem.children"
                  ></i>
                </router-link>
              
                <router-link
                  :to="childrenSubItem.path"
                  v-if="childrenSubItem.type == 'extLink'"
                  router-link-exact-active
                >
                  {{ childrenSubItem.title }}
                  <label
                    :class="
                      'badge badge-' + childrenSubItem.badgeType + ' pull-right'
                    "
                    v-if="childrenSubItem.badgeType"
                    >{{ childrenSubItem.badgeValue }}</label
                  >
                  <i
                    class="fa fa-angle-right pull-right"
                    v-if="childrenSubItem.children"
                  ></i>
                </router-link>
              </li>
            </ul>
          </li>
        </ul>
      </li>
    </ul>
  </div>

</template>
<script>

import { mapState } from "vuex";
import { layoutClasses } from "../../constants/layout";
export default {
  components: {
    
  },
  data() {
    return {
      layoutobj: {},
    };
  },
 computed: {
    ...mapState({
      menuItems: (state) => state.menu.data,
      layout: (state) => state.layout.layout,
      sidebar: (state) => state.layout.sidebarType,
      activeoverlay: (state) => state.menu.activeoverlay,
      togglesidebar: (state) => state.menu.togglesidebar,
      width: (state) => state.menu.width,
      height: (state) => state.menu.height,
      margin: (state) => state.menu.margin,
      menuWidth: (state) => state.menu.menuWidth,
    }),
    sidebarType(){
      return this.$route.query.layout || 'default'
    },
    layoutobject: {
      get: function () {
        return JSON.parse(
          JSON.stringify(
            layoutClasses.find(
              (item) => Object.keys(item).pop() === this.layout.settings.layout
            )
          )
        )[this.layout.settings.layout];
      },
      set: function () {
        this.layoutobj = layoutClasses.find(
          (item) => Object.keys(item).pop() === this.layout.settings.layout
        );
        this.layoutobj = JSON.parse(JSON.stringify(this.layoutobj))[
          this.layout.settings.layout
        ];
        return this.layoutobj;
      },
    },
  },
  watch: {
    width() {
      window.addEventListener("resize", this.handleResize);
      this.handleResize();
      window.addEventListener("scroll", this.handleScroll);
      this.handleScroll();
      if (window.innerWidth < 992) {
        const newlayout = JSON.parse(
          JSON.stringify(this.layoutobject).replace(
            "horizontal-wrapper",
            "compact-wrapper"
          )
        );
        document.querySelector(".page-wrapper").className =
          "page-wrapper " + newlayout;
        this.$store.state.menu.margin = 0;
      } else {
        document.querySelector(".page-wrapper").className =
          "page-wrapper " + this.layoutobject;
      }
    },
  },
  created() {
    window.addEventListener("resize", this.handleResize);
    this.handleResize();
  
    setTimeout(() => {
      const elmnt = document.getElementById("myDIV");
      this.$store.state.menu.menuWidth = elmnt.offsetWidth;
      this.$store.state.menu.menuWidth > window.innerWidth
        ? ((this.$store.state.menu.hideRightArrow = false),
          (this.$store.state.menu.hideLeftArrowRTL = false))
        : ((this.$store.state.menu.hideRightArrow = false),
          (this.$store.state.menu.hideLeftArrowRTL = true));
    }, 500);
    this.layoutobject = layoutClasses.find(
      (item) => Object.keys(item).pop() === this.layout.settings.layout
    );
    this.layoutobject = JSON.parse(JSON.stringify(this.layoutobject))[
      this.layout.settings.layout
    ];
  },
 
  mounted() {
    this.menuItems.filter((items) => {
      if (items.path === this.$route.path)
        this.$store.dispatch("menu/setActiveRoute", items);
      if (!items.children) return false;
      items.children.filter((subItems) => {
        if (subItems.path === this.$route.path)
          this.$store.dispatch("menu/setActiveRoute", subItems);
        if (!subItems.children) return false;
        subItems.children.filter((subSubItems) => {
          if (subSubItems.path === this.$route.path)
            this.$store.dispatch("menu/setActiveRoute", subSubItems);
        });
      });
    });
  },
  methods: {
    handleScroll() {
      if (window.scrollY > 400) {
        if (
          this.layoutobject.split(" ").pop() === "material-type" ||
          this.layoutobject.split(" ").pop() === "normal-sidebar"
        )
          document.querySelector(".sidebar-main").className =
            "sidebar-main hovered";
      } else {
        if (document.getElementById("sidebar-main"))
          document.querySelector(".sidebar-main").className = "sidebar-main";
      }
    },
    setNavActive(item) {
      this.$store.dispatch("menu/setNavActive", item);
      if (
        this.layoutobject.split(" ").includes("compact-sidebar") &&
        window.innerWidth > 991
      ) {
        if (this.menuItems.some((menuItem) => menuItem.active === true)) {
          this.$store.state.menu.activeoverlay = true;
        } else {
          this.$store.state.menu.activeoverlay = false;
        }
      }
    },
    hidesecondmenu() {
      if (window.innerWidth < 991) {
        (this.$store.state.menu.activeoverlay = false),
          (this.$store.state.menu.togglesidebar = false);
        this.menuItems.filter((menuItem) => {
          menuItem.active = false;
        });
      } else if (this.layoutobject.split(" ").includes("compact-sidebar")) {
        (this.$store.state.menu.activeoverlay = false),
          this.menuItems.filter((menuItem) => {
            menuItem.active = false;
          });
      }
    },
    handleResize() {
      this.$store.state.menu.width = window.innerWidth - 300;
    },
  },
};
</script>