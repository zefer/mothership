<script>
  import { ArrowDownWideNarrow, CalendarArrowUp, CalendarArrowDown, ArrowDownAZ, ArrowDownZA } from 'lucide-svelte';

  let { sort, direction, onChange } = $props();
  let open = $state(false);

  const options = [
    { label: 'Newest first', sort: 'date', direction: 'desc', icon: CalendarArrowUp },
    { label: 'Oldest first', sort: 'date', direction: 'asc', icon: CalendarArrowDown },
    { label: 'A to Z', sort: 'name', direction: 'asc', icon: ArrowDownAZ },
    { label: 'Z to A', sort: 'name', direction: 'desc', icon: ArrowDownZA },
  ];

  function select(opt) {
    onChange(opt.sort, opt.direction);
    open = false;
  }
</script>

<div class="browse-actions" style="position: relative;">
  <button onclick={() => open = !open}><ArrowDownWideNarrow size={18} /></button>
  {#if open}
    <div class="browse-actions-menu">
      {#each options as opt}
        <button onclick={() => select(opt)}><opt.icon size={15} /> {opt.label}</button>
      {/each}
    </div>
  {/if}
</div>

<svelte:window onclick={(e) => { if (open && !e.target.closest('.browse-actions')) open = false; }} />
