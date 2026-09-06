import { useMemo } from 'react'
import FunctionNav from '@/components/FunctionNav'
import Layout from '@/components/Layout'
import MethodsTable from '@/components/MethodsTable'
import ReportsSelector from '@/components/ReportsSelector'
import SourceCodeViewer from '@/components/SourceCodeViewer'
import SummaryMetrics from '@/components/SummaryMetrics'
import ValidationAlerts from '@/components/ValidationAlerts'
import { useReportSelection } from '@/hooks/useReportSelection'
import { applyReportSelectionToTotals, unfilterableMetrics } from '@/lib/reportSelection'
import type { DetailsV1 } from '@/lib/validation'
import { validateDetailsData } from '@/lib/validation'
import type { MetadataItem } from '@/types/summary'
import { SidebarContent, SidebarHeader } from '@/ui/sidebar'

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

    const reportSelection = useReportSelection(validatedData?.reports)

    const totals = useMemo(() => {
        if (!validatedData) return undefined
        return applyReportSelectionToTotals(
            validatedData.totals,
            validatedData.reportIndex,
            validatedData.statusBands,
            reportSelection.selection,
        )
    }, [validatedData, reportSelection.selection])

    const frozenMetricLabels = useMemo(() => {
        if (!validatedData?.reportIndex) return []
        return unfilterableMetrics(metricKeys, { file: validatedData.reportIndex }).map(
            (id) => validatedData.metricDefinitions[id]?.label ?? id,
        )
    }, [validatedData, metricKeys])

    const title = validatedData?.title ?? (rawData as Partial<DetailsV1>)?.title ?? 'Coverage Details'

    const leftSidebar =
        validatedData && totals ? (
            <>
                <SidebarHeader>
                    <div className="font-semibold text-sm">Overview</div>
                </SidebarHeader>
                <SidebarContent>
                    <SummaryMetrics
                        info={reportInfo}
                        metrics={totals}
                        metricOrder={metricKeys}
                        metricDefinitions={validatedData.metricDefinitions}
                        infoFooter={<ReportsSelector state={reportSelection} frozenMetricLabels={frozenMetricLabels} />}
                        variant="sidebar"
                    />
                </SidebarContent>
            </>
        ) : undefined

    const rightSidebar =
        validatedData?.methods && validatedData.methods.length > 0 ? (
            <FunctionNav methods={validatedData.methods} />
        ) : undefined

    return (
        <Layout title={title} showBackButton leftSidebar={leftSidebar} rightSidebar={rightSidebar}>
            {!validationResult.success && <ValidationAlerts issues={validationResult.error.issues} />}

            {validatedData ? (
                <>
                    {validatedData.methods && (
                        <MethodsTable
                            methods={validatedData.methods}
                            metricDefinitions={validatedData.metricDefinitions}
                        />
                    )}
                    <SourceCodeViewer
                        fileName={validatedData.fileName}
                        lines={validatedData.lines}
                        selection={reportSelection.selection}
                        reports={validatedData.reports}
                    />
                </>
            ) : (
                <div className="rounded-md border border-border bg-card p-10 text-center text-muted-foreground">
                    Could not render the report due to critical data errors. Please review the alerts above.
                </div>
            )}
        </Layout>
    )
}
