import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

export default defineConfig({
    base: './',

    plugins: [react(), tailwindcss(), nodePolyfills()],

    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
        },
    },

    build: {
        outDir: resolve(__dirname, 'dist'),
        rollupOptions: {
            input: {
                main: resolve(__dirname, 'index.html'),
                details: resolve(__dirname, 'details.html'),
            },
        },
    },

    server: { port: 5173 },
})
