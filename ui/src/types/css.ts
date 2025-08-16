import type { CSSProperties } from 'react'

export type CSSWithVars = CSSProperties & {
    [key: `--${string}`]: string | number | undefined
}
