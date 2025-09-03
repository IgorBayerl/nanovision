import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import type { Plugin, UserConfig } from 'vite'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

/**
 * Find the module script tag, and modify it to not use type="module"
 * Also adds defer, so its loaded after the rest of the page and it can get access to the data.js file
 * This is used so we can run the .html straight from the file, without any server
 */
function classicScriptPlugin(): Plugin {
    return {
        name: 'vite-plugin-classic-script',
        transformIndexHtml(html) {
            let transformedHtml = html.replace(/<script type="module"(.+?)><\/script>/g, '<script defer$1></script>')
            transformedHtml = transformedHtml.replace(/ crossorigin/g, '')
            return transformedHtml
        },
    }
}

export const baseConfig: UserConfig = {
    base: './',
    plugins: [react(), tailwindcss(), nodePolyfills(), classicScriptPlugin()],
    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
        },
    },
    server: { port: 5173 },
}
