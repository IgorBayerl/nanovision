import { AlertCircle, AlertTriangle } from 'lucide-react'
import { cn } from '@/lib/utils'

export default function InlineCoverage({
    percentage,
    risk,
    isFolder,
    density,
}: {
    percentage: number
    risk: 'safe' | 'warning' | 'danger'
    isFolder?: boolean
    density: 'comfortable' | 'compact'
}) {
    return (
        <div className={cn('flex w-full items-center', density === 'compact' ? 'gap-2' : 'gap-3')}>
            <span className="inline-flex items-center justify-center" style={{ width: 14 }}>
                {!isFolder ? (
                    risk === 'danger' ? (
                        <AlertCircle className="h-3 w-3 text-destructive" />
                    ) : risk === 'warning' ? (
                        <AlertTriangle className="h-3 w-3 text-accent" />
                    ) : null
                ) : null}
            </span>
            <span className="pl-1 text-right text-foreground text-xs tabular-nums" style={{ width: 32 }}>
                {Math.max(0, Math.min(100, percentage))}%
            </span>
            <div
                className={cn(
                    'relative flex-1 overflow-hidden rounded-xs bg-muted',
                    density === 'compact' ? 'h-2' : 'h-3',
                )}
            >
                <div
                    className="absolute inset-y-0 left-0 bg-primary"
                    style={{ width: `${Math.max(0, Math.min(100, percentage))}%` }}
                />
            </div>
        </div>
    )
}
