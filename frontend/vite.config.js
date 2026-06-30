import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) }
  },
  server: {
    host: '0.0.0.0',
    port: 3001,
    watch: {
      usePolling: true,
      interval: 300
    },
    // Proxy /uploads and /api to the nginx gateway so uploaded media (slides,
    // facility images, logos, PDFs etc.) resolve correctly when the SPA is
    // served from the Vite dev port. nginx serves them at port 80 inside the
    // ncsms_net docker network.
    proxy: {
      '/uploads': { target: 'http://localhost:9081', changeOrigin: true },
      '/api':     { target: 'http://localhost:9081', changeOrigin: true }
    }
  }
})
