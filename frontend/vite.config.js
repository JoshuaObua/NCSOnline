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
        '/api': { target: proxyTarget, changeOrigin: true }
      }
    }
  }
})
