import { useMemo } from 'react'
import FileExplorerBody from '@/components/FileExplorer.Body'
import FileExplorerHeader from '@/components/FileExplorer.Header'
import FileExplorerToolbar from '@/components/FileExplorer.Toolbar'
import { CHANGE_COLUMN } from '@/components/Tree.Row'
import { useFileExplorerState } from '@/hooks/useFileExplorerState'
import { useFilteredAndSortedTree } from '@/hooks/useFilteredAndSortedTree'
import { aggregateFolderDiff } from '@/lib/aggregateFolderDiff'
import { camelCaseToTitleCase } from '@/lib/utils'
import type { ConfigFile } from '@/lib/validation'
import type { FileNode, MetricDefinition, MetricDefinitions } from '@/types/summary'
import { Card, CardContent, CardHeader } from '@/ui/card'

function getShortLabel(metricId: string): string {
    const knownPrefixes = ['line', 'method', 'statement', 'function']
    const knownMatch = knownPrefixes.find((p) => metricId.toLowerCase().startsWith(p))
    if (knownMatch) return knownMatch.charAt(0).toUpperCase() + knownMatch.slice(1)
    return metricId.length > 4 ? `${metricId.slice(0, 3)}.` : metricId
}

/** Adds a last column to a metric: how far its percentage moved against the base run. */
function withChangeColumn(definition: MetricDefinition): MetricDefinition {
    const at = definition.subMetrics.findIndex((sub) => sub.id === 'percentage')
    if (at < 0) return definition
    const subMetrics = [...definition.subMetrics]
    subMetrics.splice(at + 1, 0, { id: CHANGE_COLUMN, label: 'Change', width: 76 })
    return { ...definition, subMetrics }
}

interface FileExplorerProps {
    nodes: FileNode[]
    availableMetrics: string[]
    metricDefinitions: MetricDefinitions
    /** The config files of the run; a folder with its own settings is marked. */
    configs?: ConfigFile[]
    /** Path -> metric -> change against the base run, in percentage points. */
    deltas?: Record<string, Record<string, number>>
    /** The metrics the base run was compared on; each gets a Change column. */
    deltaMetrics?: string[]
}

export default function FileExplorer({
    nodes: reportNodes,
    availableMetrics,
    metricDefinitions,
    configs,
    deltas,
    deltaMetrics,
}: FileExplorerProps) {
    // The change is put next to the numbers it belongs to, so a row reads and
    // sorts like any other column.
    const nodes = useMemo(() => {
        if (!deltas) return reportNodes
        return reportNodes.map((node) => {
            const moved = deltas[node.path]
            if (!moved || !node.metrics) return node
            const metrics = { ...node.metrics }
            for (const [id, delta] of Object.entries(moved)) {
                const metric = metrics[id]
                if (metric && 'percentage' in metric) metrics[id] = { ...metric, delta }
            }
            return { ...node, metrics }
        })
    }, [reportNodes, deltas])

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
                    definition: definition && deltaMetrics?.includes(id) ? withChangeColumn(definition) : definition,
                }
            }),
        [availableMetrics, metricDefinitions, state.enabledMetrics, deltaMetrics],
    )

    const enabledMetrics = useMemo(() => metricConfigs.filter((m) => m.enabled), [metricConfigs])

    // Diff status propagated up to folders (a folder is decorated when any
    // descendant file was added/modified).
    const folderDiffMap = useMemo(() => aggregateFolderDiff(nodes), [nodes])

    // The folder config that applies to a path: the nearest folder with its own settings.
    const configFor = useMemo(() => {
        const folders = (configs ?? []).filter((c) => c.path !== '').sort((a, b) => b.path.length - a.path.length)
        return (path: string) => folders.find((c) => path === c.path || path.startsWith(`${c.path}/`))
    }, [configs])

    const finalView = useFilteredAndSortedTree({
        nodes,
        query: state.query,
        searchMode: state.searchMode,
        riskFilter: state.riskFilter,
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
                            totalMetricsWidth={totalMetricsWidth}
                        />
                        <FileExplorerBody
                            nodes={finalView}
                            enabledMetrics={enabledMetrics}
                            expandedFolders={state.expandedFolders}
                            onToggleFolder={setters.toggleFolder}
                            viewMode={state.viewMode}
                            isPinned={state.isNameColumnPinned}
                            folderDiffMap={folderDiffMap}
                            configFor={configFor}
                        />
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
