<template>
  <LayoutDefault :title="title">
    <section class="expense-workspace">
      <header class="expense-heading"><div><h1>{{ title }}</h1><p v-if="description">{{ description }}</p></div><slot name="actions" /></header>
      <slot />
    </section>
  </LayoutDefault>
</template>
<script setup>
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
defineProps({title:String,description:String})
</script>
<style>
.expense-workspace { color:#1f2937; }
.expense-workspace .expense-heading { display:flex;flex-wrap:wrap;justify-content:space-between;align-items:center;gap:1rem;margin-bottom:1.5rem; }
.expense-workspace h1 { font-size:1.5rem;font-weight:700;margin:0; }
.expense-workspace h2 { font-size:1.1rem;font-weight:700;margin-bottom:1rem; }
.expense-workspace p { margin:.4rem 0; }
.expense-workspace .expense-card { background:white;border:1px solid #d1d5db;border-radius:.75rem;padding:1.25rem;margin-bottom:1rem; }
.expense-workspace .expense-grid { display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1rem; }
.expense-workspace .expense-input { display:block;width:100%;margin-top:.35rem;border:1px solid #9ca3af;border-radius:.5rem;padding:.65rem;background:#fff;color:#111827;font-size:.875rem; }
.expense-workspace label { display:block;font-size:.875rem;font-weight:600; }
.expense-workspace .expense-wide { grid-column:1/-1; }
.expense-workspace .expense-actions { display:flex;flex-wrap:wrap;gap:.75rem;align-items:center; }
.expense-workspace .expense-button { display:inline-flex;align-items:center;justify-content:center;padding:.65rem 1rem;border:1px solid #047857;background:#047857;color:#fff;border-radius:.5rem;font-size:.875rem;font-weight:600; }
.expense-workspace .expense-secondary { background:white;color:#065f46; }
.expense-workspace button:disabled { opacity:.5;cursor:not-allowed; }
.expense-workspace :focus-visible { outline:2px solid #047857;outline-offset:2px; }
.expense-workspace .expense-error { padding:1rem;background:#fef2f2;color:#991b1b;border:1px solid #fca5a5;border-radius:.5rem;margin:1rem 0; }
.expense-workspace .expense-success { padding:1rem;background:#ecfdf5;color:#065f46;border-radius:.5rem;margin-bottom:1rem; }
.expense-workspace .expense-muted { color:#4b5563;font-size:.875rem; }
.expense-workspace .expense-table-wrap { overflow-x:auto; }
.expense-workspace table { width:100%;font-size:.875rem;border-collapse:collapse; }
.expense-workspace th,.expense-workspace td { padding:.8rem;text-align:left;border-bottom:1px solid #e5e7eb; }
.expense-workspace th { background:#f3f4f6;white-space:nowrap; }
.expense-workspace .expense-link { color:#047857;font-weight:600;text-decoration:underline; }
.expense-workspace .expense-details dt { color:#4b5563;font-size:.875rem; }
.expense-workspace .expense-details dd { margin:.25rem 0 1rem;font-weight:600;white-space:pre-wrap;overflow-wrap:anywhere; }
@media(max-width:640px) { .expense-workspace .expense-grid { grid-template-columns:1fr; } }
</style>
