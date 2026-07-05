import { AlertCircle, AlertTriangle, ShieldCheck } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { CoverageDetail, MetricDefinition, RiskLevel, ScoreDetail } from '@/types/summary'
import { Card } from '@/ui/card'
import { Progress } from '@/ui/progress'

export const StatusIcon = ({ status, showOk = true }: { status: RiskLevel; showOk?: boolean }) => {
    if (status === 'danger') return <AlertCircle className="h-4 w-4 text-uncovered" />
    if (status === 'warning') return <AlertTriangle className="h-4 w-4 text-partial" />
    if (!showOk) return null
    return <ShieldCheck className="h-4 w-4 text-primary" />
}

const CardLabel = ({ label, status }: { label: string; status?: RiskLevel }) => (
    <div className="flex min-w-0 items-center gap-1.5">
        {status && <StatusIcon status={status} />}
        <span className="truncate font-medium text-muted-foreground text-xs uppercase tracking-wide">{label}</span>
    </div>
)

// Compact "label: value" pair shown under a percentage card. Labels stay dimmed;
// only the value uses the foreground color, and nothing changes on hover.
const SubMetric = ({ label, value }: { label: string; value: number | string }) => (
    <span className="text-muted-foreground text-xs">
        {label} <span className="font-medium font-mono text-foreground tabular-nums">{value}</span>
    </span>
)

type MetricCardProps = {
    label: string
    details: CoverageDetail | ScoreDetail | undefined
    status?: RiskLevel
    definition?: MetricDefinition
    /** Renders a denser layout suited to the narrow sidebar column. */
    compact?: boolean
}

const isScoreDetail = (d: CoverageDetail | ScoreDetail | undefined): d is ScoreDetail =>
    !!d && 'value' in d && !('percentage' in d)

export default function MetricCard({ label, details, status, definition, compact }: MetricCardProps) {
    if (definition?.kind === 'value' || isScoreDetail(details)) {
        return <ValueCardBody label={label} details={details as ScoreDetail | undefined} status={status} />
    }

    const cov = details as CoverageDetail | undefined
    const pct = cov ? Math.max(0, Math.min(100, Math.round(cov.percentage))) : undefined
    const subMetrics = cov && definition ? definition.subMetrics.filter((s) => s.id !== 'percentage') : []

    return (
        <Card className="flex w-full flex-col gap-2 rounded-md px-3 py-2.5">
            <div className="flex items-center justify-between gap-2">
                <CardLabel label={label} status={status} />
                <span
                    className={cn(
                        'font-bold text-foreground leading-none tabular-nums',
                        compact ? 'text-2xl' : 'text-3xl',
                    )}
                >
                    {pct !== undefined ? `${pct}%` : 'N/A'}
                </span>
            </div>
            <Progress value={pct} className="h-1.5" indicatorClassName="bg-primary" />
            {subMetrics.length > 0 && cov && (
                <div className="flex flex-wrap gap-x-3 gap-y-0.5">
                    {subMetrics.map((sub) => (
                        <SubMetric key={sub.id} label={sub.label} value={cov[sub.id as keyof typeof cov] ?? '-'} />
                    ))}
                </div>
            )}
        </Card>
    )
}

/**
 * Scalar-value card: a plain number with no percentage or progress bar.
 * Used for metrics like Max Cyclomatic Complexity or (future) CRAP score.
 */
function ValueCardBody({
    label,
    details,
    status,
}: {
    label: string
    details: ScoreDetail | undefined
    status?: RiskLevel
}) {
    const value = details?.value
    const display = value === undefined ? 'N/A' : Number.isInteger(value) ? String(value) : value.toFixed(2)

    return (
        <Card className="flex w-full flex-row items-center justify-between gap-2 rounded-md px-3 py-2.5">
            <CardLabel label={label} status={status} />
            <span className="font-bold text-2xl text-foreground leading-none tabular-nums">{display}</span>
        </Card>
    )
}
