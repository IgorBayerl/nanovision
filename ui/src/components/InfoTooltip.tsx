import { Info } from 'lucide-react'
import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/ui/tooltip'

interface InfoTooltipProps {
    /** Tooltip body. Nothing renders when it is empty. */
    children?: ReactNode
    label: string
    className?: string
    /** Draws attention to a caveat rather than a plain explanation. */
    tone?: 'muted' | 'warning'
}

/**
 * A hover target that explains a label or a caveat.
 *
 * It renders unconditionally wherever it is placed, so a message never appears
 * or disappears in the layout — the reason this exists instead of inline notes,
 * which pushed the surrounding UI around as the report selection changed.
 */
export default function InfoTooltip({ children, label, className, tone = 'muted' }: InfoTooltipProps) {
    if (!children) return null

    return (
        <Tooltip>
            <TooltipTrigger asChild>
                {/* A button so the explanation is reachable by keyboard; it has
                    no action of its own, the tooltip opens on hover and focus. */}
                <button
                    type="button"
                    aria-label={label}
                    className={cn(
                        'inline-flex shrink-0 cursor-help items-center justify-center rounded-sm outline-none transition-colors focus-visible:ring-1 focus-visible:ring-ring',
                        tone === 'warning'
                            ? 'text-partial hover:text-partial/80'
                            : 'text-muted-foreground/60 hover:text-foreground',
                        className,
                    )}
                >
                    <Info className="h-3.5 w-3.5" />
                </button>
            </TooltipTrigger>
            <TooltipContent className="max-w-xs text-pretty">{children}</TooltipContent>
        </Tooltip>
    )
}
