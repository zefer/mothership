<script>
  import { querystring } from 'svelte-spa-router';
  import { ls } from '../lib/mpd.js';
  import { Folder, Music, ListMusic } from 'lucide-svelte';
  import Breadcrumbs from '../components/Breadcrumbs.svelte';
  import Search from '../components/Search.svelte';
  import SortBy from '../components/SortBy.svelte';
  import Pagination from '../components/Pagination.svelte';
  import BrowseActions from '../components/BrowseActions.svelte';

  let { params = {} } = $props();

  const MAX_PER_PAGE = 200;

  let items = $state([]);
  let loading = $state(false);

  // Parse route params from wildcard and querystring.
  let uri = $derived(params?.wild || '/');
  let qsParams = $derived(new URLSearchParams($querystring || ''));
  let sort = $derived(qsParams.get('sort') || localStorage.getItem('m-sort') || 'date');
  let direction = $derived(qsParams.get('direction') || localStorage.getItem('m-direction') || 'desc');
  let filter = $derived(qsParams.get('filter') || '');
  let page = $derived(parseInt(qsParams.get('page'), 10) || 1);

  // Pagination.
  let totalPages = $derived(Math.ceil(items.length / MAX_PER_PAGE));
  let pages = $derived(Array.from({ length: totalPages }, (_, i) => i + 1));
  let pagedItems = $derived(items.slice((page - 1) * MAX_PER_PAGE, page * MAX_PER_PAGE));

  // Breadcrumbs.
  let crumbs = $derived(buildCrumbs(uri));

  function buildCrumbs(uri) {
    if (!uri || uri === '/') return [{ label: 'home', path: '' }];
    const parts = uri.split('/');
    const result = [{ label: 'home', path: '' }];
    for (let i = 0; i < parts.length; i++) {
      result.push({
        label: parts[i],
        path: parts.slice(0, i + 1).join('/'),
      });
    }
    return result;
  }

  async function fetchItems(uri, sort, direction, filter) {
    loading = true;
    try {
      items = (await ls(uri, sort, direction, filter)) || [];
    } catch {
      items = [];
    }
    loading = false;
  }

  // Save sort preferences.
  $effect(() => {
    localStorage.setItem('m-sort', sort);
    localStorage.setItem('m-direction', direction);
  });

  // Fetch on param changes.
  $effect(() => {
    fetchItems(uri, sort, direction, filter);
  });

  function navigate(path) {
    window.location.hash = `/browse/${path}`;
  }

  function buildQs(overrides) {
    const p = new URLSearchParams();
    const s = overrides.sort ?? sort;
    const d = overrides.direction ?? direction;
    const f = overrides.filter ?? filter;
    const pg = overrides.page ?? 1;
    if (s) p.set('sort', s);
    if (d) p.set('direction', d);
    if (f) p.set('filter', f);
    if (pg > 1) p.set('page', pg);
    return p.toString();
  }

  function setSort(s, d) {
    const qs = buildQs({ sort: s, direction: d, page: 1 });
    window.location.hash = `/browse/${uri === '/' ? '' : uri}?${qs}`;
  }

  function setFilter(f) {
    const qs = buildQs({ filter: f, page: 1 });
    window.location.hash = `/browse/${uri === '/' ? '' : uri}?${qs}`;
  }

  function setPage(p) {
    const qs = buildQs({ page: p });
    window.location.hash = `/browse/${uri === '/' ? '' : uri}?${qs}`;
  }

</script>

<section>
  <Breadcrumbs {crumbs} onNavigate={navigate} />

  <div class="browse-toolbar">
    <Search value={filter} onSearch={setFilter} />
    <SortBy {sort} {direction} onChange={setSort} />
    <Pagination {pages} currentPage={page} onChange={setPage} />
  </div>

  {#if loading}
    <p>Loading...</p>
  {:else if pagedItems.length === 0}
    <p>No items found.</p>
  {:else}
    <div class="browse-list">
      {#each pagedItems as item}
        <div class="browse-item" onclick={() => item.type === 'directory' ? navigate(item.path) : null}>
          <span class="browse-item-icon">
            {#if item.type === 'directory'}
              <Folder size={18} />
            {:else if item.type === 'playlist'}
              <ListMusic size={18} />
            {:else}
              <Music size={18} />
            {/if}
          </span>
          <span class="browse-item-name">{item.base}</span>
          <BrowseActions path={item.path} type={item.type} />
        </div>
      {/each}
    </div>
  {/if}

  {#if totalPages > 1}
    <div class="browse-toolbar" style="margin-top: 1rem;">
      <SortBy {sort} {direction} onChange={setSort} />
      <Pagination {pages} currentPage={page} onChange={setPage} />
    </div>
  {/if}
</section>
