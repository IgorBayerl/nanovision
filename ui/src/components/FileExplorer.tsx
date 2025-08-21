import { List, ListTree, Pin, PinOff } from 'lucide-react'
import micromatch from 'micromatch'
import { useCallback, useMemo, useRef, useState } from 'react'
import HeaderRangeSlider from '@/components/HeaderRangeSlider'
import ColumnsMenu from '@/components/Toolbar.ColumnsMenu'
import RiskSegment from '@/components/Toolbar.RiskSegment'
import SearchBox from '@/components/Toolbar.SearchBox'
import { TreeRow } from '@/components/Tree.Row'
import { SUB_METRIC_COLS } from '@/lib/consts'
import { flattenFiles } from '@/lib/metrics'
import { useKeyboardSearch } from '@/lib/useKeyboardSearch'
import { useUrlState } from '@/lib/useUrlState'
import { camelCaseToTitleCase, cn } from '@/lib/utils'
import type {
    FileNode,
    FilterRange,
    MetricConfig,
    MetricKey,
    Metrics,
    RiskFilter,
    RiskLevel,
    SortDir,
    SortKey,
} from '@/types/summary'
import { Button } from '@/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/ui/tooltip'

type RenderNode = FileNode & { depth: number }

function getShortLabel(metricId: string): string {
    const knownPrefixes = ['line', 'branch', 'method', 'statement', 'function']
    const knownMatch = knownPrefixes.find((p) => metricId.toLowerCase().startsWith(p))
    if (knownMatch) {
        return knownMatch.charAt(0).toUpperCase() + knownMatch.slice(1)
    }
    return metricId.length > 4 ? `${metricId.slice(0, 3)}.` : metricId
}

const getDefaultEnabledMetrics = (metrics: string[]) => metrics.slice(0, 3)

export default function FileExplorer({ tree, availableMetrics }: { tree: FileNode[]; availableMetrics: string[] }) {
    const [viewMode, setViewMode] = useUrlState<'tree' | 'flat'>('view', 'tree')
    const [query, setQuery] = useUrlState('q', '')
    const [searchMode, setSearchMode] = useUrlState<'glob' | 'normal'>('qMode', 'normal')
    const [riskFilter, setRiskFilter] = useUrlState<RiskFilter>('risk', 'all')
    const [isNameColumnPinned, setIsNameColumnPinned] = useUrlState('pinned', true)

    const [sortKey, setSortKey] = useUrlState<SortKey>('sortKey', 'name')
    const [sortDir, setSortDir] = useUrlState<SortDir>('sortDir', 'asc')

    const defaultEnabled = getDefaultEnabledMetrics(availableMetrics)
    const [enabledMetricsParam, setEnabledMetricsParam] = useUrlState('cols', defaultEnabled.join(','))

    const [filterRanges, setFilterRanges] = useUrlState<Record<MetricKey, FilterRange>>('ranges', {})

    const [expandedFolders, setExpandedFolders] = useState<Set<string>>(
        () => new Set(tree.filter((n) => n.type === 'folder').map((n) => n.id)),
    )
    const searchRef = useRef<HTMLInputElement>(null)
    useKeyboardSearch(searchRef)

    // Memoize a map of the tree for efficient node lookup.
    const idToNodeMap = useMemo(() => {
        const map = new Map<string, FileNode>()
        const walk = (nodes: FileNode[]) => {
            for (const node of nodes) {
                map.set(node.id, node)
                if (node.children) {
                    walk(node.children)
                }
            }
        }
        walk(tree)
        return map
    }, [tree])

    const metricConfigs: MetricConfig[] = useMemo(
        () =>
            availableMetrics.map((id) => ({
                id,
                label: camelCaseToTitleCase(id),
                shortLabel: getShortLabel(id),
                enabled: enabledMetricsParam.split(',').includes(id),
            })),
        [availableMetrics, enabledMetricsParam],
    )
    const enabledMetrics = metricConfigs.filter((m) => m.enabled)

    // Helper function to get all descendant folder IDs for recursive toggling.
    const getDescendantFolderIds = (startNode: FileNode): string[] => {
        const ids: string[] = []
        const walk = (node: FileNode) => {
            if (node.type === 'folder' && node.children) {
                ids.push(node.id)
                for (const child of node.children) {
                    walk(child)
                }
            }
        }
        if (startNode.children) {
            for (const child of startNode.children) {
                walk(child)
            }
        }
        return ids
    }

    const toggleFolder = (id: string, event: React.MouseEvent | React.KeyboardEvent) => {
        const isRecursive = event.altKey === true

        setExpandedFolders((prev) => {
            const newSet = new Set(prev)
            const shouldExpand = !newSet.has(id)

            if (isRecursive) {
                const startNode = idToNodeMap.get(id)
                if (startNode) {
                    const descendantIds = getDescendantFolderIds(startNode)
                    const allIdsToToggle = [id, ...descendantIds]

                    if (shouldExpand) {
                        allIdsToToggle.forEach((folderId) => newSet.add(folderId))
                    } else {
                        allIdsToToggle.forEach((folderId) => newSet.delete(folderId))
                    }
                }
            } else {
                // Standard non-recursive toggle
                if (shouldExpand) {
                    newSet.add(id)
                } else {
                    newSet.delete(id)
                }
            }
            return newSet
        })
    }

    const toggleMetric = (id: MetricKey) => {
        const current = enabledMetricsParam.split(',').filter(Boolean)
        const newSet = new Set(current)
        if (newSet.has(id)) {
            newSet.delete(id)
        } else {
            newSet.add(id)
        }
        const newEnabled = availableMetrics.filter((mId) => newSet.has(mId))
        setEnabledMetricsParam(newEnabled.join(','))
    }

    const updateFilterRange = (id: MetricKey, vals: [number, number]) => {
        const newRanges = { ...filterRanges }
        if (vals[0] === 0 && vals[1] === 100) {
            delete newRanges[id]
        } else {
            newRanges[id] = { min: vals[0], max: vals[1] }
        }
        setFilterRanges(newRanges)
    }

    const nameMatches = useCallback(
        (text: string) => {
            if (!query) return true
            const trimmedQuery = query.trim()
            if (searchMode === 'glob') {
                return micromatch.isMatch(text, trimmedQuery, { nocase: true })
            }
            return text.toLowerCase().includes(trimmedQuery.toLowerCase())
        },
        [query, searchMode],
    )

    const fileRiskMatches = useCallback(
        (file: FileNode) => {
            if (riskFilter === 'all') return true
            const st = file.statuses
            if (!st) return true
            const statuses: RiskLevel[] = enabledMetrics.map((cfg) => st[cfg.id]).filter(Boolean) as RiskLevel[]
            if (statuses.length === 0) return true
            const hasDanger = statuses.includes('danger')
            const hasWarning = statuses.includes('warning')
            const allSafe = statuses.every((s) => s === 'safe')
            if (riskFilter === 'danger') return hasDanger
            if (riskFilter === 'warning') return hasWarning && !hasDanger
            if (riskFilter === 'safe') return allSafe
            return true
        },
        [riskFilter, enabledMetrics],
    )

    const filterTreeForView = useCallback(
        (nodes: FileNode[]): FileNode[] => {
            const res: FileNode[] = []
            for (const node of nodes) {
                if (node.type === 'file') {
                    const passesRanges = enabledMetrics.every((cfg) => {
                        const percentage = node.metrics?.[cfg.id]?.percentage
                        const range = filterRanges[cfg.id] ?? { min: 0, max: 100 }
                        if (percentage === undefined) return true
                        return percentage >= range.min && percentage <= range.max
                    })
                    const passesName = nameMatches(node.name) || nameMatches(node.path)
                    const passesRisk = fileRiskMatches(node)
                    if (passesRanges && passesName && passesRisk) res.push(node)
                } else if (node.type === 'folder' && node.children) {
                    const keptChildren = filterTreeForView(node.children)
                    if (keptChildren.length > 0) res.push({ ...node, children: keptChildren })
                }
            }
            return res
        },
        [enabledMetrics, filterRanges, nameMatches, fileRiskMatches],
    )

    const filteredView = useMemo(() => filterTreeForView(tree), [tree, filterTreeForView])

    const sortNodes = useCallback(
        (nodes: FileNode[]): FileNode[] => {
            return [...nodes].sort((a, b) => {
                const dir = sortDir === 'asc' ? 1 : -1
                if (sortKey === 'name') {
                    if (viewMode === 'tree' && a.type !== b.type) return a.type === 'folder' ? -1 : 1
                    const nameA = viewMode === 'flat' ? a.path : a.name
                    const nameB = viewMode === 'flat' ? b.path : b.name
                    return nameA.localeCompare(nameB) * dir
                }
                if (typeof sortKey === 'object') {
                    const { metric, subMetric } = sortKey
                    const getValue = (n: FileNode) => n.metrics?.[metric]?.[subMetric] ?? -1
                    return (getValue(a) - getValue(b)) * dir
                }
                return 0
            })
        },
        [sortDir, sortKey, viewMode],
    )

    const finalView: RenderNode[] = useMemo(() => {
        if (viewMode === 'flat') {
            return sortNodes(flattenFiles(filteredView)).map((n) => ({ ...n, depth: 0 }))
        }
        const result: RenderNode[] = []
        const build = (nodes: FileNode[], depth: number) => {
            sortNodes(nodes).forEach((node) => {
                result.push({ ...node, depth })
                if (node.type === 'folder' && node.children && expandedFolders.has(node.id)) {
                    build(node.children, depth + 1)
                }
            })
        }
        build(filteredView, 0)
        return result
    }, [filteredView, sortNodes, viewMode, expandedFolders])

    const onHeaderClick = (key: SortKey) => {
        if (JSON.stringify(sortKey) === JSON.stringify(key)) {
            setSortDir(sortDir === 'asc' ? 'desc' : 'asc')
        } else {
            setSortKey(key)
            setSortDir('asc')
        }
    }

    const metricsForNode = (node: FileNode): Partial<Metrics> | undefined => {
        return node.metrics
    }

    const totalMetricsWidth = enabledMetrics.reduce((sum) => sum + SUB_METRIC_COLS.reduce((s, c) => s + c.width, 0), 0)
    const totalTableWidth = `calc(max(100%, 300px + ${totalMetricsWidth}px))`

    return (
        <Card className="rounded-md">
            <CardHeader>
                <div className="flex flex-col gap-3">
                    <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                        <CardTitle className="text-lg">File Coverage</CardTitle>
                        <div className="flex flex-wrap items-center gap-2">
                            <TooltipProvider delayDuration={100}>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="h-8 w-8 rounded-sm p-0"
                                            onClick={() => setViewMode(viewMode === 'tree' ? 'flat' : 'tree')}
                                            aria-label={`Switch to ${viewMode === 'tree' ? 'flat' : 'tree'} view`}
                                        >
                                            {viewMode === 'tree' ? (
                                                <ListTree className="h-4 w-4" />
                                            ) : (
                                                <List className="h-4 w-4" />
                                            )}
                                        </Button>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>Switch to {viewMode === 'tree' ? 'Flat' : 'Tree'} View</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                            <ColumnsMenu metricConfigs={metricConfigs} onToggleMetric={toggleMetric} />
                        </div>
                    </div>
                    <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                        <SearchBox
                            ref={searchRef}
                            value={query}
                            onChange={setQuery}
                            mode={searchMode}
                            onModeChange={setSearchMode}
                        />
                        <RiskSegment value={riskFilter} onChange={setRiskFilter} />
                    </div>
                </div>
            </CardHeader>

            <CardContent className="p-0">
                <div className="w-full overflow-x-auto">
                    <div style={{ width: totalTableWidth }}>
                        {/* HEADER */}
                        <div className="sticky top-0 z-20 grid bg-background font-semibold text-xs">
                            <div
                                className="grid"
                                style={{
                                    gridTemplateColumns: `minmax(300px, 1fr) ${totalMetricsWidth}px`,
                                }}
                            >
                                <div
                                    className={cn(
                                        'border-border border-r border-b py-2',
                                        isNameColumnPinned && 'sticky left-0 z-10 bg-background',
                                    )}
                                >
                                    <div className="flex items-center justify-between px-2">
                                        <button
                                            type="button"
                                            className="flex items-center gap-1 text-left text-muted-foreground hover:text-foreground"
                                            onClick={() => onHeaderClick('name')}
                                        >
                                            File / Folder
                                            <span className="inline-flex h-3 w-3 items-center justify-center">
                                                {sortKey === 'name' ? (sortDir === 'asc' ? '▲' : '▼') : null}
                                            </span>
                                        </button>
                                        <TooltipProvider delayDuration={100}>
                                            <Tooltip>
                                                <TooltipTrigger asChild>
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                        className="h-6 w-6 p-0 text-muted-foreground hover:text-foreground"
                                                        onClick={() => setIsNameColumnPinned(!isNameColumnPinned)}
                                                        aria-label={isNameColumnPinned ? 'Unpin column' : 'Pin column'}
                                                    >
                                                        {isNameColumnPinned ? (
                                                            <PinOff className="h-4 w-4" />
                                                        ) : (
                                                            <Pin className="h-4 w-4" />
                                                        )}
                                                    </Button>
                                                </TooltipTrigger>
                                                <TooltipContent>
                                                    <p>{isNameColumnPinned ? 'Unpin column' : 'Pin column'}</p>
                                                </TooltipContent>
                                            </Tooltip>
                                        </TooltipProvider>
                                    </div>
                                </div>
                                <div
                                    className="grid"
                                    style={{ gridTemplateColumns: `repeat(${enabledMetrics.length}, 1fr)` }}
                                >
                                    {enabledMetrics.map((m, index) => (
                                        <div
                                            key={m.id}
                                            className={cn(
                                                'flex flex-col gap-2 border-border border-b py-2 text-left',
                                                index < enabledMetrics.length - 1 && 'border-r',
                                            )}
                                        >
                                            <div className="px-2 font-bold text-foreground">{m.shortLabel}</div>

                                            <div className="px-2">
                                                <HeaderRangeSlider
                                                    range={filterRanges[m.id] ?? { min: 0, max: 100 }}
                                                    onRangeCommit={(vals) => updateFilterRange(m.id, vals)}
                                                />
                                            </div>

                                            <div
                                                className="mt-1 grid"
                                                style={{
                                                    gridTemplateColumns: SUB_METRIC_COLS.map(
                                                        (c) => `${c.width}px`,
                                                    ).join(' '),
                                                }}
                                            >
                                                {SUB_METRIC_COLS.map((sub) => (
                                                    <button
                                                        key={sub.id}
                                                        type="button"
                                                        className={cn(
                                                            'flex w-full items-center justify-end gap-1 px-2 text-left text-muted-foreground hover:text-foreground',
                                                        )}
                                                        onClick={() =>
                                                            onHeaderClick({ metric: m.id, subMetric: sub.id })
                                                        }
                                                        title={`Sort by ${m.label} ${sub.label}`}
                                                    >
                                                        <span className="inline-flex h-3 w-3 items-center justify-center">
                                                            {typeof sortKey === 'object' &&
                                                            sortKey.metric === m.id &&
                                                            sortKey.subMetric === sub.id
                                                                ? sortDir === 'asc'
                                                                    ? '▲'
                                                                    : '▼'
                                                                : null}
                                                        </span>
                                                        {sub.label}
                                                    </button>
                                                ))}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>

                        {/* BODY */}
                        <div className="w-full">
                            {finalView.length > 0 ? (
                                finalView.map((node, index) => (
                                    <TreeRow
                                        key={node.id}
                                        node={node}
                                        depth={node.depth}
                                        enabledMetrics={enabledMetrics}
                                        isExpanded={expandedFolders.has(node.id)}
                                        onToggleFolder={toggleFolder}
                                        metricsForNode={metricsForNode}
                                        viewMode={viewMode}
                                        index={index}
                                        isPinned={isNameColumnPinned}
                                    />
                                ))
                            ) : (
                                <div className="py-16 text-center text-muted-foreground">
                                    No files match the current filters/search
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
