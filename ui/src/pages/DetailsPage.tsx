import { useMemo } from 'react'
import Footer from '@/components/Layout.Footer'
import TopBar from '@/components/Layout.TopBar'
import MethodsTable from '@/components/MethodsTable'
import SourceCodeViewer from '@/components/SourceCodeViewer'
import SummaryMetrics from '@/components/SummaryMetrics'
import ValidationAlerts from '@/components/ValidationAlerts'
import { cn } from '@/lib/utils'
import type { DetailsV1 } from '@/lib/validation'
import { validateDetailsData } from '@/lib/validation'
import type { MetadataItem } from '@/types/summary'

const NON_METRIC_KEYS = new Set(['files', 'folders', 'statuses'])
const isMetricKey = (key: string): boolean => !NON_METRIC_KEYS.has(key)

export default function DetailsPage({ data: rawData }: { data: unknown }) {
    const validationResult = useMemo(() => validateDetailsData(rawData), [rawData])

    const { validatedData, metricKeys, reportInfo } = useMemo(() => {
        // <-- ADD reportInfo
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

        // Create the info for the InfoCard
        let reportInfo: { title: string; items: MetadataItem[] } | undefined
        if (data.metadata && data.metadata.length > 0) {
            reportInfo = {
                title: 'File Information',
                items: data.metadata,
            }
        }

        return { validatedData: data, metricKeys: keys, reportInfo }
    }, [validationResult, rawData])

    const title = validatedData?.title ?? (rawData as Partial<DetailsV1>)?.title ?? 'Coverage Details'

    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <TopBar title={title} />

            {!validationResult.success && <ValidationAlerts issues={validationResult.error.issues} />}

            {validatedData ? (
                <>
                    <SummaryMetrics
                        info={reportInfo} // <-- PASS INFO
                        metrics={validatedData.totals}
                        metricOrder={metricKeys}
                        metricDefinitions={validatedData.metricDefinitions}
                    />
                    {validatedData.methods && (
                        <MethodsTable
                            methods={validatedData.methods}
                            metricDefinitions={validatedData.metricDefinitions}
                        />
                    )}
                    <SourceCodeViewer fileName={validatedData.fileName} lines={validatedData.lines} />
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
