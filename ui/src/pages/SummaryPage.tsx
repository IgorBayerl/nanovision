import { useMemo } from 'react'
import ChangesTab from '@/components/ChangesTab'
import FileExplorer from '@/components/FileExplorer'
import Layout, { type LayoutProps } from '@/components/Layout'
import SummaryMetrics from '@/components/SummaryMetrics'
import Tabs from '@/components/Tabs'
import ValidationAlerts from '@/components/ValidationAlerts'
import { useReportSelection, withReportSelection } from '@/hooks/useReportSelection'
import { useUrlState } from '@/hooks/useUrlState'
import { applyReportSelection, unfilterableMetrics } from '@/lib/reportSelection'
import { useRememberSummaryQuery } from '@/lib/summaryQuery'
import { orderedMetricKeys } from '@/lib/utils'
import type { SummaryV1 } from '@/lib/validation'
import { validateSummaryData } from '@/lib/validation'
import type { MetadataItem } from '@/types/summary'
import { SidebarContent } from '@/ui/sidebar'

interface SummaryPageProps {
    data: unknown
    layout?: Pick<LayoutProps, 'nav' | 'actions'>
}

export default function SummaryPage({ data: rawData, layout }: SummaryPageProps) {
    const validationResult = useMemo(() => validateSummaryData(rawData), [rawData])
    useRememberSummaryQuery()

    const { reportInfo, metricKeys, validatedData } = useMemo(() => {
        if (!validationResult.success) {
            const partialData = rawData as Partial<SummaryV1>
            return {
                validatedData: null,
                reportInfo: undefined,
                metricKeys: orderedMetricKeys(partialData.totals, partialData.metricOrder),
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

        const keys = orderedMetricKeys(data.totals, data.metricOrder)

        return {
            validatedData: data,
            reportInfo,
            metricKeys: keys,
        }
    }, [validationResult, rawData])

    const reportSelection = useReportSelection(validatedData?.reports)

    // With some reports unticked the tree is rebuilt for the active ones: files
    // from their own index, folders and totals summed from what is below them.
    // With all ticked the report's own numbers already are the answer.
    const { nodes, totals } = useMemo(() => {
        if (!validatedData) return { nodes: [], totals: undefined }
        if (reportSelection.isAllSelected) return { nodes: validatedData.nodes, totals: validatedData.totals }

        const filtered = applyReportSelection(
            validatedData.nodes,
            validatedData.totals,
            validatedData.reportIndexes,
            validatedData.statusBands,
            reportSelection.selection,
        )

        return {
            // a changed-only report lists a part of the files, so its totals
            // cannot be summed from them and keep every report
            totals: validatedData.onlyChanged ? validatedData.totals : filtered.totals,
            nodes: filtered.nodes.map((node) => ({
                ...node,
                targetUrl: withReportSelection(node.targetUrl, reportSelection.linkQuery),
            })),
        }
    }, [validatedData, reportSelection.isAllSelected, reportSelection.selection, reportSelection.linkQuery])

    // Named on the selector so a frozen number is not mistaken for a filtered one.
    const frozenMetricLabels = useMemo(() => {
        if (!validatedData) return []
        const frozen = validatedData.onlyChanged
            ? metricKeys
            : unfilterableMetrics(metricKeys, validatedData.reportIndexes)
        return frozen.map((id) => validatedData.metricDefinitions[id]?.label ?? id)
    }, [validatedData, metricKeys])

    // A run with a diff or a base run gets the Changes tab next to the files.
    const hasChanges = !!validatedData && (!!validatedData.review || !!validatedData.comparison)
    const [tab, setTab] = useUrlState<'changes' | 'files'>('tab', 'files')
    const activeTab = hasChanges && tab === 'changes' ? 'changes' : 'files'
    const changedFiles = useMemo(
        () =>
            nodes.filter((n) => n.type === 'file' && (n.diffStatus === 'added' || n.diffStatus === 'modified')).length,
        [nodes],
    )

    // The delta counts every report, so it is shown only while all are ticked.
    const deltas = useMemo(() => {
        const comparison = validatedData?.comparison
        if (!comparison || !reportSelection.isAllSelected) return undefined
        return Object.fromEntries(comparison.metrics.map((m) => [m.key, m.delta]))
    }, [validatedData, reportSelection.isAllSelected])

    const title = validatedData?.title ?? (rawData as Partial<SummaryV1>)?.title ?? 'Coverage Report'

    const leftSidebar =
        validatedData && totals ? (
            <SidebarContent className="gap-0 p-0">
                <SummaryMetrics
                    info={reportInfo}
                    comparing={validatedData.comparing}
                    metrics={totals}
                    metricOrder={metricKeys}
                    metricDefinitions={validatedData.metricDefinitions}
                    statusBands={validatedData.statusBands}
                    reportSelection={reportSelection}
                    frozenMetricLabels={frozenMetricLabels}
                    deltas={deltas}
                />
            </SidebarContent>
        ) : undefined

    return (
        <Layout title={title} leftSidebar={leftSidebar} {...layout}>
            {!validationResult.success && <ValidationAlerts issues={validationResult.error.issues} />}
            {validatedData ? (
                <>
                    {hasChanges && (
                        <Tabs
                            tabs={[
                                { id: 'files', label: 'Files' },
                                { id: 'changes', label: 'Changes', count: changedFiles || undefined },
                            ]}
                            active={activeTab}
                            onChange={setTab}
                        />
                    )}
                    {activeTab === 'changes' ? (
                        <ChangesTab comparison={validatedData.comparison} review={validatedData.review} nodes={nodes} />
                    ) : (
                        <FileExplorer
                            nodes={nodes}
                            availableMetrics={metricKeys}
                            metricDefinitions={validatedData.metricDefinitions}
                            // like the side panel: the delta counts every report
                            deltas={reportSelection.isAllSelected ? validatedData.comparison?.deltas : undefined}
                            deltaMetrics={
                                reportSelection.isAllSelected
                                    ? validatedData.comparison?.metrics.map((m) => m.key)
                                    : undefined
                            }
                        />
                    )}
                </>
            ) : (
                <div className="rounded-md border border-border bg-card p-10 text-center text-muted-foreground">
                    Could not render the report due to critical data errors. Please review the alerts above.
                </div>
            )}
        </Layout>
    )
}
