import { ThemeProvider } from '@/components/Theme.Context'
import '@/index.css'
import SummaryPage from '@/pages/SummaryPage'
import ReactDOM from 'react-dom/client'

/**
 * Retrieves summary data from the window or a script tag.
 * The data is treated as `unknown` because it has not yet been validated.
 */
function getSummaryData(): unknown | null {
    // The primary method: data injected by Go into a global variable.
    if (window.__NANOVISION_SUMMARY__) {
        return window.__NANOVISION_SUMMARY__
    }

    // fallback
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
    const data = getSummaryData()
    if (!data) {
        ReactDOM.createRoot(rootEl).render(
            <div className="p-6 text-foreground text-sm">
                No report data found. Ensure <code>data.js</code> is loaded and valid before the main application
                script.
            </div>,
        )
    } else {
        ReactDOM.createRoot(rootEl).render(
            <ThemeProvider>
                <SummaryPage data={data} />
            </ThemeProvider>,
        )
    }
}
