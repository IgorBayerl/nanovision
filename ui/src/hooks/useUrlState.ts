import { useCallback, useState } from 'react'

const getUrlValue = (key: string) => {
    if (typeof window === 'undefined') return null
    const params = new URLSearchParams(window.location.search)
    return params.get(key)
}

const parsers = {
    string: (value: string) => value,
    boolean: (value: string) => value === 'true',
    json: <T>(value: string): T | null => {
        try {
            return JSON.parse(value) as T
        } catch {
            return null
        }
    },
}

export function useUrlState<T>(key: string, defaultValue: T): [T, (newValue: T) => void] {
    const [state, setState] = useState<T>(() => {
        const valueFromUrl = getUrlValue(key)
        if (valueFromUrl === null) {
            return defaultValue
        }

        if (typeof defaultValue === 'boolean') {
            return parsers.boolean(valueFromUrl) as unknown as T
        }
        if (typeof defaultValue === 'string') {
            return parsers.string(valueFromUrl) as T
        }
        const parsed = parsers.json<T>(valueFromUrl)
        return parsed !== null ? parsed : defaultValue
    })

    const setUrlState = useCallback(
        (newValue: T) => {
            setState(newValue)

            const params = new URLSearchParams(window.location.search)

            if (JSON.stringify(newValue) === JSON.stringify(defaultValue)) {
                params.delete(key)
            } else {
                const valueToSet = typeof newValue === 'string' ? newValue : JSON.stringify(newValue)
                params.set(key, valueToSet)
            }

            const queryString = params.toString()
            const newUrl = queryString ? `${window.location.pathname}?${queryString}` : window.location.pathname

            window.history.replaceState({ ...window.history.state, as: newUrl, url: newUrl }, '', newUrl)
        },
        [key, defaultValue],
    )

    return [state, setUrlState]
}
