<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="show"
        class="fixed inset-0 z-50 overflow-y-auto flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
      >
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm" @click="$emit('close')"></div>

        <!-- Panel -->
        <Transition name="modal-slide" appear>
          <div
            v-if="show"
            :class="sizeClass"
            class="relative w-full bg-white rounded-2xl shadow-2xl z-10"
            @click.stop
          >
            <!-- Header -->
            <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-100">
              <div v-if="icon" :class="[iconBg, 'w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0']">
                <i :class="[icon, iconColor, 'text-lg leading-none']"></i>
              </div>
              <h3 class="flex-1 text-base font-semibold text-gray-900">{{ title }}</h3>
              <button
                @click="$emit('close')"
                class="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-100 transition-colors flex-shrink-0"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Body -->
            <div class="px-6 py-5">
              <slot />
            </div>

            <!-- Footer -->
            <div v-if="$slots.footer" class="px-6 py-4 border-t border-gray-100 bg-gray-50/60 flex justify-end gap-2 rounded-b-2xl">
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  show:     { type: Boolean, default: false },
  title:    { type: String,  default: '' },
  size:     { type: String,  default: 'md', validator: v => ['sm','md','lg','xl'].includes(v) },
  icon:     { type: String,  default: '' },
  iconBg:   { type: String,  default: 'bg-primary-100' },
  iconColor:{ type: String,  default: 'text-primary-700' }
})
defineEmits(['close'])

const sizeClass = computed(() => ({
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-2xl'
}[props.size] || 'max-w-md'))
</script>

<style scoped>
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.2s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

.modal-slide-enter-active { transition: opacity 0.2s ease, transform 0.2s ease; }
.modal-slide-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.modal-slide-enter-from { opacity: 0; transform: translateY(16px) scale(0.98); }
.modal-slide-leave-to   { opacity: 0; transform: translateY(8px) scale(0.99); }
</style>
