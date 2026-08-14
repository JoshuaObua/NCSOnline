<template>
  <LayoutDefault title="Dashboard Overview">

    <!-- Loading Skeleton -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div v-for="i in 4" :key="i" class="card card-statistic-4 animate-pulse">
        <div class="h-20 bg-slate-100 dark:bg-slate-800 rounded-lg"></div>
      </div>
    </div>

    <!-- ════════════════════════════════════════════════════════════════════ -->
    <!-- A. RECEPTIONIST / FRONT DESK DASHBOARD VIEW                          -->
    <!-- ════════════════════════════════════════════════════════════════════ -->
    <div v-else-if="isReceptionist" class="space-y-6">
      
      <!-- 1. Hero Welcome Card -->
      <div class="card bg-gradient-to-r from-blue-900 via-indigo-900 to-slate-900 border-0 text-white shadow-md rounded-2xl overflow-hidden p-6 sm:p-8">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div class="text-xs uppercase tracking-widest text-blue-300 font-bold mb-1">{{ greeting }}</div>
            <h2 class="text-2xl font-extrabold">{{ welcomeName }}</h2>
            <p class="text-slate-300 text-sm mt-1 max-w-xl">
              Front Desk & Visitor Clearance Operations — National Council of Sports.
            </p>
          </div>
          <div class="flex items-center gap-2">
            <span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-500/20 text-emerald-300 text-xs font-semibold border border-emerald-500/30">
              <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
              Front Desk Active
            </span>
          </div>
        </div>
      </div>

      <!-- 2. Receptionist 4 KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Card 1: Visitors Registered -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Total Visitors</h5>
              <h2>{{ receptionVisitors.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-users"></i> Ingested Total</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-id-card"></i>
            </div>
          </div>
        </div>

        <!-- Card 2: Pending Host Approvals -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Pending Approval</h5>
              <h2>{{ pendingVisitorsCount }}</h2>
              <span class="badge" :class="pendingVisitorsCount > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-clock-time"></i> Awaiting Host
              </span>
            </div>
            <div class="banner-img bg-warning-light">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>

        <!-- Card 3: Checked In (On Premises) -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>On Premises</h5>
              <h2>{{ checkedInVisitorsCount }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Active Passes</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-badge"></i>
            </div>
          </div>
        </div>

        <!-- Card 4: Staff Leave Entitlement -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Leave Balance</h5>
              <h2>{{ remainingLeaveDays }} Days</h2>
              <span class="badge badge-info"><i class="icofont-calendar"></i> 21 Annual Days</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-sun-alt"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. Front Desk Quick Action Navigation Grid -->
      <div class="card">
        <div class="card-header">
          <h4><i class="icofont-navigation-menu text-blue-600"></i> Front Desk Quick Actions</h4>
        </div>
        <div class="card-body">
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <router-link
              to="/reception/visitors"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-blue-100 dark:bg-blue-900/50 text-blue-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-id-card"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Visitor Clearance Desk</span>
            </router-link>

            <router-link
              to="/reception/leave/apply"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-calendar"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Apply for Leave</span>
            </router-link>

            <router-link
              to="/reception/leave/status"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-purple-100 dark:bg-purple-900/50 text-purple-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-history"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Leave Status & History</span>
            </router-link>

            <router-link
              to="/reception/reports"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-amber-100 dark:bg-amber-900/50 text-amber-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-chart-bar-graph"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Front Desk Reports</span>
            </router-link>
          </div>
        </div>
      </div>

      <!-- 4. Active Visitor Inflow & Leave Status Tables -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <!-- Left: Recent Visitor Passes -->
        <div class="card lg:col-span-2">
          <div class="card-header flex items-center justify-between">
            <h4><i class="icofont-users text-blue-600"></i> Recent Visitors Ingested</h4>
            <router-link to="/reception/visitors" class="btn btn-sm btn-outline-primary">
              View All <i class="icofont-arrow-right ml-1"></i>
            </router-link>
          </div>
          <div class="card-body p-0">
            <div class="table-responsive">
              <table class="table table-striped table-hover mb-0 text-xs">
                <thead>
                  <tr>
                    <th>Pass No.</th>
                    <th>Visitor Name</th>
                    <th>Target Department</th>
                    <th>Host Officer</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="v in receptionVisitors.slice(0, 5)" :key="v.id">
                    <td class="font-bold text-blue-600 font-mono">{{ v.pass_number }}</td>
                    <td>
                      <div class="font-medium text-slate-900 dark:text-white">{{ v.visitor_name }}</div>
                      <div class="text-[10px] text-slate-400">{{ v.visitor_organization || v.visitor_phone }}</div>
                    </td>
                    <td>{{ v.target_department }}</td>
                    <td>{{ v.host_officer_name }}</td>
                    <td>
                      <span class="badge" :class="statusBadge(v.status)">{{ v.status }}</span>
                    </td>
                  </tr>
                  <tr v-if="!receptionVisitors.length">
                    <td colspan="5" class="text-center py-6 text-xs text-slate-400">
                      No visitors registered yet today.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- Right: Leave Entitlement Summary -->
        <div class="card">
          <div class="card-header">
            <h4><i class="icofont-calendar text-emerald-600"></i> Leave Entitlement</h4>
          </div>
          <div class="card-body space-y-4">
            <div class="p-4 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-center">
              <div class="text-[11px] font-bold text-slate-400 uppercase">Statutory Annual Leave</div>
              <div class="text-2xl font-extrabold text-slate-900 dark:text-white mt-1">
                21 Days
              </div>
              <div class="text-xs text-emerald-500 font-semibold mt-1">
                {{ remainingLeaveDays }} Days Available Balance
              </div>
            </div>

            <div class="space-y-2 text-xs">
              <div class="flex justify-between p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
                <span class="text-slate-500">Days Utilized:</span>
                <span class="font-bold text-slate-900 dark:text-white">{{ utilizedLeaveDays }} Days</span>
              </div>
              <div class="flex justify-between p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
                <span class="text-slate-500">Relief Officer:</span>
                <span class="font-bold text-slate-900 dark:text-white">{{ latestLeave?.relieving_officer_name || 'Designated on Apply' }}</span>
              </div>
              <div class="flex justify-between p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
                <span class="text-slate-500">Latest Request Status:</span>
                <span class="badge" :class="statusBadge(latestLeave?.status || 'NONE')">
                  {{ latestLeave?.status || 'No Active Request' }}
                </span>
              </div>
            </div>

            <router-link to="/reception/leave/apply" class="btn btn-sm btn-primary w-full justify-center">
              <i class="icofont-plus mr-1"></i> Apply for Leave
            </router-link>
          </div>
        </div>

      </div>

    </div>

    <!-- ════════════════════════════════════════════════════════════════════ -->
    <!-- B. IT & ENGINEERING WORKSTATION DASHBOARD VIEW                        -->
    <!-- ════════════════════════════════════════════════════════════════════ -->
    <div v-else-if="isDepartmentalOfficer" class="space-y-6">
      
      <!-- 1. Hero Welcome Card -->
      <div class="card bg-gradient-to-r from-blue-900 via-indigo-900 to-slate-900 border-0 text-white shadow-md rounded-2xl overflow-hidden p-6 sm:p-8">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div class="text-xs uppercase tracking-widest text-blue-300 font-bold mb-1">{{ greeting }}</div>
            <h2 class="text-2xl font-extrabold">{{ welcomeName }}</h2>
            <p class="text-slate-300 text-sm mt-1 max-w-xl">
              {{ departmentalWorkstationTitle }}
            </p>
          </div>
          <div class="flex items-center gap-2">
            <span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-500/20 text-emerald-300 text-xs font-semibold border border-emerald-500/30">
              <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
              {{ departmentalBadgeText }}
            </span>
          </div>
        </div>
      </div>

      <!-- 2. IT Officer 4 KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Card 1: Total PPDA Form 5 Requisitions -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Total Requisitions</h5>
              <h2>{{ itRequisitions.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-document-folder"></i> PPDA Form 5</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-law-order"></i>
            </div>
          </div>
        </div>

        <!-- Card 2: Pending Approvals -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Pending Approval</h5>
              <h2>{{ itPendingCount }}</h2>
              <span class="badge" :class="itPendingCount > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-clock-time"></i> In Review Pipeline
              </span>
            </div>
            <div class="banner-img bg-warning-light">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>

        <!-- Card 3: Approved / PO Issued -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Approved & Contracted</h5>
              <h2>{{ itApprovedCount }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Authorized</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-certificate"></i>
            </div>
          </div>
        </div>

        <!-- Card 4: Staff Leave Entitlement -->
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Leave Balance</h5>
              <h2>{{ remainingLeaveDays }} Days</h2>
              <span class="badge badge-info"><i class="icofont-calendar"></i> 21 Annual Days</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-sun-alt"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. IT Quick Actions Navigation Grid -->
      <div class="card">
        <div class="card-header">
          <h4><i class="icofont-navigation-menu text-blue-600"></i> ICT Officer Quick Actions</h4>
        </div>
        <div class="card-body">
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <router-link
              to="/it/ppda/new"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-blue-100 dark:bg-blue-900/50 text-blue-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-plus-circle"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">New PPDA Form 5</span>
            </router-link>

            <router-link
              to="/it/ppda/status"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-purple-100 dark:bg-purple-900/50 text-purple-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-document-folder"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">PPDA Form 5 Status</span>
            </router-link>

            <router-link
              to="/reception/leave/apply"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-calendar"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Apply for Leave</span>
            </router-link>

            <router-link
              to="/reception/leave/status"
              class="flex flex-col items-center gap-2 p-4 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-11 h-11 rounded-xl bg-amber-100 dark:bg-amber-900/50 text-amber-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-history"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Leave Status & History</span>
            </router-link>
          </div>
        </div>
      </div>

      <!-- 4. Active PPDA Form 5 Pipeline & Leave Entitlement Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <!-- Left: Recent PPDA Form 5 Requisitions -->
        <div class="card lg:col-span-2">
          <div class="card-header flex items-center justify-between">
            <h4><i class="icofont-document-folder text-blue-600"></i> Active PPDA Form 5 Requisitions</h4>
            <router-link to="/it/ppda/status" class="btn btn-sm btn-outline-primary">
              View All <i class="icofont-arrow-right ml-1"></i>
            </router-link>
          </div>
          <div class="card-body p-0">
            <div class="table-responsive">
              <table class="table table-striped table-hover mb-0 text-xs">
                <thead>
                  <tr>
                    <th>Ref No.</th>
                    <th>Subject of Procurement</th>
                    <th>Est. Amount</th>
                    <th>Current Stage</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="req in itRequisitions.slice(0, 5)" :key="req.id">
                    <td class="font-bold text-blue-600 font-mono">{{ req.reference_no }}</td>
                    <td>
                      <div class="font-medium text-slate-900 dark:text-white truncate max-w-xs">{{ req.subject_of_procurement }}</div>
                      <div class="text-[10px] text-slate-400">{{ req.procurement_category }} · {{ req.budget_vote_head }}</div>
                    </td>
                    <td class="font-mono font-bold text-slate-900 dark:text-white">
                      UGX {{ Number(req.estimated_amount_ugx || 0).toLocaleString() }}
                    </td>
                    <td class="text-[11px] text-slate-500">
                      {{ req.current_stage || 'In Review' }}
                    </td>
                    <td>
                      <span class="badge" :class="statusBadge(req.status)">{{ req.status }}</span>
                    </td>
                  </tr>
                  <tr v-if="!itRequisitions.length">
                    <td colspan="5" class="text-center py-6 text-xs text-slate-400">
                      No procurement requisitions submitted yet.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- Right: Leave Entitlement Summary -->
        <div class="card">
          <div class="card-header">
            <h4><i class="icofont-calendar text-emerald-600"></i> Leave Entitlement</h4>
          </div>
          <div class="card-body space-y-4">
            <div class="p-4 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-center">
              <div class="text-[11px] font-bold text-slate-400 uppercase">Statutory Annual Leave</div>
              <div class="text-2xl font-extrabold text-slate-900 dark:text-white mt-1">
                21 Days
              </div>
              <div class="text-xs text-emerald-500 font-semibold mt-1">
                {{ remainingLeaveDays }} Days Available Balance
              </div>
            </div>

            <div class="space-y-2 text-xs">
              <div class="flex justify-between p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
                <span class="text-slate-500">Days Utilized:</span>
                <span class="font-bold text-slate-900 dark:text-white">{{ utilizedLeaveDays }} Days</span>
              </div>
              <div class="flex justify-between p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
                <span class="text-slate-500">Relief Officer:</span>
                <span class="font-bold text-slate-900 dark:text-white">{{ latestLeave?.relieving_officer_name || 'Designated on Apply' }}</span>
              </div>
              <div class="flex justify-between p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
                <span class="text-slate-500">Latest Request Status:</span>
                <span class="badge" :class="statusBadge(latestLeave?.status || 'NONE')">
                  {{ latestLeave?.status || 'No Active Request' }}
                </span>
              </div>
            </div>

            <router-link to="/reception/leave/apply" class="btn btn-sm btn-primary w-full justify-center">
              <i class="icofont-plus mr-1"></i> Apply for Leave
            </router-link>
          </div>
        </div>

      </div>

    </div>

    <!-- ════════════════════════════════════════════════════════════════════ -->
    <!-- C. EXECUTIVE & GENERAL OVERVIEW DASHBOARD VIEW                       -->
    <!-- ════════════════════════════════════════════════════════════════════ -->
    <div v-else class="space-y-6">

      <!-- ── 1. Otika Hero Welcome Card ──────────────────────────────── -->
      <div class="card bg-gradient-to-r from-blue-900 via-indigo-900 to-slate-900 border-0 text-white shadow-md rounded-2xl overflow-hidden p-6 sm:p-8">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div class="text-xs uppercase tracking-widest text-blue-300 font-bold mb-1">{{ greeting }}</div>
            <h2 class="text-2xl font-extrabold">{{ welcomeName }}!</h2>
            <p class="text-slate-300 text-sm mt-1 max-w-xl">
              Welcome to the National Council of Sports (NCS) Intranet Command Portal.
            </p>
          </div>
          <div class="flex items-center gap-2">
            <span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-500/20 text-emerald-300 text-xs font-semibold border border-emerald-500/30">
              <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
              All Systems Operational
            </span>
          </div>
        </div>
      </div>

      <!-- ── 2. Otika 4-Statistic KPI Cards ──────────────────────────── -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Total Applications</h5>
              <h2>{{ stats?.total_applications ?? 0 }}</h2>
              <span class="badge badge-primary"><i class="icofont-files-stack"></i> Active Pipeline</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-files-stack"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Pending Review</h5>
              <h2>{{ stats?.pending_review ?? 0 }}</h2>
              <span class="badge badge-warning"><i class="icofont-clock-time"></i> Action Required</span>
            </div>
            <div class="banner-img bg-warning-light">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Active Federations</h5>
              <h2>{{ stats?.active_federations ?? 0 }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Statutory Act 2023</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-trophy"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm hover:shadow-md transition">
          <div class="card-content">
            <div>
              <h5>Portal Users</h5>
              <h2>{{ stats?.total_users ?? 0 }}</h2>
              <span class="badge badge-info"><i class="icofont-users-alt-5"></i> Verified Accounts</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-users"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- ── 3. Quick Departmental Directory ────────────────────────── -->
      <div class="card">
        <div class="card-header">
          <h4><i class="icofont-ui-office text-blue-600"></i> Departmental Operations Hub</h4>
        </div>
        <div class="card-body">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
            <router-link
              to="/executive/appraisals"
              class="flex flex-col items-center gap-2 p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-10 h-10 rounded-xl bg-blue-100 dark:bg-blue-900/50 text-blue-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-law-order"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Online Appraisal</span>
            </router-link>

            <router-link
              to="/executive/ags-technical"
              class="flex flex-col items-center gap-2 p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-10 h-10 rounded-xl bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-badge"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">AGS Technical</span>
            </router-link>

            <router-link
              to="/executive/ags-admin"
              class="flex flex-col items-center gap-2 p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-10 h-10 rounded-xl bg-purple-100 dark:bg-purple-900/50 text-purple-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-architecture-alt"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">AGS Admin</span>
            </router-link>

            <router-link
              to="/reception/visitors"
              class="flex flex-col items-center gap-2 p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-10 h-10 rounded-xl bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-id-card"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Reception & Visitors</span>
            </router-link>

            <router-link
              to="/stores/inventory"
              class="flex flex-col items-center gap-2 p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 border border-slate-200 dark:border-slate-700 transition group text-center"
            >
              <div class="w-10 h-10 rounded-xl bg-amber-100 dark:bg-amber-900/50 text-amber-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">
                <i class="icofont-box"></i>
              </div>
              <span class="text-xs font-bold text-slate-800 dark:text-slate-200 group-hover:text-blue-600">Stores & Stock</span>
            </router-link>
          </div>
        </div>
      </div>

      <!-- ── 4. Tables & Analytics Split ────────────────────────────── -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <!-- Left: Recent Applications Table -->
        <div class="card lg:col-span-2">
          <div class="card-header">
            <h4><i class="icofont-history text-blue-600"></i> Recent Licence Submissions</h4>
            <div class="card-header-action">
              <router-link to="/applications" class="btn btn-sm btn-outline-primary">
                View All <i class="icofont-arrow-right ml-1"></i>
              </router-link>
            </div>
          </div>
          <div class="card-body p-0">
            <div class="table-responsive">
              <table class="table table-striped table-hover mb-0">
                <thead>
                  <tr>
                    <th>Ref No.</th>
                    <th>Applicant / Federation</th>
                    <th>Type</th>
                    <th>Status</th>
                    <th>Submitted</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="app in recentApplicationsList" :key="app.id">
                    <td>
                      <span class="font-bold text-slate-800 dark:text-slate-200">{{ app.reference_no || app.id?.substring(0,8) }}</span>
                    </td>
                    <td>
                      <div class="font-medium text-slate-900 dark:text-white">{{ app.organisation_name || app.applicant_name || 'NCS Applicant' }}</div>
                      <div class="text-xs text-slate-400">{{ app.email || '' }}</div>
                    </td>
                    <td>
                      <span class="badge badge-light">{{ app.application_type || app.form_type || 'Licence' }}</span>
                    </td>
                    <td>
                      <span class="badge" :class="statusBadge(app.status)">
                        {{ app.status || 'Submitted' }}
                      </span>
                    </td>
                    <td class="text-xs text-slate-500">
                      {{ formatDate(app.created_at) }}
                    </td>
                  </tr>
                  <tr v-if="!recentApplicationsList.length">
                    <td colspan="5" class="text-center py-6 text-xs text-slate-400">
                      No applications submitted yet.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- Right: Fixed Asset Ledger Status Card -->
        <div class="card">
          <div class="card-header">
            <h4><i class="icofont-coins text-amber-500"></i> Asset Ledger Summary</h4>
          </div>
          <div class="card-body">
            <div class="p-4 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-center mb-4">
              <div class="text-[11px] font-bold text-slate-400 uppercase">Total Portfolio Valuation</div>
              <div class="text-xl font-extrabold text-slate-900 dark:text-white mt-1">
                UGX {{ Number(stats?.total_portfolio_valuation_ugx || 0).toLocaleString() }}
              </div>
              <div class="text-[11px] text-emerald-500 font-semibold mt-0.5">
                {{ stats?.total_asset_items_count || 0 }} Fixed Assets Registered
              </div>
            </div>

            <div class="space-y-3">
              <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs">
                <span class="text-slate-600 dark:text-slate-400">Institutional Land & Venues</span>
                <span class="font-bold text-slate-900 dark:text-white">Active</span>
              </div>
              <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs">
                <span class="text-slate-600 dark:text-slate-400">Arena, Hostels & Infrastructure</span>
                <span class="font-bold text-slate-900 dark:text-white">Active</span>
              </div>
              <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs">
                <span class="text-slate-600 dark:text-slate-400">Government Transport Fleet</span>
                <span class="font-bold text-slate-900 dark:text-white">Active</span>
              </div>
            </div>

            <div class="mt-5 pt-4 border-t border-slate-100 dark:border-slate-800">
              <router-link to="/executive/appraisals" class="btn btn-sm btn-primary w-full justify-center">
                Open Online Appraisal
              </router-link>
            </div>
          </div>
        </div>

      </div>

      <!-- ── 5. Website Traffic & Public Portal Analytics ────────────────── -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h4 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
              <i class="icofont-chart-growth text-blue-600"></i> Website & Public Portal Traffic Analytics
            </h4>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Live engagement metrics from the official National Council of Sports web portal</p>
          </div>
          <div class="flex items-center gap-2 text-xs">
            <span class="badge badge-success flex items-center gap-1"><span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span> Live Analytics</span>
          </div>
        </div>

        <div class="card-body p-5 space-y-6">
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
              <span class="text-[10px] uppercase font-bold text-slate-400 block">Total Portal Visits</span>
              <strong class="text-lg font-bold text-slate-900 dark:text-white font-mono">{{ Number(stats?.website_traffic?.total_visits || 0).toLocaleString() }}</strong>
              <div class="text-[10px] text-emerald-500 font-semibold mt-0.5"><i class="icofont-arrow-up"></i> +14.2% this month</div>
            </div>

            <div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
              <span class="text-[10px] uppercase font-bold text-slate-400 block">Unique Visitors</span>
              <strong class="text-lg font-bold text-slate-900 dark:text-white font-mono">{{ Number(stats?.website_traffic?.unique_visitors || 0).toLocaleString() }}</strong>
              <div class="text-[10px] text-blue-500 font-semibold mt-0.5">Verified IP Sessions</div>
            </div>

            <div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
              <span class="text-[10px] uppercase font-bold text-slate-400 block">Total Page Views</span>
              <strong class="text-lg font-bold text-slate-900 dark:text-white font-mono">{{ Number(stats?.website_traffic?.page_views || 0).toLocaleString() }}</strong>
              <div class="text-[10px] text-purple-500 font-semibold mt-0.5">3.8 pages / session</div>
            </div>

            <div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
              <span class="text-[10px] uppercase font-bold text-slate-400 block">Avg Session Duration</span>
              <strong class="text-lg font-bold text-slate-900 dark:text-white font-mono">{{ stats?.website_traffic?.avg_session_duration_sec || 0 }}s</strong>
              <div class="text-[10px] text-emerald-500 font-semibold mt-0.5">Bounce Rate: {{ stats?.website_traffic?.bounce_rate_pct || 0 }}%</div>
            </div>
          </div>

          <div>
            <h5 class="text-xs font-bold text-slate-800 dark:text-slate-200 mb-3">Top Visited Public Pages & Directories</h5>
            <div class="overflow-x-auto">
              <table class="table table-sm table-striped text-xs text-left mb-0">
                <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-500 uppercase tracking-wider font-semibold">
                  <tr>
                    <th class="py-2.5 px-3">Page Title / Topic</th>
                    <th class="py-2.5 px-3">Route Path</th>
                    <th class="py-2.5 px-3 font-mono">Pageviews</th>
                    <th class="py-2.5 px-3">Traffic Share</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                  <tr v-for="page in stats?.website_traffic?.top_pages || []" :key="page.path">
                    <td class="py-2.5 px-3 font-bold text-slate-900 dark:text-white">{{ page.title }}</td>
                    <td class="py-2.5 px-3 font-mono text-slate-500">{{ page.path }}</td>
                    <td class="py-2.5 px-3 font-mono font-bold text-blue-600 dark:text-blue-400">{{ Number(page.views).toLocaleString() }}</td>
                    <td class="py-2.5 px-3">
                      <div class="flex items-center gap-2">
                        <div class="w-24 bg-slate-100 dark:bg-slate-700 h-1.5 rounded-full overflow-hidden">
                          <div class="bg-blue-600 h-full rounded-full" :style="{ width: page.percentage + '%' }"></div>
                        </div>
                        <span class="text-[10px] font-mono text-slate-500">{{ page.percentage }}%</span>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'

const authStore = useAuthStore()
const loading = ref(false)
const stats = ref(null)
const recentApplicationsList = ref([])

// Receptionist state
const receptionVisitors = ref([])
const leaveApplications = ref([])

// IT Officer state
const itRequisitions = ref([])
const itDashboardStats = ref(null)

const isReceptionist = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isReception = roles.includes('receptionist') || roles.includes('helpdesk')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary')
  return isReception && !isPrivileged
})

const isITOfficer = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isIT = roles.includes('it_officer')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary') || roles.includes('ags_technical') || roles.includes('ags_admin')
  return isIT && !isPrivileged
})

const isEngineeringOfficer = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isEng = roles.includes('senior_engineer') ||
    roles.includes('assistant_engineer_civil') ||
    roles.includes('assistant_engineer_electrical') ||
    roles.includes('engineering_officer_civil') ||
    roles.includes('engineering_officer_electrical') ||
    roles.includes('plumber')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary') || roles.includes('ags_technical') || roles.includes('ags_admin')
  return isEng && !isPrivileged
})

const isDepartmentalOfficer = computed(() => {
  return isITOfficer.value || isEngineeringOfficer.value
})

const departmentalWorkstationTitle = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  
  if (roles.includes('senior_engineer')) {
    return 'Senior Infrastructure Directorate Lead — CapEx Projects, FIBA Arena Standards & HOD Vetting'
  }
  if (roles.includes('assistant_engineer_civil')) {
    return 'Civil & Structural Engineering Section — Arena Hardwood, Concrete Pitch Infrastructure & Drainage'
  }
  if (roles.includes('assistant_engineer_electrical')) {
    return 'Electrical & Power Engineering Section — 1,200 Lux Arena Floodlights, Transformers & Power Systems'
  }
  if (roles.includes('engineering_officer_civil')) {
    return 'Civil Works Operations & Inspection Unit — Maintenance, Pavements & Structural Surveys'
  }
  if (roles.includes('engineering_officer_electrical')) {
    return 'Electrical Maintenance & Lighting Unit — High-Mast Lights, Distribution Panels & Power Systems'
  }
  if (roles.includes('plumber')) {
    return 'Water Reticulation & Sanitary Infrastructure Unit — Hydration Lines, Drainage, Pumps & Pressure Systems'
  }
  return 'ICT Systems & Infrastructure Workstation — Datacenter, Network & Web Systems'
})

const departmentalBadgeText = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  
  if (roles.includes('senior_engineer')) return 'Infrastructure Directorate Active'
  if (roles.includes('assistant_engineer_civil')) return 'Civil & Structural Engineering Active'
  if (roles.includes('assistant_engineer_electrical')) return 'Electrical & Arena Lighting Active'
  if (roles.includes('engineering_officer_civil')) return 'Civil Inspections Active'
  if (roles.includes('engineering_officer_electrical')) return 'Power Systems Active'
  if (roles.includes('plumber')) return 'Water & Sanitary Reticulation Active'
  return 'Systems Operational (99.9% Uptime)'
})

const itPendingCount = computed(() => {
  return itRequisitions.value.filter(r => ['SUBMITTED', 'HOD_APPROVED', 'PDU_REVIEW', 'FINANCE_CLEARED'].includes(r.status)).length
})

const itApprovedCount = computed(() => {
  return itRequisitions.value.filter(r => ['ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED'].includes(r.status)).length
})

const pendingVisitorsCount = computed(() => {
  return receptionVisitors.value.filter(v => v.status === 'PENDING_APPROVAL').length
})

const checkedInVisitorsCount = computed(() => {
  return receptionVisitors.value.filter(v => v.status === 'CHECKED_IN').length
})

const utilizedLeaveDays = computed(() => {
  return leaveApplications.value
    .filter(l => l.status === 'APPROVED_HR' && l.leave_type === 'ANNUAL_LEAVE')
    .reduce((acc, cur) => acc + (cur.days_requested || 0), 0)
})

const remainingLeaveDays = computed(() => {
  return Math.max(0, 21 - utilizedLeaveDays.value)
})

const latestLeave = computed(() => {
  return leaveApplications.value[0] || null
})

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 12) return 'Good Morning'
  if (h < 17) return 'Good Afternoon'
  return 'Good Evening'
})

const welcomeName = computed(() => {
  const u = authStore.user
  if (!u) return 'Staff Member'
  if (u.first_name) return `${u.first_name} ${u.last_name || ''}`.trim()
  if (u.email) return u.email.split('@')[0]
  return 'Staff Member'
})

function statusBadge(status) {
  switch (status?.toUpperCase()) {
    case 'APPROVED':
    case 'APPROVED_HR':
    case 'CLEARED':
    case 'COMPLETED':
    case 'ACCOUNTING_OFFICER_APPROVED':
    case 'PO_ISSUED':
    case 'DELIVERED':
      return 'badge-success'
    case 'CHECKED_IN':
      return 'badge-info'
    case 'PENDING':
    case 'PENDING_APPROVAL':
    case 'PENDING_SUPERVISOR':
    case 'PENDING REVIEW':
    case 'SUBMITTED':
    case 'HOD_APPROVED':
    case 'PDU_REVIEW':
    case 'FINANCE_CLEARED':
      return 'badge-warning'
    case 'REJECTED':
      return 'badge-danger'
    default:
      return 'badge-primary'
  }
}

function formatDate(d) {
  if (!d) return ''
  try {
    return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })
  } catch {
    return d
  }
}

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    return
  }
  loading.value = true
  try {
    if (isReceptionist.value) {
      // Load Receptionist Front Desk KPIs
      const [visitorsRes, leaveRes] = await Promise.allSettled([
        apiClient.get('/api/v1/reception/visitors'),
        apiClient.get('/api/v1/reception/leave/status')
      ])
      if (visitorsRes.status === 'fulfilled') {
        receptionVisitors.value = visitorsRes.value.data?.data?.visitors || visitorsRes.value.data?.visitors || []
      }
      if (leaveRes.status === 'fulfilled') {
        leaveApplications.value = leaveRes.value.data?.data?.leaves || leaveRes.value.data?.leaves || []
      }
    } else if (isDepartmentalOfficer.value) {
      // Load IT & Engineering Department PPDA Form 5 and Leave KPIs
      const [ppdaRes, statsRes, leaveRes] = await Promise.allSettled([
        apiClient.get('/api/v1/it/ppda'),
        apiClient.get('/api/v1/it/dashboard/stats'),
        apiClient.get('/api/v1/reception/leave/status')
      ])
      if (ppdaRes.status === 'fulfilled') {
        itRequisitions.value = ppdaRes.value.data?.data?.requisitions || ppdaRes.value.data?.requisitions || []
      }
      if (statsRes.status === 'fulfilled') {
        itDashboardStats.value = statsRes.value.data?.data || statsRes.value.data || null
      }
      if (leaveRes.status === 'fulfilled') {
        leaveApplications.value = leaveRes.value.data?.data?.leaves || leaveRes.value.data?.leaves || []
      }
    } else {
      // Load Executive Overview
      const isStaff = authStore.hasAnyRole('super_admin', 'admin', 'general_secretary', 'ags_technical', 'ags_admin')
      const appsEndpoint = isStaff ? '/api/v1/admin/applications' : '/api/v1/applications'

      const [statsRes, appsRes] = await Promise.allSettled([
        apiClient.get('/api/v1/dashboard/stats'),
        apiClient.get(appsEndpoint, { params: { limit: 5, per_page: 5 } })
      ])
      
      if (statsRes.status === 'fulfilled') {
        stats.value = statsRes.value.data?.data || statsRes.value.data || {}
      }
      if (appsRes.status === 'fulfilled') {
        const appData = appsRes.value.data?.data?.applications || appsRes.value.data?.applications || appsRes.value.data?.data || []
        recentApplicationsList.value = Array.isArray(appData) ? appData.slice(0, 5) : []
      }
    }
  } catch {
    // Fail gracefully
  } finally {
    loading.value = false
  }
})
</script>
