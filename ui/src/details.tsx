import ReactDOM from 'react-dom/client'
import '@/index.css'
import { ThemeProvider } from '@/components/Theme.Context'
import DetailsPage from '@/pages/DetailsPage'

/**
 * Retrieves details data from the window.
 * The data is treated as `unknown` because it has not yet been validated.
 */
function getDetailsData(): unknown | null {
    if (window.__ADLERCOV_DETAILS__) {
        return window.__ADLERCOV_DETAILS__
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
        ReactDOM.createRoot(rootEl).render(
            <ThemeProvider>
                <DetailsPage data={data} />
            </ThemeProvider>,
        )
    }
}
