<script>
  import { player } from '../lib/stores.js';
  import * as mpd from '../lib/mpd.js';

  function togglePlayPause() {
    if ($player.state === 'play') mpd.pause();
    else mpd.play();
  }

  function toggleRandom() {
    if ($player.randomOn) mpd.randomOff();
    else mpd.randomOn();
  }
</script>

<div class="player-controls">
  <button onclick={() => mpd.previous()} title="Previous">⏮</button>
  {#if $player.state === 'play'}
    <button onclick={togglePlayPause} title="Pause">⏸</button>
  {:else}
    <button onclick={togglePlayPause} title="Play">▶</button>
  {/if}
  <button onclick={() => mpd.next()} title="Next">⏭</button>
  <button onclick={toggleRandom} class:active={$player.randomOn} title="Random">🔀</button>
</div>
