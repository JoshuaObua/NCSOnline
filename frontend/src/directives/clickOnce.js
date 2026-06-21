/**
 * v-click-once — prevents accidental double-clicks on action buttons / links.
 *
 *  • On any pointerdown the element receives `.click-pressed` (scale feedback).
 *  • On click, the click handler runs immediately, then the element is locked
 *    (aria-disabled, pointer-events disabled, opacity-60) for the configured
 *    cool-down so a frustrated second tap is ignored. Native `disabled`
 *    state on the element is respected and not overwritten.
 *  • Cool-down defaults to 750 ms — short enough that the legit re-click after
 *    a real network round-trip works fine, long enough to absorb a stutter.
 *
 *  Usage:
 *    <button v-click-once>Save</button>                  // default 750 ms
 *    <button v-click-once="1500">Slow action</button>    // custom ms
 *    <a v-click-once href="…">Submit</a>
 */
const TIMER_KEY = '__clickOnceTimer'
const HANDLERS_KEY = '__clickOnceHandlers'

function lock(el, ms) {
  if (el.disabled) return
  el.classList.add('click-locked')
  el.setAttribute('aria-disabled', 'true')
  const previousPointer = el.style.pointerEvents
  el.style.pointerEvents = 'none'
  el[TIMER_KEY] = setTimeout(() => {
    el.classList.remove('click-locked')
    el.removeAttribute('aria-disabled')
    el.style.pointerEvents = previousPointer || ''
    el[TIMER_KEY] = null
  }, ms)
}

function flash(el) {
  el.classList.add('click-pressed')
  setTimeout(() => el.classList.remove('click-pressed'), 180)
}

export default {
  mounted(el, binding) {
    const ms = Number(binding.value) > 0 ? Number(binding.value) : 750
    const onPointerDown = () => flash(el)
    const onClick = () => lock(el, ms)
    el.addEventListener('pointerdown', onPointerDown, { passive: true })
    el.addEventListener('click', onClick)
    el[HANDLERS_KEY] = { onPointerDown, onClick }
  },
  beforeUnmount(el) {
    const handlers = el[HANDLERS_KEY]
    if (handlers) {
      el.removeEventListener('pointerdown', handlers.onPointerDown)
      el.removeEventListener('click', handlers.onClick)
    }
    if (el[TIMER_KEY]) clearTimeout(el[TIMER_KEY])
  },
}
