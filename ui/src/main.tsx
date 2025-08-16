import ReactDOM from 'react-dom/client'
import '@/index.css'
import { ThemeProvider } from '@/components/Theme.Context'
import SummaryPage from '@/pages/SummaryPage'
import type { SummaryV1 } from '@/types/summary'

function getSummaryData(): SummaryV1 | null {
    // Prod mode: injected by Go via data.js (window.__ADLERCOV_SUMMARY__)
    if (window.__ADLERCOV_SUMMARY__) return window.__ADLERCOV_SUMMARY__

    // fallback
    const tag = document.getElementById('adl-summary')
    if (tag?.textContent) {
        try {
            return JSON.parse(tag.textContent) as SummaryV1
        } catch (err) {
            console.error('Failed to parse #adl-summary JSON:', err)
        }
    }
    return null
}

const rootEl = document.getElementById('root')
if (!rootEl) {
    console.error('Missing #root element')
} else {
    const data = getSummaryData()
    if (!data) {
        ReactDOM.createRoot(rootEl).render(
            <div className="p-6 text-foreground text-sm">
                No report data found. Ensure <code>data.js</code> runs before <code>react-islands.js</code>.
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
