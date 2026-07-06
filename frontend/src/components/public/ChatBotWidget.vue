<template>
  <!-- Toggle FAB -->
  <button
    v-if="!open"
    type="button"
    class="fixed bottom-20 right-5 z-[80] w-14 h-14 rounded-full bg-[#f5a623] hover:bg-[#e09612] text-white shadow-xl shadow-[#f5a623]/30 flex items-center justify-center transition-all duration-300 hover:scale-110 group"
    aria-label="Open NCS chatbot"
    @click="open = true"
  >
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6" aria-hidden="true"><path d="M7.9 20A9 9 0 1 0 4 16.1L2 22Z"></path></svg>
    <span class="absolute -top-1 -right-1 w-3 h-3 rounded-full bg-green-400 ring-2 ring-white animate-pulse"></span>
    <span class="absolute right-full mr-3 whitespace-nowrap bg-[#1a365d] text-white text-sm font-medium px-3 py-1.5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">Ask NCS Bot</span>
  </button>

  <!-- Chat panel -->
  <div
    v-else
    class="fixed bottom-20 right-5 z-[80] w-[calc(100vw-3rem)] sm:w-[400px] h-[600px] max-h-[calc(100vh-6rem)] bg-white rounded-2xl shadow-2xl border border-gray-200 flex flex-col overflow-hidden"
    role="dialog"
    aria-label="NCS chat assistant"
  >
    <!-- Header -->
    <div class="bg-gradient-to-br from-[#1a365d] to-[#2d4a7a] p-4 text-white flex items-center justify-between flex-shrink-0">
      <div class="flex items-center gap-3">
        <div class="relative w-10 h-10 rounded-full bg-[#f5a623] flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-white" aria-hidden="true"><path d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z"></path><path d="M20 3v4"></path><path d="M22 5h-4"></path><path d="M4 17v2"></path><path d="M5 18H3"></path></svg>
          <span class="absolute -bottom-0.5 -right-0.5 w-3 h-3 rounded-full bg-green-400 ring-2 ring-[#1a365d]"></span>
        </div>
        <div>
          <p class="font-bold leading-tight">NCS Bot</p>
          <p class="text-xs text-white/70">Online • Ask me anything</p>
        </div>
      </div>
      <div class="flex items-center gap-1">
        <button type="button" class="w-8 h-8 rounded-full hover:bg-white/10 flex items-center justify-center transition-colors" aria-label="Close chat" @click="open = false">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M18 6 6 18"></path><path d="m6 6 12 12"></path></svg>
        </button>
      </div>
    </div>

    <!-- Messages -->
    <div ref="messagesEl" class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50">
      <div v-if="!messages.length" class="text-center py-6">
        <div class="w-14 h-14 rounded-full bg-gradient-to-br from-[#f5a623] to-[#e09612] mx-auto mb-3 flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7 text-white" aria-hidden="true"><path d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z"></path><path d="M20 3v4"></path><path d="M22 5h-4"></path><path d="M4 17v2"></path><path d="M5 18H3"></path></svg>
        </div>
        <p class="font-bold text-[#1a365d] text-base mb-1">Hi! I'm NCS Bot</p>
        <p class="text-sm text-gray-600 mb-6">Ask me anything about NCS, our 52 federations, facilities, registrations, or how to invest.</p>
        <div class="space-y-2">
          <button
            v-for="(prompt, index) in suggestedPrompts"
            :key="index"
            type="button"
            class="w-full text-left text-sm bg-white border border-gray-200 hover:border-[#f5a623] hover:bg-[#f5a623]/5 rounded-lg px-3 py-2 transition-colors text-gray-700"
            @click="sendMessage(prompt)"
          >{{ prompt }}</button>
        </div>
      </div>
      <template v-else>
        <div v-for="(message, index) in messages" :key="index" :class="['flex', message.role === 'user' ? 'justify-end' : 'justify-start']">
          <div :class="['max-w-[85%] rounded-2xl px-3 py-2 text-sm whitespace-pre-line', message.role === 'user' ? 'bg-[#1a365d] text-white rounded-br-md' : 'bg-white border border-gray-200 text-gray-700 rounded-bl-md']">{{ message.text }}</div>
        </div>
      </template>
    </div>

    <!-- Input -->
    <form class="border-t border-gray-200 p-3 bg-white flex-shrink-0" @submit.prevent="sendMessage()">
      <div class="flex items-end gap-2">
        <textarea
          v-model="draft"
          placeholder="Type your question..."
          rows="1"
          class="flex-1 resize-none rounded-xl border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#f5a623] focus:border-transparent max-h-32 min-h-[40px]"
          @keydown.enter.exact.prevent="sendMessage()"
        ></textarea>
        <button
          type="submit"
          :disabled="!draft.trim()"
          class="w-10 h-10 rounded-full bg-[#f5a623] hover:bg-[#e09612] disabled:bg-gray-300 disabled:cursor-not-allowed text-white flex items-center justify-center transition-colors flex-shrink-0"
          aria-label="Send message"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z"></path><path d="m21.854 2.147-10.94 10.939"></path></svg>
        </button>
      </div>
      <p class="text-[10px] text-gray-400 text-center mt-2">Powered by Claude • Replies may occasionally be inaccurate</p>
    </form>
  </div>
</template>

<script setup>
import { nextTick, onMounted, onUnmounted, ref } from 'vue'

const open = ref(false)
const draft = ref('')
const messages = ref([])
const messagesEl = ref(null)

const suggestedPrompts = [
  'How do I register a sports federation?',
  'Who is the NCS chairperson?',
  'Tell me about Uganda Rugby Union',
  'How can I submit an investment proposal?',
]

// The live AI backend is not wired up yet, so replies are an honest
// placeholder pointing people at real contact channels rather than
// pretending to answer.
const STUB_REPLY = "Thanks for your question! I'm still being connected to my knowledge base, so I can't answer this yet.\n\nIn the meantime, the NCS team can help directly:\n• info@ncs.go.ug\n• +256 414254477\n• The contact form on the Contact Us page"

function sendMessage(text) {
  const value = (text ?? draft.value).trim()
  if (!value) return
  messages.value.push({ role: 'user', text: value })
  draft.value = ''
  scrollToBottom()
  window.setTimeout(() => {
    messages.value.push({ role: 'bot', text: STUB_REPLY })
    scrollToBottom()
  }, 450)
}

async function scrollToBottom() {
  await nextTick()
  if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight
}

function onKeydown(event) {
  if (event.key === 'Escape' && open.value) open.value = false
}
onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>
