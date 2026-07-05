import { cn } from '@/lib/utils'
import type { DiffStatus } from '@/types/summary'

/**
 * VSCode-style git decoration: a single colored letter — "A" (green) for added,
 * "M" (amber) for modified. Shared across the file tree, method table and
 * function navigation so the indicator looks identical everywhere. Pass a
 * `className` (e.g. `ml-auto`) to control placement at the call site.
 */
export default function DiffStatusBadge({ status, className }: { status?: DiffStatus; className?: string }) {
    if (status !== 'added' && status !== 'modified') return null

    const letter = status === 'added' ? 'A' : 'M'
    const colorClass = status === 'added' ? 'text-covered' : 'text-partial'
    const label = status === 'added' ? 'Added' : 'Modified'

    return (
        <span className={cn('shrink-0 select-none font-bold font-mono text-xs', colorClass, className)} title={label}>
            {letter}
        </span>
    )
}
