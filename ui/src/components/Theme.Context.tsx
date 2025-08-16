import { createContext, type ReactNode, useContext, useEffect, useState } from 'react'

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
 * Wrap the entire application with this to use theme switching.
 */
export function ThemeProvider({ children }: { children: ReactNode }) {
    const [mode, setMode] = useState<Mode>('light')

    // Effect to toggle the 'dark' class on the <html> element
    useEffect(() => {
        const root = document.documentElement
        root.classList.remove('light', 'dark')
        root.classList.add(mode)
    }, [mode])

    return <ThemeContext.Provider value={{ mode, setMode }}>{children}</ThemeContext.Provider>
}
