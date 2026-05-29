import { defineStore } from 'pinia'

export const useBreadcrumbStore = defineStore('breadcrumb', {
  state: () => ({ title: '', crumbs: [] }),
  actions: {
    set(title, crumbs = []) {
      this.title = title
      this.crumbs = crumbs
    }
  }
})
