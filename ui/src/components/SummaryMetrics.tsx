import type { ReactNode } from 'react'
import MetricRow from '@/components/MetricRow'
import ReportsSelector from '@/components/ReportsSelector'
import type { ReportSelectionState } from '@/hooks/useReportSelection'
import { camelCaseToTitleCase } from '@/lib/utils'
import type { StatusBand } from '@/lib/validation'
import type { CoverageDetail, MetadataItem, MetricDefinitions, ScoreDetail, Totals } from '@/types/summary'

type SummaryMetricsProps = {
    info?: {
        title: string
        items: MetadataItem[]
    }
    metrics: Totals
    metricOrder: string[]
    metricDefinitions: MetricDefinitions
    statusBands?: Record<string, StatusBand>
    reportSelection: ReportSelectionState
    /** Metrics that keep their merged value whatever the selection. */
    frozenMetricLabels?: string[]
}

const Section = ({ title, aside, children }: { title?: string; aside?: ReactNode; children: ReactNode }) => (
    <section className="px-5 py-4">
        {title && (
            <div className="mb-1 flex items-baseline justify-between gap-2">
                <h2 className="font-medium text-muted-foreground text-xs uppercase tracking-wide">{title}</h2>
                {aside && <span className="shrink-0 text-muted-foreground text-xs tabular-nums">{aside}</span>}
            </div>
        )}
        {children}
    </section>
)

const MetadataValue = ({ value }: { value: MetadataItem['value'] }) => {
    const display = Array.isArray(value) ? value.join(', ') : String(value ?? '')
    return (
        <span className="truncate text-right font-mono text-foreground" title={display}>
            {display || '-'}
        </span>
    )
}

/**
 * The sidebar overview: report metadata, the report selection and one row per
 * metric, laid out as a single column divided by rules rather than cards.
 */
export default function SummaryMetrics({
    info,
    metrics,
    metricOrder,
    metricDefinitions,
    statusBands,
    reportSelection,
    frozenMetricLabels,
}: SummaryMetricsProps) {
    const infoItems = info?.items ?? []
    const { reports, isSelected } = reportSelection
    const hasReportChoice = reports.length >= 2
    const selectedCount = reports.filter((_, index) => isSelected(index)).length
    const hasBands = !!statusBands && Object.keys(statusBands).length > 0

    return (
        <div className="flex flex-col divide-y divide-sidebar-border">
            {infoItems.length > 0 && (
                <Section title={info?.title}>
                    <dl className="mt-2 flex flex-col gap-1.5 text-xs">
                        {infoItems.map((item) => (
                            <div key={item.label} className="flex items-baseline justify-between gap-3">
                                <dt className="shrink-0 text-muted-foreground">{item.label}</dt>
                                <MetadataValue value={item.value} />
                            </div>
                        ))}
                    </dl>
                </Section>
            )}

            {hasReportChoice && (
                <Section>
                    <ReportsSelector state={reportSelection} frozenMetricLabels={frozenMetricLabels} />
                </Section>
            )}

            <Section
                title="Coverage"
                aside={hasReportChoice ? `${selectedCount} of ${reports.length} reports` : undefined}
            >
                <div className="flex flex-col divide-y divide-sidebar-border">
                    {metricOrder.map((metricId) => {
                        const definition = metricDefinitions[metricId]
                        return (
                            <MetricRow
                                key={metricId}
                                label={definition?.label ?? camelCaseToTitleCase(metricId)}
                                details={metrics[metricId] as CoverageDetail | ScoreDetail | undefined}
                                status={metrics.statuses?.[metricId]}
                                definition={definition}
                                band={statusBands?.[metricId]}
                                hasBands={hasBands}
                            />
                        )
                    })}
                </div>
            </Section>
        </div>
    )
}
