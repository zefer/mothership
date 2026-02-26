<script>
  import { player } from '../lib/stores.js';
  import * as mpd from '../lib/mpd.js';
  import { SkipBack, Play, Pause, SkipForward, Shuffle } from 'lucide-svelte';

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
  <button onclick={() => mpd.previous()} title="Previous"><SkipBack size={18} /></button>
  {#if $player.state === 'play'}
    <button onclick={togglePlayPause} title="Pause"><Pause size={18} /></button>
  {:else}
    <button onclick={togglePlayPause} title="Play"><Play size={18} /></button>
  {/if}
  <button onclick={() => mpd.next()} title="Next"><SkipForward size={18} /></button>
  <button onclick={toggleRandom} class:active={$player.randomOn} title="Random"><Shuffle size={18} /></button>
</div>
