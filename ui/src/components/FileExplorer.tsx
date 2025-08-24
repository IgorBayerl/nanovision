import { useMemo } from 'react'
import FileExplorerBody from '@/components/FileExplorer.Body'
import FileExplorerHeader from '@/components/FileExplorer.Header'
import FileExplorerToolbar from '@/components/FileExplorer.Toolbar'
import { useFileExplorerState } from '@/hooks/useFileExplorerState'
import { useFilteredAndSortedTree } from '@/hooks/useFilteredAndSortedTree'
import { camelCaseToTitleCase } from '@/lib/utils'
import type { FileNode, MetricDefinitions } from '@/types/summary'
import { Card, CardContent, CardHeader } from '@/ui/card'

function getShortLabel(metricId: string): string {
    const knownPrefixes = ['line', 'branch', 'method', 'statement', 'function']
    const knownMatch = knownPrefixes.find((p) => metricId.toLowerCase().startsWith(p))
    if (knownMatch) return knownMatch.charAt(0).toUpperCase() + knownMatch.slice(1)
    return metricId.length > 4 ? `${metricId.slice(0, 3)}.` : metricId
}

interface FileExplorerProps {
    tree: FileNode[]
    availableMetrics: string[]
    metricDefinitions: MetricDefinitions
}

export default function FileExplorer({ tree, availableMetrics, metricDefinitions }: FileExplorerProps) {
    const { state, setters, searchRef } = useFileExplorerState(tree, availableMetrics)

    const metricConfigs = useMemo(
        () =>
            availableMetrics.map((id) => {
                const definition = metricDefinitions[id]
                return {
                    id,
                    label: definition?.label ?? camelCaseToTitleCase(id),
                    shortLabel: definition?.shortLabel ?? getShortLabel(id),
                    enabled: state.enabledMetrics.includes(id),
                    definition: definition,
                }
            }),
        [availableMetrics, metricDefinitions, state.enabledMetrics],
    )

    const enabledMetrics = useMemo(() => metricConfigs.filter((m) => m.enabled), [metricConfigs])

    const finalView = useFilteredAndSortedTree({
        tree,
        query: state.query,
        searchMode: state.searchMode,
        riskFilter: state.riskFilter,
        filterRanges: state.filterRanges,
        sortKey: state.sortKey,
        sortDir: state.sortDir,
        viewMode: state.viewMode,
        expandedFolders: state.expandedFolders,
        enabledMetrics: enabledMetrics,
    })

    const totalMetricsWidth = enabledMetrics.reduce(
        (sum, metric) => sum + (metric.definition?.subMetrics.reduce((s, c) => s + c.width, 0) ?? 0),
        0,
    )
    const totalTableWidth = `calc(max(100%, 300px + ${totalMetricsWidth}px))`

    return (
        <Card className="rounded-md">
            <CardHeader>
                <FileExplorerToolbar
                    state={state}
                    setters={setters}
                    searchRef={searchRef}
                    metricConfigs={metricConfigs}
                />
            </CardHeader>

            <CardContent className="p-0">
                <div className="w-full overflow-x-auto">
                    <div style={{ width: totalTableWidth }}>
                        <FileExplorerHeader
                            isNameColumnPinned={state.isNameColumnPinned}
                            onPinColumn={setters.setIsNameColumnPinned}
                            enabledMetrics={enabledMetrics}
                            sortKey={state.sortKey}
                            sortDir={state.sortDir}
                            onHeaderClick={setters.handleHeaderClick}
                            filterRanges={state.filterRanges}
                            onRangeUpdate={setters.updateFilterRange}
                            totalMetricsWidth={totalMetricsWidth}
                        />
                        <FileExplorerBody
                            nodes={finalView}
                            enabledMetrics={enabledMetrics}
                            expandedFolders={state.expandedFolders}
                            onToggleFolder={setters.toggleFolder}
                            viewMode={state.viewMode}
                            isPinned={state.isNameColumnPinned}
                        />
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
