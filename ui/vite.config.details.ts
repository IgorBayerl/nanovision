import { resolve } from 'node:path'
import { defineConfig, mergeConfig } from 'vite'
import { baseConfig } from './vite.config.base'

export default defineConfig(
    mergeConfig(baseConfig, {
        build: {
            // Important to keep this false
            // We are generating 2 builds in one command, this makes sure we do not overwrite the first one.
            emptyOutDir: false,
            outDir: resolve(__dirname, '../internal/reporter/htmlreact/assets/dist'),
            rollupOptions: {
                input: {
                    details: resolve(__dirname, 'details.html'),
                },
            },
        },
    }),
)
