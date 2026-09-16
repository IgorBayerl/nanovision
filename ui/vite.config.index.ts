import { resolve } from 'node:path'
import { defineConfig, mergeConfig } from 'vite'
import { baseConfig } from './vite.config.base.ts'

export default defineConfig(
    mergeConfig(baseConfig, {
        build: {
            emptyOutDir: true,
            outDir: resolve(import.meta.dirname, '../internal/reporter/htmlreact/assets/dist'),
            rolldownOptions: {
                input: {
                    main: resolve(import.meta.dirname, 'index.html'),
                },
            },
        },
    }),
)
