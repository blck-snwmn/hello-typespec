import { defineConfig } from 'vite'
import devServer from '@hono/vite-dev-server'

export default defineConfig({
  server: {
    port: 3000,
  },
  build: {
    lib: {
      entry: 'src/index.ts',
      formats: ['es'],
      fileName: 'index',
    },
    rollupOptions: {
      external: ['hono', 'hono/cors', '@hono/node-server'],
    },
  },
  plugins: [
    devServer({
      entry: 'src/index.ts',
    })
  ],
})