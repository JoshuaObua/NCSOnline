<template>
  <AccountShell title="Profile Overview" subtitle="Core identity, public display information, and account roles.">
    <section class="account-grid">
      <article class="account-panel profile-hero">
        <div class="avatar-frame">
          <img v-if="profile.avatar_url" :src="profile.avatar_url" alt="" referrerpolicy="no-referrer" />
          <span v-else>{{ initials }}</span>
        </div>
        <div>
          <h2>{{ fullName }}</h2>
          <p>{{ profile.email }}</p>
          <div class="role-row">
            <span v-for="role in profile.roles || []" :key="role.id || role.name">{{ role.name }}</span>
          </div>
        </div>
      </article>

      <article class="account-panel">
        <h3>Account Details</h3>
        <dl class="detail-list">
          <div><dt>First name</dt><dd>{{ profile.first_name || '-' }}</dd></div>
          <div><dt>Last name</dt><dd>{{ profile.last_name || '-' }}</dd></div>
          <div><dt>Phone</dt><dd>{{ profile.phone || '-' }}</dd></div>
          <div><dt>Email verified</dt><dd>{{ profile.is_email_verified ? 'Verified' : 'Not verified' }}</dd></div>
          <div><dt>Provider</dt><dd>{{ profile.auth_provider || 'password' }}</dd></div>
          <div><dt>Last login</dt><dd>{{ formatDate(profile.last_login_at) }}</dd></div>
        </dl>
      </article>
    </section>
  </AccountShell>
</template>

<script setup>
import { computed, onMounted, reactive } from 'vue'
import AccountShell from './AccountShell.vue'
import { getProfile } from '@/api/account.js'

const profile = reactive({})
const fullName = computed(() => `${profile.first_name || ''} ${profile.last_name || ''}`.trim() || profile.email || 'Account')
const initials = computed(() => fullName.value.split(/\s+/).slice(0, 2).map(part => part[0] || '').join('').toUpperCase() || 'U')

function formatDate(value) {
  return value ? new Date(value).toLocaleString() : '-'
}

onMounted(async () => {
  const res = await getProfile()
  Object.assign(profile, res.data?.data || {})
})
</script>
