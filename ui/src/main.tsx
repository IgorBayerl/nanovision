import { useEffect, useState } from 'react'
import { ThemeProvider } from '@/components/Theme.Context'
import '@/index.css'
import ReactDOM from 'react-dom/client'
import { applyDefaultFilters } from '@/lib/applyDefaultFilters'
import DetailsPage from '@/pages/DetailsPage'
import SummaryPage from '@/pages/SummaryPage'
import { TooltipProvider } from '@/ui/tooltip'

// Declare types for global variables attached to the window
declare global {
    interface Window {
        __NANOVISION_SUMMARY__?: unknown
        __NANOVISION_MODE__?: 'single'
        __NANOVISION_FULL_DATA__?: {
            summary: unknown
            details: Record<string, unknown>
        }
    }
}

/**
 * Normal router for the single file mode that watches window.location.hash
 */
function SingleFileApp({ data }: { data: NonNullable<Window['__NANOVISION_FULL_DATA__']> }) {
    const [currentHash, setCurrentHash] = useState(window.location.hash)

    useEffect(() => {
        const handler = () => setCurrentHash(window.location.hash)
        window.addEventListener('hashchange', handler)
        return () => window.removeEventListener('hashchange', handler)
    }, [])

    if (currentHash.startsWith('#/details/')) {
        const filePath = currentHash.replace('#/details/', '')
        const detailsData = data.details[filePath]

        if (!detailsData) {
            return (
                <div className="p-6 text-foreground text-sm">
                    <h2>File Not Found</h2>
                    <p>No details found for path: {filePath}</p>
                    <button
                        type="button"
                        onClick={() => {
                            window.location.hash = ''
                        }}
                        className="mt-4 text-primary underline"
                    >
                        Return to Summary
                    </button>
                </div>
            )
        }

        // The DetailsPage expects its data prop so we can just render it natively here
        return <DetailsPage data={detailsData} />
    }

    return <SummaryPage data={data.summary} />
}

/**
 * Retrieves summary data from the window or a script tag.
 */
function getSummaryData(): unknown | null {
    if (window.__NANOVISION_SUMMARY__) {
        return window.__NANOVISION_SUMMARY__
    }

    const tag = document.getElementById('adl-summary')
    if (tag?.textContent) {
        try {
            return JSON.parse(tag.textContent)
        } catch (err) {
            console.error('Failed to parse #adl-summary JSON:', err)
        }
    }
    return null
}

const rootEl = document.getElementById('root')
if (!rootEl) {
    console.error('Fatal: Missing #root element in HTML.')
} else {
    // Branch on NANOVISION_MODE
    if (window.__NANOVISION_MODE__ === 'single' && window.__NANOVISION_FULL_DATA__) {
        applyDefaultFilters(window.__NANOVISION_FULL_DATA__.summary)
        ReactDOM.createRoot(rootEl).render(
            <ThemeProvider>
                <TooltipProvider delayDuration={300} skipDelayDuration={300}>
                    <SingleFileApp data={window.__NANOVISION_FULL_DATA__} />
                </TooltipProvider>
            </ThemeProvider>,
        )
    } else {
        // Fallback to traditional mode
        const data = getSummaryData()
        if (!data) {
            ReactDOM.createRoot(rootEl).render(
                <div className="p-6 text-foreground text-sm">
                    No report data found. Ensure <code>data.js</code> is loaded and valid before the main application
                    script.
                </div>,
            )
        } else {
            applyDefaultFilters(data)
            ReactDOM.createRoot(rootEl).render(
                <ThemeProvider>
                    <TooltipProvider delayDuration={300} skipDelayDuration={300}>
                        <SummaryPage data={data} />
                    </TooltipProvider>
                </ThemeProvider>,
            )
        }
    }
}
