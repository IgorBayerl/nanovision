import { useMemo } from 'react'
import FileExplorer from '@/components/FileExplorer'
import Footer from '@/components/Layout.Footer'
import TopBar from '@/components/Layout.TopBar'
import SummaryMetrics from '@/components/SummaryMetrics'
import ValidationAlerts from '@/components/ValidationAlerts'
import { cn } from '@/lib/utils'
import type { SummaryV1 } from '@/lib/validation' // Corrected: This import now works
import { validateSummaryData } from '@/lib/validation'
import type { MetadataItem } from '@/types/summary'

const NON_METRIC_KEYS = new Set(['files', 'folders', 'statuses'])

const isMetricKey = (key: string): boolean => {
    return !NON_METRIC_KEYS.has(key)
}

export default function SummaryPage({ data: rawData }: { data: unknown }) {
    const validationResult = useMemo(() => validateSummaryData(rawData), [rawData])

    const { reportInfo, metricKeys, validatedData } = useMemo(() => {
        if (!validationResult.success) {
            const partialData = rawData as Partial<SummaryV1>
            return {
                validatedData: null,
                reportInfo: undefined,
                metricKeys: partialData.totals ? Object.keys(partialData.totals).filter(isMetricKey) : [],
            }
        }

        const data = validationResult.data
        let reportInfo: { title: string; items: MetadataItem[] } | undefined

        if (data.metadata) {
            const validItems = data.metadata.filter(
                (item) => item.value !== undefined && (!Array.isArray(item.value) || item.value.length > 0),
            )
            if (validItems.length > 0) {
                reportInfo = {
                    title: 'Report Information',
                    items: validItems,
                }
            }
        }

        const keys = Object.keys(data.totals).filter(isMetricKey)

        return {
            validatedData: data,
            reportInfo,
            metricKeys: keys,
        }
    }, [validationResult, rawData])

    const title = validatedData?.title ?? (rawData as Partial<SummaryV1>)?.title ?? 'Coverage Report'

    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <TopBar title={title} />

            {!validationResult.success && <ValidationAlerts issues={validationResult.error.issues} />}

            {validatedData ? (
                <>
                    <SummaryMetrics
                        info={reportInfo}
                        metrics={validatedData.totals}
                        metricOrder={metricKeys}
                        metricDefinitions={validatedData.metricDefinitions}
                    />
                    <FileExplorer
                        tree={validatedData.tree}
                        availableMetrics={metricKeys}
                        metricDefinitions={validatedData.metricDefinitions}
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
