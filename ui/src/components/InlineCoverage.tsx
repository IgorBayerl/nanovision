import { AlertCircle, AlertTriangle } from 'lucide-react'

export default function InlineCoverage({
    percentage,
    risk,
    isFolder,
}: {
    percentage: number
    risk: 'safe' | 'warning' | 'danger'
    isFolder?: boolean
}) {
    return (
        <div className="flex w-full items-center gap-2 pl-2">
            <span className="inline-flex items-center justify-center" style={{ width: 14 }}>
                {!isFolder ? (
                    risk === 'danger' ? (
                        <AlertCircle className="h-3.5 w-3.5 text-destructive" />
                    ) : risk === 'warning' ? (
                        <AlertTriangle className="h-3.5 w-3.5 text-accent" />
                    ) : null
                ) : null}
            </span>
            <span className="pl-1 text-right text-foreground text-xs tabular-nums" style={{ width: 32 }}>
                {Math.max(0, Math.min(100, percentage))}%
            </span>
            <div className="relative h-2 flex-1 overflow-hidden rounded-xs bg-uncovered">
                <div
                    className="absolute inset-y-0 left-0 rounded-xs bg-covered"
                    style={{ width: `${Math.max(0, Math.min(100, percentage))}%` }}
                />
            </div>
        </div>
    )
}
