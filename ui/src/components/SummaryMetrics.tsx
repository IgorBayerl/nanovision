import InfoCard from '@/components/InfoCard'
import MetricCard from '@/components/MetricCard'
import { camelCaseToTitleCase, cn } from '@/lib/utils'
import type { CoverageDetail, MetadataItem, MetricDefinitions, ScoreDetail, Totals } from '@/types/summary'

type SummaryMetricsProps = {
    info?: {
        title: string
        items: MetadataItem[]
    }
    metrics: Totals
    metricOrder: string[]
    metricDefinitions: MetricDefinitions
    /** 'grid' = wrapping cards (default); 'sidebar' = stacked compact column. */
    variant?: 'grid' | 'sidebar'
}

export default function SummaryMetrics({
    info,
    metrics,
    metricOrder,
    metricDefinitions,
    variant = 'grid',
}: SummaryMetricsProps) {
    const isSidebar = variant === 'sidebar'

    return (
        <div className={cn(isSidebar ? 'flex flex-col gap-2' : 'flex flex-wrap gap-4')}>
            {info && info.items.length > 0 && (
                <div className={cn(isSidebar ? '' : 'flex-grow rounded-lg')}>
                    <InfoCard title={info.title} items={info.items} />
                </div>
            )}

            {metricOrder.map((metricId) => {
                const metricDetails = metrics[metricId] as CoverageDetail | ScoreDetail | undefined
                const status = metrics.statuses?.[metricId]
                const definition = metricDefinitions[metricId]
                const label = definition?.label ?? camelCaseToTitleCase(metricId)

                return (
                    <div key={metricId} className={cn(isSidebar ? '' : 'min-w-sm flex-grow lg:max-w-1/2')}>
                        <MetricCard
                            label={label}
                            details={metricDetails}
                            status={status}
                            definition={definition}
                            compact={isSidebar}
                        />
                    </div>
                )
            })}
        </div>
    )
}
