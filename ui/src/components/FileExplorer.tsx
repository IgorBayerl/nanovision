import { useMemo } from 'react'
import FileExplorerBody from '@/components/FileExplorer.Body'
import FileExplorerHeader from '@/components/FileExplorer.Header'
import FileExplorerToolbar from '@/components/FileExplorer.Toolbar'
import { useFileExplorerState } from '@/hooks/useFileExplorerState'
import { useFilteredAndSortedTree } from '@/hooks/useFilteredAndSortedTree'
import { aggregateFolderDiff } from '@/lib/aggregateFolderDiff'
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
    nodes: FileNode[]
    availableMetrics: string[]
    metricDefinitions: MetricDefinitions
}

export default function FileExplorer({ nodes, availableMetrics, metricDefinitions }: FileExplorerProps) {
    const { state, setters, searchRef } = useFileExplorerState(nodes, availableMetrics)

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

    // Diff status propagated up to folders (a folder is decorated when any
    // descendant file was added/modified).
    const folderDiffMap = useMemo(() => aggregateFolderDiff(nodes), [nodes])

    // Data-derived upper bound for value (non-percentage) metrics, so their range
    // sliders span the actual observed values instead of a fixed 0–100.
    const metricMaxes = useMemo(() => {
        const maxes: Record<string, number> = {}
        for (const cfg of metricConfigs) {
            if (cfg.definition?.kind !== 'value') continue
            let max = 0
            for (const node of nodes) {
                const metric = node.metrics?.[cfg.id]
                const val = metric && 'value' in metric ? metric.value : undefined
                if (typeof val === 'number' && val > max) max = val
            }
            maxes[cfg.id] = max
        }
        return maxes
    }, [metricConfigs, nodes])

    const finalView = useFilteredAndSortedTree({
        nodes,
        query: state.query,
        searchMode: state.searchMode,
        riskFilter: state.riskFilter,
        diffFilter: state.diffFilter,
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
    const totalTableWidth = `calc(max(99.9%, 450px + ${totalMetricsWidth}px))`

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
                            metricMaxes={metricMaxes}
                        />
                        <FileExplorerBody
                            nodes={finalView}
                            enabledMetrics={enabledMetrics}
                            expandedFolders={state.expandedFolders}
                            onToggleFolder={setters.toggleFolder}
                            viewMode={state.viewMode}
                            isPinned={state.isNameColumnPinned}
                            folderDiffMap={folderDiffMap}
                        />
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
