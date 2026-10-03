import { useEffect, useReducer } from 'react'
import { parseRoute } from '@/app/router'

declare global {
    interface Window {
        /** The server build; view URLs carry it, so the browser may cache them for good. */
        __NANOVISION_BUILD__?: string
    }
}

const build = () => encodeURIComponent(window.__NANOVISION_BUILD__ ?? '')

async function fetchJSON(url: string): Promise<unknown> {
    const resp = await fetch(url, { headers: { Accept: 'application/json' } })
    if (!resp.ok) {
        let message = `${resp.status} ${resp.statusText}`
        try {
            const body = (await resp.json()) as { error?: string }
            if (body.error) message = body.error
        } catch {
            // not JSON
        }
        throw new Error(message)
    }
    return resp.json()
}

export const api = {
    summary: (id: number) => `/api/v1/runs/${id}/summary?b=${build()}`,
    file: (id: number, path: string) => `/api/v1/runs/${id}/file?path=${encodeURIComponent(path)}&b=${build()}`,
}

interface Entry {
    done: Promise<void>
    value?: unknown
    error?: Error
}

// A saved run never changes, so an answer is kept for the life of the page.
// ponytail: nothing is evicted and a failed load is not tried again; reload
// the page for either. Add a size limit if one tab opens many big runs.
const entries = new Map<string, Entry>()

function load(url: string): Entry {
    let entry = entries.get(url)
    if (!entry) {
        const created: Entry = {
            done: fetchJSON(url).then(
                (value) => {
                    created.value = value
                },
                (error: Error) => {
                    created.error = error
                },
            ),
        }
        entry = created
        entries.set(url, entry)
    }
    return entry
}

export type Resource<T> = { data: T | undefined; error: Error | undefined }

/**
 * Reads an API URL. Loaded data renders at once, so a prefetched page shows
 * without waiting; pass null to load nothing.
 */
export function useApi<T>(url: string | null): Resource<T> {
    const entry = url ? load(url) : undefined
    const [, rerender] = useReducer((n: number) => n + 1, 0)
    useEffect(() => {
        let mounted = true
        entry?.done.then(() => {
            if (mounted) rerender()
        })
        return () => {
            mounted = false
        }
    }, [entry])
    return { data: entry?.value as T | undefined, error: entry?.error }
}

/** Loads the data of the page behind href, for example when a link is hovered. */
export function prefetch(href: string) {
    const url = new URL(href, window.location.href)
    if (url.origin !== window.location.origin) return
    const route = parseRoute(url.pathname)
    if (route.page === 'run') load(api.summary(route.id))
    else if (route.page === 'file') load(api.file(route.id, route.path))
}

/** Prefetches the page of any link the pointer rests on. */
export function prefetchOnHover(): () => void {
    let timer: number | undefined
    const onOver = (event: MouseEvent) => {
        const anchor = (event.target as Element | null)?.closest?.('a')
        const href = anchor?.getAttribute('href')
        if (!href || href.startsWith('#') || href.startsWith('http')) return
        window.clearTimeout(timer)
        // a short pause skips links the pointer only crosses
        timer = window.setTimeout(() => prefetch(href), 60)
    }
    document.addEventListener('mouseover', onOver)
    return () => {
        document.removeEventListener('mouseover', onOver)
        window.clearTimeout(timer)
    }
}
