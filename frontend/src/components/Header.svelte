<script>
  import { location } from 'svelte-spa-router';
  import { connected, player } from '../lib/stores.js';
  import PlayerControls from './PlayerControls.svelte';

  let menuOpen = $state(false);

  function toggleMenu() {
    menuOpen = !menuOpen;
  }

  function closeMenu() {
    menuOpen = false;
  }
</script>

<header>
  <nav>
    <PlayerControls />
    <div class="nav-links">
      <a href="#/playing" class:active={$location === '/playing' || $location === '/'}>Playing</a>
      <a href="#/browse" class:active={$location?.startsWith('/browse')}>Browse</a>
      <a href="#/playlist" class:active={$location === '/playlist'}>
        Playlist
        {#if $player.playlistLength > 0}
          <span class="nav-badge">{$player.position}/{$player.playlistLength}</span>
        {/if}
      </a>
      {#if !$connected}
        <span class="connection-lost">disconnected</span>
      {/if}
    </div>
    <button class="hamburger" onclick={toggleMenu}>☰</button>
  </nav>
  <div class="mobile-menu" class:open={menuOpen}>
    <a href="#/playing" onclick={closeMenu}>Playing</a>
    <a href="#/browse" onclick={closeMenu}>Browse</a>
    <a href="#/playlist" onclick={closeMenu}>
      Playlist
      {#if $player.playlistLength > 0}
        <span class="nav-badge">{$player.position}/{$player.playlistLength}</span>
      {/if}
    </a>
    {#if !$connected}
      <span class="connection-lost">disconnected</span>
    {/if}
  </div>
</header>
