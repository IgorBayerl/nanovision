import { useMemo } from 'react'
import InfoCard from '@/components/InfoCard'
import MetricCard from '@/components/MetricCard'
import { camelCaseToTitleCase } from '@/lib/utils'
import type { CoverageDetail, MetadataItem } from '@/types/summary'

// This threshold determines when the card should expand to take up more grid space.
const LINE_WEIGHT_THRESHOLD = 5

type SummaryMetricsProps = {
    info?: {
        title: string
        items: MetadataItem[]
    }
    metrics: Record<string, CoverageDetail | number>
    metricOrder: string[]
}

export default function SummaryMetrics({ info, metrics, metricOrder }: SummaryMetricsProps) {
    const infoCardColSpanClass = useMemo(() => {
        if (!info?.items || info.items.length === 0) {
            return ''
        }

        // Calculate a "weight" to estimate the vertical space needed.
        const totalWeight = info.items.reduce((weight, _item) => {
            // Arrays used to take more space, but now they are single-line.
            // We can treat every item as a single line for simplicity.
            return weight + 1
        }, 0)

        // Decide how many grid slots the card should span based on the weight.
        if (totalWeight > LINE_WEIGHT_THRESHOLD * 2) {
            return 'sm:col-span-2 lg:col-span-3'
        }
        if (totalWeight > LINE_WEIGHT_THRESHOLD) {
            return 'sm:col-span-2 lg:col-span-2'
        }
        return 'lg:col-span-1'
    }, [info])

    return (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {info && info.items.length > 0 && (
                <div className={infoCardColSpanClass}>
                    <InfoCard title={info.title} items={info.items} />
                </div>
            )}

            {metricOrder.map((metricId) => {
                const metricDetails = metrics[metricId] as CoverageDetail | undefined
                const label = camelCaseToTitleCase(metricId)
                return <MetricCard key={metricId} label={label} details={metricDetails} />
            })}
        </div>
    )
}
