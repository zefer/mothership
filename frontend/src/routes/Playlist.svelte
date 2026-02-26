<script>
  import { playlist, player } from '../lib/stores.js';
</script>

<section>
  <h3>Playlist</h3>
  {#if $playlist.length === 0}
    <p>Playlist is empty.</p>
  {:else}
    <div>
      {#each $playlist as item}
        <div class="playlist-item" class:active={$player.position === item.pos}>
          {item.pos}. {item.name}
        </div>
      {/each}
    </div>
  {/if}
</section>
