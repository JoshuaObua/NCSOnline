<template>
  <div class="account-page">
    <div class="account-wrap">
      <aside class="account-nav" aria-label="Account sections">
        <router-link to="/account/profile">Profile</router-link>
        <router-link to="/account/settings">Settings & Security</router-link>
        <router-link to="/account/activities">My Audit Activities</router-link>
      </aside>
      <main class="account-main">
        <header class="account-heading">
          <p>Account Workspace</p>
          <h1>{{ title }}</h1>
          <span>{{ subtitle }}</span>
        </header>
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
})
</script>

<style>
.account-page{background:#f5f7fb;min-height:calc(100vh - 7rem);padding:2rem 1rem 4rem}.account-wrap{max-width:1180px;margin:0 auto;display:grid;grid-template-columns:250px minmax(0,1fr);gap:1rem}.account-nav{background:white;border:1px solid #e5e7eb;border-radius:8px;padding:.75rem;height:max-content;position:sticky;top:8.5rem}.account-nav a{display:block;border-radius:6px;padding:.75rem .85rem;color:#1f2937;font-weight:700;font-size:.9rem}.account-nav a.router-link-active{background:#112b4e;color:white}.account-main{display:grid;gap:1rem}.account-heading{background:white;border:1px solid #e5e7eb;border-radius:8px;padding:1.25rem}.account-heading p{margin:0 0 .35rem;color:#f48c06;font-weight:800;text-transform:uppercase;font-size:.75rem}.account-heading h1{margin:0;color:#112b4e;font-size:1.75rem;font-weight:900}.account-heading span{display:block;margin-top:.35rem;color:#64748b}.account-grid{display:grid;grid-template-columns:1fr 1fr;gap:1rem}.account-panel{background:white;border:1px solid #e5e7eb;border-radius:8px;padding:1rem}.account-panel h3{margin:0 0 1rem;color:#112b4e;font-weight:900}.profile-hero{display:flex;gap:1rem;align-items:center}.avatar-frame{width:86px;height:86px;border-radius:50%;background:#112b4e;color:white;display:grid;place-items:center;font-size:1.5rem;font-weight:900;overflow:hidden;flex:0 0 auto}.avatar-frame img{width:100%;height:100%;object-fit:cover}.profile-hero h2{font-size:1.35rem;font-weight:900;color:#112b4e;margin:0}.profile-hero p{margin:.2rem 0;color:#64748b}.role-row{display:flex;gap:.4rem;flex-wrap:wrap;margin-top:.65rem}.role-row span{background:#eef2ff;color:#3730a3;border-radius:999px;padding:.25rem .55rem;font-size:.75rem;font-weight:800}.detail-list{display:grid;gap:.7rem;margin:0}.detail-list div{display:grid;grid-template-columns:130px minmax(0,1fr);gap:1rem;border-bottom:1px solid #f1f5f9;padding-bottom:.6rem}.detail-list dt{color:#64748b;font-size:.8rem}.detail-list dd{margin:0;color:#111827;font-weight:700;word-break:break-word}@media(max-width:900px){.account-wrap{grid-template-columns:1fr}.account-nav{position:static}.account-grid{grid-template-columns:1fr}}
</style>
