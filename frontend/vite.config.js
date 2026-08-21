import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig(({ mode }) => {
  const env = { ...process.env, ...loadEnv(mode, process.cwd(), '') }
  const devHost = env.VITE_DEV_HOST
  const devPort = Number(env.VITE_DEV_PORT)
  const watchInterval = Number(env.VITE_DEV_WATCH_INTERVAL)
  const proxyTarget = env.VITE_DEV_PROXY_TARGET
  const usePolling = String(env.VITE_DEV_WATCH_POLLING).toLowerCase() === 'true'

  return {
    plugins: [vue()],
    build: {
      target: 'esnext',
      cssCodeSplit: true,
      chunkSizeWarningLimit: 1000,
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('node_modules')) {
              if (id.includes('vue') || id.includes('vue-router') || id.includes('pinia')) {
                return 'vendor-vue'
              }
              if (id.includes('sweetalert2') || id.includes('gsap') || id.includes('tippy')) {
                return 'vendor-ui'
              }
              return 'vendor-deps'
            }
          },
          assetFileNames: assetInfo => assetInfo.names?.includes('icofont.woff2')
            ? 'assets/icofont.woff2'
            : assetInfo.names?.includes('icofont.woff')
              ? 'assets/icofont.woff'
              : 'assets/[name]-[hash][extname]'
        }
      }
    },
    resolve: {
      alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) }
    },
    server: {
      host: devHost,
      port: devPort,
      watch: {
        usePolling,
        interval: watchInterval
      },
      proxy: {
        '/uploads': { target: proxyTarget, changeOrigin: true },
        '/api': { target: proxyTarget, changeOrigin: true },
        '/portal-api': {
          target: env.VITE_PORTAL_API_URL || proxyTarget,
          changeOrigin: true,
          rewrite: path => path.replace(/^\/portal-api/, '')
        }
      }
    }
  }
})
