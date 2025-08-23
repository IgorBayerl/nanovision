import { AlertCircle, AlertTriangle } from 'lucide-react'
import { Progress } from '@/ui/progress' // 1. Import the Progress component

export default function InlineCoverage({
    percentage,
    risk,
    isFolder,
}: {
    percentage: number
    risk: 'safe' | 'warning' | 'danger'
    isFolder?: boolean
}) {
    const clampedPercentage = Math.round(Math.max(0, Math.min(100, percentage)))

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
                {clampedPercentage}%
            </span>
            <Progress
                value={clampedPercentage}
                className="h-2 flex-1 rounded-xs bg-uncovered"
                indicatorClassName="bg-covered rounded-xs"
            />
        </div>
    )
}
