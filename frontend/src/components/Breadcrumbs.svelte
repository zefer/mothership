<script>
  let { crumbs, onNavigate } = $props();
</script>

<ol class="breadcrumbs">
  {#each crumbs as crumb, i}
    <li>
      {#if i < crumbs.length - 1}
        <a href="#/browse/{crumb.path}" onclick={(e) => { e.preventDefault(); onNavigate(crumb.path); }}>
          {crumb.label}
        </a>
      {:else}
        {crumb.label}
      {/if}
    </li>
  {/each}
</ol>
