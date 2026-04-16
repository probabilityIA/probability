// @ts-check
import { defineConfig } from 'astro/config';
import node from '@astrojs/node';

import tailwindcss from '@tailwindcss/vite';

import react from '@astrojs/react';

// https://astro.build/config
export default defineConfig({
  output: 'server',
  adapter: node({
    mode: 'standalone'
  }),
  base: '/website/',
  vite: {
    plugins: [tailwindcss()],
    ssr: {
      external: ['framer-motion']
    }
  },

  integrations: [react()]
});