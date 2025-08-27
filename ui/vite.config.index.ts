import { resolve } from 'node:path'
import { defineConfig, mergeConfig } from 'vite'
import { baseConfig } from './vite.config.base'

export default defineConfig(
    mergeConfig(baseConfig, {
        build: {
            emptyOutDir: true,
            outDir: resolve(__dirname, '../internal/reporter/htmlreact/assets/dist'),
            rollupOptions: {
                input: {
                    main: resolve(__dirname, 'index.html'),
                },
            },
        },
    }),
)
