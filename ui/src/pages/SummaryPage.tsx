import { useMemo } from 'react'
import FileExplorer from '@/components/FileExplorer'
import Layout from '@/components/Layout'
import ProblemsPanel from '@/components/ProblemsPanel'
import ReportsSelector from '@/components/ReportsSelector'
import ReviewSummary from '@/components/ReviewSummary'
import SummaryMetrics from '@/components/SummaryMetrics'
import ValidationAlerts from '@/components/ValidationAlerts'
import { useReportSelection, withReportSelection } from '@/hooks/useReportSelection'
import { applyReportSelection, unfilterableMetrics } from '@/lib/reportSelection'
import type { SummaryV1 } from '@/lib/validation'
import { validateSummaryData } from '@/lib/validation'
import type { MetadataItem } from '@/types/summary'
import { SidebarContent, SidebarHeader } from '@/ui/sidebar'

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

    const reportSelection = useReportSelection(validatedData?.reports)

    // The whole tree is rebuilt for the active reports: files from their own
    // index, folders and totals summed from what is below them.
    const { nodes, totals } = useMemo(() => {
        if (!validatedData) return { nodes: [], totals: undefined }

        const filtered = applyReportSelection(
            validatedData.nodes,
            validatedData.totals,
            validatedData.reportIndexes,
            validatedData.statusBands,
            reportSelection.selection,
        )

        if (!reportSelection.linkQuery) return filtered

        return {
            ...filtered,
            nodes: filtered.nodes.map((node) => ({
                ...node,
                targetUrl: withReportSelection(node.targetUrl, reportSelection.linkQuery),
            })),
        }
    }, [validatedData, reportSelection.selection, reportSelection.linkQuery])

    // Named on the selector so a frozen number is not mistaken for a filtered one.
    const frozenMetricLabels = useMemo(() => {
        if (!validatedData) return []
        return unfilterableMetrics(metricKeys, validatedData.reportIndexes).map(
            (id) => validatedData.metricDefinitions[id]?.label ?? id,
        )
    }, [validatedData, metricKeys])

    const title = validatedData?.title ?? (rawData as Partial<SummaryV1>)?.title ?? 'Coverage Report'

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

    return (
        <Layout title={title} leftSidebar={leftSidebar}>
            {!validationResult.success && <ValidationAlerts issues={validationResult.error.issues} />}
            {validatedData ? (
                <>
                    {validatedData.review && <ReviewSummary review={validatedData.review} nodes={nodes} />}
                    {validatedData.diagnostics && validatedData.diagnostics.length > 0 && (
                        <ProblemsPanel diagnostics={validatedData.diagnostics} nodes={nodes} />
                    )}
                    <FileExplorer
                        nodes={nodes}
                        availableMetrics={metricKeys}
                        metricDefinitions={validatedData.metricDefinitions}
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
