import MetricCard from '@/components/MetricCard'
import type { MetricKey, Metrics } from '@/types/summary'

const METRICS_TO_DISPLAY: { id: MetricKey; label: string }[] = [
    { id: 'lineCoverage', label: 'Line Coverage' },
    { id: 'branchCoverage', label: 'Branch Coverage' },
    { id: 'methodCoverage', label: 'Method Coverage' },
    { id: 'statementCoverage', label: 'Statement Coverage' },
    { id: 'functionCoverage', label: 'Function Coverage' },
]

export default function SummaryMetrics({ averages }: { averages?: Metrics }) {
    return (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {METRICS_TO_DISPLAY.map((m) => {
                const percentageValue = averages?.[m.id as keyof Metrics]?.percentage

                return <MetricCard key={m.id} label={m.label} value={percentageValue} />
            })}
        </div>
    )
}
