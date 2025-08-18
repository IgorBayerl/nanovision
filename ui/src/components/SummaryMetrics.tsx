import InfoCard from '@/components/InfoCard'
import MetricCard from '@/components/MetricCard'
import type { SummaryV1 } from '@/types/summary'

const METRICS_TO_DISPLAY = [
    { id: 'lineCoverage', label: 'Line Coverage' },
    { id: 'branchCoverage', label: 'Branch Coverage' },
    { id: 'methodCoverage', label: 'Method Coverage' },
    { id: 'statementCoverage', label: 'Statement Coverage' },
    { id: 'functionCoverage', label: 'Function Coverage' },
] as const

export default function SummaryMetrics({ data }: { data: SummaryV1 }) {
    const { totals, generatedAt, parsers, configFiles, importedReports } = data

    return (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div className="sm:col-span-2 lg:col-span-1">
                <InfoCard
                    generatedAt={generatedAt}
                    files={totals.files}
                    folders={totals.folders}
                    parsers={parsers}
                    configFiles={configFiles}
                    importedReports={importedReports}
                />
            </div>
            {METRICS_TO_DISPLAY.map((m) => {
                const metricDetails = totals[m.id]
                return <MetricCard key={m.id} label={m.label} details={metricDetails} />
            })}
        </div>
    )
}
