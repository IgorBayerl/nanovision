import Footer from '@/components/Layout.Footer'
import TopBar from '@/components/Layout.TopBar'
import MethodsTable from '@/components/MethodsTable'
import ReportsSelector from '@/components/ReportsSelector'
import SourceCodeViewer from '@/components/SourceCodeViewer'
import SummaryMetrics from '@/components/SummaryMetrics'
import ValidationAlerts from '@/components/ValidationAlerts'
import { cn } from '@/lib/utils'
import type { DetailsV1 } from '@/lib/validation'
import { validateDetailsData } from '@/lib/validation'
import type { MetadataItem } from '@/types/summary'
import { useMemo, useState } from 'react'

const NON_METRIC_KEYS = new Set(['files', 'folders', 'statuses'])
const isMetricKey = (key: string): boolean => !NON_METRIC_KEYS.has(key)

export default function DetailsPage({ data: rawData }: { data: unknown }) {
    const validationResult = useMemo(() => validateDetailsData(rawData), [rawData])

    const { validatedData, metricKeys, reportInfo } = useMemo(() => {
        if (!validationResult.success) {
            const partialData = rawData as Partial<DetailsV1>
            return {
                validatedData: null,
                metricKeys: partialData.totals ? Object.keys(partialData.totals).filter(isMetricKey) : [],
                reportInfo: undefined,
            }
        }
        const data = validationResult.data
        const keys = Object.keys(data.totals).filter(isMetricKey)

        let reportInfo: { title: string; items: MetadataItem[] } | undefined
        if (data.metadata && data.metadata.length > 0) {
            reportInfo = {
                title: 'File Information',
                items: data.metadata,
            }
        }

        return { validatedData: data, metricKeys: keys, reportInfo }
    }, [validationResult, rawData])

    const [activeReportIndices, setActiveReportIndices] = useState<Set<number>>(
        () => new Set(validatedData?.reports?.map((_, i) => i) ?? []),
    )

    const handleToggleReport = (index: number) => {
        setActiveReportIndices((prev) => {
            const newSet = new Set(prev)
            if (newSet.has(index)) {
                newSet.delete(index)
            } else {
                newSet.add(index)
            }
            return newSet
        })
    }

    const title = validatedData?.title ?? (rawData as Partial<DetailsV1>)?.title ?? 'Coverage Details'

    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <TopBar title={title} showBackButton />

            {!validationResult.success && <ValidationAlerts issues={validationResult.error.issues} />}

            {validatedData ? (
                <>
                    <SummaryMetrics
                        info={reportInfo} // <-- PASS INFO
                        metrics={validatedData.totals}
                        metricOrder={metricKeys}
                        metricDefinitions={validatedData.metricDefinitions}
                    />
                    {validatedData.reports && validatedData.reports.length > 0 && (
                        <ReportsSelector
                            reports={validatedData.reports}
                            activeReportIndices={activeReportIndices}
                            onToggleReport={handleToggleReport}
                        />
                    )}
                    {validatedData.methods && (
                        <MethodsTable
                            methods={validatedData.methods}
                            metricDefinitions={validatedData.metricDefinitions}
                        />
                    )}
                    <SourceCodeViewer
                        fileName={validatedData.fileName}
                        lines={validatedData.lines}
                        activeReportIndices={activeReportIndices}
                    />
                </>
            ) : (
                <div className="rounded-md border border-border bg-card p-10 text-center text-muted-foreground">
                    Could not render the report due to critical data errors. Please review the alerts above.
                </div>
            )}

            <Footer />
        </div>
    )
}
