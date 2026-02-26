<script>
  import { addToPlaylist, updateLibrary } from '../lib/mpd.js';
  import { EllipsisVertical, Plus, Play, CornerUpRight, RefreshCw } from 'lucide-svelte';

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
  <button onclick={toggle}><EllipsisVertical size={18} /></button>
  {#if open}
    <div class="browse-actions-menu">
      <button onclick={add}><Plus size={15} /> Add</button>
      <button onclick={addPlay}><Play size={15} /> Add & play</button>
      <button onclick={addReplacePlay}><CornerUpRight size={15} /> Add, replace & play</button>
      {#if type === 'directory'}
        <button onclick={update}><RefreshCw size={15} /> Update this folder</button>
      {/if}
    </div>
  {/if}
</div>

<svelte:window onclick={() => { if (open) open = false; }} />
