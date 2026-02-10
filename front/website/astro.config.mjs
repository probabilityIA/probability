// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';

import preact from '@astrojs/preact';

// https://astro.build/config
export default defineConfig({
<<<<<<< HEAD
=======
  base: '/website/',
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
  vite: {
    plugins: [tailwindcss()]
  },

  integrations: [preact()]
});