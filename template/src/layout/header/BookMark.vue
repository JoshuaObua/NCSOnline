<template>
     <div class="bookmark-box"><i data-feather="star"><vue-feather type="star"></vue-feather></i></div>
     <div class="onhover-show-div bookmark-flip" :class="bookmarkSearchBox ? 'active' : ''">
        <div class="flip-card">
            <div class="flip-card-inner" :class="bookmarkSearchBox ? 'flipped' : ''">
                <div class="front dropdown-title bookmark-div ">
                    <h6 class="f-18 mb-0">Bookmark</h6>
                    <ul class="bookmark-dropdown pb-0">
                        <li class="custom-scrollbar">
                            <div class="row">
                                <div class="col-4 text-center" v-for="(menuItem, index) in bookmarkItems.slice(0, 8)"
                                    :key="index">
                                    <div class="bookmark-content">
                                        <div class="bookmark-ico">
                                            <vue-feather :type="menuItem.icon" class="bookmark-icons" ></vue-feather>
                                        </div>
                                        <h5 class="mt-2"> <router-link :to="{ path: menuItem.path }" class="realname">{{
                                            menuItem.title
                                        }}
                                        </router-link></h5>

                                    </div>
                                </div>
                            </div>
                        </li>
                        <li class="text-center">
                            <a class="flip-btn btn btn-primary" id="flip-btn" href="javascript:void(0)"
                                @click="openbookmark"> Add New
                                Bookmark</a>
                        </li>
                    </ul>
                </div>
                <div class="back dropdown-title bookmark-icon">
                    <ul>
                        <li>
                            <div class="bookmark-dropdown flip-back-content">
                                <input type="text" placeholder="search..." @keyup="searchTerm" v-model="terms" />
                            </div>
                            <div class="bookmark-search custom-scrollbar"
                                :class="!bookmarkSearchResultEmpty ? 'Typeahead-menu is-open' : 'Typeahead-menu'"
                                v-if="searchMenuItems.length">
                                <div class="ProfileCard u-cf" v-for="(menuItem, index) in searchMenuItems.slice(0, 8)"
                                    :key="index">
                                    <div class="ProfileCard-avatar header-search">
                                      <vue-feather :type="menuItem.icon"></vue-feather>
                                    </div>
                                    <div class="ProfileCard-details">
                                        <div class="ProfileCard-realName bookmark-realName">
                                            <span @click="removeFix()" style="float:left">
                                                <router-link :to="{ path: menuItem.path }" class="realname">{{
                                                    menuItem.title }}</router-link>
                                            </span>
                                            <span class="float-right "><a href="JavaScript:void(0);"
                                                    @click="addToBookmark(menuItem)"><i
                                                        class="fa fa-star-o f-18 bookmark-search f-right"
                                                        :class="menuItem.bookmark ? 'text-warning' : ''"></i></a></span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div :class="bookmarkSearchResultEmpty ? 'Typeahead-menu is-open' : 'Typeahead-menu'">
                                <div class="tt-dataset tt-dataset-0">
                                    <div class="EmptyMessage"> Your search turned up 0 results. Opps There are no result
                                        found.</div>
                                </div>
                            </div>
                        </li>
                        <li>
                            <a href="javascript:void(0)" class="d-block flip-back f-w-700" id="flip-back"
                                @click="openbookmark"> Back
                            </a>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import { mapState } from 'vuex';
export default {
  data() {
    return {
      filtered: false,
      terms: '',
      bookmarkSearchBox: false,
      bookmarkSearchResult: false,
      bookmarkSearchResultEmpty: false,
      bookmarkItems: [],
    }
  },
  computed: {
    ...mapState({
      menuItems: (state) => state.menu.data,
      searchMenuItems: (state) => state.menu.searchData,
    }),
  },
  watch: {
    searchMenuItems: function () {
      this.terms ? this.addFix() : this.removeFix();
      if (!this.searchMenuItems.length) this.bookmarkSearchResultEmpty = true;
      else this.bookmarkSearchResultEmpty = false;
    },
  },
  mounted() {
    this.menuItems.filter((items) => {
      if (items.bookmark) {
        this.bookmarkItems.push(items);
      }
    });
  },
  methods: {
    collapseFilter() {
      this.filtered = !this.filtered;
    },
    openbookmark() {
      this.bookmarkSearchBox = !this.bookmarkSearchBox;
      if (!this.bookmarkSearchBox) this.removeFix();
    },
    searchTerm: function () {
      this.addFix();
      this.$store.dispatch('menu/searchTerm', this.terms);
    },
    addFix() {
      this.bookmarkSearchResult = true;
    },
    removeFix() {
      this.bookmarkSearchResult = false;
      this.text = '';
    },
    addToBookmark(items) {
      const index = this.bookmarkItems.indexOf(items);
      if (index === -1 && !items.bookmark) {
        items.bookmark = true;
        this.bookmarkItems.push(items);
        this.text = '';
      } else {
        this.bookmarkItems.splice(index, 1);
        items.bookmark = false;
      }
    },
  }
}
</script>