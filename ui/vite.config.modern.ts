import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig, type Plugin } from 'vite'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

/**
 * A custom plugin to handle development server rewrites.
 * This makes links like `...go.html` correctly serve the `details.html` page
 * during development, mimicking the production environment.
 */
function devServerRewrites(): Plugin {
    return {
        name: 'dev-server-rewrites',
        configureServer(server) {
            // This middleware runs for every request
            server.middlewares.use((req, _res, next) => {
                // If the URL matches the pattern of a details page link
                if (req.url?.endsWith('.go.html')) {
                    // Rewrite the request to point to /details.html instead
                    req.url = '/details.html'
                }
                // Continue to the next middleware
                next()
            })
        },
    }
}

export default defineConfig({
    base: './',

    plugins: [react(), tailwindcss(), nodePolyfills(), devServerRewrites()],

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
