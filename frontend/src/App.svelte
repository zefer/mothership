<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import Router, { location } from 'svelte-spa-router';
  import { connectWebSocket, play, pause, next, previous, randomOn, randomOff } from './lib/mpd.js';
  import { player } from './lib/stores.js';
  import Header from './components/Header.svelte';
  import Playing from './routes/Playing.svelte';
  import Playlist from './routes/Playlist.svelte';
  import Browse from './routes/Browse.svelte';

  const routes = {
    '/': Playing,
    '/playing': Playing,
    '/playlist': Playlist,
    '/browse': Browse,
    '/browse/*': Browse,
  };

  onMount(() => {
    connectWebSocket();

    document.addEventListener('keydown', (e) => {
      const tag = e.target.tagName.toLowerCase();
      if (tag === 'input' || tag === 'textarea' || tag === 'select') return;
      if (e.metaKey || e.ctrlKey) return;

      switch (e.key) {
        case ' ':
        case 'p':
          e.preventDefault();
          if (get(player).state === 'play') pause();
          else play();
          break;
        case 'ArrowRight':
          e.preventDefault();
          next();
          break;
        case 'ArrowLeft':
          e.preventDefault();
          previous();
          break;
        case 'r':
          if (get(player).randomOn) randomOff();
          else randomOn();
          break;
        case '1':
          window.location.hash = '/playing';
          break;
        case '2':
          window.location.hash = '/browse';
          break;
        case '3':
          window.location.hash = '/playlist';
          break;
        case '?':
          alert(
            'Keyboard shortcuts:\n\n' +
            'Space / P - Play/Pause\n' +
            '\u2192 - Next track\n' +
            '\u2190 - Previous track\n' +
            'R - Toggle random\n' +
            '1 - Playing view\n' +
            '2 - Browse view\n' +
            '3 - Playlist view\n' +
            '? - Show this help'
          );
          break;
      }
    });
  });
</script>

<Header />
<main class="container">
  <Router {routes} />
</main>
