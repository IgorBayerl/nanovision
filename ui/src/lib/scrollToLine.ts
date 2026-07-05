/**
 * Scrolls the source viewer to a given line number and briefly pulses it.
 * Shared by the methods table and the function-navigation sidebar.
 */
export function scrollToLine(lineNumber: number): void {
    const selector = `[data-line-number="${lineNumber}"]`
    const lineElement = document.querySelector(selector) as HTMLElement | null
    if (!lineElement) return

    lineElement.scrollIntoView({ behavior: 'smooth', block: 'center' })
    lineElement.classList.add('animate-pulse-bg')
    setTimeout(() => {
        lineElement.classList.remove('animate-pulse-bg')
    }, 1500)
}
