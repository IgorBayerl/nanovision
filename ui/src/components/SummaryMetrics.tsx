import InfoCard from '@/components/InfoCard'
import MetricCard from '@/components/MetricCard'
import { camelCaseToTitleCase } from '@/lib/utils'
import type { CoverageDetail, MetadataItem, MetricDefinitions, Totals } from '@/types/summary'

type SummaryMetricsProps = {
    info?: {
        title: string
        items: MetadataItem[]
    }
    metrics: Totals
    metricOrder: string[]
    metricDefinitions: MetricDefinitions
}

export default function SummaryMetrics({ info, metrics, metricOrder, metricDefinitions }: SummaryMetricsProps) {
    return (
        <div className="flex flex-wrap gap-4">
            {info && info.items.length > 0 && (
                <div className="flex-grow rounded-lg">
                    <InfoCard title={info.title} items={info.items} />
                </div>
            )}

            {metricOrder.map((metricId) => {
                const metricDetails = metrics[metricId] as CoverageDetail | undefined
                const status = metrics.statuses?.[metricId]
                const definition = metricDefinitions[metricId]
                const label = definition?.label ?? camelCaseToTitleCase(metricId)

                return (
                    <div key={metricId} className="min-w-sm flex-grow lg:max-w-1/2">
                        <MetricCard label={label} details={metricDetails} status={status} definition={definition} />
                    </div>
                )
            })}
        </div>
    )
}
