import { useEffect, useMemo } from 'react'
import { prefetchOnHover } from '@/app/api'
import { AppShell } from '@/app/components/Chrome'
import FilePage from '@/app/pages/FilePage'
import RunPage from '@/app/pages/RunPage'
import { interceptLinks, parseRoute, type Route, useLocationKey } from '@/app/router'

function pageTitle(route: Route): string {
    switch (route.page) {
        case 'run':
            return `Run ${route.id} · nanovision`
        case 'file':
            return `${route.path.slice(route.path.lastIndexOf('/') + 1)} · run ${route.id} · nanovision`
        default:
            return 'nanovision'
    }
}

/**
 * The web UI of `nanovision serve`: the report pages, fed by the API instead
 * of data.js. The server sends its root to the newest run.
 */
export default function ServerApp() {
    const location = useLocationKey()
    const route = useMemo(() => parseRoute(new URL(location, window.location.origin).pathname), [location])

    useEffect(() => {
        const stopLinks = interceptLinks()
        const stopPrefetch = prefetchOnHover()
        return () => {
            stopLinks()
            stopPrefetch()
        }
    }, [])

    useEffect(() => {
        document.title = pageTitle(route)
    }, [route])

    switch (route.page) {
        case 'run':
            return <RunPage key={`run:${route.id}`} id={route.id} />
        case 'file':
            return <FilePage key={`file:${route.id}:${route.path}`} id={route.id} path={route.path} />
        case 'home':
            // the server only serves this page while its store is empty
            return (
                <AppShell>
                    <p className="text-muted-foreground text-sm">
                        This server has no runs yet. A CI build uploads one with{' '}
                        <code className="font-mono">nanovision -store {window.location.origin} -run-kind submit</code>.
                    </p>
                </AppShell>
            )
        default:
            return (
                <AppShell>
                    <p className="text-muted-foreground text-sm">
                        There is no page at <code className="font-mono">{window.location.pathname}</code>.{' '}
                        <a href="/" className="text-primary hover:underline">
                            Open the newest run
                        </a>
                    </p>
                </AppShell>
            )
    }
}
