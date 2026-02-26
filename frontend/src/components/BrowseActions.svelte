<script>
  import { addToPlaylist, updateLibrary } from '../lib/mpd.js';

  let { path, type } = $props();
  let open = $state(false);

  function stop(e) {
    e.stopPropagation();
  }

  function toggle(e) {
    e.stopPropagation();
    open = !open;
  }

  function add(e) {
    e.stopPropagation();
    addToPlaylist(path, type, false, false);
    open = false;
  }

  function addPlay(e) {
    e.stopPropagation();
    addToPlaylist(path, type, false, true);
    open = false;
  }

  function addReplacePlay(e) {
    e.stopPropagation();
    addToPlaylist(path, type, true, true);
    open = false;
  }

  function update(e) {
    e.stopPropagation();
    updateLibrary(path);
    open = false;
  }
</script>

<div class="browse-actions" role="group" onclick={stop} onkeydown={stop}>
  <button onclick={toggle}>☰</button>
  {#if open}
    <div class="browse-actions-menu">
      <button onclick={add}>Add</button>
      <button onclick={addPlay}>Add & play</button>
      <button onclick={addReplacePlay}>Add, replace & play</button>
      {#if type === 'directory'}
        <button onclick={update}>Update this folder</button>
      {/if}
    </div>
  {/if}
</div>

<svelte:window onclick={() => { if (open) open = false; }} />
