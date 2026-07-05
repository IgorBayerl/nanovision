/**
 * Seeds the URL query string from a report's `defaultFilters` (a raw query
 * string like "diff=changed&risk=danger") on first load. Because filter/column
 * state is read from the URL by `useUrlState`, seeding the URL before React
 * mounts makes every configured filter apply automatically.
 *
 * Only runs when the page has no existing query string, so a user's own
 * navigation (shared links, manual filters) always wins over the defaults.
 */
export function applyDefaultFilters(data: unknown): void {
    if (typeof window === 'undefined') return
    if (window.location.search && window.location.search !== '?') return
    if (!data || typeof data !== 'object') return

    const defaultFilters = (data as { defaultFilters?: unknown }).defaultFilters
    if (typeof defaultFilters !== 'string' || defaultFilters.trim() === '') return

    const query = defaultFilters.replace(/^\?/, '').trim()
    if (!query) return

    const newUrl = `${window.location.pathname}?${query}${window.location.hash}`
    window.history.replaceState(window.history.state, '', newUrl)
}
