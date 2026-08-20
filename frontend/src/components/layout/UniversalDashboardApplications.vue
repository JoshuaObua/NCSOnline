<template>
  <section class="universal-dashboard-applications rounded-2xl border border-slate-200 bg-white p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900">
    <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
      <div>
        <p class="text-[11px] font-black uppercase tracking-[0.18em] text-blue-600 dark:text-blue-400">Universal dashboard modules</p>
        <h2 class="mt-1 text-base font-black text-slate-900 dark:text-white">Leave Application & PPDA Application Desk</h2>
        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
          Available from every dashboard for staff self-service, procurement initiation, and Form 5 follow-up.
        </p>
      </div>

      <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:min-w-[34rem]">
        <article class="rounded-xl border border-emerald-100 bg-emerald-50/70 p-3 dark:border-emerald-900/60 dark:bg-emerald-950/20">
          <div class="flex items-start gap-3">
            <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-100 text-emerald-700 dark:bg-emerald-900/70 dark:text-emerald-300">
              <i class="icofont-calendar text-lg"></i>
            </div>
            <div class="min-w-0 flex-1">
              <h3 class="text-sm font-black text-slate-900 dark:text-white">Leave Application Module</h3>
              <p class="mt-0.5 text-[11px] text-slate-500 dark:text-slate-400">Apply for statutory leave and check approval history.</p>
              <div class="mt-3 flex flex-wrap gap-2">
                <router-link to="/reception/leave/apply" class="btn btn-sm btn-success text-[11px]">
                  <i class="icofont-plus-circle mr-1"></i> Apply Leave
                </router-link>
                <router-link to="/reception/leave/status" class="btn btn-sm btn-outline-success text-[11px]">
                  Status
                </router-link>
              </div>
            </div>
          </div>
        </article>

        <article class="rounded-xl border border-blue-100 bg-blue-50/70 p-3 dark:border-blue-900/60 dark:bg-blue-950/20">
          <div class="flex items-start gap-3">
            <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-100 text-blue-700 dark:bg-blue-900/70 dark:text-blue-300">
              <i class="icofont-document-folder text-lg"></i>
            </div>
            <div class="min-w-0 flex-1">
              <h3 class="text-sm font-black text-slate-900 dark:text-white">PPDA Application Module</h3>
              <p class="mt-0.5 text-[11px] text-slate-500 dark:text-slate-400">Create and track PPDA Form 5 procurement requests.</p>
              <div class="mt-3 flex flex-wrap gap-2">
                <router-link to="/it/ppda/new" class="btn btn-sm btn-primary text-[11px]">
                  <i class="icofont-plus-circle mr-1"></i> New PPDA
                </router-link>
                <router-link to="/it/ppda/status" class="btn btn-sm btn-outline-primary text-[11px]">
                  Status
                </router-link>
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
