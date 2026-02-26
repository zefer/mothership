<script>
  let { pages = [], currentPage = 1, onChange } = $props();
</script>

{#if pages.length > 1}
  <div class="pagination">
    <button disabled={currentPage <= 1} onclick={() => onChange(currentPage - 1)}>«</button>
    {#each pages as p}
      <button class:active={p === currentPage} onclick={() => onChange(p)}>{p}</button>
    {/each}
    <button disabled={currentPage >= pages.length} onclick={() => onChange(currentPage + 1)}>»</button>
  </div>
{/if}
