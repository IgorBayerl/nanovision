import { ThemeProvider } from '@/components/Theme.Context'
import '@/index.css'
import ReactDOM from 'react-dom/client'
import { applyDefaultFilters } from '@/lib/applyDefaultFilters'
import DetailsPage from '@/pages/DetailsPage'
import { TooltipProvider } from '@/ui/tooltip'

/**
 * Retrieves details data from the window.
 * The data is treated as `unknown` because it has not yet been validated.
 */
function getDetailsData(): unknown | null {
    if (window.__NANOVISION_DETAILS__) {
        return window.__NANOVISION_DETAILS__
    }
    // Data for the details page must be embedded directly by the generator.
    return null
}

const rootEl = document.getElementById('root')
if (!rootEl) {
    console.error('Fatal: Missing #root element in HTML.')
} else {
    const data = getDetailsData()
    if (!data) {
        ReactDOM.createRoot(rootEl).render(
            <div className="p-6 text-foreground text-sm">
                No report data found. Ensure the details data object is embedded in the HTML before the main script.
            </div>,
        )
    } else {
        applyDefaultFilters(data)
        ReactDOM.createRoot(rootEl).render(
            <ThemeProvider>
                <TooltipProvider delayDuration={300} skipDelayDuration={300}>
                    <DetailsPage data={data} />
                </TooltipProvider>
            </ThemeProvider>,
        )
    }
}
