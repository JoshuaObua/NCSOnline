let motionGeneration = 0
let activeTriggers = []
let mediaObserver

function enhanceLazyMedia(root) {
  if (!root) return
  root.querySelectorAll('img').forEach((image, index) => {
    const eager = image.closest('#section-hero, .public-carousel') || index === 0
    if (!image.hasAttribute('loading')) image.loading = eager ? 'eager' : 'lazy'
    if (!image.hasAttribute('decoding')) image.decoding = 'async'
    if (!eager && !image.hasAttribute('fetchpriority')) image.setAttribute('fetchpriority', 'low')
  })
  root.querySelectorAll('iframe').forEach(frame => {
    if (!frame.hasAttribute('loading')) frame.loading = 'lazy'
  })
  root.querySelectorAll('video').forEach(video => {
    if (!video.hasAttribute('preload')) video.preload = video.closest('#section-hero, .public-carousel') ? 'metadata' : 'none'
  })
}

function observeLazyMedia(root) {
  mediaObserver?.disconnect()
  enhanceLazyMedia(root)
  mediaObserver = new MutationObserver(() => enhanceLazyMedia(root))
  mediaObserver.observe(root, { childList: true, subtree: true })
}

function topLevelSections(root) {
  return Array.from(root.querySelectorAll('section:not(#section-hero)')).filter(section => !section.parentElement?.closest('section'))
}

export async function animatePublicPage(root) {
  if (!root) return
  cleanupPublicMotion(false)
  const generation = ++motionGeneration
  observeLazyMedia(root)

  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return

  const [gsapPackage, triggerPackage] = await Promise.all([
    import('gsap'),
    import('gsap/ScrollTrigger'),
  ])
  if (generation !== motionGeneration || !root.isConnected) return

  const gsap = gsapPackage.gsap || gsapPackage.default
  const ScrollTrigger = triggerPackage.ScrollTrigger || triggerPackage.default
  gsap.registerPlugin(ScrollTrigger)

  const reveal = (target, options = {}) => {
    const trigger = ScrollTrigger.create({
      trigger: target,
      start: options.start || 'top 88%',
      once: true,
      onEnter: () => {
        const items = options.items || [target]
        gsap.fromTo(items, {
          autoAlpha: 0,
          y: options.y ?? 28,
          scale: options.scale ?? 0.985,
        }, {
          autoAlpha: 1,
          y: 0,
          scale: 1,
          duration: options.duration || 0.72,
          stagger: options.stagger || 0,
          ease: options.ease || 'power3.out',
          clearProps: 'opacity,visibility,transform',
        })
      },
    })
    activeTriggers.push(trigger)
  }

  topLevelSections(root).forEach((section, index) => reveal(section, {
    y: index % 2 === 0 ? 28 : 36,
    duration: 0.78,
  }))

  const animatedChildren = new Set()
  root.querySelectorAll('.card-grid, .association-grid, .help-grid, .logo-strip, [data-motion-stagger]').forEach(group => {
    const items = Array.from(group.children).filter(item => !animatedChildren.has(item)).slice(0, 16)
    if (items.length < 2) return
    items.forEach(item => animatedChildren.add(item))
    reveal(group, { items, y: 24, scale: 0.97, stagger: 0.07, duration: 0.62 })
  })

  if (root.querySelector('.home-redesign')) {
    root.querySelectorAll('.section-heading, .core-functions, .quick-contact').forEach(target => reveal(target, {
      y: 20,
      duration: 0.65,
    }))
    root.querySelectorAll('#section-about img, .content-card > img, .image-card > img').forEach(image => {
      const tween = gsap.fromTo(image, { yPercent: -2 }, {
        yPercent: 3,
        ease: 'none',
        scrollTrigger: { trigger: image, start: 'top bottom', end: 'bottom top', scrub: 0.6 },
      })
      if (tween.scrollTrigger) activeTriggers.push(tween.scrollTrigger)
    })
  }

  ScrollTrigger.refresh()
}

export function cleanupPublicMotion(disconnectMedia = true) {
  motionGeneration += 1
  activeTriggers.forEach(trigger => trigger?.kill())
  activeTriggers = []
  if (disconnectMedia) {
    mediaObserver?.disconnect()
    mediaObserver = undefined
  }
}
