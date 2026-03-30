import { defineConfig } from 'astro/config';
import sitemap from '@astrojs/sitemap';

export default defineConfig({
  site: 'https://needle-bench.cc',
  output: 'static',
  integrations: [sitemap()],
  build: {
    format: 'directory',
  },
});
