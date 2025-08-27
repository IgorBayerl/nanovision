import ReactDOM from 'react-dom/client'
import '@/index.css'
import { ThemeProvider } from '@/components/Theme.Context'
import { mountIslands } from '@/islands'

const rootEl = document.getElementById('root')

function App() {
    // Mount the interactive islands
    mountIslands()

    // The rest of the page is static HTML, so we don't need to render a full React app.
    // If you wanted a root component for this page, you could render it here.
    // For now, returning null is fine as `mountIslands` does the work.
    return null
}

if (!rootEl) {
    // If there's no #root, we can just mount the theme provider and island mounter
    // without a dedicated root element.
    mountIslands()
} else {
    ReactDOM.createRoot(rootEl).render(
        <ThemeProvider>
            <App />
        </ThemeProvider>,
    )
}
