import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/scrape': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/results': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/schedules': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
