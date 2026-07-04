<template>
  <main class="otika-cms">
    <div class="otika-app">
      <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
        <div class="navbar-bg"></div>
        <nav class="navbar navbar-expand-lg main-navbar sticky">
          <div class="form-inline mr-auto">
            <ul class="navbar-nav mr-3">
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon collapse-btn" @click="toggleSidebar" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'">
                  <i class="icofont-navigation-menu"></i>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon fullscreen-btn" @click="loadAll" title="Refresh">
                  <i class="icofont-refresh"></i>
                </button>
              </li>
              <li>
                <form class="form-inline mr-auto" @submit.prevent="runGlobalSearch">
                  <div class="search-element">
                    <input v-model="globalSearch" class="form-control" type="search" placeholder="Search CMS" aria-label="Search" />
                    <button class="btn" type="submit"><i class="fas fa-search"></i></button>
                  </div>
                </form>
              </li>
            </ul>
          </div>
          <ul class="navbar-nav navbar-right">
            <li class="dropdown dropdown-list-toggle" :class="{ show: messagesOpen }">
              <button type="button" class="nav-link nav-link-lg cms-top-icon message-toggle" @click="messagesOpen = !messagesOpen; notificationsOpen = false; profileOpen = false">
                <i class="icofont-envelope"></i><span class="badge headerBadge1">{{ unreadMessageCount }}</span>
              </button>
              <div class="dropdown-menu dropdown-list dropdown-menu-right pullDown" :class="{ show: messagesOpen }">
                <div class="dropdown-header">Messages <div class="float-right"><button type="button" class="link-button" @click="markAllMessagesRead">Mark All As Read</button></div></div>
                <div class="dropdown-list-content dropdown-list-message">
                  <button v-for="item in topMessages" :key="item.id" type="button" class="dropdown-item" @click="openTopMessage(item)">
                    <span class="dropdown-item-avatar text-white"><img :src="item.avatar" alt="" class="rounded-circle" /></span>
                    <span class="dropdown-item-desc"><span class="message-user">{{ item.name }}</span><span class="time messege-text">{{ item.text }}</span><span class="time">{{ item.time }}</span></span>
                  </button>
                  <p v-if="!topMessages.length" class="cms-empty">No messages yet.</p>
                </div>
                <div class="dropdown-footer text-center"><button type="button" class="link-button" @click="active = 'messages'; messagesOpen = false">View All <i class="fas fa-chevron-right"></i></button></div>
              </div>
            </li>
            <li class="dropdown dropdown-list-toggle" :class="{ show: notificationsOpen }">
              <button type="button" class="nav-link notification-toggle nav-link-lg cms-top-icon" @click="notificationsOpen = !notificationsOpen; messagesOpen = false; profileOpen = false">
                <i class="icofont-notification"></i><span class="badge headerBadge2">{{ unreadNotificationCount }}</span>
              </button>
              <div class="dropdown-menu dropdown-list dropdown-menu-right pullDown" :class="{ show: notificationsOpen }">
                <div class="dropdown-header">Notifications <div class="float-right"><button type="button" class="link-button" @click="markAllNotifications">Mark All As Read</button></div></div>
                <div class="dropdown-list-content dropdown-list-icons">
                  <button v-for="item in notificationItems" :key="item.id" type="button" class="dropdown-item" :class="{ 'dropdown-item-unread': item.unread }" @click="openNotification(item)">
                    <span class="dropdown-item-icon text-white" :class="item.color"><i :class="item.icon"></i></span>
                    <span class="dropdown-item-desc">{{ item.text }} <span class="time">{{ item.time }}</span></span>
                  </button>
                  <p v-if="!notificationItems.length" class="cms-empty">No notifications yet.</p>
                </div>
                <div class="dropdown-footer text-center"><button type="button" class="link-button" @click="active = 'notifications'; notificationsOpen = false">View All <i class="fas fa-chevron-right"></i></button></div>
              </div>
            </li>
            <li class="dropdown" :class="{ show: profileOpen }">
              <button type="button" class="nav-link dropdown-toggle nav-link-lg nav-link-user cms-top-icon" @click="profileOpen = !profileOpen; notificationsOpen = false; messagesOpen = false">
                <img v-if="currentUserAvatar" :src="currentUserAvatar" alt="" class="user-img-radious-style" />
                <img v-else alt="" src="/otika-assets/img/user.png" class="user-img-radious-style" />
                <span class="d-sm-none d-lg-inline-block"></span>
              </button>
              <div class="dropdown-menu dropdown-menu-right pullDown" :class="{ show: profileOpen }">
                <div class="dropdown-title">Hello {{ currentUserName }}</div>
                <button type="button" class="dropdown-item has-icon" @click="active = 'my-profile'; profileOpen = false"><i class="far fa-user"></i> Profile</button>
                <button type="button" class="dropdown-item has-icon" @click="active = 'audit-logs'; profileOpen = false"><i class="fas fa-bolt"></i> Activities</button>
                <button type="button" class="dropdown-item has-icon" @click="active = 'settings'; profileOpen = false"><i class="fas fa-cog"></i> Site Settings</button>
                <div class="dropdown-divider"></div>
                <button type="button" class="dropdown-item has-icon text-danger" @click="logout"><i class="fas fa-sign-out-alt"></i> Logout</button>
              </div>
            </li>
          </ul>
        </nav>

        <div class="main-sidebar sidebar-style-2">
          <aside id="sidebar-wrapper">
            <div class="sidebar-brand">
              <router-link to="/cms">
                <img alt="NCS" src="/main-logo.png" class="header-logo" />
              </router-link>
            </div>
            <ul class="sidebar-menu">
              <li class="menu-header">Main</li>
              <li v-for="item in visibleItems(topSections)" :key="item.id" :class="{ active: active === item.id }">
                <button type="button" class="nav-link" @click="selectSection(item.id)"><i :class="item.icon"></i><span>{{ item.label }}</span></button>
              </li>
              <li class="menu-header">Website</li>
              <li v-if="visibleItems(homepageSections).length" class="dropdown" :class="{ active: homepageGroupOpen || groupHasActive(homepageSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="homepageGroupOpen = !homepageGroupOpen"><i class="icofont-home"></i><span>Homepage Management</span></button>
                <ul class="dropdown-menu" :style="{ display: homepageGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(homepageSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(slideshowSections).length" class="dropdown" :class="{ active: slideshowGroupOpen || groupHasActive(slideshowSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="slideshowGroupOpen = !slideshowGroupOpen"><i class="icofont-image"></i><span>Slideshow Manager</span></button>
                <ul class="dropdown-menu" :style="{ display: slideshowGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(slideshowSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(blogSections).length" class="dropdown" :class="{ active: blogsGroupOpen || groupHasActive(blogSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="blogsGroupOpen = !blogsGroupOpen"><i class="icofont-newspaper"></i><span>Blogs Management</span></button>
                <ul class="dropdown-menu" :style="{ display: blogsGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(blogSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(staticPageSections).length" class="dropdown" :class="{ active: staticPagesGroupOpen || groupHasActive(staticPageSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="staticPagesGroupOpen = !staticPagesGroupOpen"><i class="icofont-page"></i><span>Static Pages</span></button>
                <ul class="dropdown-menu" :style="{ display: staticPagesGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(staticPageSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(projectSections).length" class="dropdown" :class="{ active: projectsGroupOpen || groupHasActive(projectSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="projectsGroupOpen = !projectsGroupOpen"><i class="icofont-briefcase"></i><span>Projects</span></button>
                <ul class="dropdown-menu" :style="{ display: projectsGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(projectSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(caseStudySections).length" class="dropdown" :class="{ active: caseStudiesGroupOpen || groupHasActive(caseStudySections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="caseStudiesGroupOpen = !caseStudiesGroupOpen"><i class="icofont-read-book"></i><span>Case Studies</span></button>
                <ul class="dropdown-menu" :style="{ display: caseStudiesGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(caseStudySections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(faqSections).length" class="dropdown" :class="{ active: faqsGroupOpen || groupHasActive(faqSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="faqsGroupOpen = !faqsGroupOpen"><i class="icofont-question-circle"></i><span>FAQs</span></button>
                <ul class="dropdown-menu" :style="{ display: faqsGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(faqSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(resourceSections).length" class="dropdown" :class="{ active: resourcesGroupOpen || groupHasActive(resourceSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="resourcesGroupOpen = !resourcesGroupOpen"><i class="icofont-download"></i><span>Resource Centre</span></button>
                <ul class="dropdown-menu" :style="{ display: resourcesGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(resourceSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(careerSections).length" class="dropdown" :class="{ active: careersGroupOpen || groupHasActive(careerSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="careersGroupOpen = !careersGroupOpen"><i class="icofont-businessman"></i><span>Careers</span></button>
                <ul class="dropdown-menu" :style="{ display: careersGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(careerSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(teamSections).length" class="dropdown" :class="{ active: teamGroupOpen || groupHasActive(teamSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="teamGroupOpen = !teamGroupOpen"><i class="icofont-users-alt-5"></i><span>Team Members</span></button>
                <ul class="dropdown-menu" :style="{ display: teamGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(teamSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(roleSections).length" class="dropdown" :class="{ active: rolesGroupOpen || groupHasActive(roleSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="rolesGroupOpen = !rolesGroupOpen"><i class="icofont-key"></i><span>Roles</span></button>
                <ul class="dropdown-menu" :style="{ display: rolesGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(roleSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(userSections).length" class="dropdown" :class="{ active: usersGroupOpen || groupHasActive(userSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="usersGroupOpen = !usersGroupOpen"><i class="icofont-users-social"></i><span>Users</span></button>
                <ul class="dropdown-menu" :style="{ display: usersGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(userSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(facilitySections).length" class="dropdown" :class="{ active: facilitiesGroupOpen || groupHasActive(facilitySections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="facilitiesGroupOpen = !facilitiesGroupOpen"><i class="icofont-building-alt"></i><span>Facilities</span></button>
                <ul class="dropdown-menu" :style="{ display: facilitiesGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(facilitySections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(eventSections).length" class="dropdown" :class="{ active: eventsGroupOpen || groupHasActive(eventSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="eventsGroupOpen = !eventsGroupOpen"><i class="icofont-calendar"></i><span>Events</span></button>
                <ul class="dropdown-menu" :style="{ display: eventsGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(eventSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(investSections).length" class="dropdown" :class="{ active: investGroupOpen || groupHasActive(investSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="investGroupOpen = !investGroupOpen"><i class="icofont-money-bag"></i><span>Invest With Us</span></button>
                <ul class="dropdown-menu" :style="{ display: investGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(investSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(federationSections).length" class="dropdown" :class="{ active: federationsGroupOpen || groupHasActive(federationSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="federationsGroupOpen = !federationsGroupOpen"><i class="icofont-trophy"></i><span>Federations</span></button>
                <ul class="dropdown-menu" :style="{ display: federationsGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(federationSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(funFactSections).length" class="dropdown" :class="{ active: funFactsGroupOpen || groupHasActive(funFactSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="funFactsGroupOpen = !funFactsGroupOpen"><i class="icofont-chart-bar-graph"></i><span>Fun Facts</span></button>
                <ul class="dropdown-menu" :style="{ display: funFactsGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(funFactSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li v-if="visibleItems(newsletterSections).length" class="dropdown" :class="{ active: newsletterGroupOpen || groupHasActive(newsletterSections) }">
                <button type="button" class="menu-toggle nav-link has-dropdown" @click="newsletterGroupOpen = !newsletterGroupOpen"><i class="icofont-envelope"></i><span>Newsletter</span></button>
                <ul class="dropdown-menu" :style="{ display: newsletterGroupOpen ? 'block' : 'none' }">
                  <li v-for="item in visibleItems(newsletterSections)" :key="item.id" :class="{ active: active === item.id }"><button type="button" class="nav-link" @click="selectSection(item.id)">{{ item.label }}</button></li>
                </ul>
              </li>
              <li class="menu-header">Modules</li>
              <li v-for="item in visibleItems(contentSections)" :key="item.id" :class="{ active: active === item.id }">
                <button type="button" class="nav-link" @click="selectSection(item.id)"><i :class="item.icon"></i><span>{{ item.label }}</span></button>
              </li>
            </ul>
          </aside>
        </div>

        <div class="main-content">
          <section class="section">
            <div class="section-body">
              <div class="section-header">
                <h1>{{ currentSection.label }}</h1>
                <div class="section-header-breadcrumb">
                  <div class="breadcrumb-item active"><router-link to="/cms">CMS</router-link></div>
                  <div class="breadcrumb-item">{{ currentSection.label }}</div>
                </div>
              </div>
              <div class="cms-actions otika-page-actions">
                <ThemeToggle />
                <router-link to="/" target="_blank" class="btn btn-icon icon-left btn-primary"><i class="fas fa-external-link-alt"></i> Preview website</router-link>
                <button type="button" class="btn btn-icon icon-left btn-info" @click="loadAll"><i class="fas fa-sync"></i> Refresh</button>
              </div>

        <p v-if="message" class="cms-message">{{ message }}</p>
        <p v-if="error" class="cms-error">{{ error }}</p>

        <section v-if="active === 'overview'" class="otika-dashboard">
          <div class="row">
            <div v-for="card in dashboardStatCards" :key="card.label" class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
              <div class="card">
                <div class="card-statistic-4">
                  <div class="align-items-center justify-content-between">
                    <div class="row">
                      <div class="col-lg-6 col-md-6 col-sm-6 col-xs-6 pr-0 pt-3">
                        <div class="card-content">
                          <h5 class="font-15">{{ card.label }}</h5>
                          <h2 class="mb-3 font-18">{{ card.value }}</h2>
                          <p class="mb-0"><span :class="card.trendClass">{{ card.trend }}</span> {{ card.caption }}</p>
                        </div>
                      </div>
                      <div class="col-lg-6 col-md-6 col-sm-6 col-xs-6 pl-0">
                        <div class="banner-img">
                          <i :class="card.icon"></i>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-12 col-sm-12 col-lg-8">
              <div class="card">
                <div class="card-header">
                  <h4>Page Visit Statistics</h4>
                  <div class="card-header-action">
                    <span class="badge badge-primary">Last 7 days</span>
                  </div>
                </div>
                <div class="card-body">
                  <div class="cms-chart-bars">
                    <div v-for="point in visitSeries" :key="point.label" class="cms-chart-bar">
                      <span class="cms-chart-bar__value">{{ point.value }}</span>
                      <div class="cms-chart-bar__track">
                        <span :style="{ height: `${point.percent}%` }"></span>
                      </div>
                      <strong>{{ point.label }}</strong>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-12 col-sm-12 col-lg-4">
              <div class="card">
                <div class="card-header">
                  <h4>Website KPI Summary</h4>
                </div>
                <div class="card-body">
                  <div v-for="item in websiteKpis" :key="item.label" class="cms-kpi-row">
                    <span>{{ item.label }}</span>
                    <strong>{{ item.value }}</strong>
                    <sup :class="item.trendClass">{{ item.trend }}</sup>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-12 col-sm-12 col-lg-4">
              <div class="card">
                <div class="card-header">
                  <h4>Platforms</h4>
                </div>
                <div class="card-body">
                  <div v-for="item in platformStats" :key="item.label" class="cms-progress-item">
                    <div><span>{{ item.label }}</span><strong>{{ item.value }}%</strong></div>
                    <div class="progress" data-height="6">
                      <div class="progress-bar" :class="item.color" :style="{ width: `${item.value}%` }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-12 col-sm-12 col-lg-4">
              <div class="card">
                <div class="card-header">
                  <h4>Traffic Sources</h4>
                </div>
                <div class="card-body">
                  <div v-for="item in trafficSources" :key="item.label" class="cms-source-row">
                    <span><i :class="item.icon"></i> {{ item.label }}</span>
                    <strong>{{ item.value }}%</strong>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-12 col-sm-12 col-lg-4">
              <div class="card">
                <div class="card-header">
                  <h4>Content Analytics</h4>
                </div>
                <div class="card-body">
                  <div class="summary">
                    <div class="summary-chart active">
                      <div class="cms-donut" :style="{ '--first': `${contentMix[0].value}%`, '--second': `${contentMix[0].value + contentMix[1].value}%` }">
                        <span>{{ totalContentItems }}</span>
                        <small>items</small>
                      </div>
                    </div>
                    <div class="cms-donut-legend">
                      <span v-for="item in contentMix" :key="item.label"><i :class="item.color"></i>{{ item.label }} {{ item.value }}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-12">
              <div class="card">
                <div class="card-header">
                  <h4>Website Operations Table</h4>
                  <div class="card-header-form">
                    <div class="input-group">
                      <input v-model="globalSearch" type="text" class="form-control" placeholder="Search KPI">
                      <div class="input-group-btn">
                        <button class="btn btn-primary" type="button" @click="runGlobalSearch"><i class="fas fa-search"></i></button>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="card-body p-0">
                  <div class="table-responsive">
                    <table class="table table-striped table-hover table-sm mb-0">
                      <thead>
                        <tr>
                          <th>Task Name</th>
                          <th>Owner</th>
                          <th>Task Status</th>
                          <th>Due Date</th>
                          <th>Priority</th>
                          <th>Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="task in dashboardTasks" :key="task.name">
                          <td>{{ task.name }}</td>
                          <td class="text-truncate">
                            <ul class="list-unstyled order-list m-b-0">
                              <li v-for="avatar in task.avatars" :key="avatar" class="team-member team-member-sm">
                                <img class="rounded-circle" :src="avatar" alt="user">
                              </li>
                              <li class="avatar avatar-sm"><span class="badge badge-primary">+{{ task.extra }}</span></li>
                            </ul>
                          </td>
                          <td class="align-middle">
                            <div class="progress-text">{{ task.progress }}%</div>
                            <div class="progress" data-height="6">
                              <div class="progress-bar" :class="task.barClass" :style="{ width: `${task.progress}%` }"></div>
                            </div>
                          </td>
                          <td>{{ task.due }}</td>
                          <td><div class="badge" :class="task.priorityClass">{{ task.priority }}</div></td>
                          <td><button type="button" class="btn btn-outline-primary" @click="selectSection(task.target)">Detail</button></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-md-6 col-lg-12 col-xl-6">
              <div class="card">
                <div class="card-header">
                  <h4>Website Activity</h4>
                  <form class="card-header-form" @submit.prevent="runGlobalSearch">
                    <input v-model="globalSearch" type="text" class="form-control" placeholder="Search">
                  </form>
                </div>
                <div class="card-body">
                  <div v-for="item in websiteActivity" :key="item.id" class="support-ticket media pb-1 mb-3">
                    <img :src="item.avatar" class="user-img mr-2" alt="">
                    <div class="media-body ml-3">
                      <div class="badge badge-pill mb-1 float-right" :class="item.badgeClass">{{ item.type }}</div>
                      <span class="font-weight-bold">#{{ item.id }}</span>
                      <button type="button" class="cms-inline-link" @click="selectSection(item.target)">{{ item.title }}</button>
                      <p class="my-1">{{ item.summary }}</p>
                      <small class="text-muted">Created by <span class="font-weight-bold font-13">{{ item.owner }}</span> - {{ item.time }}</small>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-6 col-lg-12 col-xl-6">
              <div class="card">
                <div class="card-header">
                  <h4>Analytics Snapshot</h4>
                </div>
                <div class="card-body">
                  <div class="cms-analytics-snapshot">
                    <article v-for="item in analyticsSnapshot" :key="item.label">
                      <i :class="item.icon"></i>
                      <div>
                        <span>{{ item.label }}</span>
                        <strong>{{ item.value }}</strong>
                      </div>
                    </article>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section v-else-if="active === 'analytics'" class="otika-dashboard">
          <div class="row">
            <div v-for="card in dashboardStatCards" :key="`analytics-${card.label}`" class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
              <div class="card">
                <div class="card-statistic-4">
                  <div class="row">
                    <div class="col-7">
                      <div class="card-content">
                        <h5 class="font-15">{{ card.label }}</h5>
                        <h2 class="mb-3 font-18">{{ card.value }}</h2>
                        <p class="mb-0"><span :class="card.trendClass">{{ card.trend }}</span> {{ card.caption }}</p>
                      </div>
                    </div>
                    <div class="col-5">
                      <i :class="card.icon"></i>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-12 col-lg-8">
              <div class="card">
                <div class="card-header">
                  <h4>Traffic Timeline</h4>
                  <div class="card-header-action">
                    <button type="button" class="btn btn-primary btn-sm" @click="loadAnalytics"><i class="fas fa-sync"></i> Refresh</button>
                  </div>
                </div>
                <div class="card-body">
                  <div class="cms-chart-bars">
                    <div v-for="point in visitSeries" :key="`analytics-${point.label}`" class="cms-chart-bar">
                      <span class="cms-chart-bar__value">{{ point.value }}</span>
                      <div class="cms-chart-bar__track">
                        <span :style="{ height: `${point.percent}%` }"></span>
                      </div>
                      <strong>{{ point.label }}</strong>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-12 col-lg-4">
              <div class="card">
                <div class="card-header"><h4>Acquisition Channels</h4></div>
                <div class="card-body">
                  <div v-for="item in trafficSources" :key="`channel-${item.label}`" class="cms-source-row">
                    <span><i :class="item.icon"></i> {{ item.label }}</span>
                    <strong>{{ item.value }}%</strong>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="row">
            <div class="col-12 col-lg-4">
              <div class="card">
                <div class="card-header"><h4>Device Types</h4></div>
                <div class="card-body">
                  <div v-for="item in platformStats" :key="`device-${item.label}`" class="cms-progress-item">
                    <div><span>{{ item.label }}</span><strong>{{ item.value }}%</strong></div>
                    <div class="progress" data-height="6">
                      <div class="progress-bar" :class="item.color" :style="{ width: `${item.value}%` }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-12 col-lg-4">
              <div class="card">
                <div class="card-header"><h4>Top Cities</h4></div>
                <div class="card-body p-0">
                  <table class="table table-striped table-sm mb-0">
                    <tbody>
                      <tr v-for="row in analyticsRows('top_cities')" :key="`city-${row.label}`">
                        <td>{{ row.label }}</td>
                        <td class="text-right">{{ row.value }}</td>
                        <td class="text-right">{{ row.share }}%</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            <div class="col-12 col-lg-4">
              <div class="card">
                <div class="card-header"><h4>Top Pages</h4></div>
                <div class="card-body p-0">
                  <table class="table table-striped table-sm mb-0">
                    <tbody>
                      <tr v-for="row in analyticsRows('top_pages')" :key="`page-${row.label}`">
                        <td>{{ row.label }}</td>
                        <td class="text-right">{{ row.value }}</td>
                        <td class="text-right">{{ row.share }}%</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section v-else-if="active === 'homepage'" class="cms-panel">
          <div class="cms-panel-head"><h2>Homepage Sections</h2><button @click="saveHomepage">Save homepage</button></div>
          <div class="cms-two">
            <label>Facilities eyebrow<input v-model="homepage.facilities.eyebrow" /></label>
            <label>Facilities title<input v-model="homepage.facilities.title" /></label>
            <label class="wide">Facilities intro<textarea v-model="homepage.facilities.intro"></textarea></label>
            <label>Events eyebrow<input v-model="homepage.events.eyebrow" /></label>
            <label>Events title<input v-model="homepage.events.title" /></label>
            <label class="wide">Events intro<textarea v-model="homepage.events.intro"></textarea></label>
          </div>
        </section>

        <section v-else-if="active === 'homepage-about'" class="cms-panel">
          <div class="cms-panel-head"><h2>About NCS Section</h2><button @click="saveHomepage">Save about section</button></div>
          <div class="cms-two">
            <label>Eyebrow<input v-model="homepage.about.eyebrow" /></label>
            <label>Title<input v-model="homepage.about.title" /></label>
            <label class="wide">Intro<textarea v-model="homepage.about.intro"></textarea></label>
            <label class="wide">Body<textarea v-model="homepage.about.body"></textarea></label>
            <label>Leadership button label<input v-model="homepage.about.leadership_label" /></label>
            <label>Leadership button URL<input v-model="homepage.about.leadership_url" /></label>
            <label>Core functions title<input v-model="homepage.about.core_title" /></label>
            <label>Mandate button label<input v-model="homepage.about.mandate_label" /></label>
            <label class="wide">Core functions intro<textarea v-model="homepage.about.core_intro"></textarea></label>
            <label>Mandate button URL<input v-model="homepage.about.mandate_url" /></label>
            <label>Chairperson name<input v-model="homepage.leadership.chairperson_name" /></label>
            <label>General Secretary name<input v-model="homepage.leadership.secretary_name" /></label>
            <label class="wide">Chairperson image URL<input v-model="homepage.leadership.chairperson_image" /></label>
          </div>
          <article class="cms-subpanel">
            <div class="cms-panel-head"><h3>About Value Cards</h3><button type="button" @click="addHomepageValue">Add value card</button></div>
            <div class="homepage-card-list">
              <div v-for="(item, index) in homepage.values" :key="index" class="homepage-dynamic-card">
                <div class="cms-two">
                  <label>Title<input v-model="item.title" /></label>
                  <label>Icon class<input v-model="item.icon" placeholder="icofont-dart" /></label>
                  <label class="wide">Text<textarea v-model="item.text"></textarea></label>
                  <label class="cms-check"><input v-model="item.featured" type="checkbox" /> Featured navy card</label>
                </div>
                <button type="button" class="btn btn-sm btn-danger" @click="homepage.values.splice(index, 1)">Remove</button>
              </div>
            </div>
          </article>
        </section>

        <section v-else-if="active === 'homepage-topbar'" class="cms-panel">
          <div class="cms-panel-head"><h2>Topbar Marquee & Social Links</h2><button @click="saveHomepage">Save topbar</button></div>
          <article class="cms-subpanel">
            <div class="cms-panel-head"><h3>Marquee Messages</h3><button type="button" @click="homepage.topbar.marquee.push('')">Add message</button></div>
            <div class="cms-list-editor">
              <div v-for="(_, index) in homepage.topbar.marquee" :key="index" class="cms-row">
                <input v-model="homepage.topbar.marquee[index]" class="form-control" placeholder="Topbar announcement" />
                <button type="button" @click="homepage.topbar.marquee.splice(index, 1)">Remove</button>
              </div>
            </div>
          </article>
          <div class="cms-two">
            <label>Facebook URL<input v-model="homepage.topbar.social.facebook" /></label>
            <label>X / Twitter URL<input v-model="homepage.topbar.social.twitter" /></label>
            <label>LinkedIn URL<input v-model="homepage.topbar.social.linkedin" /></label>
            <label>YouTube URL<input v-model="homepage.topbar.social.youtube" /></label>
            <label>Instagram URL<input v-model="homepage.topbar.social.instagram" /></label>
            <label>Webmail URL<input v-model="homepage.topbar.webmail_url" /></label>
          </div>
        </section>

        <section v-else-if="active === 'core'" class="cms-panel">
          <div class="cms-panel-head"><h2>Core Functions</h2><button @click="saveHomepage">Save functions</button></div>
          <div class="cms-list-editor">
            <div v-for="(_, index) in homepage.core_functions" :key="index" class="cms-row">
              <textarea v-model="homepage.core_functions[index]"></textarea>
              <button type="button" @click="homepage.core_functions.splice(index, 1)">Remove</button>
            </div>
            <button type="button" @click="homepage.core_functions.push('')">Add function</button>
          </div>
        </section>

        <section v-else-if="active === 'homepage-sports-excellence'" class="cms-panel">
          <div class="cms-panel-head"><h2>Sports Excellence Counters</h2><button @click="saveHomepage">Save counters</button></div>
          <div class="cms-two">
            <label>Section title<input v-model="homepage.stats_title" /></label>
            <label>Section intro<input v-model="homepage.stats_intro" /></label>
          </div>
          <div class="homepage-card-list">
            <div v-for="(counter, index) in homepage.milestones" :key="index" class="homepage-dynamic-card">
              <div class="cms-panel-head"><h3>Counter {{ index + 1 }}</h3><button type="button" @click="homepage.milestones.splice(index, 1)">Remove</button></div>
              <div class="cms-two">
                <label>Number / value<input v-model="counter.value" placeholder="60+" /></label>
                <label>Title<input v-model="counter.label" /></label>
                <label>Icon class<input v-model="counter.icon" placeholder="icofont-award" /></label>
                <label class="wide">Description<textarea v-model="counter.description"></textarea></label>
              </div>
            </div>
          </div>
          <button type="button" class="btn btn-icon icon-left btn-primary" @click="addSportsCounter"><i class="fas fa-plus"></i> Add counter</button>
        </section>


        <section v-else-if="active === 'slideshow-manager'" class="cms-panel">
          <SlideshowManager :slideshow="slideshow" :slides="slides" @refresh="loadAll" @message="setMsg" @error="setErr" />
        </section>

        <section v-else-if="active === 'create-post'" class="cms-panel">
          <BlogPostEditor :model="postForm" :categories="blogCategories" @save="savePost" />
        </section>

        <section v-else-if="active === 'manage-posts'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Posts</h2><button type="button" @click="resetPostForm(); active = 'create-post'">New post</button></div>
          <ContentTable :items="posts" title-key="title" subtitle-key="status" @edit="editPost" @delete="removePost" />
        </section>

        <section v-else-if="active === 'static-pages'" class="cms-panel">
          <StaticPageBuilder :model="pageForm" @save="savePage" />
        </section>

        <section v-else-if="active === 'manage-pages'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Static Pages</h2><button type="button" @click="resetPageForm(); active = 'static-pages'">New page</button></div>
          <ContentTable :items="pages" title-key="title" subtitle-key="slug" @edit="editPage" @delete="removePost" />
        </section>

        <section v-else-if="active === 'projects'" class="cms-panel">
          <BlogPostEditor :model="projectForm" :categories="projectCategories" category-field="category_tag" @save="saveProject" />
        </section>

        <section v-else-if="active === 'manage-projects'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Projects</h2><button type="button" @click="resetProjectForm(); active = 'projects'">New project</button></div>
          <ContentTable :items="projects" title-key="title" subtitle-key="status" @edit="editProject" @delete="removePost" />
        </section>

        <section v-else-if="active === 'create-project-categories'" class="cms-panel">
          <EditorForm title="Create Project Categories" :model="projectCategoryForm" :fields="projectCategoryFields" @save="saveProjectCategory" />
        </section>

        <section v-else-if="active === 'manage-project-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Project Categories</h2><button type="button" @click="resetProjectCategoryForm(); active = 'create-project-categories'">New category</button></div>
          <ContentTable :items="projectCategories" title-key="name" subtitle-key="slug" @edit="editProjectCategory" @delete="removeProjectCategory" />
        </section>

        <section v-else-if="active === 'case-studies'" class="cms-panel">
          <BlogPostEditor :model="caseStudyForm" :categories="caseStudyCategories" category-field="category_tag" @save="saveCaseStudy" />
        </section>

        <section v-else-if="active === 'manage-case-studies'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Case Study Posts</h2><button type="button" @click="resetCaseStudyForm(); active = 'case-studies'">New case study</button></div>
          <ContentTable :items="caseStudies" title-key="title" subtitle-key="status" @edit="editCaseStudy" @delete="removePost" />
        </section>

        <section v-else-if="active === 'create-case-study-categories'" class="cms-panel">
          <EditorForm title="Create Case Study Categories" :model="caseStudyCategoryForm" :fields="caseStudyCategoryFields" @save="saveCaseStudyCategory" />
        </section>

        <section v-else-if="active === 'manage-case-study-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Case Study Categories</h2><button type="button" @click="resetCaseStudyCategoryForm(); active = 'create-case-study-categories'">New category</button></div>
          <ContentTable :items="caseStudyCategories" title-key="name" subtitle-key="slug" @edit="editCaseStudyCategory" @delete="removeCaseStudyCategory" />
        </section>

        <section v-else-if="active === 'create-blog-categories'" class="cms-panel">
          <EditorForm title="Create Blog Categories" :model="blogCategoryForm" :fields="blogCategoryFields" @save="saveBlogCategory" />
        </section>

        <section v-else-if="active === 'manage-blog-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Categories</h2><button type="button" @click="resetBlogCategoryForm(); active = 'create-blog-categories'">New category</button></div>
          <ContentTable :items="blogCategories" title-key="name" subtitle-key="slug" @edit="editBlogCategory" @delete="removeBlogCategory" />
        </section>

        <section v-else-if="active === 'comments'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Comment Moderation Queue</h2>
            <select v-model="commentStatus" @change="loadComments">
              <option value="">All comments</option>
              <option value="pending">Pending</option>
              <option value="approved">Approved</option>
              <option value="flagged">Flagged</option>
            </select>
          </div>
          <div class="cms-table">
            <article v-for="comment in comments" :key="comment.id" class="cms-table-row comment-row">
              <div>
                <strong>{{ comment.user_name }} on {{ comment.post_title }}</strong>
                <span>{{ comment.status }} · {{ comment.body }}</span>
              </div>
              <div>
                <button type="button" @click="moderateComment(comment, 'approve')">Approve</button>
                <button type="button" @click="moderateComment(comment, 'flag')">Flag</button>
                <button type="button" @click="deleteComment(comment)">Delete</button>
              </div>
            </article>
            <p v-if="!comments.length" class="cms-empty">No comments in this queue.</p>
          </div>
        </section>

        <section v-else-if="active === 'events'" class="cms-panel">
          <form class="otika-form-card event-form-card" @submit.prevent="saveEvent">
            <div class="card">
              <div class="card-header"><h4>{{ eventForm.id ? 'Edit Event' : 'Add New Event' }}</h4></div>
              <div class="card-body">
                <div class="section-title mt-0">Event Details</div>
                <div class="row">
                  <div class="form-group col-lg-8"><label>Title</label><input v-model="eventForm.title" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Status</label><select v-model="eventForm.status" class="form-control selectric"><option value="draft">Draft</option><option value="published">Published</option><option value="archived">Archived</option></select></div>
                  <div class="form-group col-lg-6"><label>Slug</label><input v-model="eventForm.slug" class="form-control" placeholder="event-url-slug" /></div>
                  <div class="form-group col-lg-6"><label>Category</label><select v-model="eventForm.category" class="form-control selectric"><option value="">Select category</option><option v-for="category in eventCategories" :key="category.slug || category.id" :value="category.slug || category.name">{{ category.name }}</option></select></div>
                  <div class="form-group col-lg-6"><label>Location</label><input v-model="eventForm.location" class="form-control" /></div>
                  <div class="form-group col-lg-3"><label>Start Date</label><input v-model="eventForm.event_date" type="datetime-local" class="form-control" /></div>
                  <div class="form-group col-lg-3"><label>End Date</label><input v-model="eventForm.end_date" type="datetime-local" class="form-control" /></div>
                  <div class="form-group col-12">
                    <label>Cover Image</label>
                    <DropzoneUpload v-model="eventForm.cover_image_url" label="event cover image" accept="image/png,image/jpeg,image/webp" @error="setErr" />
                  </div>
                </div>
                <div class="section-title">Description</div>
                <CmsRichTextEditor v-model="eventForm.description" />
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary">{{ eventForm.id ? 'Update Event' : 'Create Event' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-events'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Events</h2><button type="button" @click="resetEventForm(); active = 'events'">New event</button></div>
          <ContentTable :items="events" title-key="title" subtitle-key="category" @edit="editEvent" @delete="removeEvent" />
        </section>

        <section v-else-if="active === 'create-event-categories'" class="cms-panel">
          <form class="otika-form-card event-category-form-card" @submit.prevent="saveEventCategory">
            <div class="card">
              <div class="card-header"><h4>{{ eventCategoryForm.id ? 'Edit Event Category' : 'Add Event Category' }}</h4></div>
              <div class="card-body">
                <div class="row">
                  <div class="form-group col-lg-6"><label>Name</label><input v-model="eventCategoryForm.name" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Slug</label><input v-model="eventCategoryForm.slug" class="form-control" /></div>
                  <div class="form-group col-lg-2"><label>Sort Order</label><input v-model.number="eventCategoryForm.sort_order" type="number" class="form-control" /></div>
                  <div class="form-group col-12"><label>Description</label><textarea v-model="eventCategoryForm.description" class="form-control otika-textarea"></textarea></div>
                  <div class="form-group col-12"><div class="custom-control custom-checkbox"><input id="event-category-active" v-model="eventCategoryForm.is_active" type="checkbox" class="custom-control-input" /><label class="custom-control-label" for="event-category-active">Active category</label></div></div>
                </div>
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary">{{ eventCategoryForm.id ? 'Update Category' : 'Create Category' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-event-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Event Categories</h2><button type="button" @click="resetEventCategoryForm(); active = 'create-event-categories'">New category</button></div>
          <ContentTable :items="eventCategories" title-key="name" subtitle-key="slug" @edit="editEventCategory" @delete="removeEventCategory" />
        </section>

        <section v-else-if="active === 'facilities'" class="cms-panel">
          <EditorForm title="Add Facility" :model="facilityForm" :fields="facilityFields" @save="saveFacility" />
        </section>

        <section v-else-if="active === 'manage-facilities'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Facilities</h2><button type="button" @click="resetFacilityForm(); active = 'facilities'">New facility</button></div>
          <ContentTable :items="facilities" title-key="name" subtitle-key="category" @edit="editFacility" @delete="removeFacility" />
        </section>

        <section v-else-if="active === 'create-facility-categories'" class="cms-panel">
          <EditorForm title="Add Facility Category" :model="facilityCategoryForm" :fields="facilityCategoryFields" @save="saveFacilityCategory" />
        </section>

        <section v-else-if="active === 'manage-facility-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Facility Categories</h2><button type="button" @click="resetFacilityCategoryForm(); active = 'create-facility-categories'">New category</button></div>
          <ContentTable :items="facilityCategories" title-key="name" subtitle-key="slug" @edit="editFacilityCategory" @delete="removeFacilityCategory" />
        </section>

        <section v-else-if="active === 'associations'" class="cms-panel">
          <EditorForm title="Add New Federation" :model="associationForm" :fields="associationFields" @save="saveAssociation" />
        </section>

        <section v-else-if="active === 'manage-federations'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Federations</h2><button type="button" @click="resetFederationForm(); active = 'associations'">New federation</button></div>
          <ContentTable :items="associations" title-key="name" subtitle-key="category" @edit="editAssociation" @delete="removeAssociation" />
        </section>

        <section v-else-if="active === 'create-federation-categories'" class="cms-panel">
          <EditorForm title="Create Federation Category" :model="federationCategoryForm" :fields="federationCategoryFields" @save="saveFederationCategory" />
        </section>

        <section v-else-if="active === 'manage-federation-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Federation Categories</h2><button type="button" @click="resetFederationCategoryForm(); active = 'create-federation-categories'">New category</button></div>
          <ContentTable :items="federationCategories" title-key="name" subtitle-key="slug" @edit="editFederationCategory" @delete="removeFederationCategory" />
        </section>

        <section v-else-if="active === 'facts'" class="cms-panel">
          <EditorForm title="Create Fun Fact Article" :model="factForm" :fields="factFields" @save="saveFact" />
        </section>

        <section v-else-if="active === 'manage-facts'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Fun Fact Articles</h2><button type="button" @click="resetFactForm(); active = 'facts'">New article</button></div>
          <ContentTable :items="facts" title-key="label" subtitle-key="value" @edit="editFact" @delete="removeFact" />
        </section>

        <section v-else-if="active === 'faqs'" class="cms-panel">
          <form class="otika-form-card faq-form-card" @submit.prevent="saveFAQ">
            <div class="card">
              <div class="card-header"><h4>{{ faqForm.id ? 'Edit FAQ Article' : 'Create FAQ Article' }}</h4></div>
              <div class="card-body">
                <div class="section-title mt-0">Question Details</div>
                <div class="row">
                  <div class="form-group col-lg-8"><label>Question</label><input v-model="faqForm.question" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Category</label><select v-model="faqForm.category" class="form-control selectric"><option value="General">General</option><option v-for="category in faqCategories" :key="category.slug || category.id" :value="category.name">{{ category.name }}</option></select></div>
                  <div class="form-group col-lg-3"><label>Sort Order</label><input v-model.number="faqForm.sort_order" type="number" class="form-control" /></div>
                  <div class="form-group col-lg-3"><label class="d-block">Visibility</label><div class="custom-control custom-checkbox mt-2"><input id="faq-active" v-model="faqForm.is_active" type="checkbox" class="custom-control-input" /><label class="custom-control-label" for="faq-active">Active article</label></div></div>
                </div>
                <div class="section-title">Answer</div>
                <CmsRichTextEditor v-model="faqForm.answer" />
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary">{{ faqForm.id ? 'Update FAQ Article' : 'Create FAQ Article' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-faqs'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage FAQ Articles</h2><button type="button" @click="resetFAQForm(); active = 'faqs'">New article</button></div>
          <ContentTable :items="faqs" title-key="question" subtitle-key="category" @edit="editFAQ" @delete="removeFAQ" />
        </section>

        <section v-else-if="active === 'create-faq-categories'" class="cms-panel">
          <form class="otika-form-card" @submit.prevent="saveFAQCategory">
            <div class="card">
              <div class="card-header"><h4>{{ faqCategoryForm.id ? 'Edit FAQ Category' : 'Create FAQ Category' }}</h4></div>
              <div class="card-body">
                <div class="row">
                  <div class="form-group col-lg-6"><label>Name</label><input v-model="faqCategoryForm.name" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Slug</label><input v-model="faqCategoryForm.slug" class="form-control" /></div>
                  <div class="form-group col-lg-2"><label>Sort Order</label><input v-model.number="faqCategoryForm.sort_order" type="number" class="form-control" /></div>
                  <div class="form-group col-12"><label>Description</label><textarea v-model="faqCategoryForm.description" class="form-control otika-textarea"></textarea></div>
                  <div class="form-group col-12"><div class="custom-control custom-checkbox"><input id="faq-category-active" v-model="faqCategoryForm.is_active" type="checkbox" class="custom-control-input" /><label class="custom-control-label" for="faq-category-active">Active category</label></div></div>
                </div>
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary">{{ faqCategoryForm.id ? 'Update Category' : 'Create Category' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-faq-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage FAQ Categories</h2><button type="button" @click="resetFAQCategoryForm(); active = 'create-faq-categories'">New category</button></div>
          <ContentTable :items="faqCategories" title-key="name" subtitle-key="slug" @edit="editFAQCategory" @delete="removeFAQCategory" />
        </section>

        <section v-else-if="active === 'messages'" class="cms-panel message-center">
          <div class="cms-panel-head">
            <h2>Messages</h2>
            <div class="cms-actions-inline">
              <select v-model="messageStatus" class="form-control" @change="loadContactMessages"><option value="">All</option><option value="unread">Unread</option><option value="read">Read</option></select>
              <button type="button" @click="loadContactMessages">Refresh</button>
              <button type="button" @click="clearAllMessages">Clear All</button>
            </div>
          </div>
          <div class="notification-list">
            <article v-for="item in contactMessages" :key="item.id" class="notification-item" :class="{ unread: item.status === 'unread' }">
              <span class="notification-dot"></span>
              <div>
                <strong>{{ item.subject || 'Contact Us Message' }}</strong>
                <span>{{ item.name }} · {{ item.email }} · {{ formatDateTime(item.created_at) }}</span>
                <p>{{ item.message }}</p>
              </div>
              <div class="cms-actions-inline">
                <button type="button" @click="selectedMessage = item; markMessage(item, 'read')">View</button>
                <button type="button" @click="markMessage(item, item.status === 'read' ? 'unread' : 'read')">{{ item.status === 'read' ? 'Unread' : 'Read' }}</button>
                <button type="button" @click="removeMessage(item)">Delete</button>
              </div>
            </article>
            <p v-if="!contactMessages.length" class="cms-empty">No messages found.</p>
          </div>
          <article v-if="selectedMessage" class="cms-subpanel">
            <div class="cms-panel-head"><h2>{{ selectedMessage.subject || 'Contact Us Message' }}</h2><button type="button" @click="selectedMessage = null">Close</button></div>
            <p><strong>From:</strong> {{ selectedMessage.name }} &lt;{{ selectedMessage.email }}&gt;</p>
            <p><strong>Received:</strong> {{ formatDateTime(selectedMessage.created_at) }}</p>
            <p>{{ selectedMessage.message }}</p>
          </article>
        </section>

        <section v-else-if="active === 'notifications'" class="cms-panel notification-center">
          <div class="cms-panel-head">
            <h2>Notifications</h2>
            <div class="cms-actions-inline">
              <select v-model="notificationStatus" class="form-control" @change="loadNotifications"><option value="">All</option><option value="unread">Unread</option><option value="read">Read</option></select>
              <button type="button" @click="markAllNotifications">Mark All as Read</button>
              <button type="button" @click="clearAllNotifications">Clear All</button>
            </div>
          </div>
          <div class="notification-list">
            <article v-for="item in cmsNotifications" :key="item.id" class="notification-item" :class="{ unread: item.status === 'unread' }">
              <span class="notification-icon" :class="notificationThemeClass(item)"><i :class="notificationIcon(item)"></i></span>
              <div>
                <strong>{{ item.title }}</strong>
                <span>{{ item.type }} · {{ formatDateTime(item.created_at) }}</span>
                <p>{{ item.message }}</p>
              </div>
              <div class="cms-actions-inline">
                <button type="button" @click="toggleNotificationRead(item)">{{ item.status === 'read' ? 'Unread' : 'Read' }}</button>
                <button type="button" @click="dismissNotification(item)">Dismiss</button>
              </div>
            </article>
            <p v-if="!cmsNotifications.length" class="cms-empty">No new notifications.</p>
          </div>
        </section>

        <section v-else-if="active === 'resources'" class="cms-panel">
          <EditorForm title="Create Resource Centre Article" :model="resourceForm" :fields="resourceFields" @save="saveResource" @error="setErr" />
        </section>

        <section v-else-if="active === 'manage-resources'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Resource Centre Articles</h2><button type="button" @click="resetResourceForm(); active = 'resources'">New article</button></div>
          <ContentTable :items="resources" title-key="title" subtitle-key="category" @edit="editResource" @delete="removeResource" />
        </section>

        <section v-else-if="active === 'create-resource-categories'" class="cms-panel">
          <EditorForm title="Create Resource Article Category" :model="resourceCategoryForm" :fields="resourceCategoryFields" @save="saveResourceCategory" />
        </section>

        <section v-else-if="active === 'manage-resource-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Resource Article Categories</h2><button type="button" @click="resetResourceCategoryForm(); active = 'create-resource-categories'">New category</button></div>
          <ContentTable :items="resourceCategories" title-key="name" subtitle-key="slug" @edit="editResourceCategory" @delete="removeResourceCategory" />
        </section>

        <section v-else-if="active === 'careers'" class="cms-panel">
          <form class="otika-form-card" @submit.prevent="saveCareer">
            <div class="card">
              <div class="card-header"><h4>{{ careerForm.id ? 'Edit Career Post' : 'Add Career Post' }}</h4></div>
              <div class="card-body">
                <div class="section-title mt-0">Post Details</div>
                <div class="row">
                  <div class="form-group col-lg-8"><label>Title</label><input v-model="careerForm.title" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Status</label><select v-model="careerForm.status" class="form-control selectric"><option value="draft">Draft</option><option value="published">Published</option><option value="closed">Closed</option></select></div>
                  <div class="form-group col-lg-4"><label>Department</label><select v-model="careerForm.department_id" class="form-control selectric" @change="syncCareerDepartment"><option value="">Select department</option><option v-for="dept in departmentOptions" :key="dept.id" :value="dept.id">{{ dept.name }}</option></select></div>
                  <div class="form-group col-lg-4"><label>Category</label><select v-model="careerForm.category" class="form-control selectric"><option value="jobs">Jobs</option><option v-for="category in careerCategories" :key="category.slug || category.id" :value="category.slug || category.name">{{ category.name }}</option></select></div>
                  <div class="form-group col-lg-4"><label>Job Type</label><select v-model="careerForm.job_type" class="form-control selectric"><option value="full_time">Full Time</option><option value="part_time">Part Time</option><option value="contract">Contract</option><option value="internship">Internship</option></select></div>
                  <div class="form-group col-lg-4"><label>Location</label><input v-model="careerForm.location" class="form-control" /></div>
                  <div class="form-group col-lg-4"><label>Salary Range</label><input v-model="careerForm.salary_range" class="form-control" /></div>
                  <div class="form-group col-lg-4"><label>Deadline</label><input v-model="careerForm.deadline_at" type="datetime-local" class="form-control" /></div>
                </div>
                <div class="section-title">Description</div>
                <CmsRichTextEditor v-model="careerForm.description" />
                <div class="section-title">Requirements</div>
                <CmsRichTextEditor v-model="careerForm.requirements" />
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary mr-1">{{ careerForm.id ? 'Update Career Post' : 'Create Career Post' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-careers'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Career Posts</h2><button type="button" @click="resetCareerForm(); active = 'careers'">New post</button></div>
          <ContentTable :items="careers" title-key="title" subtitle-key="status" @edit="editCareer" @delete="removeCareer" />
        </section>

        <section v-else-if="active === 'create-career-categories'" class="cms-panel">
          <form class="otika-form-card" @submit.prevent="saveCareerCategory">
            <div class="card">
              <div class="card-header"><h4>{{ careerCategoryForm.id ? 'Edit Career Category' : 'Create Career Category' }}</h4></div>
              <div class="card-body">
                <div class="row">
                  <div class="form-group col-lg-6"><label>Name</label><input v-model="careerCategoryForm.name" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Slug</label><input v-model="careerCategoryForm.slug" class="form-control" /></div>
                  <div class="form-group col-lg-2"><label>Sort Order</label><input v-model.number="careerCategoryForm.sort_order" type="number" class="form-control" /></div>
                  <div class="form-group col-12"><label>Description</label><textarea v-model="careerCategoryForm.description" class="form-control otika-textarea"></textarea></div>
                  <div class="form-group col-12"><div class="custom-control custom-checkbox"><input id="career-category-active" v-model="careerCategoryForm.is_active" type="checkbox" class="custom-control-input" /><label class="custom-control-label" for="career-category-active">Active category</label></div></div>
                </div>
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary">{{ careerCategoryForm.id ? 'Update Category' : 'Create Category' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-career-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Career Categories</h2><button type="button" @click="resetCareerCategoryForm(); active = 'create-career-categories'">New category</button></div>
          <ContentTable :items="careerCategories" title-key="name" subtitle-key="slug" @edit="editCareerCategory" @delete="removeCareerCategory" />
        </section>

        <section v-else-if="active === 'invest'" class="cms-panel">
          <EditorForm title="Create Investment Post" :model="investForm" :fields="investFields" @save="saveInvest" @error="setErr" />
        </section>

        <section v-else-if="active === 'manage-invest'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Investment Posts</h2><button type="button" @click="resetInvestForm(); active = 'invest'">New post</button></div>
          <ContentTable :items="investItems" title-key="title" subtitle-key="subtitle" @edit="editInvest" @delete="removeInvest" />
        </section>

        <section v-else-if="active === 'create-invest-categories'" class="cms-panel">
          <EditorForm title="Create Investment Category" :model="investCategoryForm" :fields="investCategoryFields" @save="saveInvestCategory" />
        </section>

        <section v-else-if="active === 'manage-invest-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Investment Categories</h2><button type="button" @click="resetInvestCategoryForm(); active = 'create-invest-categories'">New category</button></div>
          <ContentTable :items="investCategories" title-key="name" subtitle-key="slug" @edit="editInvestCategory" @delete="removeInvestCategory" />
        </section>

        <section v-else-if="active === 'newsletter'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Newsletter Subscribers</h2>
            <div class="cms-actions-inline">
              <button type="button" @click="loadNewsletterSubscribers">Refresh</button>
              <button type="button" @click="exportNewsletter('csv')">Export CSV</button>
              <button type="button" @click="exportNewsletter('excel')">Export Excel</button>
            </div>
          </div>
          <div class="table-responsive newsletter-table">
            <table class="table table-striped table-hover table-sm mb-0">
              <thead><tr><th>Email</th><th>Date Submitted</th><th>Day</th><th>Status</th></tr></thead>
              <tbody>
                <tr v-for="subscriber in newsletterSubscribers" :key="subscriber.id || subscriber.email">
                  <td>{{ subscriber.email }}</td>
                  <td>{{ formatDateTime(subscriber.created_at || subscriber.submitted_at) }}</td>
                  <td>{{ formatDay(subscriber.created_at || subscriber.submitted_at) }}</td>
                  <td><span class="badge badge-success">{{ subscriber.status || 'subscribed' }}</span></td>
                </tr>
              </tbody>
            </table>
            <p v-if="!newsletterSubscribers.length" class="cms-empty">No newsletter subscribers found.</p>
          </div>
        </section>

        <section v-else-if="active === 'team'" class="cms-panel">
          <form class="otika-form-card" @submit.prevent="saveTeam">
            <div class="card">
              <div class="card-header"><h4>{{ teamForm.id ? 'Edit Team Member' : 'Add Team Member' }}</h4></div>
              <div class="card-body">
                <div class="section-title mt-0">Member Details</div>
                <div class="row">
                  <div class="form-group col-lg-6"><label>Full Name</label><input v-model="teamForm.full_name" class="form-control" required /></div>
                  <div class="form-group col-lg-6"><label>Designation</label><input v-model="teamForm.designation" class="form-control" /></div>
                  <div class="form-group col-lg-6"><label>Department</label><select v-model="teamForm.department_id" class="form-control selectric"><option value="">Select department</option><option v-for="dept in departmentOptions" :key="dept.id" :value="dept.id">{{ dept.name }}</option></select></div>
                  <div class="form-group col-lg-3"><label>Sort Order</label><input v-model.number="teamForm.sort_order" type="number" class="form-control" /></div>
                  <div class="form-group col-lg-3"><label class="d-block">Visibility</label><div class="custom-control custom-checkbox mt-2"><input id="team-active" v-model="teamForm.is_active" type="checkbox" class="custom-control-input" /><label class="custom-control-label" for="team-active">Active member</label></div></div>
                  <div class="form-group col-12">
                    <label>Profile Image</label>
                    <DropzoneUpload v-model="teamForm.image_url" label="team member image" accept="image/png,image/jpeg,image/webp" @error="setErr" />
                  </div>
                </div>
                <div class="section-title">Bio</div>
                <CmsRichTextEditor v-model="teamForm.bio" />
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary mr-1">{{ teamForm.id ? 'Update Team Member' : 'Create Team Member' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-team'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Team Members</h2><button type="button" @click="resetTeamForm(); active = 'team'">New team member</button></div>
          <ContentTable :items="teamMembers" title-key="full_name" subtitle-key="designation" @edit="editTeam" @delete="removeTeam" />
        </section>

        <section v-else-if="active === 'create-team-departments'" class="cms-panel">
          <form class="otika-form-card" @submit.prevent="saveTeamDepartment">
            <div class="card">
              <div class="card-header"><h4>{{ teamDepartmentForm.id ? 'Edit Department' : 'Create Department' }}</h4></div>
              <div class="card-body">
                <div class="row">
                  <div class="form-group col-lg-6"><label>Name</label><input v-model="teamDepartmentForm.name" class="form-control" required /></div>
                  <div class="form-group col-lg-4"><label>Slug</label><input v-model="teamDepartmentForm.slug" class="form-control" /></div>
                  <div class="form-group col-lg-2"><label>Sort Order</label><input v-model.number="teamDepartmentForm.sort_order" type="number" class="form-control" /></div>
                  <div class="form-group col-12"><label>Description</label><textarea v-model="teamDepartmentForm.description" class="form-control otika-textarea"></textarea></div>
                  <div class="form-group col-12"><div class="custom-control custom-checkbox"><input id="team-department-active" v-model="teamDepartmentForm.is_active" type="checkbox" class="custom-control-input" /><label class="custom-control-label" for="team-department-active">Active department</label></div></div>
                </div>
              </div>
              <div class="card-footer text-right"><button type="submit" class="btn btn-primary">{{ teamDepartmentForm.id ? 'Update Department' : 'Create Department' }}</button></div>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-team-departments'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Departments</h2><button type="button" @click="resetTeamDepartmentForm(); active = 'create-team-departments'">New department</button></div>
          <ContentTable :items="teamDepartments" title-key="name" subtitle-key="slug" @edit="editTeamDepartment" @delete="removeTeamDepartment" />
        </section>

        <section v-else-if="active === 'users'" class="cms-panel">
          <div class="cms-panel-head"><h2>Create New User</h2><button type="button" @click="resetUserForm">Clear</button></div>
          <form class="cms-editor" @submit.prevent="saveUser">
            <div class="cms-two">
              <label>First name<input v-model="userForm.first_name" required /></label>
              <label>Last name<input v-model="userForm.last_name" required /></label>
              <label>Email<input v-model="userForm.email" type="email" required :disabled="!!userForm.id" /></label>
              <label>Phone<input v-model="userForm.phone" /></label>
              <label v-if="!userForm.id">Password<input v-model="userForm.password" type="password" required minlength="8" autocomplete="new-password" /></label>
            </div>
            <div class="cms-actions-inline">
              <button type="submit">{{ userForm.id ? 'Update user' : 'Create user' }}</button>
              <button type="button" @click="active = 'manage-users'; loadUsers()">Manage users</button>
            </div>
          </form>
        </section>

        <section v-else-if="active === 'manage-users'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Manage Users</h2>
            <div class="cms-actions-inline">
              <button type="button" @click="resetUserForm(); active = 'users'">New user</button>
              <button type="button" @click="loadUsers">Refresh</button>
            </div>
          </div>
          <form class="cms-search" @submit.prevent="loadUsers">
            <input v-model="userSearch" placeholder="Search users by name or email..." />
            <button type="submit">Search</button>
          </form>
          <div class="card otika-basic-table-card">
            <div class="card-body p-0">
              <div class="table-responsive">
                <table class="table table-striped table-hover table-sm mb-0">
                  <thead><tr><th>User</th><th>Status</th><th>Roles</th><th class="text-right">Actions</th></tr></thead>
                  <tbody>
                    <tr v-for="user in users" :key="user.id">
                      <td><strong>{{ user.first_name }} {{ user.last_name }}</strong><br /><span>{{ user.email }}</span></td>
                      <td><span class="badge" :class="user.is_active ? 'badge-success' : 'badge-danger'">{{ user.account_status || (user.is_active ? 'ACTIVE' : 'DISABLED') }}</span></td>
                      <td>
                        <span v-for="role in user.roles || []" :key="role.id || role.name" class="badge badge-primary mr-1">{{ role.name || role }}</span>
                        <span v-if="!(user.roles || []).length" class="text-muted">No roles</span>
                      </td>
                      <td class="table-actions text-right">
                        <button v-if="hasPermission('users:update')" type="button" class="btn btn-sm btn-primary" @click="editUser(user)">Edit</button>
                        <button v-if="hasPermission('users:roles') || hasPermission('roles:assign')" type="button" class="btn btn-sm btn-info" @click="openUserRoles(user)">Roles</button>
                        <button v-if="hasPermission('users:reset_password')" type="button" class="btn btn-sm btn-warning" @click="promptResetPassword(user)">Reset</button>
                        <button v-if="hasPermission('users:activate')" type="button" class="btn btn-sm btn-secondary" @click="toggleUserActive(user)">{{ user.is_active ? 'Disable' : 'Activate' }}</button>
                        <button v-if="hasPermission('users:delete')" type="button" class="btn btn-sm btn-danger" @click="removeUser(user)">Delete</button>
                      </td>
                    </tr>
                    <tr v-if="!users.length"><td colspan="4" class="text-center text-muted py-4">No users found.</td></tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
          <article v-if="selectedUser" class="cms-subpanel">
            <div class="cms-panel-head">
              <h2>Roles for {{ selectedUser.first_name }} {{ selectedUser.last_name }}</h2>
              <span>{{ selectedUser.email }}</span>
            </div>
            <div class="permission-grid">
              <label v-for="role in roles" :key="role.id" class="permission-row">
                <input type="checkbox" :checked="userHasRole(selectedUser, role)" @change="toggleUserRole(selectedUser, role, $event.target.checked)" />
                <span><strong>{{ role.name }}</strong><small>{{ role.description || (role.is_system ? 'System role' : 'Custom role') }}</small></span>
              </label>
            </div>
          </article>
        </section>

        <section v-else-if="active === 'roles'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Role Management</h2>
            <button type="button" @click="resetRoleForm">New role</button>
          </div>
          <div class="cms-admin-grid">
            <form class="cms-editor" @submit.prevent="saveRole">
              <div class="cms-two">
                <label>Role name<input v-model="roleForm.name" required placeholder="content_editor_events" /></label>
                <label class="wide">Description<textarea v-model="roleForm.description" placeholder="Describe what this custom content manager can do."></textarea></label>
              </div>
              <div class="cms-actions-inline">
                <button type="submit">{{ roleForm.id ? 'Update role' : 'Create role' }}</button>
                <button v-if="roleForm.id && !roleForm.is_system" type="button" @click="removeRole(roleForm)">Delete role</button>
              </div>
            </form>
            <article class="cms-subpanel">
              <h3>Custom roles</h3>
              <div class="cms-table compact">
                <article v-for="role in roles" :key="role.id" class="cms-table-row" :class="{ selected: role.id === roleForm.id }">
                  <div>
                    <strong>{{ role.name }}</strong>
                    <span>{{ role.description || (role.is_system ? 'System role' : 'Custom role') }}</span>
                  </div>
                  <div><button type="button" @click="selectRole(role)">Customize</button></div>
                </article>
              </div>
            </article>
          </div>
          <article class="cms-subpanel">
            <div class="cms-panel-head">
              <h2>Permission Matrix</h2>
              <span>{{ roleForm.id ? roleForm.name : 'Select or create a role' }}</span>
            </div>
            <div class="permission-grid">
              <section v-for="group in permissionGroups" :key="group.resource" class="permission-group">
                <h3>{{ group.label }}</h3>
                <label v-for="permission in group.permissions" :key="permission.id" class="permission-row">
                  <input
                    type="checkbox"
                    :checked="selectedPermissionIds.has(permission.id)"
                    :disabled="!roleForm.id"
                    @change="togglePermission(permission, $event.target.checked)"
                  />
                  <span>
                    <strong>{{ permission.action }}</strong>
                    <small>{{ permission.description || permission.name }}</small>
                  </span>
                </label>
              </section>
            </div>
          </article>
        </section>

        <section v-else-if="active === 'manage-roles'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Manage Roles</h2>
            <button type="button" @click="resetRoleForm(); active = 'roles'">New role</button>
          </div>
          <div class="cms-admin-grid">
            <article class="cms-subpanel">
              <h3>Custom roles</h3>
              <div class="cms-table compact">
                <article v-for="role in roles" :key="role.id" class="cms-table-row" :class="{ selected: role.id === roleForm.id }">
                  <div>
                    <strong>{{ role.name }}</strong>
                    <span>{{ role.description || (role.is_system ? 'System role' : 'Custom role') }}</span>
                  </div>
                  <div><button type="button" @click="selectRole(role)">Customize</button></div>
                </article>
              </div>
            </article>
            <article class="cms-subpanel">
              <div class="cms-panel-head">
                <h2>Permission Matrix</h2>
                <span>{{ roleForm.id ? roleForm.name : 'Select a role' }}</span>
              </div>
              <div class="permission-grid">
                <section v-for="group in permissionGroups" :key="group.resource" class="permission-group">
                  <h3>{{ group.label }}</h3>
                  <label v-for="permission in group.permissions" :key="permission.id" class="permission-row">
                    <input type="checkbox" :checked="selectedPermissionIds.has(permission.id)" :disabled="!roleForm.id" @change="togglePermission(permission, $event.target.checked)" />
                    <span>
                      <strong>{{ permission.action }}</strong>
                      <small>{{ permission.description || permission.name }}</small>
                    </span>
                  </label>
                </section>
              </div>
            </article>
          </div>
        </section>

        <section v-else-if="active === 'audit-logs'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Audit Logs</h2>
            <form class="cms-search" @submit.prevent="loadAuditLogs">
              <input v-model="auditSearch" placeholder="Search user, IP, endpoint, action..." />
              <button type="submit">Search</button>
            </form>
          </div>
          <div class="audit-summary">
            <article><span>Total events</span><strong>{{ auditMeta.total || auditLogs.length }}</strong></article>
            <article><span>Failed/blocked</span><strong>{{ auditRiskCount }}</strong></article>
            <article><span>High threat</span><strong>{{ auditHighThreatCount }}</strong></article>
          </div>
          <div class="audit-table">
            <article v-for="log in auditLogs" :key="log.id" class="audit-row">
              <div>
                <strong>{{ log.method || log.action }} {{ log.endpoint || log.resource }}</strong>
                <span>{{ formatDateTime(log.created_at) }} · {{ log.username || log.first_name || 'Anonymous' }} · {{ log.ip_address || log.forwarded_ip || 'No IP' }}</span>
              </div>
              <div>
                <span :class="['audit-code', statusClass(log)]">{{ log.response_code || log.event_status || 'n/a' }}</span>
                <span>{{ log.response_time_ms || 0 }}ms</span>
                <button type="button" @click="viewAuditLog(log)">Details</button>
              </div>
            </article>
            <p v-if="!auditLogs.length" class="cms-empty">No audit entries found.</p>
          </div>
          <article v-if="selectedAuditLog" class="cms-subpanel audit-detail">
            <div class="cms-panel-head">
              <h2>Audit Detail</h2>
              <button type="button" @click="selectedAuditLog = null">Close</button>
            </div>
            <pre>{{ selectedAuditLog }}</pre>
          </article>
        </section>

        <section v-else-if="active === 'menus'" class="cms-panel menu-panel">
          <div class="cms-panel-head"><h2>Main Menu Builder</h2></div>
          <div class="cms-menu-builder-grid">
            <article>
              <h3>Main menu</h3>
              <MenuBuilder v-model="mainMenuTree" :pages="internalPageOptions" />
            </article>
            <article>
              <h3>Footer menu</h3>
              <MenuBuilder v-model="footerMenuTree" :pages="internalPageOptions" />
            </article>
          </div>
          <p class="footer-columns-hint">The footer's link columns (title + links) are managed under Settings → Footer Columns, not here.</p>
          <div class="menu-save-bar">
            <span v-if="menusDirty" class="menu-save-bar__status menu-save-bar__status--dirty">Unsaved changes</span>
            <span v-else class="menu-save-bar__status">Up to date</span>
            <button type="button" :disabled="!menusDirty || savingMenus" @click="discardMenuChanges">Discard changes</button>
            <button type="button" class="menu-save-bar__save" :disabled="savingMenus" @click="saveMenus">{{ savingMenus ? 'Saving…' : 'Save menus' }}</button>
          </div>
        </section>

        <section v-else-if="active === 'settings'" class="cms-panel">
          <div class="cms-panel-head"><h2>Contact Details</h2><button @click="saveSettings">Save contact details</button></div>
          <div class="cms-two">
            <label>Phone<input v-model="contact.phone" class="form-control" placeholder="+256..." /></label>
            <label>WhatsApp number<input v-model="contact.whatsapp" class="form-control" placeholder="+256..." /></label>
            <label>Email<input v-model="contact.email" class="form-control" placeholder="info@ncs.go.ug" /></label>
            <label>Fax<input v-model="contact.fax" class="form-control" /></label>
            <label>Location<input v-model="contact.location" class="form-control" placeholder="Lugogo, Kampala" /></label>
            <label>P.O. Box<input v-model="contact.postal_address" class="form-control" placeholder="P.O. Box..." /></label>
            <label class="wide">Physical address<textarea v-model="contact.address" class="form-control"></textarea></label>
            <label class="wide">Google Maps link<textarea v-model="contact.mapUrl" class="form-control" placeholder="https://maps.google.com/..."></textarea></label>
          </div>
          <article class="cms-subpanel">
            <div class="cms-panel-head"><h2>Social Media Links</h2></div>
            <div class="cms-two">
              <label>Facebook<input v-model="contact.social.facebook" class="form-control" /></label>
              <label>X / Twitter<input v-model="contact.social.twitter" class="form-control" /></label>
              <label>LinkedIn<input v-model="contact.social.linkedin" class="form-control" /></label>
              <label>Instagram<input v-model="contact.social.instagram" class="form-control" /></label>
              <label>YouTube<input v-model="contact.social.youtube" class="form-control" /></label>
            </div>
          </article>
          <article class="cms-subpanel">
            <div class="cms-panel-head"><h2>Footer About NCS</h2></div>
            <div class="cms-two">
              <label class="wide">About NCS footer description<textarea v-model="footer.about" class="form-control"></textarea></label>
              <label class="wide">Copyright text<input v-model="footer.copyright" class="form-control" placeholder="National Council of Sports, Uganda. All rights reserved." /></label>
            </div>
          </article>
          <article class="cms-subpanel">
            <div class="cms-panel-head"><h2>Footer Columns</h2><button type="button" @click="addFooterColumn">Add column</button></div>
            <p class="footer-columns-hint">These columns render side-by-side in the site footer, in this order. Both the title and the links in each column are fully editable.</p>
            <div v-for="(column, cIdx) in footer.columns" :key="cIdx" class="footer-column-editor">
              <div class="cms-panel-head">
                <label class="footer-column-title">Column title<input v-model="column.title" class="form-control" placeholder="e.g. Quick Links" /></label>
                <button type="button" class="footer-remove-btn" @click="footer.columns.splice(cIdx, 1)">Remove column</button>
              </div>
              <div v-for="(link, lIdx) in column.links" :key="lIdx" class="footer-link-row">
                <input v-model="link.label" class="form-control" placeholder="Link label" />
                <input v-model="link.url" class="form-control" placeholder="/path or https://..." />
                <button type="button" class="footer-remove-btn" @click="column.links.splice(lIdx, 1)">Remove</button>
              </div>
              <button type="button" @click="column.links.push({ label:'', url:'/' })">Add link</button>
              <p v-if="!column.links.length" class="cms-empty">No links in this column yet.</p>
            </div>
            <p v-if="!footer.columns.length" class="cms-empty">No footer columns yet — add one above.</p>
          </article>
        </section>

              <WebsiteSettingsPanel v-else-if="active === 'website-settings'" @message="setMsg" @error="setErr" />
              <SitemapPanel v-else-if="active === 'sitemap'" @message="setMsg" @error="setErr" />
              <AppearanceSettingsPanel v-else-if="active === 'appearance'" @message="setMsg" @error="setErr" />
              <ProfileSettingsPanel v-else-if="active === 'my-profile'" @message="setMsg" @error="setErr" @profile-updated="refreshStoredUser" />
              <StorageSettingsPanel v-else-if="active === 'storage'" @message="setMsg" @error="setErr" />
              <SystemCommandCenterPanel v-else-if="active === 'command-center'" @message="setMsg" @error="setErr" />
              <MaintenanceModePanel v-else-if="active === 'maintenance'" @message="setMsg" @error="setErr" />
              <SmartUpdatesPanel v-else-if="active === 'smart-updates'" @message="setMsg" @error="setErr" />
            </div>
          </section>
        </div>
        <footer class="main-footer cms-main-footer">
          <div class="footer-left">
            Design By: Ateni Media Technologies LLC
          </div>
          <div class="footer-right">
            National Council of Sports CMS
          </div>
        </footer>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, defineComponent, h, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { API_BASE_URL, mediaUrl } from '@/api/client.js'
import * as cms from '@/api/cms.js'
import BlogPostEditor from '@/components/cms/BlogPostEditor.vue'
import CmsRichTextEditor from '@/components/cms/CmsRichTextEditor.vue'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'
import MenuBuilder from '@/components/cms/MenuBuilder.vue'
import SlideshowManager from '@/components/cms/SlideshowManager.vue'
import StaticPageBuilder from '@/components/cms/StaticPageBuilder.vue'
import StorageSettingsPanel from '@/components/cms/StorageSettingsPanel.vue'
import AppearanceSettingsPanel from '@/components/cms/AppearanceSettingsPanel.vue'
import WebsiteSettingsPanel from '@/components/cms/WebsiteSettingsPanel.vue'
import SitemapPanel from '@/components/cms/SitemapPanel.vue'
import ProfileSettingsPanel from '@/components/cms/ProfileSettingsPanel.vue'
import SystemCommandCenterPanel from '@/components/cms/SystemCommandCenterPanel.vue'
import MaintenanceModePanel from '@/components/cms/MaintenanceModePanel.vue'
import SmartUpdatesPanel from '@/components/cms/SmartUpdatesPanel.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'
import { normalizeMenuTree, toCmsMenuItems } from '@/utils/menuTree.js'

const router = useRouter()
const active = ref('overview')
const message = ref('')
const error = ref('')
const apiAvailable = ref(null)
const sidebarCollapsed = ref(localStorage.getItem('ncsms_sidebar_collapsed') === 'true')
const homepageGroupOpen = ref(false)
const blogsGroupOpen = ref(false)
const slideshowGroupOpen = ref(false)
const staticPagesGroupOpen = ref(false)
const projectsGroupOpen = ref(false)
const caseStudiesGroupOpen = ref(false)
const faqsGroupOpen = ref(false)
const resourcesGroupOpen = ref(false)
const careersGroupOpen = ref(false)
const teamGroupOpen = ref(false)
const rolesGroupOpen = ref(false)
const usersGroupOpen = ref(false)
const facilitiesGroupOpen = ref(false)
const eventsGroupOpen = ref(false)
const investGroupOpen = ref(false)
const federationsGroupOpen = ref(false)
const funFactsGroupOpen = ref(false)
const newsletterGroupOpen = ref(false)
const messagesOpen = ref(false)
const notificationsOpen = ref(false)
const profileOpen = ref(false)
const globalSearch = ref('')
const analytics = ref(null)

ensureOtikaStyles()

const storedUser = reactive((() => {
  try { return JSON.parse(localStorage.getItem('ncsms_user') || '{}') } catch { return {} }
})())
function refreshStoredUser(patch) {
  Object.assign(storedUser, patch)
}
const currentUserName = computed(() => {
  const fullName = storedUser.first_name ? `${storedUser.first_name} ${storedUser.last_name || ''}`.trim() : ''
  return storedUser.username || storedUser.name || storedUser.display_name || fullName || emailLocalPart(storedUser.email) || 'cms-admin'
})
const currentUserAvatar = computed(() => storedUser.avatar_url ? mediaUrl(storedUser.avatar_url) : '')
const sidebarUserName = computed(() => {
  const value = currentUserName.value || 'cms-admin'
  return value.length > 18 ? `${value.slice(0, 15)}...` : value
})
function timeAgo(value) {
  if (!value) return ''
  const diffMs = Date.now() - new Date(value).getTime()
  const mins = Math.floor(diffMs / 60000)
  if (mins < 1) return 'Just now'
  if (mins < 60) return `${mins} Min${mins === 1 ? '' : 's'} Ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} Hour${hours === 1 ? '' : 's'} Ago`
  const days = Math.floor(hours / 24)
  return `${days} Day${days === 1 ? '' : 's'} Ago`
}
const NOTIFICATION_THEME_BADGE = { warn:'bg-warning', success:'bg-success', info:'bg-info' }
const topMessages = computed(() => contactMessages.value.slice(0, 5).map(item => ({
  id: item.id,
  name: item.subject || 'Contact Us Message',
  text: `${item.name || 'Visitor'}: ${item.message || ''}`,
  time: timeAgo(item.created_at),
  avatar: '/otika-assets/img/users/user-1.png',
  target: 'messages',
})))
const notificationItems = computed(() => cmsNotifications.value.slice(0, 5).map(item => ({
  id: item.id,
  text: item.title || item.message || 'Notification',
  time: timeAgo(item.created_at),
  icon: notificationIcon(item),
  color: NOTIFICATION_THEME_BADGE[notificationThemeClass(item)] || 'bg-info',
  unread: item.status === 'unread',
  target: 'notifications',
})))
const unreadMessageCount = computed(() => contactMessages.value.filter(m => m.status === 'unread').length)
const unreadNotificationCount = computed(() => cmsNotifications.value.filter(n => n.status === 'unread').length)

const topSections = [
  { id:'overview', label:'Dashboard', icon:'icofont-dashboard-web' },
  { id:'analytics', label:'Analytics', icon:'icofont-chart-histogram' },
]
const homepageSections = [
  { id:'homepage', label:'General Sections', icon:'icofont-home' },
  { id:'homepage-about', label:'About NCS Section', icon:'icofont-info-circle' },
  { id:'homepage-topbar', label:'Topbar Section', icon:'icofont-megaphone' },
  { id:'core', label:'Core Functions', icon:'icofont-check-circled' },
  { id:'homepage-sports-excellence', label:'Sports Excellence', icon:'icofont-chart-growth' },
]
const slideshowSections = [
  { id:'slideshow-manager', label:'Homepage Hero', icon:'icofont-slidshare' },
]
const blogSections = [
  { id:'create-post', label:'Create New Post', icon:'icofont-edit' },
  { id:'manage-posts', label:'Manage Posts', icon:'icofont-list' },
  { id:'create-blog-categories', label:'Create Blog Categories', icon:'icofont-folder-open' },
  { id:'manage-blog-categories', label:'Manage Categories', icon:'icofont-tags' },
]
const staticPageSections = [
  { id:'static-pages', label:'Create Static Page', icon:'icofont-page' },
  { id:'manage-pages', label:'Manage Static Pages', icon:'icofont-copy' },
]
const projectSections = [
  { id:'projects', label:'Add New Project', icon:'icofont-plus-circle' },
  { id:'manage-projects', label:'Manage Projects', icon:'icofont-list' },
  { id:'create-project-categories', label:'Create Project Categories', icon:'icofont-folder-open' },
  { id:'manage-project-categories', label:'Manage Project Categories', icon:'icofont-tags' },
]
const caseStudySections = [
  { id:'case-studies', label:'Create New Case Study Post', icon:'icofont-edit' },
  { id:'manage-case-studies', label:'Manage Case Study Posts', icon:'icofont-list' },
  { id:'create-case-study-categories', label:'Create Case Study Category', icon:'icofont-folder-open' },
  { id:'manage-case-study-categories', label:'Manage Case Study Categories', icon:'icofont-tags' },
]
const faqSections = [
  { id:'faqs', label:'Create New Article', icon:'icofont-edit' },
  { id:'manage-faqs', label:'Manage Articles', icon:'icofont-list' },
  { id:'create-faq-categories', label:'Create New FAQ Category', icon:'icofont-folder-open' },
  { id:'manage-faq-categories', label:'Manage Categories', icon:'icofont-tags' },
]
const resourceSections = [
  { id:'resources', label:'Create New Article', icon:'icofont-edit' },
  { id:'manage-resources', label:'Manage Articles', icon:'icofont-list' },
  { id:'create-resource-categories', label:'Create Article Category', icon:'icofont-folder-open' },
  { id:'manage-resource-categories', label:'Manage Article Categories', icon:'icofont-tags' },
]
const careerSections = [
  { id:'careers', label:'Add New Post', icon:'icofont-plus-circle' },
  { id:'manage-careers', label:'Manage Posts', icon:'icofont-list' },
  { id:'create-career-categories', label:'Create Category', icon:'icofont-folder-open' },
  { id:'manage-career-categories', label:'Manage Categories', icon:'icofont-tags' },
]
const teamSections = [
  { id:'team', label:'Add Team Member', icon:'icofont-user-suited' },
  { id:'manage-team', label:'Manage Team Members', icon:'icofont-list' },
  { id:'create-team-departments', label:'Create Department', icon:'icofont-building-alt' },
  { id:'manage-team-departments', label:'Manage Departments', icon:'icofont-tags' },
]
const roleSections = [
  { id:'roles', label:'Add New Role', icon:'icofont-plus-circle' },
  { id:'manage-roles', label:'Manage Roles', icon:'icofont-list' },
]
const userSections = [
  { id:'users', label:'Create New User', icon:'icofont-user-alt-3' },
  { id:'manage-users', label:'Manage Users', icon:'icofont-users-social' },
]
const facilitySections = [
  { id:'facilities', label:'Add New Facility', icon:'icofont-plus-circle' },
  { id:'manage-facilities', label:'Manage Facilities', icon:'icofont-list' },
  { id:'create-facility-categories', label:'Add Facility Category', icon:'icofont-folder-open' },
  { id:'manage-facility-categories', label:'Manage Facility Categories', icon:'icofont-tags' },
]
const eventSections = [
  { id:'events', label:'Add New Event', icon:'icofont-plus-circle' },
  { id:'manage-events', label:'Manage Events', icon:'icofont-list' },
  { id:'create-event-categories', label:'Add Event Category', icon:'icofont-folder-open' },
  { id:'manage-event-categories', label:'Manage Event Categories', icon:'icofont-tags' },
]
const investSections = [
  { id:'create-invest-categories', label:'Create Investment Category', icon:'icofont-folder-open' },
  { id:'manage-invest-categories', label:'Manage Investment Categories', icon:'icofont-tags' },
  { id:'invest', label:'Create Investment Post', icon:'icofont-edit' },
  { id:'manage-invest', label:'Manage Investment Posts', icon:'icofont-list' },
]
const federationSections = [
  { id:'associations', label:'Add New Federation', icon:'icofont-plus-circle' },
  { id:'manage-federations', label:'Manage Federations', icon:'icofont-list' },
  { id:'create-federation-categories', label:'Create Federation Category', icon:'icofont-folder-open' },
  { id:'manage-federation-categories', label:'Manage Federation Categories', icon:'icofont-tags' },
]
const funFactSections = [
  { id:'facts', label:'Create Article', icon:'icofont-edit' },
  { id:'manage-facts', label:'Manage Articles', icon:'icofont-list' },
]
const newsletterSections = [
  { id:'newsletter', label:'Subscribers', icon:'icofont-email' },
]
const contentSections = [
  { id:'messages', label:'Messages', icon:'icofont-envelope' },
  { id:'notifications', label:'Notifications', icon:'icofont-notification' },
  { id:'comments', label:'Comment Moderation', icon:'icofont-speech-comments' },
  { id:'audit-logs', label:'Audit Logs', icon:'icofont-shield-alt' },
  { id:'menus', label:'Main Menu', icon:'icofont-navigation-menu' },
  { id:'settings', label:'Contact Details', icon:'icofont-contacts' },
  { id:'website-settings', label:'Website Settings', icon:'icofont-globe' },
  { id:'sitemap', label:'Sitemap', icon:'icofont-site-map' },
  { id:'appearance', label:'Appearance', icon:'icofont-paint' },
  { id:'storage', label:'Storage Settings', icon:'icofont-cloud-upload' },
  { id:'command-center', label:'Command Center', icon:'icofont-layers' },
  { id:'maintenance', label:'Maintenance & Backups', icon:'icofont-shield' },
  { id:'smart-updates', label:'Smart Updates', icon:'icofont-refresh' },
]
// Reachable only via the navbar profile dropdown, not the sidebar — kept out
// of contentSections so it doesn't clutter Modules, but still included below
// so currentSection resolves its label for the page header/breadcrumb.
const profileSections = [
  { id:'my-profile', label:'Profile Settings', icon:'far fa-user' },
]
const sections = [...topSections, ...homepageSections, ...slideshowSections, ...blogSections, ...staticPageSections, ...projectSections, ...caseStudySections, ...faqSections, ...resourceSections, ...careerSections, ...teamSections, ...roleSections, ...userSections, ...facilitySections, ...eventSections, ...investSections, ...federationSections, ...funFactSections, ...newsletterSections, ...contentSections, ...profileSections]
const currentSection = computed(() => sections.find(s => s.id === active.value) || sections[0])

const sectionPermissionMap = {
  overview:['dashboard:read'],
  analytics:['analytics:read','dashboard:read'],
  homepage:['homepage:read','homepage:update'],
  'homepage-about':['homepage:read','homepage:update'],
  'homepage-topbar':['homepage:read','homepage:update'],
  core:['homepage:read','homepage:update'],
  'homepage-sports-excellence':['homepage:read','homepage:update'],
  'slideshow-manager':['slideshows:read'],
  'create-post':['blog_posts:create'],
  'manage-posts':['blog_posts:read'],
  'create-blog-categories':['blog_categories:create'],
  'manage-blog-categories':['blog_categories:read'],
  'static-pages':['static_pages:create'],
  'manage-pages':['static_pages:read'],
  projects:['projects:create'],
  'manage-projects':['projects:read'],
  'create-project-categories':['project_categories:create'],
  'manage-project-categories':['project_categories:read'],
  'case-studies':['case_studies:create'],
  'manage-case-studies':['case_studies:read'],
  'create-case-study-categories':['case_study_categories:create'],
  'manage-case-study-categories':['case_study_categories:read'],
  faqs:['faqs:create'],
  'manage-faqs':['faqs:read'],
  'create-faq-categories':['faq_categories:create'],
  'manage-faq-categories':['faq_categories:read'],
  resources:['resources:create'],
  'manage-resources':['resources:read'],
  'create-resource-categories':['resource_categories:create'],
  'manage-resource-categories':['resource_categories:read'],
  careers:['careers:create'],
  'manage-careers':['careers:read'],
  'create-career-categories':['career_categories:create'],
  'manage-career-categories':['career_categories:read'],
  team:['team_members:create'],
  'manage-team':['team_members:read'],
  'create-team-departments':['team_departments:create'],
  'manage-team-departments':['team_departments:read'],
  roles:['roles:write'],
  'manage-roles':['roles:read'],
  users:['users:create'],
  'manage-users':['users:read'],
  facilities:['facilities:create'],
  'manage-facilities':['facilities:read'],
  'create-facility-categories':['facility_categories:create'],
  'manage-facility-categories':['facility_categories:read'],
  events:['events:create'],
  'manage-events':['events:read'],
  'create-event-categories':['event_categories:create'],
  'manage-event-categories':['event_categories:read'],
  invest:['investments:create'],
  'manage-invest':['investments:read'],
  'create-invest-categories':['investment_categories:create'],
  'manage-invest-categories':['investment_categories:read'],
  associations:['federations:create'],
  'manage-federations':['federations:read'],
  'create-federation-categories':['federation_categories:create'],
  'manage-federation-categories':['federation_categories:read'],
  facts:['fun_facts:create'],
  'manage-facts':['fun_facts:read'],
  newsletter:['newsletter:read'],
  messages:['messages:read'],
  notifications:['notifications:read'],
  comments:['comments:read'],
  'audit-logs':['audit:read'],
  menus:['menus:read'],
  settings:['settings:read'],
  'website-settings':['settings:read'],
  sitemap:['settings:read'],
  storage:['storage:read'],
  'command-center':['dashboard:read'],
  maintenance:['dashboard:read'],
  'smart-updates':['dashboard:read'],
}
const resourceLabels = {
  homepage:'Homepage Management', slideshows:'Slideshows', blog_posts:'Blog Posts', blog_categories:'Blog Categories',
  static_pages:'Static Pages', projects:'Projects', project_categories:'Project Categories', case_studies:'Case Studies',
  case_study_categories:'Case Study Categories', faqs:'FAQ Articles', faq_categories:'FAQ Categories', resources:'Resource Centre',
  resource_categories:'Resource Categories', careers:'Career Posts', career_categories:'Career Categories', team_members:'Team Members',
  team_departments:'Team Departments', facilities:'Facilities', facility_categories:'Facility Categories', events:'Events',
  event_categories:'Event Categories', investments:'Investment Posts', investment_categories:'Investment Categories',
  federations:'Federations', federation_categories:'Federation Categories', fun_facts:'Fun Facts', newsletter:'Newsletter',
  comments:'Comments', menus:'Main Menu', settings:'Contact Details', storage:'Storage Settings', users:'Users',
  roles:'Roles', audit:'Audit Logs', dashboard:'Dashboard', cms:'Legacy CMS', applications:'Applications',
  messages:'Messages', notifications:'Notifications',
}
const actionOrder = ['read', 'create', 'update', 'write', 'delete', 'activate', 'assign', 'roles', 'reset_password', 'export']

const homepageDefaults = {
  stats_title:'Sports Excellence in Numbers',
  stats_intro:'Driving the development of sports across Uganda through dedicated programs and world-class facilities',
  about:{eyebrow:'About NCS',title:'Developing Sports Excellence Since 1964',intro:'The National Council of Sports (NCS) is a statutory body established to develop, promote, and control sports in Uganda under the Ministry of Education and Sports.',body:'Established under the <strong>National Council of Sports Act (Chapter 48)</strong>, assented on 22 June 1964 and commenced on 25 June 1964, NCS serves as the apex regulator for sports development in Uganda, now updated by the <strong>National Sports Act, 2023</strong>.',leadership_label:'View Current Membership',leadership_url:'/team',core_title:'Core Functions of NCS',core_intro:'As mandated by the National Sports Act, NCS performs the following key functions:',mandate_label:'Read Full Mandate',mandate_url:'/pages/the-mandate'},
  leadership:{chairperson_name:'Mr. Ambrose Tashobya',chairperson_image:'',secretary_name:'Dr. Bernard Patrick Ogwel'},
  milestones:[
    {value:'60+',label:'Years of Excellence',description:'Established in 1964 and still driving national sport.',icon:'icofont-award'},
    {value:'54+',label:'Sports Associations',description:'Recognised bodies supported across Uganda.',icon:'icofont-trophy'},
    {value:'32+',label:'Sports Facilities',description:'Facilities and venues supporting athletes and federations.',icon:'icofont-building-alt'},
    {value:'100K+',label:'Athletes Reached',description:'Athletes, coaches, and administrators served through NCS programs.',icon:'icofont-users-alt-5'},
  ],
  values:[
    {title:'Our Mission',text:'Maximizing opportunities for all Ugandans to participate and excel in Sports.',icon:'icofont-dart',featured:true},
    {title:'Our Vision',text:'A centre of excellence for promotion and development of Sports.',icon:'icofont-eye'},
    {title:'Integrity',text:'Upholding the highest standards of ethics and fair play in all sporting activities.',icon:'icofont-shield'},
    {title:'Inclusivity',text:'Ensuring sports opportunities are accessible to all Ugandans regardless of background.',icon:'icofont-people'},
    {title:'Excellence',text:'Striving for the highest standards in athlete development and sports administration.',icon:'icofont-award'},
    {title:'Global Recognition',text:'Positioning Uganda as a leading sports nation on the African and world stage.',icon:'icofont-globe',featured:true},
  ],
  facilities:{eyebrow:'World-Class Infrastructure', title:'Our Sports Facilities', intro:'', button_label:'Explore All Facilities'},
  events:{eyebrow:'Upcoming Events', title:'NCS Calendar', intro:''},
  core_functions:['Developing, promoting, and controlling sports on a national basis, including training and staffing','Recognizing sports disciplines and registering national sports organizations','Regulating associations and federations, awarding medals, certificates, trophies, and incentives','Approving international and national competitions and festivals'],
  topbar:{
    marquee:['Welcome to National Council of Sports Uganda','A centre of excellence for promotion and development of Sports.','Maximizing opportunities for all Ugandans to participate and excel in Sports.','Established 1964'],
    webmail_url:'https://mail.umcs.go.ug/',
    social:{facebook:'',twitter:'',linkedin:'',instagram:'',youtube:''},
  },
}
const homepage = reactive(JSON.parse(JSON.stringify(homepageDefaults)))
const contact = reactive({
  phone:'',
  whatsapp:'',
  email:'',
  fax:'',
  location:'',
  address:'',
  postal_address:'',
  hours:'',
  mapUrl:'',
  social:{ facebook:'', twitter:'', linkedin:'', instagram:'', youtube:'' },
})
const footer = reactive({ about:'', copyright:'', columns:[] })
function addFooterColumn() {
  footer.columns.push({ title:'New column', links:[] })
}

const posts = ref([]), pages = ref([]), projects = ref([]), caseStudies = ref([]), events = ref([]), slides = ref([]), facilities = ref([]), associations = ref([]), facts = ref([]), faqs = ref([]), resources = ref([]), careers = ref([]), investItems = ref([]), teamMembers = ref([])
const blogCategories = ref([])
const projectCategories = ref([])
const caseStudyCategories = ref([])
const faqCategories = ref([])
const resourceCategories = ref([])
const careerCategories = ref([])
const teamDepartments = ref([])
const institutionalDepartments = ref([])
const facilityCategories = ref([])
const eventCategories = ref([])
const investCategories = ref([])
const federationCategories = ref([])
const newsletterSubscribers = ref([])
const contactMessages = ref([])
const cmsNotifications = ref([])
const comments = ref([])
const users = ref([])
const roles = ref([])
const permissions = ref([])
const selectedPermissionIds = ref(new Set())
const selectedUser = ref(null)
const userSearch = ref('')
const auditLogs = ref([])
const auditMeta = ref({})
const auditSearch = ref('')
const selectedAuditLog = ref(null)
const commentStatus = ref('pending')
const messageStatus = ref('')
const notificationStatus = ref('')
const selectedMessage = ref(null)
const slideshow = reactive({ id:'homepage-hero', name:'Homepage Hero', slug:'homepage-hero', transition_effect:'fade', transition_duration:700, autoplay_speed:6500, pause_on_hover:true, is_active:true })
const mainMenuTree = ref([]), footerMenuTree = ref([])
const savingMenus = ref(false)
const menusSavedSnapshot = ref(null)

const STATIC_SITE_PAGES = [
  { url: '/', label: 'Home' },
  { url: '/news', label: 'News' },
  { url: '/events', label: 'Events' },
  { url: '/careers', label: 'Careers' },
  { url: '/projects', label: 'Projects' },
  { url: '/case-studies', label: 'Case Studies' },
  { url: '/resource-centre', label: 'Resource Centre' },
  { url: '/facilities', label: 'Facilities' },
  { url: '/associations', label: 'Associations' },
  { url: '/invest', label: 'Invest with NCS' },
  { url: '/faqs', label: 'FAQs' },
  { url: '/contact-us', label: 'Contact Us' },
  { url: '/team', label: 'Team' },
]

const internalPageOptions = computed(() => [
  ...STATIC_SITE_PAGES.map(page => ({ ...page, group: 'Site Pages' })),
  ...posts.value.filter(p => p.slug).map(p => ({ url: `/news/${p.slug}`, label: p.title, group: 'News Posts' })),
  ...pages.value.filter(p => p.slug).map(p => ({ url: `/pages/${p.slug}`, label: p.title, group: 'Pages' })),
])

const menusDirty = computed(() => {
  if (menusSavedSnapshot.value === null) return false
  return JSON.stringify([toCmsMenuItems(mainMenuTree.value), toCmsMenuItems(footerMenuTree.value)]) !== menusSavedSnapshot.value
})

function snapshotMenus() {
  menusSavedSnapshot.value = JSON.stringify([toCmsMenuItems(mainMenuTree.value), toCmsMenuItems(footerMenuTree.value)])
}

const postForm = reactive({ id:'', title:'', slug:'', category:'news', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
const pageForm = reactive({ id:'', title:'', slug:'', category:'page', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
const projectForm = reactive({ id:'', title:'', slug:'', category:'project', category_tag:'', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
const caseStudyForm = reactive({ id:'', title:'', slug:'', category:'case_study', category_tag:'', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
const blogCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true })
const projectCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'project' })
const caseStudyCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'case_study' })
const faqCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'faq' })
const resourceCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'resource' })
const careerCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'career' })
const teamDepartmentForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'team_department' })
const facilityCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'facility' })
const eventCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'event' })
const investCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'investment' })
const federationCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'federation' })
const eventForm = reactive({ id:'', title:'', slug:'', category:'', location:'', event_date:'', end_date:'', cover_image_url:'', description:'', status:'published' })
const facilityForm = reactive({ id:'', name:'', slug:'', category:'', description:'', image_url:'', sort_order:0, is_active:true })
const associationForm = reactive({ id:'', name:'', slug:'', category:'', president:'', secretary:'', phone:'', website_url:'', description:'', logo_url:'', sort_order:0, is_active:true })
const factForm = reactive({ id:'', label:'', value:'', icon:'icofont-chart-growth', sort_order:0, is_active:true })
const faqForm = reactive({ id:'', question:'', answer:'', category:'General', sort_order:0, is_active:true })
const resourceForm = reactive({ id:'', title:'', category:'Guidelines', file_url:'', description:'', sort_order:0, is_active:true })
const careerForm = reactive({ id:'', title:'', department:'', department_id:'', location:'', job_type:'full_time', category:'jobs', description:'', requirements:'', salary_range:'', status:'draft', deadline_at:'' })
const investForm = reactive({ id:'', title:'', subtitle:'', content:'', image_url:'', sort_order:0, is_active:true })
const teamForm = reactive({ id:'', full_name:'', designation:'', department_id:'', image_url:'', bio:'', sort_order:0, is_active:true })
const roleForm = reactive({ id:'', name:'', description:'', is_system:false })
const userForm = reactive({ id:'', first_name:'', last_name:'', email:'', phone:'', password:'' })

const blogCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const projectCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const caseStudyCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const faqCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const resourceCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const careerCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const teamDepartmentFields = fields(['name','slug','sort_order','is_active'], ['description'])
const facilityCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const eventCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const investCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const federationCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
// Renders the given field as a category dropdown sourced from that content
// type's own scoped category list, instead of a free-text input that has no
// connection to the category management screens for that type.
function withCategoryOptions(fieldDefs, categoriesRef) {
  return computed(() => fieldDefs.map(f => f.name === 'category' ? { ...f, type: 'select', options: categoriesRef.value } : f))
}
const facilityFields = withCategoryOptions(fields(['name','slug','category','image_url','sort_order','is_active'], ['description']), facilityCategories)
const associationFields = withCategoryOptions(fields(['name','slug','category','president','secretary','phone','website_url','logo_url','sort_order','is_active'], ['description']), federationCategories)
const factFields = fields(['label','value','sort_order','is_active'])
const resourceFields = withCategoryOptions(fields(['title','category','file_url','sort_order','is_active'], ['description']), resourceCategories)
const investFields = fields(['title','subtitle','image_url','sort_order','is_active'], ['content'])
const teamFields = fields(['full_name','designation','image_url','sort_order','is_active'], ['bio'])

const overviewCards = computed(() => [
  { label:'Posts', value:posts.value.length },
  { label:'Pages', value:pages.value.length },
  { label:'Careers', value:careers.value.length },
  { label:'Resources', value:resources.value.length },
  { label:'Roles', value:roles.value.length },
  { label:'Audit events', value:auditMeta.value.total || auditLogs.value.length },
  { label:'Events', value:events.value.length },
  { label:'Facilities', value:facilities.value.length },
  { label:'Associations', value:associations.value.length },
  { label:'Facts', value:facts.value.length },
  { label:'FAQs', value:faqs.value.length },
])
const totalContentItems = computed(() => posts.value.length + pages.value.length + projects.value.length + caseStudies.value.length + resources.value.length + careers.value.length)
const dashboardStatCards = computed(() => [
  { label:'Page Visits', value:formatNumber(analytics.value?.total_views || 0), trend:`${formatNumber(analytics.value?.unique_visitors || 0)}`, caption:'unique visitors', trendClass:'col-green', icon:'icofont-eye-alt col-blue cms-stat-icon' },
  { label:'Active Content', value:formatNumber(totalContentItems.value), trend:'08%', caption:'Increase', trendClass:'col-green', icon:'icofont-page col-green cms-stat-icon' },
  { label:'Engagement', value:analytics.value ? `${formatNumber(Math.round(analytics.value.average_session_seconds || 0))}s` : '0s', trend:`${analytics.value?.bounce_rate || 0}%`, caption:'bounce rate', trendClass:'col-orange', icon:'icofont-chart-growth col-orange cms-stat-icon' },
  { label:'Audit Events', value:formatNumber(auditMeta.value.total || auditLogs.value.length), trend:'16%', caption:'Tracked', trendClass:'col-green', icon:'icofont-shield-alt col-cyan cms-stat-icon' },
])
const visitSeries = computed(() => {
  if (analytics.value?.timeline?.length) {
    const points = analytics.value.timeline.slice(-7)
    const max = Math.max(1, ...points.map(point => point.views || 0))
    return points.map(point => ({
      label: point.label,
      value: formatNumber(point.views || 0),
      percent: Math.max(8, Math.round(((point.views || 0) / max) * 100)),
    }))
  }
  const values = [320, 460, 390, 680, 740, 610, 830].map((value, index) => value + Math.min(260, totalContentItems.value * (index + 1)))
  const max = Math.max(...values)
  return ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((label, index) => ({
    label,
    value: formatNumber(values[index]),
    percent: Math.max(18, Math.round((values[index] / max) * 100)),
  }))
})
const websiteKpis = computed(() => [
  { label:'Total visitors', value:formatNumber(analytics.value?.unique_visitors || 0), trend:'+live', trendClass:'col-green' },
  { label:'Published items', value:formatNumber(posts.value.length + pages.value.length + projects.value.length), trend:'+16%', trendClass:'col-green' },
  { label:'Top city', value:analytics.value?.top_cities?.[0]?.label || 'Unknown', trend:`${analytics.value?.top_cities?.[0]?.share || 0}%`, trendClass:'col-green' },
  { label:'Pending comments', value:formatNumber(comments.value.length), trend:comments.value.length ? '+04%' : '0%', trendClass:comments.value.length ? 'col-orange' : 'col-green' },
  { label:'High risk events', value:formatNumber(auditHighThreatCount.value), trend:auditHighThreatCount.value ? '-02%' : '0%', trendClass:auditHighThreatCount.value ? 'text-danger' : 'col-green' },
])
const platformStats = computed(() => (analytics.value?.device_types?.length ? analytics.value.device_types : [
  { label:'Desktop', share:0 }, { label:'Mobile', share:0 }, { label:'Tablet', share:0 },
]).map((item, index) => ({ label:item.label, value:Math.round(item.share || 0), color:['bg-primary','bg-success','bg-warning','bg-info'][index % 4] })))
const trafficSources = computed(() => (analytics.value?.acquisition_channels?.length ? analytics.value.acquisition_channels : [
  { label:'Direct', share:0 },
]).map(item => ({ label:item.label, value:Math.round(item.share || 0), icon:trafficIcon(item.label) })))
const contentMix = computed(() => {
  const total = Math.max(1, totalContentItems.value)
  const postShare = Math.round((posts.value.length / total) * 100) || 44
  const resourceShare = Math.round((resources.value.length / total) * 100) || 28
  return [
    { label:'Posts', value:postShare, color:'bg-primary' },
    { label:'Resources', value:resourceShare, color:'bg-success' },
    { label:'Pages', value:Math.max(0, 100 - postShare - resourceShare), color:'bg-warning' },
  ]
})
const dashboardTasks = computed(() => [
  { name:'Publish homepage updates', avatars:userAvatars(1, 2, 3), extra:2, progress:Math.min(95, 45 + posts.value.length), due:'2026-07-08', priority:'High', priorityClass:'badge-danger', barClass:'bg-danger', target:'homepage' },
  { name:'Review resource centre uploads', avatars:userAvatars(4, 5), extra:3, progress:Math.min(90, 35 + resources.value.length * 8), due:'2026-07-10', priority:'Average', priorityClass:'badge-info', barClass:'bg-cyan', target:'resources' },
  { name:'Moderate public comments', avatars:userAvatars(6, 7, 8), extra:1, progress:comments.value.length ? 50 : 100, due:'2026-07-11', priority:comments.value.length ? 'High' : 'Low', priorityClass:comments.value.length ? 'badge-danger' : 'badge-success', barClass:comments.value.length ? 'bg-warning' : 'bg-success', target:'comments' },
  { name:'Audit access and role changes', avatars:userAvatars(9, 10), extra:4, progress:Math.min(100, 60 + auditLogs.value.length), due:'2026-07-12', priority:'Low', priorityClass:'badge-success', barClass:'bg-success', target:'audit-logs' },
])
const websiteActivity = computed(() => [
  { id:'89754', type:'Content', badgeClass:'badge-success', title:'Homepage hero requires review', summary:'Hero slideshow and homepage sections are ready for CMS approval.', owner:currentUserName.value, time:'1 day ago', avatar:'/otika-assets/img/users/user-1.png', target:'slideshow-manager' },
  { id:'89755', type:'Analytics', badgeClass:'badge-info', title:'Traffic source report updated', summary:'Website KPI platform and source split has been refreshed.', owner:'Analytics Monitor', time:'2 days ago', avatar:'/otika-assets/img/users/user-2.png', target:'overview' },
  { id:'89756', type:'Audit', badgeClass:'badge-warning', title:'Audit trail needs inspection', summary:'Recent admin requests are available in the audit log module.', owner:'Audit Monitor', time:'3 days ago', avatar:'/otika-assets/img/users/user-5.png', target:'audit-logs' },
])
const analyticsSnapshot = computed(() => [
  { label:'Bounce Rate', value:`${analytics.value?.bounce_rate || 0}%`, icon:'icofont-exit col-orange' },
  { label:'Avg. Session', value:formatDuration(analytics.value?.average_session_seconds || 0), icon:'icofont-clock-time col-green' },
  { label:'Top Country', value:analytics.value?.top_countries?.[0]?.label || 'Unknown', icon:'icofont-globe col-blue' },
  { label:'Top Page', value:analytics.value?.top_pages?.[0]?.label || '/', icon:'icofont-page col-cyan' },
])
const departmentOptions = computed(() => institutionalDepartments.value.length ? institutionalDepartments.value : teamDepartments.value)
const currentRoleNames = computed(() => normalizeRoleNames(storedUser.roles || storedUser.role || []))
const currentPermissionNames = computed(() => {
  const roleNames = currentRoleNames.value
  if (roleNames.includes('super_admin')) return new Set(['*'])
  const names = new Set()
  collectPermissionNames(storedUser.permissions, names)
  collectPermissionNames((storedUser.roles || []).flatMap(role => role.permissions || []), names)
  for (const role of roles.value || []) {
    if (roleNames.includes(role.name)) collectPermissionNames(role.permissions, names)
  }
  for (const roleName of roleNames) {
    for (const permission of builtinRolePermissions(roleName)) names.add(permission)
  }
  return names
})
const permissionGroups = computed(() => {
  const groups = new Map()
  for (const permission of Array.isArray(permissions.value) ? permissions.value : []) {
    const resource = permission.resource || 'general'
    if (!groups.has(resource)) groups.set(resource, [])
    groups.get(resource).push(permission)
  }
  return Array.from(groups, ([resource, groupPermissions]) => ({
    resource,
    label: resourceLabels[resource] || titleize(resource),
    permissions: groupPermissions.sort((a, b) => actionOrder.indexOf(a.action) - actionOrder.indexOf(b.action) || a.action.localeCompare(b.action)),
  })).sort((a, b) => a.label.localeCompare(b.label))
})
const auditRiskCount = computed(() => auditLogs.value.filter(log => Number(log.response_code) >= 400 || log.event_status === 'failure' || log.anomaly_detected).length)
const auditHighThreatCount = computed(() => auditLogs.value.filter(log => Number(log.threat_score || 0) >= 70 || log.severity_level === 'high' || log.severity_level === 'critical').length)
const sampleAuditLogs = [
  { id:'preview-login', user_name:'cms-admin', action:'preview_login', method:'POST', endpoint:'/api/v1/auth/login', ip_address:'127.0.0.1', response_code:200, response_time_ms:31, event_status:'success', severity_level:'low', created_at:new Date().toISOString(), details:{ mode:'local-preview' } },
  { id:'preview-audit', user_name:'audit-monitor', action:'view_audit_logs', method:'GET', endpoint:'/api/v1/admin/audit-logs', ip_address:'127.0.0.1', response_code:200, response_time_ms:18, event_status:'success', severity_level:'low', created_at:new Date().toISOString(), details:{ source:'sample-data' } },
]

onMounted(() => {
  if (!localStorage.getItem('ncsms_access_token')) {
    router.replace('/login')
    return
  }
  const params = new URLSearchParams(window.location.search)
  // Returning from the Google Drive OAuth consent screen lands back on this
  // same route with ?code=&state= — land on the Storage Settings tab so the
  // admin sees the connect flow finish instead of the Overview tab.
  if (params.get('code')) {
    active.value = 'storage'
  } else {
    // Restore the active panel from the URL (?section=) so a browser refresh
    // keeps the admin on the same section instead of resetting to Overview.
    const requested = params.get('section')
    if (requested && requested !== active.value && sections.some(s => s.id === requested) && canAccessSection(requested)) {
      active.value = requested
    }
  }
  loadAll()
})

// Mirror the active panel into the URL (?section=) so it survives a refresh and
// is shareable. router.replace avoids polluting history; existing query params
// (e.g. the OAuth code/state) are preserved. Overview keeps a clean URL.
watch(active, (id) => {
  const currentQuery = router.currentRoute.value.query
  if ((currentQuery.section || 'overview') === id) return
  const query = { ...currentQuery }
  if (id === 'overview') delete query.section
  else query.section = id
  router.replace({ query }).catch(() => {})
})

function fields(short = [], long = []) {
  return [...short.map(name => ({
    name,
    type: name === 'is_active'
      ? 'checkbox'
      : (name === 'sort_order' || name.endsWith('_order') ? 'number' : 'input'),
  })), ...long.map(name => ({ name, type: 'textarea' }))]
}
function formatNumber(value) {
  return new Intl.NumberFormat('en-UG').format(Number(value || 0))
}
function userAvatars(...ids) {
  return ids.map(id => `/otika-assets/img/users/user-${id}.png`)
}
function emailLocalPart(email) {
  return typeof email === 'string' && email.includes('@') ? email.split('@')[0] : ''
}
function data(res) { return res?.data?.data ?? {} }
function isLocalPreviewSession() {
  return localStorage.getItem('ncsms_access_token') === 'local-cms-preview-token'
}
async function canReachApi(force = false) {
  if (isLocalPreviewSession()) return false
  if (!force && apiAvailable.value !== null) return apiAvailable.value
  let timer
  try {
    const controller = new AbortController()
    timer = window.setTimeout(() => controller.abort(), 1800)
    const res = await fetch(`${API_BASE_URL}/health`, { signal: controller.signal, cache: 'no-store' })
    apiAvailable.value = res.ok
  } catch {
    apiAvailable.value = false
  } finally {
    if (timer) window.clearTimeout(timer)
  }
  return apiAvailable.value
}
function listData(res, fallback = []) {
  const value = data(res)
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.items)) return value.items
  return fallback
}
function pageMeta(res) { return res?.data?.meta || {} }
function setMsg(text) { message.value = text; error.value = ''; setTimeout(() => { message.value = '' }, 2500) }
function setErr(err) { error.value = err.response?.data?.error?.message || err.message || 'Action failed' }
async function confirmAction(title, text = 'This cannot be undone.') {
  const result = await Swal.fire({
    title,
    text,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Yes, continue',
    cancelButtonText: 'Cancel',
    confirmButtonColor: '#6777ef',
    cancelButtonColor: '#fc544b',
    reverseButtons: true,
  })
  return result.isConfirmed
}
function copyInto(target, source) { Object.keys(target).forEach(k => { target[k] = source?.[k] ?? (typeof target[k] === 'boolean' ? false : '') }) }
function clean(payload) {
  // Fields the backend decodes as Go numbers. HTML inputs yield strings, so we
  // coerce here — sending "2" where an int is expected fails json decode with
  // "Invalid JSON body" (400). This is the single serialization point for every
  // saveEntity create/update, so coercing here fixes them all at once.
  const isNumericField = k => k === 'sort_order' || k.endsWith('_order') || k === 'transition_duration' || k === 'autoplay_speed'
  return Object.fromEntries(
    Object.entries(payload)
      .filter(([k, v]) => k !== 'id' && !(k.endsWith('_id') && v === ''))
      .map(([k, v]) => {
        if (isNumericField(k)) { const n = Number(v); return [k, Number.isFinite(n) ? n : 0] }
        return [k, v]
      }),
  )
}
function normalizeSlug(value) { return String(value || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '') }
function titleize(value) { return String(value || '').replace(/[_:-]+/g, ' ').replace(/\b\w/g, char => char.toUpperCase()) }
function normalizeRoleNames(input) {
  const rolesList = Array.isArray(input) ? input : [input]
  return rolesList.map(role => typeof role === 'string' ? role : role?.name).filter(Boolean)
}
function collectPermissionNames(input, target) {
  for (const permission of Array.isArray(input) ? input : []) {
    if (typeof permission === 'string') target.add(permission)
    else if (permission?.name) target.add(permission.name)
  }
}
function builtinRolePermissions(roleName) {
  if (roleName === 'admin') return ['dashboard:read','audit:read','users:read','users:create','users:update','users:delete','users:activate','users:reset_password','users:roles','roles:read','roles:assign','cms:read','cms:write','cms:delete','menus:read','menus:update','settings:read','settings:update','storage:read','storage:update','messages:read','messages:update','messages:delete','notifications:read','notifications:update','notifications:delete']
  if (roleName === 'content_manager') return ['cms:read','cms:write','cms:delete','messages:read','messages:update','notifications:read','notifications:update']
  return []
}
function hasPermission(permissionName) {
  const names = currentPermissionNames.value
  if (names.has('*') || names.has(permissionName)) return true
  const [resource, action] = String(permissionName).split(':')
  if (resource && action && names.has(`${resource}:write`) && ['create', 'update'].includes(action)) return true
  if (resource !== 'users' && resource !== 'roles' && names.has('cms:read') && action === 'read') return true
  if (resource !== 'users' && resource !== 'roles' && names.has('cms:write') && ['create', 'update'].includes(action)) return true
  if (resource !== 'users' && resource !== 'roles' && names.has('cms:delete') && action === 'delete') return true
  return false
}
function canAccessSection(id) {
  const required = sectionPermissionMap[id] || []
  return !required.length || required.some(hasPermission)
}
function visibleItems(items) { return items.filter(item => canAccessSection(item.id)) }
function syncCareerDepartment() {
  const selected = departmentOptions.value.find(dept => dept.id === careerForm.department_id)
  if (selected) careerForm.department = selected.name
  else if (!careerForm.department_id) careerForm.department = ''
}
function selectSection(id) {
  if (!canAccessSection(id)) {
    setErr(new Error('You do not have permission to access that CMS section.'))
    return
  }
  active.value = id
  messagesOpen.value = false
  notificationsOpen.value = false
  profileOpen.value = false
  if (id === 'audit-logs' && !auditLogs.value.length) loadAuditLogs()
  if (id === 'messages' && !contactMessages.value.length) loadContactMessages()
  if (id === 'notifications' && !cmsNotifications.value.length) loadNotifications()
  if (['roles', 'manage-roles'].includes(id) && (!roles.value.length || !permissions.value.length)) loadRoles()
  if (['users', 'manage-users'].includes(id) && !users.value.length) loadUsers()
  if (id === 'analytics' && !analytics.value) loadAnalytics()
}

function trafficIcon(label) {
  const value = String(label || '').toLowerCase()
  if (value.includes('search')) return 'fas fa-search'
  if (value.includes('social')) return 'fas fa-share-alt'
  if (value.includes('direct')) return 'fas fa-link'
  return 'fas fa-project-diagram'
}

function formatDuration(seconds) {
  const total = Math.round(Number(seconds) || 0)
  const minutes = Math.floor(total / 60)
  const remainder = total % 60
  return minutes ? `${minutes}m ${remainder}s` : `${remainder}s`
}

function analyticsRows(key) {
  const rows = analytics.value?.[key]
  return Array.isArray(rows) && rows.length ? rows : [{ label:'No data yet', value:0, share:0 }]
}

function groupHasActive(items) {
  return items.some(item => item.id === active.value)
}

function runGlobalSearch() {
  const query = globalSearch.value.trim().toLowerCase()
  if (!query) return
  const match = sections.find(section => section.label.toLowerCase().includes(query) || section.id.includes(query))
  if (match) {
    selectSection(match.id)
    return
  }
  setErr(new Error(`No CMS module found for "${globalSearch.value}".`))
}

function openTopMessage(item) {
  selectSection(item.target || 'overview')
  messagesOpen.value = false
}

function openNotification(item) {
  selectSection(item.target || 'overview')
  notificationsOpen.value = false
}

function mergeHomepageDefaults(target, defaults = homepageDefaults) {
  for (const [key, value] of Object.entries(defaults)) {
    if (Array.isArray(value)) {
      if (!Array.isArray(target[key]) || !target[key].length) target[key] = JSON.parse(JSON.stringify(value))
    } else if (value && typeof value === 'object') {
      if (!target[key] || typeof target[key] !== 'object' || Array.isArray(target[key])) target[key] = {}
      mergeHomepageDefaults(target[key], value)
    } else if (target[key] === undefined || target[key] === null || target[key] === '') {
      target[key] = value
    }
  }
  // Root-level-only invariants. These must NOT run on the nested recursive
  // calls (e.g. target = homepage.about / homepage.topbar.social), where
  // target.topbar / target.milestones don't exist — accessing
  // target.topbar.marquee there throws "Cannot read properties of undefined".
  if (defaults === homepageDefaults) {
    if (!Array.isArray(target.milestones)) target.milestones = JSON.parse(JSON.stringify(homepageDefaults.milestones))
    while (target.milestones.length < 4) target.milestones.push(JSON.parse(JSON.stringify(homepageDefaults.milestones[target.milestones.length] || homepageDefaults.milestones[0])))
    if (!target.topbar || typeof target.topbar !== 'object') target.topbar = JSON.parse(JSON.stringify(homepageDefaults.topbar))
    if (!Array.isArray(target.topbar.marquee)) target.topbar.marquee = JSON.parse(JSON.stringify(homepageDefaults.topbar.marquee))
  }
}

function addHomepageValue() {
  homepage.values.push({ title:'New value', text:'', icon:'icofont-star', featured:false })
}

function addSportsCounter() {
  homepage.milestones.push({ value:'0+', label:'New counter', description:'', icon:'icofont-chart-growth' })
}

async function loadAll() {
  try {
    apiAvailable.value = null
    if (isLocalPreviewSession()) {
      error.value = 'CMS preview mode: backend API is offline, so live content is paused.'
      snapshotMenus()
      return
    }
    if (!(await canReachApi())) {
      error.value = 'Backend API is not reachable on port 9080. Start the native Go API to load live CMS data.'
      snapshotMenus()
      return
    }
    const results = await Promise.allSettled([
      cms.getSettings('homepage'), cms.getSettings('contact'), cms.getSettings('footer'), cms.getMenu('main'), cms.getMenu('footer'), cms.adminGetSlideshow('homepage-hero'),
      cms.adminListPosts({ per_page:200 }), cms.adminListEvents({ per_page:200 }), cms.adminListSlides(), cms.adminListFacilities(), cms.adminListAssociations(),
      cms.adminListFunFacts(), cms.adminListFAQs(), cms.adminListResources({ per_page:200 }), cms.adminListBlogCategories(), cms.adminListComments({ status: commentStatus.value, per_page:50 }),
      cms.adminListPosts({ category:'page', per_page:200 }), cms.adminListPosts({ category:'project', per_page:200 }), cms.adminListPosts({ category:'case_study', per_page:200 }),
      cms.adminListCareers({ per_page:200 }), cms.adminListInvest(), cms.adminListTeam(),
      cms.adminListRoles(), cms.adminListPermissions(), cms.adminListAuditLogs({ page:1, per_page:50 }),
      cms.adminListProjectCategories(), cms.adminListCaseStudyCategories(),
      cms.adminListFAQCategories(), cms.adminListResourceCategories(), cms.adminListCareerCategories(), cms.adminListTeamDepartments(),
      cms.adminListFacilityCategories(), cms.adminListEventCategories(), cms.adminListInvestCategories(), cms.adminListFederationCategories(), cms.adminListNewsletterSubscribers({ per_page:200 }),
      cms.listInstitutionalDepartments(), cms.adminGetAnalytics({ days: 30 }),
      cms.listMessages({ status:'', per_page:100 }), cms.listNotifications({ status:'', per_page:100 }),
    ])
    Object.assign(homepage, data(results[0].value)?.value || {})
    mergeHomepageDefaults(homepage)
    Object.assign(contact, data(results[1].value)?.value || {})
    Object.assign(footer, data(results[2].value)?.value || {})
    mainMenuTree.value = normalizeMenuTree(data(results[3].value)?.items || [])
    footerMenuTree.value = normalizeMenuTree(data(results[4].value)?.items || [])
    snapshotMenus()
    const show = data(results[5].value) || {}
    Object.assign(slideshow, show)
    posts.value = listData(results[6].value)
    events.value = listData(results[7].value)
    slides.value = Array.isArray(show.slides) ? show.slides : listData(results[8].value)
    facilities.value = listData(results[9].value)
    associations.value = listData(results[10].value)
    facts.value = listData(results[11].value)
    faqs.value = listData(results[12].value)
    resources.value = listData(results[13].value)
    blogCategories.value = listData(results[14].value)
    comments.value = listData(results[15].value)
    pages.value = listData(results[16].value)
    projects.value = listData(results[17].value)
    caseStudies.value = listData(results[18].value)
    careers.value = listData(results[19].value)
    investItems.value = listData(results[20].value)
    teamMembers.value = listData(results[21].value)
    roles.value = listData(results[22].value)
    permissions.value = listData(results[23].value)
    auditLogs.value = listData(results[24].value)
    auditMeta.value = pageMeta(results[24].value)
    projectCategories.value = listData(results[25].value)
    caseStudyCategories.value = listData(results[26].value)
    faqCategories.value = listData(results[27].value)
    resourceCategories.value = listData(results[28].value)
    careerCategories.value = listData(results[29].value)
    teamDepartments.value = listData(results[30].value)
    facilityCategories.value = listData(results[31].value)
    eventCategories.value = listData(results[32].value)
    investCategories.value = listData(results[33].value)
    federationCategories.value = listData(results[34].value)
    newsletterSubscribers.value = listData(results[35].value, sampleNewsletterSubscribers())
    institutionalDepartments.value = listData(results[36].value)
    analytics.value = data(results[37].value) || null
    contactMessages.value = listData(results[38].value)
    cmsNotifications.value = listData(results[39].value)
  } catch (err) { setErr(err) }
}

async function loadAnalytics() {
  try {
    const res = await cms.adminGetAnalytics({ days: 30 })
    analytics.value = data(res)
  } catch (err) { setErr(err) }
}

async function saveHomepage() {
  try {
    mergeHomepageDefaults(homepage)
    const payload = JSON.parse(JSON.stringify(homepage))
    await cms.adminUpdateSettings('homepage', payload)
    if (payload.topbar) {
      await Promise.allSettled([
        cms.adminUpdateSettings('header', { marquee: payload.topbar.marquee, webmail_url: payload.topbar.webmail_url, social: payload.topbar.social }),
        cms.adminUpdateSettings('contact', { ...JSON.parse(JSON.stringify(contact)), social: { ...(contact.social || {}), ...(payload.topbar.social || {}) } }),
      ])
    }
    setMsg('Homepage saved')
  } catch (err) { setErr(err) }
}
async function saveMenus() {
  savingMenus.value = true
  try {
    await cms.adminUpdateMenu('main', toCmsMenuItems(mainMenuTree.value))
    await cms.adminUpdateMenu('footer', toCmsMenuItems(footerMenuTree.value))
    snapshotMenus()
    setMsg('Menus saved')
  } catch (err) {
    setErr(err)
  } finally {
    savingMenus.value = false
  }
}
async function discardMenuChanges() {
  try {
    const [mainRes, footerRes] = await Promise.all([cms.getMenu('main'), cms.getMenu('footer')])
    mainMenuTree.value = normalizeMenuTree(data(mainRes)?.items || [])
    footerMenuTree.value = normalizeMenuTree(data(footerRes)?.items || [])
    snapshotMenus()
    setMsg('Menu changes discarded')
  } catch (err) { setErr(err) }
}
async function saveSettings() { try { await cms.adminUpdateSettings('contact', { ...contact }); await cms.adminUpdateSettings('footer', { ...footer }); setMsg('Settings saved') } catch (err) { setErr(err) } }
async function createSlide() { try { await cms.adminCreateSlide({ title:'New slide', subtitle:'National Council of Sports', description:'', image_url:'', button_text:'Learn More', button_url:'/', sort_order:slides.value.length + 1, is_active:true }); await loadAll(); setMsg('Slide added') } catch (err) { setErr(err) } }
function editSlide(item) { active.value = 'homepage'; Object.assign(homepage, { hero_quick_edit: item.title }) }
async function removeSlide(item) {
  if (!(await confirmAction('Delete this slide?', 'The slide will be removed from the public slideshow.'))) return
  try {
    await cms.adminDeleteSlide(item.id)
    await loadAll()
    setMsg('Slide deleted')
  } catch (err) { setErr(err) }
}

async function savePost() { await saveEntity(postForm, cms.adminCreatePost, cms.adminUpdatePost, 'Post saved') }
async function savePage() {
  pageForm.category = 'page'
  pageForm.slug = normalizeSlug(pageForm.slug || pageForm.title)
  await saveEntity(pageForm, cms.adminCreatePost, cms.adminUpdatePost, 'Page saved')
}
async function saveProject() { projectForm.category = 'project'; await saveEntity(projectForm, cms.adminCreatePost, cms.adminUpdatePost, 'Project saved') }
async function saveCaseStudy() { caseStudyForm.category = 'case_study'; await saveEntity(caseStudyForm, cms.adminCreatePost, cms.adminUpdatePost, 'Case study saved') }
async function saveBlogCategory() { await saveEntity(blogCategoryForm, cms.adminCreateBlogCategory, cms.adminUpdateBlogCategory, 'Category saved') }
async function saveProjectCategory() {
  projectCategoryForm.content_type = 'project'
  await saveEntity(projectCategoryForm, cms.adminCreateProjectCategory, cms.adminUpdateProjectCategory, 'Project category saved')
}
async function saveCaseStudyCategory() {
  caseStudyCategoryForm.content_type = 'case_study'
  await saveEntity(caseStudyCategoryForm, cms.adminCreateCaseStudyCategory, cms.adminUpdateCaseStudyCategory, 'Case study category saved')
}
async function saveFAQCategory() {
  faqCategoryForm.content_type = 'faq'
  await saveEntity(faqCategoryForm, cms.adminCreateFAQCategory, cms.adminUpdateFAQCategory, 'FAQ category saved')
}
async function saveResourceCategory() {
  resourceCategoryForm.content_type = 'resource'
  await saveEntity(resourceCategoryForm, cms.adminCreateResourceCategory, cms.adminUpdateResourceCategory, 'Resource category saved')
}
async function saveCareerCategory() {
  careerCategoryForm.content_type = 'career'
  await saveEntity(careerCategoryForm, cms.adminCreateCareerCategory, cms.adminUpdateCareerCategory, 'Career category saved')
}
async function saveTeamDepartment() {
  teamDepartmentForm.content_type = 'team_department'
  await saveEntity(teamDepartmentForm, cms.adminCreateTeamDepartment, cms.adminUpdateTeamDepartment, 'Department saved')
}
async function saveFacilityCategory() {
  facilityCategoryForm.content_type = 'facility'
  await saveEntity(facilityCategoryForm, cms.adminCreateFacilityCategory, cms.adminUpdateFacilityCategory, 'Facility category saved')
}
async function saveEventCategory() {
  eventCategoryForm.content_type = 'event'
  await saveEntity(eventCategoryForm, cms.adminCreateEventCategory, cms.adminUpdateEventCategory, 'Event category saved')
}
async function saveInvestCategory() {
  investCategoryForm.content_type = 'investment'
  await saveEntity(investCategoryForm, cms.adminCreateInvestCategory, cms.adminUpdateInvestCategory, 'Investment category saved')
}
async function saveFederationCategory() {
  federationCategoryForm.content_type = 'federation'
  await saveEntity(federationCategoryForm, cms.adminCreateFederationCategory, cms.adminUpdateFederationCategory, 'Federation category saved')
}
async function saveEvent() { await saveEntity(eventForm, cms.adminCreateEvent, cms.adminUpdateEvent, 'Event saved') }
async function saveFacility() { await saveEntity(facilityForm, cms.adminCreateFacility, cms.adminUpdateFacility, 'Facility saved') }
async function saveAssociation() { await saveEntity(associationForm, cms.adminCreateAssociation, cms.adminUpdateAssociation, 'Association saved') }
async function saveFact() { await saveEntity(factForm, cms.adminCreateFunFact, cms.adminUpdateFunFact, 'Fact saved') }
async function saveFAQ() { await saveEntity(faqForm, cms.adminCreateFAQ, cms.adminUpdateFAQ, 'FAQ saved') }
async function saveResource() { await saveEntity(resourceForm, cms.adminCreateResource, cms.adminUpdateResource, 'Resource saved') }
async function saveCareer() { syncCareerDepartment(); await saveEntity(careerForm, cms.adminCreateCareer, cms.adminUpdateCareer, 'Career post saved') }
async function saveInvest() { await saveEntity(investForm, cms.adminCreateInvest, cms.adminUpdateInvest, 'Investment item saved') }
async function saveTeam() { await saveEntity(teamForm, cms.adminCreateTeam, cms.adminUpdateTeam, 'Team member saved') }
async function saveEntity(form, createFn, updateFn, ok) { try { form.id ? await updateFn(form.id, clean(form)) : await createFn(clean(form)); await loadAll(); setMsg(ok) } catch (err) { setErr(err) } }

function editPost(item) { copyInto(postForm, item); active.value = 'create-post' }
function editPage(item) { copyInto(pageForm, item); pageForm.category = 'page'; active.value = 'static-pages' }
function editProject(item) { copyInto(projectForm, item); projectForm.category = 'project'; active.value = 'projects' }
function editCaseStudy(item) { copyInto(caseStudyForm, item); caseStudyForm.category = 'case_study'; active.value = 'case-studies' }
function editBlogCategory(item) { copyInto(blogCategoryForm, item); active.value = 'create-blog-categories' }
function editProjectCategory(item) { copyInto(projectCategoryForm, item); projectCategoryForm.content_type = 'project'; active.value = 'create-project-categories' }
function editCaseStudyCategory(item) { copyInto(caseStudyCategoryForm, item); caseStudyCategoryForm.content_type = 'case_study'; active.value = 'create-case-study-categories' }
function editFAQCategory(item) { copyInto(faqCategoryForm, item); faqCategoryForm.content_type = 'faq'; active.value = 'create-faq-categories' }
function editResourceCategory(item) { copyInto(resourceCategoryForm, item); resourceCategoryForm.content_type = 'resource'; active.value = 'create-resource-categories' }
function editCareerCategory(item) { copyInto(careerCategoryForm, item); careerCategoryForm.content_type = 'career'; active.value = 'create-career-categories' }
function editTeamDepartment(item) { copyInto(teamDepartmentForm, item); teamDepartmentForm.content_type = 'team_department'; active.value = 'create-team-departments' }
function editFacilityCategory(item) { copyInto(facilityCategoryForm, item); facilityCategoryForm.content_type = 'facility'; active.value = 'create-facility-categories' }
function editEventCategory(item) { copyInto(eventCategoryForm, item); eventCategoryForm.content_type = 'event'; active.value = 'create-event-categories' }
function editInvestCategory(item) { copyInto(investCategoryForm, item); investCategoryForm.content_type = 'investment'; active.value = 'create-invest-categories' }
function editFederationCategory(item) { copyInto(federationCategoryForm, item); federationCategoryForm.content_type = 'federation'; active.value = 'create-federation-categories' }
function editEvent(item) {
  copyInto(eventForm, item)
  if (item.event_date) eventForm.event_date = item.event_date.slice(0, 16)
  if (item.end_date) eventForm.end_date = item.end_date.slice(0, 16)
  active.value = 'events'
}
function editFacility(item) { copyInto(facilityForm, item); active.value = 'facilities' }
function editAssociation(item) { copyInto(associationForm, item); active.value = 'associations' }
function editFact(item) { copyInto(factForm, item); active.value = 'facts' }
function editFAQ(item) { copyInto(faqForm, item); active.value = 'faqs' }
function editResource(item) { copyInto(resourceForm, item); active.value = 'resources' }
function editCareer(item) { copyInto(careerForm, item); if (item.deadline_at) careerForm.deadline_at = item.deadline_at; active.value = 'careers' }
function editInvest(item) { copyInto(investForm, item); active.value = 'invest' }
function editTeam(item) { copyInto(teamForm, item); active.value = 'team' }

async function removePost(item) { await removeEntity(item, cms.adminDeletePost) }
async function removeBlogCategory(item) { await removeEntity(item, cms.adminDeleteBlogCategory) }
async function removeProjectCategory(item) { await removeEntity(item, cms.adminDeleteProjectCategory) }
async function removeCaseStudyCategory(item) { await removeEntity(item, cms.adminDeleteCaseStudyCategory) }
async function removeFAQCategory(item) { await removeEntity(item, cms.adminDeleteFAQCategory) }
async function removeResourceCategory(item) { await removeEntity(item, cms.adminDeleteResourceCategory) }
async function removeCareerCategory(item) { await removeEntity(item, cms.adminDeleteCareerCategory) }
async function removeTeamDepartment(item) { await removeEntity(item, cms.adminDeleteTeamDepartment) }
async function removeFacilityCategory(item) { await removeEntity(item, cms.adminDeleteFacilityCategory) }
async function removeEventCategory(item) { await removeEntity(item, cms.adminDeleteEventCategory) }
async function removeInvestCategory(item) { await removeEntity(item, cms.adminDeleteInvestCategory) }
async function removeFederationCategory(item) { await removeEntity(item, cms.adminDeleteFederationCategory) }
async function removeEvent(item) { await removeEntity(item, cms.adminDeleteEvent) }
async function removeFacility(item) { await removeEntity(item, cms.adminDeleteFacility) }
async function removeAssociation(item) { await removeEntity(item, cms.adminDeleteAssociation) }
async function removeFact(item) { await removeEntity(item, cms.adminDeleteFunFact) }
async function removeFAQ(item) { await removeEntity(item, cms.adminDeleteFAQ) }
async function removeResource(item) { await removeEntity(item, cms.adminDeleteResource) }
async function removeCareer(item) { await removeEntity(item, cms.adminDeleteCareer) }
async function removeInvest(item) { await removeEntity(item, cms.adminDeleteInvest) }
async function removeTeam(item) { await removeEntity(item, cms.adminDeleteTeam) }
async function removeEntity(item, fn) {
  if (!(await confirmAction('Delete this item?'))) return
  try {
    await fn(item.id)
    await loadAll()
    setMsg('Deleted')
  } catch (err) { setErr(err) }
}

async function saveRole() {
  try {
    if (!(await canReachApi())) {
      error.value = 'Backend API is not reachable on port 9080. Start the native Go API to create or update live roles.'
      return
    }
    const payload = { name: roleForm.name.trim(), description: roleForm.description.trim() }
    const res = roleForm.id ? await cms.adminUpdateRole(roleForm.id, payload) : await cms.adminCreateRole(payload)
    const saved = data(res)
    await loadRoles()
    Object.assign(roleForm, { id: saved.id || roleForm.id, name: saved.name || payload.name, description: saved.description || payload.description, is_system: !!saved.is_system })
    if (roleForm.id) await selectRole(roleForm)
    setMsg(roleForm.id ? 'Role saved' : 'Role created')
  } catch (err) { setErr(err) }
}

async function loadRoles() {
  if (!(await canReachApi())) {
    const previewPerms = previewPermissions()
    roles.value = roles.value.length ? roles.value : [
      { id:'role_super_admin', name:'super_admin', description:'Full CMS access', is_system:true, permissions:previewPerms },
      { id:'role_content_manager', name:'content_manager', description:'Preview content manager role. Start the native Go API to manage live roles.', is_system:false, permissions:previewPerms.filter(permission => !['users','roles','audit','storage'].includes(permission.resource)) },
    ]
    permissions.value = permissions.value.length ? permissions.value : previewPerms
    error.value = 'Backend API is not reachable on port 9080. Role management is showing preview data.'
    return
  }
  const [roleRes, permRes] = await Promise.all([cms.adminListRoles(), cms.adminListPermissions()])
  roles.value = listData(roleRes)
  permissions.value = listData(permRes)
}

function previewPermissions() {
  const seen = new Set()
  return Object.values(sectionPermissionMap).flat().filter(name => {
    if (seen.has(name)) return false
    seen.add(name)
    return true
  }).map(name => {
    const [resource, action] = name.split(':')
    return {
      id: `preview-${resource}-${action}`,
      name,
      resource,
      action,
      description: `${titleize(action)} ${resourceLabels[resource] || titleize(resource)}`,
    }
  })
}

async function selectRole(role) {
  try {
    if (!(await canReachApi())) {
      Object.assign(roleForm, {
        id: role.id,
        name: role.name,
        description: role.description || '',
        is_system: !!role.is_system,
      })
      selectedPermissionIds.value = new Set((role.permissions || []).map(permission => permission.id))
      active.value = 'manage-roles'
      error.value = 'Backend API is not reachable on port 9080. Start the native Go API to customize live permissions.'
      return
    }
    const res = await cms.adminGetRole(role.id)
    const fullRole = data(res) || role
    Object.assign(roleForm, {
      id: fullRole.id,
      name: fullRole.name,
      description: fullRole.description || '',
      is_system: !!fullRole.is_system,
    })
    selectedPermissionIds.value = new Set((fullRole.permissions || []).map(permission => permission.id))
    active.value = 'manage-roles'
  } catch (err) { setErr(err) }
}

async function togglePermission(permission, checked) {
  if (!roleForm.id) return
  try {
    if (!(await canReachApi())) {
      const next = new Set(selectedPermissionIds.value)
      checked ? next.add(permission.id) : next.delete(permission.id)
      selectedPermissionIds.value = next
      setMsg(checked ? 'Preview permission selected' : 'Preview permission removed')
      return
    }
    checked
      ? await cms.adminAssignRolePermission(roleForm.id, permission.id)
      : await cms.adminRemoveRolePermission(roleForm.id, permission.id)
    const next = new Set(selectedPermissionIds.value)
    checked ? next.add(permission.id) : next.delete(permission.id)
    selectedPermissionIds.value = next
    setMsg(checked ? 'Permission assigned' : 'Permission removed')
  } catch (err) { setErr(err) }
}

async function removeRole(role) {
  if (!role.id || role.is_system) return
  if (!(await confirmAction('Delete this custom role?', 'Content managers assigned to this role may lose access.'))) return
  try {
    if (!(await canReachApi())) {
      error.value = 'Backend API is not reachable on port 9080. Start the native Go API to delete live roles.'
      return
    }
    await cms.adminDeleteRole(role.id)
    resetRoleForm()
    await loadRoles()
    setMsg('Role deleted')
  } catch (err) { setErr(err) }
}

function resetRoleForm() {
  Object.assign(roleForm, { id:'', name:'', description:'', is_system:false })
  selectedPermissionIds.value = new Set()
}

function sampleContactMessages() {
  return [
    { id:'preview-message-1', name:'Public Visitor', email:'visitor@example.com', subject:'Facility inquiry', message:'I would like to know more about booking a sports facility.', status:'unread', created_at:new Date().toISOString() },
  ]
}

function sampleNotificationsList() {
  return [
    { id:'preview-notification-1', type:'new_comment', title:'New Comment', message:'A public comment is waiting for moderation.', status:'unread', icon_key:'chat', created_at:new Date().toISOString() },
    { id:'preview-notification-2', type:'unusual_activity', title:'Unusual Activity', message:'A new sign-in was detected for the CMS.', status:'unread', icon_key:'shield-alert', created_at:new Date().toISOString() },
  ]
}

async function loadContactMessages() {
  try {
    if (!(await canReachApi())) {
      contactMessages.value = sampleContactMessages()
      error.value = 'Backend API is not reachable on port 9080. Messages is showing preview data.'
      return
    }
    const res = await cms.listMessages({ status:messageStatus.value, per_page:100 })
    contactMessages.value = listData(res)
  } catch (err) { setErr(err) }
}

async function markMessage(item, status) {
  try {
    if (!(await canReachApi())) {
      item.status = status
      return
    }
    await cms.updateMessage(item.id, { status })
    await loadContactMessages()
  } catch (err) { setErr(err) }
}

async function removeMessage(item) {
  if (!(await confirmAction('Delete this message?', 'The contact message will be removed from the CMS inbox.'))) return
  try {
    if (!(await canReachApi())) {
      contactMessages.value = contactMessages.value.filter(row => row.id !== item.id)
      return
    }
    await cms.deleteMessage(item.id)
    await loadContactMessages()
    setMsg('Message deleted')
  } catch (err) { setErr(err) }
}

async function markAllMessagesRead() {
  try {
    if (!(await canReachApi())) {
      contactMessages.value.forEach(item => item.status = 'read')
      return
    }
    const unread = contactMessages.value.filter(item => item.status === 'unread')
    await Promise.all(unread.map(item => cms.updateMessage(item.id, { status: 'read' })))
    await loadContactMessages()
    setMsg('Messages marked read')
  } catch (err) { setErr(err) }
}

async function clearAllMessages() {
  if (!(await confirmAction('Clear all messages?', 'All contact messages will be removed from the inbox.'))) return
  try {
    if (!(await canReachApi())) {
      contactMessages.value = []
      return
    }
    await cms.clearMessages()
    await loadContactMessages()
    setMsg('Messages cleared')
  } catch (err) { setErr(err) }
}

async function loadNotifications() {
  try {
    if (!(await canReachApi())) {
      cmsNotifications.value = sampleNotificationsList()
      error.value = 'Backend API is not reachable on port 9080. Notifications is showing preview data.'
      return
    }
    const res = await cms.listNotifications({ status:notificationStatus.value, per_page:100 })
    cmsNotifications.value = listData(res)
  } catch (err) { setErr(err) }
}

async function toggleNotificationRead(item) {
  const status = item.status === 'read' ? 'unread' : 'read'
  try {
    if (!(await canReachApi())) {
      item.status = status
      return
    }
    await cms.updateNotification(item.id, { status })
    await loadNotifications()
  } catch (err) { setErr(err) }
}

async function dismissNotification(item) {
  try {
    if (!(await canReachApi())) {
      cmsNotifications.value = cmsNotifications.value.filter(row => row.id !== item.id)
      return
    }
    await cms.deleteNotification(item.id)
    await loadNotifications()
  } catch (err) { setErr(err) }
}

async function clearAllNotifications() {
  if (!(await confirmAction('Clear all notifications?', 'All notifications will be dismissed.'))) return
  try {
    if (!(await canReachApi())) {
      cmsNotifications.value = []
      return
    }
    await cms.clearNotifications()
    await loadNotifications()
    setMsg('Notifications cleared')
  } catch (err) { setErr(err) }
}

async function markAllNotifications() {
  try {
    if (!(await canReachApi())) {
      cmsNotifications.value.forEach(item => item.status = 'read')
      return
    }
    await cms.markAllNotificationsRead()
    await loadNotifications()
    setMsg('Notifications marked read')
  } catch (err) { setErr(err) }
}

function notificationIcon(item) {
  const key = item.icon_key || item.type
  if (key === 'chat' || item.type === 'new_comment') return 'icofont-speech-comments'
  if (key === 'shield-alert' || item.type === 'unusual_activity' || item.type === 'new_sign_in') return 'icofont-shield-alt'
  if (key === 'key' || item.type === 'password_reset') return 'icofont-key'
  return 'icofont-check-circled'
}

function notificationThemeClass(item) {
  if (['unusual_activity', 'new_sign_in', 'password_reset'].includes(item.type)) return 'warn'
  if (item.type === 'system_success') return 'success'
  return 'info'
}

function sampleUsers() {
  return [
    { id:'preview-admin', first_name:'CMS', last_name:'Admin', email:'admin@ncs.go.ug', phone:'', is_active:true, account_status:'ACTIVE', roles:[{ id:'role_super_admin', name:'super_admin' }] },
    { id:'preview-editor', first_name:'Content', last_name:'Manager', email:'content@ncs.go.ug', phone:'', is_active:true, account_status:'ACTIVE', roles:[{ id:'role_content_manager', name:'content_manager' }] },
  ]
}

async function loadUsers() {
  try {
    if (!(await canReachApi())) {
      users.value = sampleUsers()
      selectedUser.value = null
      error.value = 'Backend API is not reachable on port 9080. Users is showing preview data.'
      return
    }
    const res = await cms.adminListUsers({ page:1, per_page:100, search:userSearch.value })
    users.value = listData(res)
    selectedUser.value = selectedUser.value ? users.value.find(user => user.id === selectedUser.value.id) || null : null
  } catch (err) { setErr(err) }
}

async function saveUser() {
  try {
    if (!(await canReachApi())) {
      setMsg(userForm.id ? 'Preview user updated' : 'Preview user created')
      resetUserForm()
      active.value = 'manage-users'
      users.value = users.value.length ? users.value : sampleUsers()
      return
    }
    if (userForm.id) {
      await cms.adminUpdateUser(userForm.id, {
        first_name:userForm.first_name,
        last_name:userForm.last_name,
        phone:userForm.phone,
      })
      setMsg('User updated')
    } else {
      await cms.adminCreateUser({
        first_name:userForm.first_name,
        last_name:userForm.last_name,
        email:userForm.email,
        phone:userForm.phone,
        password:userForm.password,
      })
      setMsg('User created')
    }
    resetUserForm()
    active.value = 'manage-users'
    await loadUsers()
  } catch (err) { setErr(err) }
}

function editUser(user) {
  Object.assign(userForm, {
    id:user.id,
    first_name:user.first_name || '',
    last_name:user.last_name || '',
    email:user.email || '',
    phone:user.phone || '',
    password:'',
  })
  active.value = 'users'
}

function resetUserForm() {
  Object.assign(userForm, { id:'', first_name:'', last_name:'', email:'', phone:'', password:'' })
}

async function removeUser(user) {
  if (!(await confirmAction('Delete this user?', `${user.email} will be disabled and removed from active CMS access.`))) return
  try {
    if (!(await canReachApi())) {
      users.value = users.value.filter(row => row.id !== user.id)
      setMsg('Preview user deleted')
      return
    }
    await cms.adminDeleteUser(user.id)
    await loadUsers()
    setMsg('User deleted')
  } catch (err) { setErr(err) }
}

async function toggleUserActive(user) {
  try {
    if (!(await canReachApi())) {
      user.is_active = !user.is_active
      user.account_status = user.is_active ? 'ACTIVE' : 'DISABLED'
      setMsg(user.is_active ? 'Preview user activated' : 'Preview user disabled')
      return
    }
    user.is_active ? await cms.adminDeactivateUser(user.id) : await cms.adminActivateUser(user.id)
    await loadUsers()
    setMsg(user.is_active ? 'User disabled' : 'User activated')
  } catch (err) { setErr(err) }
}

async function promptResetPassword(user) {
  const result = await Swal.fire({
    title: `Reset password for ${user.email}`,
    input: 'password',
    inputLabel: 'New password',
    inputAttributes: { minlength: 12, autocomplete: 'new-password' },
    showCancelButton: true,
    confirmButtonText: 'Reset password',
    confirmButtonColor: '#6777ef',
    cancelButtonColor: '#fc544b',
    inputValidator: value => !value || value.length < 12 ? 'Use at least 12 characters.' : undefined,
  })
  if (!result.isConfirmed) return
  try {
    if (!(await canReachApi())) {
      setMsg('Preview password reset')
      return
    }
    await cms.adminResetUserPassword(user.id, result.value)
    setMsg('Password reset')
  } catch (err) { setErr(err) }
}

function openUserRoles(user) {
  selectedUser.value = user
  if (!roles.value.length) loadRoles()
}

function userHasRole(user, role) {
  return (user.roles || []).some(item => (item.id && item.id === role.id) || (item.name || item) === role.name)
}

async function toggleUserRole(user, role, checked) {
  try {
    if (!(await canReachApi())) {
      const current = user.roles || []
      user.roles = checked
        ? [...current, { id:role.id, name:role.name, description:role.description }]
        : current.filter(item => (item.id && item.id !== role.id) || ((item.name || item) !== role.name))
      setMsg(checked ? 'Preview role assigned' : 'Preview role revoked')
      return
    }
    checked ? await cms.adminAssignUserRole(user.id, role.id) : await cms.adminRemoveUserRole(user.id, role.id)
    await loadUsers()
    selectedUser.value = users.value.find(row => row.id === user.id) || null
    setMsg(checked ? 'Role assigned' : 'Role revoked')
  } catch (err) { setErr(err) }
}

async function loadAuditLogs() {
  try {
    if (!(await canReachApi())) {
      const query = auditSearch.value.trim().toLowerCase()
      auditLogs.value = query
        ? sampleAuditLogs.filter(log => JSON.stringify(log).toLowerCase().includes(query))
        : sampleAuditLogs
      auditMeta.value = { total: auditLogs.value.length }
      selectedAuditLog.value = null
      error.value = 'Backend API is not reachable on port 9080. Audit Logs is showing preview data.'
      return
    }
    const res = await cms.adminListAuditLogs({ page:1, per_page:50, search:auditSearch.value })
    auditLogs.value = listData(res)
    auditMeta.value = pageMeta(res)
    selectedAuditLog.value = null
  } catch (err) { setErr(err) }
}

async function viewAuditLog(log) {
  try {
    if (!(await canReachApi()) || String(log.id).startsWith('preview-')) {
      selectedAuditLog.value = log
      return
    }
    const res = await cms.adminGetAuditLog(log.id)
    selectedAuditLog.value = data(res) || log
  } catch (err) { setErr(err) }
}

function statusClass(log) {
  const code = Number(log.response_code || 0)
  if (code >= 500 || log.severity_level === 'critical') return 'danger'
  if (code >= 400 || log.anomaly_detected || log.severity_level === 'high') return 'warn'
  return 'ok'
}

function formatDateTime(value) {
  if (!value) return ''
  return new Date(value).toLocaleString('en-UG', { dateStyle:'medium', timeStyle:'short' })
}
function formatDay(value) {
  if (!value) return ''
  return new Date(value).toLocaleDateString('en-UG', { weekday:'long' })
}

function sampleNewsletterSubscribers() {
  return [
    { id:'preview-1', email:'communications@ncs.go.ug', created_at:new Date().toISOString(), status:'subscribed' },
    { id:'preview-2', email:'sportsdesk@example.com', created_at:new Date(Date.now() - 86400000).toISOString(), status:'subscribed' },
  ]
}

async function loadNewsletterSubscribers() {
  try {
    if (!(await canReachApi())) {
      newsletterSubscribers.value = sampleNewsletterSubscribers()
      error.value = 'Backend API is not reachable on port 9080. Newsletter is showing preview subscribers.'
      return
    }
    const res = await cms.adminListNewsletterSubscribers({ per_page:200 })
    newsletterSubscribers.value = listData(res)
  } catch (err) { setErr(err) }
}

async function exportNewsletter(format) {
  if (await canReachApi()) {
    try {
      const res = await cms.adminExportNewsletterSubscribers(format)
      downloadBlob(res.data, `newsletter-subscribers.${format === 'excel' ? 'xls' : 'csv'}`)
      return
    } catch {
      // Fall back to client-side export from the loaded rows.
    }
  }
  const rows = newsletterSubscribers.value.length ? newsletterSubscribers.value : sampleNewsletterSubscribers()
  const csv = toNewsletterCsv(rows)
  downloadBlob(new Blob([csv], { type: format === 'excel' ? 'application/vnd.ms-excel' : 'text/csv;charset=utf-8' }), `newsletter-subscribers.${format === 'excel' ? 'xls' : 'csv'}`)
}

function toNewsletterCsv(rows) {
  const escape = value => `"${String(value ?? '').replaceAll('"', '""')}"`
  return [
    ['Email', 'Date Submitted', 'Day', 'Status'].map(escape).join(','),
    ...rows.map(row => [row.email, formatDateTime(row.created_at || row.submitted_at), formatDay(row.created_at || row.submitted_at), row.status || 'subscribed'].map(escape).join(',')),
  ].join('\n')
}

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}

function logout() {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  router.push('/login')
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('ncsms_sidebar_collapsed', String(sidebarCollapsed.value))
}

function expandSidebarOnNav() {
  if (sidebarCollapsed.value) {
    sidebarCollapsed.value = false
    localStorage.setItem('ncsms_sidebar_collapsed', 'false')
  }
}

function handleHomepageGroupFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) homepageGroupOpen.value = false
}

function handleBlogsGroupFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) blogsGroupOpen.value = false
}

function handleSlideshowGroupFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) slideshowGroupOpen.value = false
}

function resetPostForm() {
  Object.assign(postForm, { id:'', title:'', slug:'', category:'news', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
}

function resetPageForm() {
  Object.assign(pageForm, { id:'', title:'', slug:'', category:'page', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
}

function resetProjectForm() {
  Object.assign(projectForm, { id:'', title:'', slug:'', category:'project', category_tag:'', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
}

function resetCaseStudyForm() {
  Object.assign(caseStudyForm, { id:'', title:'', slug:'', category:'case_study', category_tag:'', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
}

function resetBlogCategoryForm() {
  Object.assign(blogCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true })
}

function resetProjectCategoryForm() {
  Object.assign(projectCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'project' })
}

function resetCaseStudyCategoryForm() {
  Object.assign(caseStudyCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'case_study' })
}

function resetFAQForm() {
  Object.assign(faqForm, { id:'', question:'', answer:'', category:'General', sort_order:0, is_active:true })
}

function resetFAQCategoryForm() {
  Object.assign(faqCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'faq' })
}

function resetResourceForm() {
  Object.assign(resourceForm, { id:'', title:'', category:'Guidelines', file_url:'', description:'', sort_order:0, is_active:true })
}

function resetResourceCategoryForm() {
  Object.assign(resourceCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'resource' })
}

function resetCareerForm() {
  Object.assign(careerForm, { id:'', title:'', department:'', department_id:'', location:'', job_type:'full_time', category:'jobs', description:'', requirements:'', salary_range:'', status:'draft', deadline_at:'' })
}

function resetCareerCategoryForm() {
  Object.assign(careerCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'career' })
}

function resetTeamForm() {
  Object.assign(teamForm, { id:'', full_name:'', designation:'', department_id:'', image_url:'', bio:'', sort_order:0, is_active:true })
}

function resetTeamDepartmentForm() {
  Object.assign(teamDepartmentForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'team_department' })
}

function resetFacilityForm() {
  Object.assign(facilityForm, { id:'', name:'', slug:'', category:'', description:'', image_url:'', sort_order:0, is_active:true })
}

function resetFacilityCategoryForm() {
  Object.assign(facilityCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'facility' })
}

function resetEventForm() {
  Object.assign(eventForm, { id:'', title:'', slug:'', category:'', location:'', event_date:'', end_date:'', cover_image_url:'', description:'', status:'published' })
}

function resetEventCategoryForm() {
  Object.assign(eventCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'event' })
}

function resetInvestForm() {
  Object.assign(investForm, { id:'', title:'', subtitle:'', content:'', image_url:'', sort_order:0, is_active:true })
}

function resetInvestCategoryForm() {
  Object.assign(investCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'investment' })
}

function resetFederationForm() {
  Object.assign(associationForm, { id:'', name:'', slug:'', category:'', president:'', secretary:'', phone:'', website_url:'', description:'', logo_url:'', sort_order:0, is_active:true })
}

function resetFederationCategoryForm() {
  Object.assign(federationCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true, content_type:'federation' })
}

function resetFactForm() {
  Object.assign(factForm, { id:'', label:'', value:'', icon:'icofont-chart-growth', sort_order:0, is_active:true })
}

async function loadComments() {
  try {
    if (!(await canReachApi())) {
      error.value = 'Backend API is not reachable on port 9080. Start the native Go API to load comments.'
      return
    }
    const res = await cms.adminListComments({ status: commentStatus.value, per_page:50 })
    comments.value = data(res)?.items || []
  } catch (err) { setErr(err) }
}

async function moderateComment(comment, action) {
  try {
    action === 'approve' ? await cms.adminApproveComment(comment.id) : await cms.adminFlagComment(comment.id)
    await loadComments()
    setMsg(`Comment ${action === 'approve' ? 'approved' : 'flagged'}`)
  } catch (err) { setErr(err) }
}

async function deleteComment(comment) {
  if (!(await confirmAction('Delete this comment permanently?'))) return
  try {
    await cms.adminDeleteComment(comment.id)
    await loadComments()
    setMsg('Comment deleted')
  } catch (err) { setErr(err) }
}

const ContentTable = defineComponent({
  props: { items:Array, titleKey:String, subtitleKey:String },
  emits: ['edit', 'delete'],
  setup(props, { emit }) {
    return () => h('div', { class:'card otika-basic-table-card' }, [
      h('div', { class:'card-body p-0' }, [
        h('div', { class:'table-responsive' }, [
          h('table', { class:'table table-striped table-hover table-sm mb-0' }, [
            h('thead', [
              h('tr', [
                h('th', { scope:'col' }, '#'),
                h('th', { scope:'col' }, 'Title'),
                h('th', { scope:'col' }, 'Details'),
                h('th', { scope:'col', class:'text-right' }, 'Action'),
              ]),
            ]),
            h('tbody', (props.items || []).length
              ? (props.items || []).map((item, index) => h('tr', { key:item.id || item.slug || index }, [
                  h('th', { scope:'row' }, index + 1),
                  h('td', [h('strong', item[props.titleKey] || 'Untitled')]),
                  h('td', [h('span', item[props.subtitleKey] || '')]),
                  h('td', { class:'text-right table-actions' }, [
                    h('button', { type:'button', class:'btn btn-sm btn-primary mr-1', onClick:() => emit('edit', item) }, 'Edit'),
                    h('button', { type:'button', class:'btn btn-sm btn-danger', onClick:() => emit('delete', item) }, 'Delete'),
                  ]),
                ]))
              : [h('tr', [h('td', { colspan:4, class:'text-center text-muted' }, 'No records found.')])]),
          ]),
        ]),
      ]),
    ])
  },
})

const EditorForm = defineComponent({
  props: { title:String, model:Object, fields:Array },
  emits: ['save', 'error'],
  setup(props, { emit }) {
    const humanize = name => name.replaceAll('_', ' ')
    const control = (field) => {
      if (field.name === 'file_url') {
        return h(DropzoneUpload, {
          modelValue: props.model[field.name],
          'onUpdate:modelValue': value => props.model[field.name] = value,
          mediaType: 'file',
          label: 'resource file',
          accept: 'application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
          hint: 'Upload PDF, Word, or Excel files for the Resource Centre.',
          onError: error => emit('error', error),
        })
      }
      if (['image_url', 'cover_image_url', 'logo_url'].includes(field.name)) {
        return h(DropzoneUpload, {
          modelValue: props.model[field.name],
          'onUpdate:modelValue': value => props.model[field.name] = value,
          label: 'image',
          accept: 'image/png,image/jpeg,image/webp,image/svg+xml',
          onError: error => emit('error', error),
        })
      }
      if (field.type === 'select') {
        return h('select', { class: 'form-control selectric', value: props.model[field.name], onChange: e => props.model[field.name] = e.target.value }, [
          h('option', { value: '' }, 'None'),
          ...(field.options || []).map(o => h('option', { key: o.id || o.slug, value: o.slug || o.value }, o.name || o.label)),
        ])
      }
      if (field.type === 'textarea') {
        return h('textarea', { class: 'form-control', rows: 4, value: props.model[field.name], onInput: e => props.model[field.name] = e.target.value })
      }
      if (field.type === 'number') {
        // Coerce to a real number so the JSON payload sends e.g. sort_order:2
        // (an int the backend accepts) rather than the string "2" (which fails
        // Go's int decode and returns 400 Bad Request).
        return h('input', { class: 'form-control', type: 'number', value: props.model[field.name], onInput: e => props.model[field.name] = e.target.value === '' ? 0 : Number(e.target.value) })
      }
      return h('input', { class: 'form-control', value: props.model[field.name], onInput: e => props.model[field.name] = e.target.value })
    }
    // Mirrors the Otika "basic-form.html" markup: card > card-body > form-group
    // (label + .form-control), checkboxes as Bootstrap .form-check, and a
    // card-footer submit button — so all EditorForm sections match the template.
    return () => h('form', { class: 'otika-form-card', onSubmit: e => { e.preventDefault(); emit('save') } }, [
      h('div', { class: 'card' }, [
        h('div', { class: 'card-header' }, [h('h4', props.title)]),
        h('div', { class: 'card-body' }, props.fields.map(field => {
          if (field.type === 'checkbox') {
            const id = `ef-${field.name}`
            return h('div', { class: 'form-group' }, [
              h('div', { class: 'form-check' }, [
                h('input', { class: 'form-check-input', type: 'checkbox', id, checked: !!props.model[field.name], onChange: e => props.model[field.name] = e.target.checked }),
                h('label', { class: 'form-check-label text-capitalize', for: id }, humanize(field.name)),
              ]),
            ])
          }
          return h('div', { class: 'form-group' }, [
            h('label', { class: 'text-capitalize' }, humanize(field.name)),
            control(field),
          ])
        })),
        h('div', { class: 'card-footer text-right' }, [
          h('button', { class: 'btn btn-primary mr-1', type: 'submit' }, props.model.id ? 'Update' : 'Create'),
        ]),
      ]),
    ])
  },
})
</script>

<style scoped>
.cms-shell{display:grid;grid-template-columns:18rem 1fr;min-height:100vh;transition:grid-template-columns 200ms ease}.cms-sidebar{position:sticky;top:0;height:100vh;background:#10233f;color:white;padding:1rem;display:flex;flex-direction:column;min-width:0;overflow:hidden}.cms-sidebar-head{display:flex;align-items:center;justify-content:space-between;gap:.5rem;margin-bottom:1rem}.cms-brand{display:flex;align-items:center;gap:.7rem;color:white;font-weight:800;min-width:0;overflow:hidden}.cms-brand img{flex:0 0 auto;width:2.6rem;height:2.6rem;object-fit:contain;background:white;border-radius:.35rem}.cms-collapse-toggle{flex:0 0 auto;display:flex;align-items:center;justify-content:center;width:2.25rem;height:2.25rem;border-radius:.4rem;background:rgb(255 255 255/.08);color:white;border:1px solid rgb(255 255 255/.14)}.cms-collapse-toggle:hover,.cms-collapse-toggle:focus-visible{background:#f5a623;color:#10233f}.cms-nav{display:grid;gap:.25rem;overflow:auto}.cms-nav button,.cms-logout{display:flex;align-items:center;gap:.65rem;border-radius:.45rem;padding:.65rem .75rem;color:rgb(255 255 255/.78);text-align:left}.cms-nav button.active,.cms-nav button:hover,.cms-nav button:focus-visible{background:#f5a623;color:#10233f}.cms-nav-group{display:grid}.cms-nav-parent{width:100%;justify-content:space-between}.cms-nav-submenu{display:grid;gap:.25rem;max-height:0;opacity:0;overflow:hidden;transform:translateZ(0);transition:max-height 220ms ease,opacity 180ms ease,padding 220ms ease;padding-left:.75rem}.cms-nav-group.expanded .cms-nav-submenu{max-height:24rem;opacity:1;padding-top:.25rem;padding-bottom:.25rem}.cms-nav-submenu button{font-size:.86rem;padding-left:1rem}.cms-logout{margin-top:auto;background:rgb(255 255 255/.08)}.cms-main{padding:2rem;min-width:0}.cms-header{display:flex;justify-content:space-between;gap:1rem;align-items:center;margin-bottom:1.5rem}.cms-kicker{text-transform:uppercase;letter-spacing:.18em;color:#d88700;font-size:.72rem;font-weight:800}.cms-header h1{font-size:2rem;font-weight:850;color:#1a365d}.cms-actions{display:flex;gap:.5rem;align-items:center}.cms-actions a,.cms-actions button,.cms-panel-head button,.cms-list-editor>button,.cms-editor button{border-radius:.45rem;background:#1a365d;color:white;padding:.6rem .9rem;font-size:.85rem;font-weight:700}.cms-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:1rem}.cms-card,.cms-panel{background:white;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem}.metric span{color:#64748b;font-size:.8rem}.metric strong{display:block;color:#1a365d;font-size:2rem}.cms-panel{display:grid;gap:1rem}.cms-panel-head{display:flex;justify-content:space-between;align-items:center;gap:1rem}.cms-panel h2{font-size:1.15rem;font-weight:800;color:#1a365d}.menu-save-bar{position:sticky;bottom:0;display:flex;align-items:center;justify-content:flex-end;gap:.75rem;margin-top:1rem;padding:.85rem 1rem;background:#fff;border-top:1px solid #e5e7eb;box-shadow:0 -4px 16px rgba(15,23,42,.06)}.menu-save-bar__status{margin-right:auto;font-size:.8rem;font-weight:700;color:#64748b}.menu-save-bar__status--dirty{color:#b45309}.menu-save-bar button{border-radius:.45rem;border:1px solid #d1d5db;padding:.6rem 1rem;font-weight:700;font-size:.85rem}.menu-save-bar button:disabled{opacity:.5;cursor:not-allowed}.menu-save-bar__save{background:#1a365d;color:#fff;border-color:#1a365d}:global(.dark) .menu-save-bar{background:#111827;border-color:#334155}:global(.dark) .menu-save-bar__status{color:#94a3b8}:global(.dark) .menu-save-bar__status--dirty{color:#fdba74}:global(.dark) .menu-save-bar button{border-color:#475569;color:#e5e7eb}.cms-menu-builder-grid{display:grid;gap:1rem}.cms-menu-builder-grid article{display:grid;gap:.75rem;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem}.cms-menu-builder-grid h3{font-weight:800;color:#1a365d}.cms-two{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1rem}.cms-two label{display:grid;gap:.35rem;font-size:.76rem;font-weight:800;text-transform:capitalize;color:#475569}.cms-two .wide{grid-column:1/-1}.cms-two input,.cms-two textarea,.cms-list-editor textarea{width:100%;border:1px solid #cbd5e1;border-radius:.45rem;padding:.65rem;text-transform:none;font-weight:500;color:#111827}.cms-two textarea,.cms-list-editor textarea{min-height:6rem}.cms-list-editor{display:grid;gap:.75rem}.cms-row{display:grid;grid-template-columns:1fr auto;gap:.75rem}.cms-row button,.cms-table-row button{border:1px solid #d1d5db;border-radius:.4rem;padding:.45rem .7rem}.cms-table{display:grid;gap:.5rem}.cms-table-row{display:flex;justify-content:space-between;gap:1rem;align-items:center;border:1px solid #e5e7eb;border-radius:.45rem;padding:.75rem}.cms-table-row strong{display:block;color:#1a365d}.cms-table-row span{font-size:.8rem;color:#64748b}.cms-table-row div:last-child{display:flex;gap:.4rem}.cms-message{background:#ecfdf5;color:#047857;border:1px solid #a7f3d0;border-radius:.45rem;padding:.75rem}.cms-error{background:#fef2f2;color:#b91c1c;border:1px solid #fecaca;border-radius:.45rem;padding:.75rem}:global(.dark .cms-main){background:#0f172a;color:#e5e7eb}:global(.dark .cms-card),:global(.dark .cms-panel),:global(.dark .cms-menu-builder-grid article){background:#111827;border-color:#334155}:global(.dark .cms-header h1),:global(.dark .cms-panel h2),:global(.dark .cms-menu-builder-grid h3),:global(.dark .cms-table-row strong),:global(.dark .metric strong){color:#f8fafc}:global(.dark .cms-two input),:global(.dark .cms-two textarea),:global(.dark .cms-list-editor textarea){background:#0f172a;color:#f8fafc;border-color:#475569}@media(max-width:900px){.cms-shell{grid-template-columns:1fr}.cms-sidebar{position:relative;height:auto}.cms-grid,.cms-two{grid-template-columns:1fr}.cms-header{align-items:flex-start;flex-direction:column}}
@media(min-width:901px){.cms-shell.collapsed{grid-template-columns:4.5rem 1fr}.cms-shell.collapsed .cms-sidebar{padding:1rem .6rem}.cms-shell.collapsed .cms-brand span,.cms-shell.collapsed .cms-nav button span,.cms-shell.collapsed .cms-logout span,.cms-shell.collapsed .cms-nav-parent .icofont-rounded-down,.cms-shell.collapsed .cms-nav-submenu{display:none}.cms-shell.collapsed .cms-nav button,.cms-shell.collapsed .cms-nav-parent,.cms-shell.collapsed .cms-logout{justify-content:center}.cms-shell.collapsed .cms-sidebar-head{justify-content:center;flex-direction:column;gap:.6rem}}
.cms-panel-head select,.cms-two select{border:1px solid #cbd5e1;border-radius:.45rem;padding:.65rem;color:#111827;background:white}.cms-empty{color:#64748b;border:1px dashed #cbd5e1;border-radius:.45rem;padding:1rem;text-align:center}.comment-row span{max-width:54rem;display:block;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}:deep(.blog-editor){display:grid;gap:1rem}:deep(.inline-field){display:flex;gap:.5rem;align-items:center;text-transform:none}:deep(.inline-field input){width:auto}:deep(.counter){float:right;color:#64748b;font-weight:700}:deep(.counter.warn){color:#b45309}:deep(.editor-toolbar button){border-radius:.4rem;border:1px solid #cbd5e1;padding:.55rem .8rem;font-weight:800}:deep(.editor-toolbar){display:flex;flex-wrap:wrap;gap:.4rem}:deep(.editor-toolbar .active){background:#1a365d;color:white}:deep(.rich-editor){border:1px solid #cbd5e1;border-radius:.5rem;background:white;padding:1rem;min-height:18rem}:deep(.rich-editor .ProseMirror){min-height:16rem;outline:none}:deep(.rich-editor h2){font-size:1.5rem}:deep(.rich-editor ul){list-style:disc;padding-left:1.25rem}:deep(.rich-editor img){max-width:100%;border-radius:.5rem}:deep(.draft-state){font-size:.8rem;color:#64748b}
.cms-admin-grid{display:grid;grid-template-columns:minmax(0,1fr) minmax(18rem,.85fr);gap:1rem}.cms-subpanel{display:grid;gap:.85rem;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem}.cms-subpanel h3{font-weight:800;color:#1a365d}.footer-columns-hint{color:#6c757d;font-size:.82rem;margin:0}.footer-column-editor{display:grid;gap:.6rem;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem;margin-bottom:1rem}.footer-column-editor .cms-panel-head{border:0;margin:0;padding:0;align-items:flex-end}.footer-column-title{display:grid;gap:.4rem;flex:1 1 auto;font-size:.78rem;font-weight:700;color:#34395e}.footer-column-title input{border:1px solid #e4e6fc;border-radius:.35rem;padding:.5rem .7rem;font-weight:600}.footer-link-row{display:grid;grid-template-columns:1fr 1fr auto;gap:.6rem;align-items:center}.footer-link-row input{border:1px solid #e4e6fc;border-radius:.35rem;padding:.5rem .7rem}.footer-remove-btn{border:0;border-radius:.4rem;background:#fc544b;color:#fff;padding:.45rem .8rem;font-size:.78rem;font-weight:700;white-space:nowrap}:global(.dark .footer-column-editor){background:#111827;border-color:#334155}:global(.dark .footer-column-title){color:#f8fafc}:global(.dark .footer-column-title input),:global(.dark .footer-link-row input){background:#0f172a;color:#f8fafc;border-color:#475569}@media(max-width:700px){.footer-link-row{grid-template-columns:1fr}}.cms-actions-inline{display:flex;flex-wrap:wrap;gap:.6rem}.cms-actions-inline button,.cms-search button,.audit-row button{border-radius:.45rem;background:#1a365d;color:white;padding:.55rem .8rem;font-weight:800;font-size:.82rem}.cms-actions-inline button:last-child{background:#991b1b}.cms-table.compact{max-height:28rem;overflow:auto}.cms-table-row.selected{border-color:#f5a623;background:#fff7ed}.permission-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:.75rem}.permission-group{display:grid;align-content:start;gap:.5rem;border:1px solid #e5e7eb;border-radius:.5rem;padding:.85rem}.permission-group h3{text-transform:capitalize}.permission-row{display:grid;grid-template-columns:auto 1fr;gap:.55rem;align-items:start;text-transform:none!important;font-size:.82rem!important}.permission-row input{margin-top:.2rem}.permission-row strong{display:block;color:#1a365d;text-transform:capitalize}.permission-row small{display:block;color:#64748b;font-weight:500;text-transform:none}.cms-search{display:flex;gap:.5rem;align-items:center}.cms-search input{min-width:18rem;border:1px solid #cbd5e1;border-radius:.45rem;padding:.6rem;color:#111827}.audit-summary{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:.75rem}.audit-summary article{border:1px solid #e5e7eb;border-radius:.5rem;padding:.85rem}.audit-summary span{display:block;color:#64748b;font-size:.76rem}.audit-summary strong{display:block;color:#1a365d;font-size:1.5rem}.audit-table{display:grid;gap:.5rem}.audit-row{display:grid;grid-template-columns:minmax(0,1fr) auto;gap:1rem;align-items:center;border:1px solid #e5e7eb;border-radius:.45rem;padding:.75rem}.audit-row strong{display:block;color:#1a365d}.audit-row span{font-size:.78rem;color:#64748b}.audit-row>div:last-child{display:flex;align-items:center;gap:.5rem}.audit-code{border-radius:999px;padding:.25rem .55rem;font-weight:800}.audit-code.ok{background:#ecfdf5;color:#047857}.audit-code.warn{background:#fff7ed;color:#b45309}.audit-code.danger{background:#fef2f2;color:#b91c1c}.audit-detail pre{max-height:26rem;overflow:auto;border-radius:.45rem;background:#0f172a;color:#e5e7eb;padding:1rem;font-size:.78rem;white-space:pre-wrap}:global(.dark .cms-subpanel),:global(.dark .permission-group),:global(.dark .audit-summary article),:global(.dark .audit-row){background:#111827;border-color:#334155}:global(.dark .permission-row strong),:global(.dark .audit-row strong),:global(.dark .audit-summary strong),:global(.dark .cms-subpanel h3){color:#f8fafc}:global(.dark .cms-search input){background:#0f172a;color:#f8fafc;border-color:#475569}:global(.dark .cms-table-row.selected){background:#422006;border-color:#f5a623}@media(max-width:1100px){.cms-admin-grid,.permission-grid{grid-template-columns:1fr}.audit-row{grid-template-columns:1fr}.audit-row>div:last-child{flex-wrap:wrap}.audit-summary{grid-template-columns:1fr}.cms-search{width:100%;flex-wrap:wrap}.cms-search input{min-width:0;flex:1 1 14rem}}
.otika-cms{min-height:100vh;background:#f4f6f9;color:#34395e}.otika-cms :deep(.main-navbar){box-shadow:0 4px 25px 0 rgba(0,0,0,.1)}.otika-cms :deep(.navbar-bg){background:#6777ef;position:fixed}.otika-cms button{font-family:inherit}.otika-cms .main-sidebar{position:fixed;width:260px;height:100vh;overflow:hidden}.otika-cms .sidebar-mini .main-sidebar{width:65px}.otika-cms #sidebar-wrapper{display:flex;flex-direction:column;height:100vh;min-height:0}.otika-cms .main-content{padding-top:80px}.otika-cms .sidebar-brand{display:flex;align-items:center;justify-content:center;flex:0 0 70px;height:70px}.otika-cms .sidebar-brand a{display:flex;align-items:center;justify-content:center;width:100%}.otika-cms .header-logo{max-width:84px;max-height:48px;object-fit:contain}.sidebar-user{display:flex;gap:10px;align-items:center;flex:0 0 auto;min-width:0;margin:8px 14px 12px;padding:12px;border-radius:8px;background:#f8f9fa}.sidebar-user img{width:38px;height:38px;flex:0 0 38px}.sidebar-user div{min-width:0}.sidebar-user strong{display:block;max-width:150px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:#34395e;font-size:13px}.sidebar-user span{display:block;color:#98a6ad;font-size:11px}.sidebar-menu{flex:1 1 auto;min-height:0;overflow-y:auto;overflow-x:hidden;padding-bottom:18px;scrollbar-width:thin}.sidebar-menu::-webkit-scrollbar{width:6px}.sidebar-menu::-webkit-scrollbar-thumb{background:#d7dbea;border-radius:999px}.sidebar-menu button{border:0;background:transparent;width:100%;min-width:0;text-align:left}.sidebar-menu .nav-link{display:flex!important;align-items:center;gap:10px;width:100%;height:auto!important;min-height:38px;padding:9px 18px!important;line-height:1.2!important}.sidebar-menu .nav-link i{flex:0 0 18px;width:18px;text-align:center}.sidebar-menu .nav-link span,.sidebar-menu .dropdown-menu .nav-link{min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-size:12px}.sidebar-menu .dropdown-menu{position:static!important;width:100%!important;min-width:0!important;float:none!important;box-shadow:none!important;background:transparent!important;padding:0 0 6px 32px!important}.sidebar-menu .menu-header{white-space:nowrap}.sidebar-menu li.active>button,.sidebar-menu li.active>a{color:#6777ef!important;font-weight:600}.sidebar-menu li.active>button i{color:#6777ef!important}.main-navbar button.nav-link{border:0;background:transparent}.cms-top-icon{display:inline-flex!important;align-items:center;justify-content:center;color:#111827!important}.cms-top-icon i{color:#111827!important;font-size:18px}.search-element .btn i{color:#111827!important}.navbar .dropdown-menu.show{display:block;position:absolute}.navbar .dropdown-item{border:0;background:transparent;text-align:left}.navbar .dropdown-item.has-icon{width:100%;text-align:left;border:0;background:transparent}.link-button{border:0;background:transparent;color:#6777ef;font:inherit;font-weight:600;padding:0}.search-element{display:flex;align-items:center;gap:8px}.search-element .form-control{border-radius:30px!important;margin-right:0}.search-element .btn{border-radius:30px;-webkit-appearance:none;appearance:none;transition:background-color .15s ease,color .15s ease}.search-element .btn:hover,.search-element .btn:focus-visible{background:#f5a623}.search-element .btn i{color:#111827!important;transition:color .15s ease}.search-element .btn:hover i,.search-element .btn:focus-visible i{color:#10233f!important}.headerBadge1,.headerBadge2{position:absolute;top:8px;right:4px;background:#ffa426;color:#fff}.otika-page-actions{display:flex!important;justify-content:flex-end;gap:10px;margin:-12px 0 18px}.cms-grid{display:grid!important;grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:25px!important}.cms-card,.cms-panel,.cms-subpanel,.permission-group,.audit-summary article,.audit-row{background:#fff!important;border:0!important;border-radius:3px!important;box-shadow:0 4px 25px 0 rgba(0,0,0,.1)!important}.cms-card.metric{position:relative;padding:25px!important;min-height:118px}.metric span{color:#98a6ad!important;font-size:13px!important;font-weight:500}.metric strong{color:#34395e!important;font-size:30px!important;font-weight:700}.cms-panel{display:grid;gap:20px;padding:25px!important}.cms-panel-head{border-bottom:1px solid #f9f9f9;margin:-4px -4px 8px;padding-bottom:15px}.cms-panel h2,.cms-subpanel h3,.permission-group h3{color:#34395e!important;font-size:16px!important;font-weight:700}.cms-panel-head button,.cms-list-editor>button,.cms-editor button,.cms-actions-inline button,.cms-search button,.audit-row button,.cms-table-row button{border:0!important;border-radius:30px!important;background:#6777ef!important;color:#fff!important;padding:8px 18px!important;font-size:12px!important;font-weight:600!important;box-shadow:0 2px 6px #acb5f6}.cms-actions-inline button:last-child,.cms-table-row button:last-child{background:#fc544b!important;box-shadow:0 2px 6px #fd9b96!important}.cms-two label,.cms-list-editor label{color:#34395e!important;font-size:12px!important;font-weight:600!important;text-transform:none!important}.cms-two input,.cms-two textarea,.cms-list-editor textarea,.cms-panel-head select,.cms-two select,.cms-search input{height:auto!important;border:1px solid #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#495057!important;padding:10px 15px!important;box-shadow:none!important}.cms-two textarea,.cms-list-editor textarea{min-height:120px!important}.cms-two input:focus,.cms-two textarea:focus,.cms-search input:focus{border-color:#6777ef!important;box-shadow:0 2px 6px #acb5f6!important}.cms-table,.audit-table{display:block!important;overflow-x:auto}.cms-table-row,.audit-row{display:grid!important;grid-template-columns:minmax(0,1fr) auto;align-items:center;margin-bottom:12px;padding:15px 18px!important}.cms-table-row strong,.audit-row strong,.permission-row strong{color:#34395e!important;font-size:14px}.cms-table-row span,.audit-row span,.permission-row small{color:#6c757d!important}.cms-message{border:0!important;border-radius:3px!important;background:#e8f7f0!important;color:#3abaf4!important;box-shadow:0 4px 25px rgba(0,0,0,.05)}.cms-error{border:0!important;border-radius:3px!important;background:#fdeaea!important;color:#fc544b!important;box-shadow:0 4px 25px rgba(0,0,0,.05)}.permission-grid{gap:25px!important}.permission-row input[type=checkbox]{width:18px;height:18px;accent-color:#6777ef}.audit-summary article{padding:18px 22px!important}.audit-code{border-radius:30px!important}.audit-code.ok{background:#e8f7f0!important;color:#47c363!important}.audit-code.warn{background:#fff4e6!important;color:#ffa426!important}.audit-code.danger{background:#fdeaea!important;color:#fc544b!important}:deep(.blog-editor),:deep(.menu-builder){background:#fff;border-radius:3px}:deep(.editor-toolbar button){border:0!important;border-radius:30px!important;background:#f4f6f9!important;color:#34395e!important;padding:8px 15px!important}:deep(.editor-toolbar .active){background:#6777ef!important;color:white!important}:deep(.rich-editor){border:1px solid #e4e6fc!important;border-radius:3px!important;box-shadow:none!important}.cms-menu-builder-grid article{border:0!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.cms-empty{border:1px dashed #e4e6fc!important;border-radius:3px!important;background:#fdfdff;color:#98a6ad!important}.cms-panel :deep(.otika-form-card .card){border:0!important;box-shadow:none!important;border-radius:3px!important;margin-bottom:0!important}.cms-panel :deep(.blog-editor-card){box-shadow:none!important;border-radius:3px!important;margin-bottom:0!important}.sidebar-mini .sidebar-user strong,.sidebar-mini .sidebar-user span,.sidebar-mini .sidebar-menu .nav-link span,.sidebar-mini .sidebar-menu .menu-header,.sidebar-mini .sidebar-menu .dropdown-menu{display:none!important}.sidebar-mini .sidebar-user{justify-content:center;margin:8px 8px 12px;padding:10px 6px}.sidebar-mini .sidebar-menu .nav-link{justify-content:center;padding:10px!important}.sidebar-mini .sidebar-menu .nav-link i{margin:0!important}@media(max-width:1200px){.cms-grid{grid-template-columns:repeat(2,minmax(0,1fr))!important}.permission-grid{grid-template-columns:repeat(2,minmax(0,1fr))!important}}@media(max-width:991px){.otika-cms .navbar-bg{height:116px}.otika-cms .main-navbar{left:0;right:0;min-height:70px;padding:8px 14px;flex-wrap:wrap}.otika-cms .main-navbar .form-inline.mr-auto{flex:1 1 auto;min-width:0}.otika-cms .main-navbar .navbar-nav{flex-direction:row;align-items:center;gap:4px;margin-right:0!important}.otika-cms .navbar-right{margin-left:auto;flex-direction:row}.otika-cms .main-sidebar{position:relative!important;width:100%!important;height:auto!important;left:0!important;top:0!important;z-index:1;box-shadow:none}.otika-cms #sidebar-wrapper{height:auto;max-height:none}.otika-cms .main-content{margin-left:0!important;padding:24px 16px!important;padding-top:24px!important}.otika-cms .sidebar-menu{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:4px 8px;max-height:60vh;padding-bottom:18px}.otika-cms .sidebar-menu .menu-header{grid-column:1/-1}.otika-cms .sidebar-user{margin:8px 16px 14px}.section-header{align-items:flex-start!important;gap:10px;flex-direction:column}.section-header-breadcrumb{margin-left:0!important}.cms-admin-grid,.cms-two,.cms-menu-builder-grid{grid-template-columns:1fr!important}.otika-page-actions{justify-content:flex-start!important;flex-wrap:wrap}.cms-panel-head{align-items:flex-start;flex-direction:column}.cms-search{width:100%;flex-wrap:wrap}.cms-search input{min-width:0!important;flex:1 1 12rem}}@media(max-width:768px){.navbar .search-element{display:none}.cms-grid,.permission-grid,.audit-summary{grid-template-columns:1fr!important}.cms-table-row,.audit-row{grid-template-columns:1fr!important}.cms-table-row>div:last-child,.audit-row>div:last-child{justify-content:flex-start;flex-wrap:wrap}.cms-row{grid-template-columns:1fr}.cms-panel{padding:18px!important}.cms-card.metric{min-height:auto}.navbar .dropdown-menu.show{position:fixed!important;top:62px!important;left:12px!important;right:12px!important;width:auto!important;max-width:none}}@media(max-width:575px){.otika-cms .main-navbar{padding:6px 10px}.otika-cms .sidebar-menu{grid-template-columns:1fr}.otika-cms .main-content{padding:18px 12px!important}.otika-cms .section-header h1{font-size:20px!important;line-height:1.25}.cms-panel,.cms-card.metric{padding:16px!important}.cms-actions a,.cms-actions button,.cms-panel-head button,.cms-list-editor>button,.cms-editor button{width:100%;justify-content:center;text-align:center}.cms-two input,.cms-two textarea,.cms-list-editor textarea,.cms-panel-head select,.cms-two select,.cms-search input{font-size:14px!important}.permission-row{grid-template-columns:24px 1fr!important}.comment-row span{white-space:normal}}
.otika-cms button:focus,.otika-cms a:focus,.otika-cms input:focus,.otika-cms textarea:focus,.otika-cms select:focus{outline:none!important}.otika-cms button:focus-visible,.otika-cms a:focus-visible,.otika-cms input:focus-visible,.otika-cms textarea:focus-visible,.otika-cms select:focus-visible{outline:2px solid rgba(103,119,239,.45)!important;outline-offset:2px!important}.otika-cms .sidebar-menu button:focus-visible,.otika-cms .main-navbar button:focus-visible{outline:0!important;box-shadow:0 0 0 3px rgba(103,119,239,.22)!important}.otika-cms .cms-table-row button:focus,.otika-cms .audit-row button:focus,.otika-cms .cms-actions button:focus,.otika-cms .cms-actions a:focus{box-shadow:0 2px 6px #acb5f6!important}.otika-cms .sidebar-menu .nav-link{position:relative}.otika-cms .sidebar-mini .main-sidebar,.otika-cms .sidebar-mini #sidebar-wrapper{overflow:visible!important}.otika-cms .sidebar-mini .sidebar-menu{overflow-y:auto!important;overflow-x:visible!important}.otika-cms .sidebar-mini .sidebar-user strong,.otika-cms .sidebar-mini .sidebar-user span,.otika-cms .sidebar-mini .sidebar-menu .menu-header,.otika-cms .sidebar-mini .sidebar-menu .dropdown-menu{display:none!important}.otika-cms .sidebar-mini .sidebar-menu .nav-link span{display:block!important;position:absolute;left:62px;top:50%;z-index:1200;max-width:220px;padding:8px 12px;border-radius:4px;background:#111827;color:#fff!important;box-shadow:0 8px 24px rgba(15,23,42,.22);font-size:12px;line-height:1;opacity:0;pointer-events:none;transform:translateY(-50%) translateX(-6px);transition:opacity 140ms ease,transform 140ms ease;visibility:hidden}.otika-cms .sidebar-mini .sidebar-menu .nav-link:hover span,.otika-cms .sidebar-mini .sidebar-menu .nav-link:focus-visible span{opacity:1;transform:translateY(-50%) translateX(0);visibility:visible}.otika-cms .sidebar-mini .sidebar-menu .nav-link span::before{content:"";position:absolute;left:-5px;top:50%;width:10px;height:10px;background:#111827;transform:translateY(-50%) rotate(45deg)}:global(.dark .otika-cms){background:#0f172a!important;color:#e5e7eb!important}:global(.dark .otika-cms .main-content),:global(.dark .otika-cms .section-body){background:#0f172a!important;color:#e5e7eb!important}:global(.dark .otika-cms :deep(.navbar-bg)){background:#111827!important}:global(.dark .otika-cms :deep(.main-navbar)){background:#111827!important;box-shadow:0 4px 24px rgba(0,0,0,.35)!important}:global(.dark .otika-cms .main-sidebar),:global(.dark .otika-cms #sidebar-wrapper){background:#111827!important}:global(.dark .otika-cms .sidebar-brand),:global(.dark .otika-cms .sidebar-user){background:#1f2937!important}:global(.dark .otika-cms .sidebar-user strong),:global(.dark .otika-cms .section-header h1),:global(.dark .otika-cms .cms-panel h2),:global(.dark .otika-cms .cms-subpanel h3),:global(.dark .otika-cms .permission-group h3),:global(.dark .otika-cms .cms-table-row strong),:global(.dark .otika-cms .audit-row strong),:global(.dark .otika-cms .permission-row strong),:global(.dark .otika-cms .metric strong),:global(.dark .otika-cms label){color:#f8fafc!important}:global(.dark .otika-cms .sidebar-user span),:global(.dark .otika-cms .cms-table-row span),:global(.dark .otika-cms .audit-row span),:global(.dark .otika-cms .permission-row small),:global(.dark .otika-cms .metric span),:global(.dark .otika-cms .cms-empty),:global(.dark .otika-cms .section-header-breadcrumb),:global(.dark .otika-cms .breadcrumb-item){color:#cbd5e1!important}:global(.dark .otika-cms .cms-card),:global(.dark .otika-cms .cms-panel),:global(.dark .otika-cms .cms-subpanel),:global(.dark .otika-cms .permission-group),:global(.dark .otika-cms .audit-summary article),:global(.dark .otika-cms .audit-row),:global(.dark .otika-cms .cms-table-row),:global(.dark .otika-cms .cms-menu-builder-grid article){background:#1f2937!important;border-color:#334155!important;box-shadow:0 4px 25px rgba(0,0,0,.28)!important}:global(.dark .otika-cms input),:global(.dark .otika-cms textarea),:global(.dark .otika-cms select),:global(.dark .otika-cms .form-control),:global(.dark .otika-cms :deep(.rich-editor)),:global(.dark .otika-cms :deep(.ProseMirror)){background:#111827!important;color:#f8fafc!important;border-color:#475569!important}:global(.dark .otika-cms input::placeholder),:global(.dark .otika-cms textarea::placeholder){color:#94a3b8!important}:global(.dark .otika-cms .cms-top-icon),:global(.dark .otika-cms .cms-top-icon i),:global(.dark .otika-cms .search-element .btn i),:global(.dark .otika-cms .navbar .nav-link){color:#f8fafc!important}:global(.dark .otika-cms .dropdown-menu){background:#1f2937!important;border-color:#334155!important;color:#e5e7eb!important}:global(.dark .otika-cms .dropdown-header),:global(.dark .otika-cms .dropdown-title),:global(.dark .otika-cms .dropdown-item),:global(.dark .otika-cms .dropdown-item span){color:#e5e7eb!important}:global(.dark .otika-cms .dropdown-item:hover){background:#111827!important}:global(.dark .otika-cms .sidebar-menu li.active>button),:global(.dark .otika-cms .sidebar-menu li.active>a),:global(.dark .otika-cms .sidebar-menu button:hover){background:#1f2937!important;color:#93c5fd!important}:global(.dark .otika-cms .sidebar-menu li.active>button i),:global(.dark .otika-cms .sidebar-menu button:hover i){color:#93c5fd!important}:global(.dark .otika-cms .cms-message){background:#052e2b!important;color:#6ee7b7!important}:global(.dark .otika-cms .cms-error){background:#3b1218!important;color:#fca5a5!important}:global(.dark .otika-cms .cms-empty){background:#111827!important;border-color:#334155!important}:global(.dark .otika-cms .audit-code.ok){background:#052e2b!important;color:#86efac!important}:global(.dark .otika-cms .audit-code.warn){background:#422006!important;color:#fdba74!important}:global(.dark .otika-cms .audit-code.danger){background:#450a0a!important;color:#fca5a5!important}:global(.dark .otika-cms :deep(.theme-toggle)){color:#f8fafc!important;background:#1f2937!important}:global(.dark .otika-cms :deep(.editor-toolbar button)){background:#111827!important;color:#e5e7eb!important}:global(.dark .otika-cms :deep(.editor-toolbar .active)){background:#6777ef!important;color:#fff!important}
.otika-dashboard .card{border:0!important;border-radius:3px!important;box-shadow:0 4px 25px 0 rgba(0,0,0,.1)!important}.otika-dashboard .card-header{border-bottom-color:#f9f9f9!important}.otika-dashboard .card-header h4{font-size:16px!important;font-weight:700!important;color:#34395e!important}.otika-dashboard .card-statistic-4{position:relative;color:#34395e;padding:15px;border-radius:3px;overflow:hidden}.otika-dashboard .card-statistic-4 .card-content{padding:8px 0 8px 10px}.otika-dashboard .card-statistic-4 h5{color:#6c757d;font-weight:600}.otika-dashboard .card-statistic-4 h2{color:#34395e;font-weight:700}.cms-stat-icon{display:flex!important;align-items:center;justify-content:center;width:72px;height:72px;margin:22px auto 0;border-radius:50%;background:#f4f6f9;font-size:36px}.cms-chart-bars{display:flex;align-items:flex-end;justify-content:space-between;gap:16px;min-height:260px;padding:12px 4px}.cms-chart-bar{display:grid;grid-template-rows:auto 1fr auto;gap:8px;min-width:44px;height:250px;text-align:center;color:#6c757d}.cms-chart-bar__value{font-size:12px;font-weight:700;color:#34395e}.cms-chart-bar__track{display:flex;align-items:flex-end;width:100%;height:190px;border-radius:30px;background:#f4f6f9;overflow:hidden}.cms-chart-bar__track span{display:block;width:100%;border-radius:30px 30px 0 0;background:#6777ef}.cms-chart-bar strong{font-size:12px}.cms-kpi-row,.cms-source-row,.cms-progress-item>div:first-child{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:12px;color:#6c757d}.cms-kpi-row strong,.cms-source-row strong,.cms-progress-item strong{color:#34395e}.cms-progress-item{margin-bottom:18px}.cms-progress-item .progress,.otika-dashboard .progress{height:6px!important;border-radius:30px;background:#f4f6f9}.cms-source-row{padding:10px 0;border-bottom:1px solid #f4f6f9}.cms-source-row i{width:22px;color:#6777ef}.cms-donut{--first:44%;--second:72%;display:grid;place-content:center;width:190px;height:190px;margin:0 auto 18px;border-radius:50%;background:conic-gradient(#6777ef 0 var(--first),#47c363 var(--first) var(--second),#ffa426 var(--second) 100%);color:#34395e;position:relative}.cms-donut::before{content:"";position:absolute;inset:28px;border-radius:50%;background:#fff}.cms-donut span,.cms-donut small{position:relative;z-index:1;text-align:center}.cms-donut span{font-size:28px;font-weight:800}.cms-donut small{font-size:12px;color:#6c757d}.cms-donut-legend{display:grid;gap:8px}.cms-donut-legend span{display:flex;align-items:center;gap:8px;color:#6c757d}.cms-donut-legend i{display:inline-block;width:10px;height:10px;border-radius:50%}.otika-dashboard .table td,.otika-dashboard .table th{vertical-align:middle}.otika-dashboard .order-list{display:flex;align-items:center}.otika-dashboard .team-member img{width:32px;height:32px;object-fit:cover}.cms-inline-link{border:0;background:transparent;color:#6777ef;font-weight:700;padding:0 0 0 8px}.cms-analytics-snapshot{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:18px}.cms-analytics-snapshot article{display:flex;align-items:center;gap:12px;padding:18px;border-radius:3px;background:#f4f6f9}.cms-analytics-snapshot i{font-size:30px}.cms-analytics-snapshot span{display:block;color:#6c757d;font-size:12px}.cms-analytics-snapshot strong{display:block;color:#34395e;font-size:20px}.cms-main-footer{margin-top:0!important}:global(.dark) .otika-dashboard .card-header h4,:global(.dark) .otika-dashboard .card-statistic-4 h2,:global(.dark) .cms-kpi-row strong,:global(.dark) .cms-source-row strong,:global(.dark) .cms-progress-item strong,:global(.dark) .cms-chart-bar__value,:global(.dark) .cms-donut span,:global(.dark) .cms-analytics-snapshot strong{color:#f8fafc!important}:global(.dark) .otika-dashboard .card-statistic-4 h5,:global(.dark) .cms-kpi-row,:global(.dark) .cms-source-row,:global(.dark) .cms-progress-item>div:first-child,:global(.dark) .cms-chart-bar,:global(.dark) .cms-donut small,:global(.dark) .cms-donut-legend span,:global(.dark) .cms-analytics-snapshot span{color:#cbd5e1!important}:global(.dark) .otika-dashboard .card{background:#1f2937!important;border-color:#334155!important;box-shadow:0 4px 25px rgba(0,0,0,.28)!important}:global(.dark) .cms-donut::before,:global(.dark) .cms-stat-icon,:global(.dark) .cms-chart-bar__track,:global(.dark) .cms-analytics-snapshot article{background:#111827!important}@media(max-width:768px){.cms-chart-bars{gap:8px;overflow-x:auto}.cms-chart-bar{min-width:40px}.cms-analytics-snapshot{grid-template-columns:1fr}.cms-main-footer{display:block;text-align:center}.cms-main-footer .footer-right{float:none;margin-top:6px}}
.otika-cms .navbar-bg{position:fixed!important;top:0!important;right:0!important;left:260px!important;z-index:1030!important;height:70px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.otika-cms .main-navbar{position:fixed!important;top:0!important;right:0!important;left:260px!important;z-index:1040!important;min-height:70px!important;background:#fff!important;color:#111827!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.otika-cms .main-navbar .nav-link,.otika-cms .main-navbar i{color:#111827!important}.otika-cms .main-wrapper{min-height:100vh!important;display:flex!important;flex-direction:column!important}.otika-cms .main-content{flex:1 0 auto!important;padding-top:92px!important}.otika-cms .cms-main-footer{flex:0 0 auto!important;margin-top:auto!important;border-top:1px solid #e4e6fc!important;background:#fff!important;color:#6c757d!important}.otika-cms .cms-main-footer .footer-left{font-weight:700;color:#34395e!important}.otika-cms .main-wrapper.sidebar-mini .navbar-bg,.otika-cms .main-wrapper.sidebar-mini .main-navbar{left:65px!important}.otika-cms .navbar .dropdown-menu.show{z-index:1060!important}.otika-basic-table-card{border:0!important;border-radius:3px!important;box-shadow:0 4px 25px rgba(0,0,0,.1)!important}.otika-basic-table-card .table th,.otika-basic-table-card .table td,.newsletter-table .table th,.newsletter-table .table td,.otika-dashboard .table th,.otika-dashboard .table td{vertical-align:middle!important}.otika-basic-table-card .table thead th,.newsletter-table .table thead th,.otika-dashboard .table thead th{border-bottom:1px solid #f4f6f9!important;color:#34395e!important;font-weight:700!important}.otika-basic-table-card .table strong{color:#34395e!important}.otika-basic-table-card .table span{color:#6c757d!important}.table-actions{white-space:nowrap}.cms-table{display:block!important;overflow-x:auto!important}.cms-table-row{display:grid!important;grid-template-columns:minmax(0,1fr) auto!important;align-items:center!important;margin-bottom:0!important;border:0!important;border-bottom:1px solid #f4f6f9!important;border-radius:0!important;background:#fff!important;padding:14px 18px!important;box-shadow:none!important}.cms-table-row:nth-child(odd){background:#fbfbfd!important}.cms-table.compact{border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important;overflow:auto!important}:global(.dark .otika-cms .navbar-bg),:global(.dark .otika-cms .main-navbar){background:#111827!important;color:#f8fafc!important;box-shadow:0 4px 24px rgba(0,0,0,.35)!important}:global(.dark .otika-cms .main-navbar .nav-link),:global(.dark .otika-cms .main-navbar i){color:#f8fafc!important}:global(.dark .otika-cms .cms-main-footer),:global(.dark) .otika-basic-table-card,:global(.dark) .cms-table.compact{background:#1f2937!important;color:#cbd5e1!important;border-color:#334155!important}:global(.dark .otika-cms .cms-main-footer .footer-left),:global(.dark) .otika-basic-table-card .table thead th,:global(.dark) .newsletter-table .table thead th,:global(.dark) .otika-dashboard .table thead th,:global(.dark) .otika-basic-table-card .table strong{color:#f8fafc!important}:global(.dark) .otika-basic-table-card .table span{color:#cbd5e1!important}:global(.dark) .cms-table-row{background:#1f2937!important;border-color:#334155!important}:global(.dark) .cms-table-row:nth-child(odd){background:#111827!important}@media(max-width:991px){.otika-cms .navbar-bg,.otika-cms .main-navbar,.otika-cms .main-wrapper.sidebar-mini .navbar-bg,.otika-cms .main-wrapper.sidebar-mini .main-navbar{left:0!important}.otika-cms .main-content{padding-top:116px!important}}@media(max-width:575px){.table-actions .btn{display:inline-flex;margin-bottom:4px}.otika-cms .cms-main-footer{text-align:center!important}.otika-cms .cms-main-footer .footer-right{float:none!important;margin-top:6px!important}}
.homepage-card-list{display:grid;gap:16px}.homepage-dynamic-card{display:grid;gap:12px;border:1px solid #e4e6fc;border-radius:3px;background:#fdfdff;padding:18px;box-shadow:0 4px 25px rgba(0,0,0,.04)}.homepage-dynamic-card .cms-panel-head{margin-bottom:0}.homepage-dynamic-card h3,.cms-subpanel h3{font-size:15px!important;color:#34395e!important}.cms-row input.form-control{height:42px;border:1px solid #e4e6fc;border-radius:3px;background:#fff;color:#34395e;padding:10px 15px}.cms-check{display:flex!important;grid-column:1/-1;align-items:center!important;gap:10px!important;text-transform:none!important}.cms-check input{width:18px!important;height:18px!important;accent-color:#6777ef}:global(.dark) .homepage-dynamic-card{background:#111827!important;border-color:#334155!important}:global(.dark) .homepage-dynamic-card h3,:global(.dark) .cms-subpanel h3{color:#f8fafc!important}:global(.dark) .cms-row input.form-control{background:#111827!important;color:#f8fafc!important;border-color:#475569!important}
.otika-form-card .card{border:0!important;border-radius:3px!important;box-shadow:none!important;margin-bottom:0!important}.otika-form-card .card-header{border-bottom:1px solid #f9f9f9!important;padding:18px 25px!important}.otika-form-card .card-header h4{font-size:16px!important;font-weight:700!important;color:#34395e!important;margin:0!important}.otika-form-card .card-body{padding:25px!important}.otika-form-card .card-footer{border-top:1px solid #f9f9f9!important;background:#fff!important;padding:18px 25px!important}.otika-form-card .section-title{margin:18px 0 16px!important;font-size:13px!important;font-weight:700!important;color:#34395e!important}.otika-form-card .form-group{margin-bottom:18px!important}.otika-form-card label{font-size:12px!important;font-weight:600!important;color:#34395e!important;margin-bottom:7px!important}.otika-form-card .form-control{height:42px!important;border:1px solid #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#495057!important;padding:10px 15px!important;box-shadow:none!important}.otika-form-card .form-control:focus{border-color:#6777ef!important;box-shadow:0 2px 6px #acb5f6!important}.otika-form-card textarea.form-control,.otika-form-card .otika-textarea{height:auto!important;min-height:110px!important}.otika-form-card .custom-control-label{line-height:1.8!important}.otika-form-card .btn-primary{border:0!important;border-radius:30px!important;background:#6777ef!important;box-shadow:0 2px 6px #acb5f6!important;font-size:12px!important;font-weight:600!important;padding:8px 18px!important}.otika-form-card .form-check{padding-left:1.5rem;min-height:auto;margin-bottom:.35rem}.otika-form-card .form-check-input{width:16px;height:16px;accent-color:#6777ef}.otika-form-card .form-check-label{font-weight:500!important;margin-bottom:0!important;line-height:1.6}:global(.dark) .otika-form-card .card,:global(.dark) .otika-form-card .card-footer{background:#1f2937!important;border-color:#334155!important}:global(.dark) .otika-form-card .card-header h4,:global(.dark) .otika-form-card .section-title,:global(.dark) .otika-form-card label{color:#f8fafc!important}:global(.dark) .otika-form-card .form-control{background:#111827!important;color:#f8fafc!important;border-color:#475569!important}
.notification-list{display:grid;gap:12px}.notification-item{display:grid;grid-template-columns:auto minmax(0,1fr) auto;gap:14px;align-items:flex-start;padding:16px 18px;border-radius:3px;background:#fff;box-shadow:0 4px 25px rgba(0,0,0,.08);border-left:3px solid transparent;color:#6c757d}.notification-item.unread{background:#f4f7ff;border-left-color:#6777ef}.notification-item strong{display:block;color:#34395e;font-size:14px}.notification-item span{display:block;color:#98a6ad;font-size:12px}.notification-item p{margin:6px 0 0;color:#6c757d}.notification-dot{width:10px;height:10px;margin-top:7px;border-radius:50%;background:#6777ef}.notification-icon{display:inline-flex;align-items:center;justify-content:center;width:38px;height:38px;border-radius:50%;background:#eaf4ff;color:#3abaf4}.notification-icon.warn{background:#fff4e6;color:#ffa426}.notification-icon.success{background:#e8f7f0;color:#47c363}.notification-icon.info{background:#eaf4ff;color:#3abaf4}.notification-item:not(.unread){opacity:.82}:global(.dark) .notification-item{background:#1f2937;box-shadow:0 4px 25px rgba(0,0,0,.28)}:global(.dark) .notification-item.unread{background:#111827}:global(.dark) .notification-item strong{color:#f8fafc}:global(.dark) .notification-item p{color:#cbd5e1}@media(max-width:768px){.notification-item{grid-template-columns:auto minmax(0,1fr)}.notification-item>.cms-actions-inline{grid-column:1/-1}}
</style>
