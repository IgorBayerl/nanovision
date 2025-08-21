import { createContext, type ReactNode, useContext, useEffect, useState } from 'react'

const THEME_STORAGE_KEY = 'adlercov-ui-theme'

export type Mode = 'light' | 'dark'

type ThemeContextValue = {
    mode: Mode
    setMode: (m: Mode) => void
}

export const ThemeContext = createContext<ThemeContextValue | null>(null)

/**
 * Hook to access the current theme (light/dark) and a function to update it.
 */
export function useTheme() {
    const ctx = useContext(ThemeContext)
    if (!ctx) {
        throw new Error('useTheme must be used within a ThemeProvider')
    }
    return ctx
}

/**
 * A helper function to determine the initial theme based on storage or system preference.
 * This runs only once when the application starts.
 */
function getInitialTheme(): Mode {
    return 'light'
    // // Check for a user's explicit preference in localStorage.
    // const storedTheme = window.localStorage.getItem(THEME_STORAGE_KEY)
    // if (storedTheme === 'light' || storedTheme === 'dark') {
    //     return storedTheme
    // }

    // // If no preference is stored, fall back to the system's color scheme.
    // const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches
    // return prefersDark ? 'dark' : 'light'
}

/**
 * Wrap the entire application with this to use theme switching.
 */
export function ThemeProvider({ children }: { children: ReactNode }) {
    // 2. Initialize the state by calling our new function.
    // Using a function initializer ensures it runs only on the first render.
    const [mode, setMode] = useState<Mode>(getInitialTheme)

    // 3. This effect runs whenever the `mode` changes, SAVING the new value to localStorage.
    useEffect(() => {
        try {
            window.localStorage.setItem(THEME_STORAGE_KEY, mode)
        } catch (e) {
            console.error('Failed to save theme to localStorage', e)
        }
    }, [mode])

    // This effect remains the same. It applies the 'dark' or 'light' class to the <html> tag.
    useEffect(() => {
        const root = document.documentElement
        root.classList.remove('light', 'dark')
        root.classList.add(mode)
    }, [mode])

    return <ThemeContext.Provider value={{ mode, setMode }}>{children}</ThemeContext.Provider>
}
