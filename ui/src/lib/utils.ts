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
