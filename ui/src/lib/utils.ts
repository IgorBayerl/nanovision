import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

/**
 * Converts a camelCase string to a Title Case string.
 * Example: "lineCoverage" -> "Line Coverage"
 */
export function camelCaseToTitleCase(text: string): string {
    const result = text.replace(/([A-Z])/g, ' $1')
    return result.charAt(0).toUpperCase() + result.slice(1)
}

const NON_METRIC_KEYS = new Set(['files', 'folders', 'statuses'])

/**
 * The metrics present in `totals`, in the report's configured order. Metrics
 * the order does not name (older reports carry no order) follow in data order.
 */
export function orderedMetricKeys(totals: object | undefined, order: string[] = []): string[] {
    if (!totals) return []
    const present = Object.keys(totals).filter((key) => !NON_METRIC_KEYS.has(key))
    const ordered = order.filter((key) => present.includes(key))
    return [...ordered, ...present.filter((key) => !ordered.includes(key))]
}
