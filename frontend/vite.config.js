import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/websocket': {
        target: 'ws://127.0.0.1:8080',
        ws: true,
      },
      '/play': 'http://127.0.0.1:8080',
      '/pause': 'http://127.0.0.1:8080',
      '/next': 'http://127.0.0.1:8080',
      '/previous': 'http://127.0.0.1:8080',
      '/randomOn': 'http://127.0.0.1:8080',
      '/randomOff': 'http://127.0.0.1:8080',
      '/files': 'http://127.0.0.1:8080',
      '/playlist': 'http://127.0.0.1:8080',
      '/library': 'http://127.0.0.1:8080',
    },
  },
});
