import { useEffect, useState } from 'react'

/**
 * A custom hook that debounces a value.
 * @param value The value to debounce.
 * @param delay The debounce delay in milliseconds.
 * @returns The debounced value.
 */
export function useDebounce<T>(value: T, delay: number): T {
    const [debouncedValue, setDebouncedValue] = useState<T>(value)

    useEffect(() => {
        // Set up a timer to update the debounced value after the specified delay.
        const handler = setTimeout(() => {
            setDebouncedValue(value)
        }, delay)

        // Clean up the timer if the value or delay changes, or if the component unmounts.
        // This prevents the debounced value from updating if a new value is provided
        // before the delay has passed.
        return () => {
            clearTimeout(handler)
        }
    }, [value, delay])

    return debouncedValue
}
