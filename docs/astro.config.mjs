// @ts-check
import { defineConfig } from 'astro/config';

import react from '@astrojs/react';
import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  site: 'https://igorbayerl.github.io',
  base: process.env.GITHUB_ACTIONS ? '/nanovision/docs/' : '/',
  integrations: [react()],

  vite: {
    plugins: [tailwindcss()]
  }
});