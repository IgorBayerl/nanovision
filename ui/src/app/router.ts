import { useSyncExternalStore } from 'react'

// The server UI is one document. Routes are plain paths, so every view has a
// short URL people can paste into a review or a chat.
export type Route =
    | { page: 'home' }
    | { page: 'run'; id: number }
    | { page: 'file'; id: number; path: string }
    | { page: 'notFound' }

export function parseRoute(pathname: string): Route {
    const parts = pathname.split('/').filter(Boolean)
    const id = (value: string | undefined) => {
        const n = Number(value)
        return Number.isInteger(n) && n > 0 ? n : null
    }

    if (parts.length === 0) return { page: 'home' }
    if (parts[0] === 'runs') {
        const runId = id(parts[1])
        if (runId === null) return { page: 'notFound' }
        if (parts.length === 2) return { page: 'run', id: runId }
        if (parts[2] === 'files' && parts.length > 3) {
            return { page: 'file', id: runId, path: parts.slice(3).map(decodeURIComponent).join('/') }
        }
    }
    return { page: 'notFound' }
}

export const runHref = (id: number) => `/runs/${id}`

export const fileHref = (id: number, path: string) =>
    `/runs/${id}/files/${path.split('/').map(encodeURIComponent).join('/')}`

// --- navigation ---

const NAVIGATE_EVENT = 'nanovision:navigate'

export function navigate(href: string, options: { replace?: boolean } = {}) {
    const url = new URL(href, window.location.href)
    const same = url.pathname === window.location.pathname && url.search === window.location.search
    if (options.replace) window.history.replaceState(null, '', url)
    else window.history.pushState(null, '', url)
    window.dispatchEvent(new Event(NAVIGATE_EVENT))
    if (url.hash) {
        requestAnimationFrame(() => document.getElementById(url.hash.slice(1))?.scrollIntoView())
    } else if (!same) {
        window.scrollTo(0, 0)
    }
}

function subscribe(onChange: () => void) {
    window.addEventListener('popstate', onChange)
    window.addEventListener(NAVIGATE_EVENT, onChange)
    return () => {
        window.removeEventListener('popstate', onChange)
        window.removeEventListener(NAVIGATE_EVENT, onChange)
    }
}

const snapshot = () => window.location.pathname + window.location.search

/** The current path and query; re-renders on every navigation. */
export function useLocationKey(): string {
    return useSyncExternalStore(subscribe, snapshot)
}

/**
 * Turns clicks on same-origin links into client-side navigation, so the report
 * components, which render plain anchors, move between pages without a reload.
 */
export function interceptLinks(onHref?: (href: string) => void): () => void {
    const onClick = (event: MouseEvent) => {
        if (event.defaultPrevented || event.button !== 0) return
        if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return
        const anchor = (event.target as Element | null)?.closest?.('a')
        if (!anchor || anchor.target === '_blank' || anchor.hasAttribute('download')) return
        const href = anchor.getAttribute('href')
        if (!href || href.startsWith('#')) return
        const url = new URL(href, window.location.href)
        // the root is the server's redirect to the newest run, not a page of the app
        if (url.origin !== window.location.origin || url.pathname === '/' || url.pathname.startsWith('/api/')) return
        event.preventDefault()
        onHref?.(url.pathname + url.search + url.hash)
        navigate(url.pathname + url.search + url.hash)
    }
    document.addEventListener('click', onClick)
    return () => document.removeEventListener('click', onClick)
}
