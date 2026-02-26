<script>
  import { player, connected } from '../lib/stores.js';
</script>

<section>
  {#if $player.now}
    <h3>{$player.now}</h3>
    {#if $player.sub}
      <p>{$player.sub}</p>
    {/if}
  {/if}

  {#if $player.error}
    <div class="alert alert-error">{$player.error}</div>
  {/if}

  {#if !$connected}
    <div class="alert alert-error">Connection lost</div>
  {/if}

  {#if $player.state === 'pause'}
    <div class="alert alert-warning">Paused</div>
  {:else if $player.state === 'stop'}
    <div class="alert alert-warning">Stopped</div>
  {/if}

  {#if $player.percentage > 0 || $player.state === 'play'}
    <div class="progress-bar">
      <div class="progress-fill" style="width: {$player.percentage}%"></div>
    </div>
    <div class="progress-text">
      <span>{$player.elapsed}</span>
      <span>{$player.total}</span>
    </div>
  {/if}

  {#if $player.quality}
    <p><small>{$player.quality}</small></p>
  {/if}
</section>
