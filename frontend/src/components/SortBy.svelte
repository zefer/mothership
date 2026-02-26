<script>
  let { sort, direction, onChange } = $props();
  let open = $state(false);

  const options = [
    { label: 'Newest first', sort: 'date', direction: 'desc' },
    { label: 'Oldest first', sort: 'date', direction: 'asc' },
    { label: 'A to Z', sort: 'name', direction: 'asc' },
    { label: 'Z to A', sort: 'name', direction: 'desc' },
  ];

  function select(opt) {
    onChange(opt.sort, opt.direction);
    open = false;
  }

  let currentLabel = $derived(
    options.find(o => o.sort === sort && o.direction === direction)?.label || 'Sort'
  );
</script>

<div class="browse-actions" style="position: relative;">
  <button onclick={() => open = !open}>{currentLabel} ▾</button>
  {#if open}
    <div class="browse-actions-menu">
      {#each options as opt}
        <button onclick={() => select(opt)}>{opt.label}</button>
      {/each}
    </div>
  {/if}
</div>

<svelte:window onclick={(e) => { if (open && !e.target.closest('.browse-actions')) open = false; }} />
