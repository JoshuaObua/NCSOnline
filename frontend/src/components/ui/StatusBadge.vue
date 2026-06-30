<template>
  <span :class="badgeClasses" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true
  }
})

const statusMap = {
  DRAFT: { label: 'Draft', classes: 'bg-gray-100 text-gray-700' },
  SUBMITTED: { label: 'Submitted', classes: 'bg-blue-100 text-blue-700' },
  UNDER_REVIEW: { label: 'Under Review', classes: 'bg-yellow-100 text-yellow-700' },
  APPROVED: { label: 'Approved', classes: 'bg-green-100 text-green-700' },
  REJECTED: { label: 'Rejected', classes: 'bg-red-100 text-red-700' },
  NEEDS_INFORMATION: { label: 'Queried', classes: 'bg-orange-100 text-orange-700' },
  PENDING_PAYMENT: { label: 'Pending Payment', classes: 'bg-purple-100 text-purple-700' },
  // User status variants
  active: { label: 'Active', classes: 'bg-green-100 text-green-700' },
  inactive: { label: 'Inactive', classes: 'bg-gray-100 text-gray-500' },
  ACTIVE: { label: 'Active', classes: 'bg-green-100 text-green-700' },
  INACTIVE: { label: 'Inactive', classes: 'bg-gray-100 text-gray-500' }
}

const badgeClasses = computed(() => {
  return statusMap[props.status]?.classes || 'bg-gray-100 text-gray-700'
})

const label = computed(() => {
  return statusMap[props.status]?.label || props.status
})
</script>
