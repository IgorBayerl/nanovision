import { AlertCircle, AlertTriangle, ShieldCheck } from 'lucide-react'
import InfoTooltip from '@/components/InfoTooltip'
import { cn } from '@/lib/utils'
import type { StatusBand } from '@/lib/validation'
import type { CoverageDetail, MetricDefinition, RiskLevel, ScoreDetail } from '@/types/summary'

export const StatusIcon = ({ status, showOk = true }: { status: RiskLevel; showOk?: boolean }) => {
    if (status === 'danger') return <AlertCircle className="h-4 w-4 text-uncovered" />
    if (status === 'warning') return <AlertTriangle className="h-4 w-4 text-partial" />
    if (!showOk) return null
    return <ShieldCheck className="h-4 w-4 text-primary" />
}

const formatCount = (value: number) => value.toLocaleString('en-US')
const formatPct = (value: number) => `${Number.isInteger(value) ? value : value.toFixed(1)}%`

// Only a metric that needs attention carries an icon, so a healthy list reads quietly.
const RowLabel = ({ label, status, description }: { label: string; status?: RiskLevel; description?: string }) => (
    <div className="flex min-w-0 items-center gap-1.5">
        {status && <StatusIcon status={status} showOk={false} />}
        <span className="truncate font-medium text-foreground text-sm">{label}</span>
        <InfoTooltip label={`What ${label} means`}>{description}</InfoTooltip>
    </div>
)

type MetricRowProps = {
    label: string
    details: CoverageDetail | ScoreDetail | undefined
    status?: RiskLevel
    definition?: MetricDefinition
    /** Risk band of this metric; its upper bound is drawn as the target. */
    band?: StatusBand
    /** The report configures bands at all, so a metric without one says so. */
    hasBands?: boolean
}

const isScoreDetail = (d: CoverageDetail | ScoreDetail | undefined): d is ScoreDetail =>
    !!d && 'value' in d && !('percentage' in d)

export default function MetricRow({ label, details, status, definition, band, hasBands }: MetricRowProps) {
    if (definition?.kind === 'value' || isScoreDetail(details)) {
        const value = (details as ScoreDetail | undefined)?.value
        const display = value === undefined ? 'N/A' : Number.isInteger(value) ? String(value) : value.toFixed(2)

        return (
            <div className="flex items-center justify-between gap-3 py-3.5 first:pt-2 last:pb-0">
                <RowLabel label={label} status={status} description={definition?.description} />
                <span className="font-semibold text-foreground text-xl tabular-nums leading-none">{display}</span>
            </div>
        )
    }

    const cov = details as CoverageDetail | undefined
    const pct = cov ? Math.max(0, Math.min(100, cov.percentage)) : undefined
    const target = band ? Math.max(0, Math.min(100, band.max)) : undefined

    return (
        <div className="flex flex-col gap-2 py-3.5 first:pt-2 last:pb-0">
            <div className="flex items-center justify-between gap-3">
                <RowLabel label={label} status={status} description={definition?.description} />
                <span className="font-semibold text-foreground text-xl tabular-nums leading-none">
                    {pct !== undefined ? `${Math.round(pct)}%` : 'N/A'}
                </span>
            </div>

            <div className="relative h-1.5 w-full rounded-full bg-foreground/10">
                <div
                    className={cn(
                        'h-full rounded-full transition-[width]',
                        status === 'danger' ? 'bg-uncovered' : status === 'warning' ? 'bg-partial' : 'bg-primary',
                    )}
                    style={{ width: `${pct ?? 0}%` }}
                />
                {target !== undefined && (
                    <div
                        aria-hidden
                        className="-top-0.5 -bottom-0.5 absolute w-0.5 rounded-full bg-foreground/50"
                        style={{ left: `calc(${target}% - 1px)` }}
                    />
                )}
            </div>

            {cov && (
                <div className="flex items-baseline justify-between gap-3 text-muted-foreground text-xs">
                    <span className="font-mono tabular-nums">
                        {formatCount(cov.covered)} / {formatCount(cov.total)}
                    </span>
                    {target !== undefined ? (
                        <span className="tabular-nums">target {formatPct(target)}</span>
                    ) : (
                        hasBands && <span>no threshold</span>
                    )}
                </div>
            )}
        </div>
    )
}
