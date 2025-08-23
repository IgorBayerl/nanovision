import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { nodePolyfills } from 'vite-plugin-node-polyfills'
import { viteSingleFile } from 'vite-plugin-singlefile'

export default defineConfig(() => ({
    base: './',
    plugins: [react(), tailwindcss(), viteSingleFile(), nodePolyfills()].filter(Boolean),

    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
        },
    },

    build: {
        outDir: resolve(__dirname, '../internal/reporter/htmlreact/assets/dist'),
        assetsDir: 'assets',
        sourcemap: false,
        modulePreload: false,
        rollupOptions: {
            input: resolve(__dirname, 'index.html'),
            output: {
                entryFileNames: 'assets/react-islands.js',
                assetFileNames: 'assets/[name][extname]',
                chunkFileNames: 'assets/[name]-[hash].js',
            },
        },
    },
    server: { port: 5173 },
}))
