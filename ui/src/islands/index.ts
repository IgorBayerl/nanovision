import React from 'react'
import { createRoot } from 'react-dom/client'
import { FileDetailsIsland, type IslandProps } from './placeholder'

// All islands must accept a generic bag of props (string->unknown)
type IslandComponent = React.ComponentType<IslandProps>

// Register islands here
const registry: Record<string, IslandComponent> = {
    FileDetails: FileDetailsIsland,
}

export function mountIslands(): void {
    document.querySelectorAll<HTMLElement>('[data-island]').forEach((el) => {
        const name = el.dataset.island
        if (!name) return
        const Comp = registry[name]
        if (!Comp) return

        const props = parseProps(el.dataset.props) // IslandProps
        const holder = document.createElement('div')
        el.appendChild(holder)

        // Use React.createElement to avoid JSX spread typing headaches
        createRoot(holder).render(React.createElement(Comp, props))
    })
}

function parseProps(raw?: string): IslandProps {
    if (!raw) return {}
    try {
        const v = JSON.parse(raw)
        return v && typeof v === 'object' ? (v as Record<string, unknown>) : {}
    } catch {
        return {}
    }
}
