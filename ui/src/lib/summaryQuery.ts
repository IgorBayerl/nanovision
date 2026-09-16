import { useEffect } from 'react'

// The report selection is shared by every screen, so the details page's value wins.
const SHARED_PARAM = 'reports'

// Keyed by the summary document, so two reports open in one tab keep their own filters.
const storageKey = () => `nanovision-summary-query:${new URL('./index.html', window.location.href).pathname}`

const remember = () => {
    try {
        window.sessionStorage.setItem(storageKey(), window.location.search)
    } catch {
        // storage blocked: the back button falls back to the current query
    }
}

/**
 * Saves the summary page's query string (its filters, sort and columns) when
 * the user leaves it, so the details page's back button can restore them. The
 * browser's own back button restores the URL by itself; this link does not.
 *
 * Clicks are caught too: a link opened in a new tab never fires `pagehide`,
 * and single-file reports move to details by hash without leaving the page.
 */
export function useRememberSummaryQuery() {
    useEffect(() => {
        window.addEventListener('pagehide', remember)
        document.addEventListener('click', remember, true)
        document.addEventListener('auxclick', remember, true)
        return () => {
            window.removeEventListener('pagehide', remember)
            document.removeEventListener('click', remember, true)
            document.removeEventListener('auxclick', remember, true)
        }
    }, [])
}

/** Link back to the summary with its saved filters and the current report selection. */
export function summaryHref(): string {
    let saved: string | null = null
    try {
        saved = window.sessionStorage.getItem(storageKey())
    } catch {
        // storage blocked
    }
    if (saved === null) return `./index.html${window.location.search}`

    const params = new URLSearchParams(saved)
    const shared = new URLSearchParams(window.location.search).get(SHARED_PARAM)
    if (shared === null) params.delete(SHARED_PARAM)
    else params.set(SHARED_PARAM, shared)

    // URLSearchParams escapes commas; the filters read them back either way.
    const query = params.toString().replaceAll('%2C', ',')
    return `./index.html${query ? `?${query}` : ''}`
}
