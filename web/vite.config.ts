import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/v1': {
        target: 'http://localhost:9002',
        changeOrigin: true,
      },
      '/gateway': {
        target: 'http://localhost:9002',
        changeOrigin: true,
      },
    },
  },
})
