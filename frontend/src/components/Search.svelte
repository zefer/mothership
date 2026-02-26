<script>
  import { X } from 'lucide-svelte';

  let { value = '', onSearch } = $props();
  let input = $state('');
  let timeout;

  $effect(() => {
    input = value;
  });

  function handleInput(e) {
    input = e.target.value;
    clearTimeout(timeout);
    timeout = setTimeout(() => onSearch(input), 200);
  }

  function clear() {
    input = '';
    onSearch('');
  }
</script>

<div class="search-wrapper">
  <input
    type="text"
    placeholder="filter"
    value={input}
    oninput={handleInput}
  />
  {#if input}
    <button class="search-clear" onclick={clear}><X size={16} /></button>
  {/if}
</div>
