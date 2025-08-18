import { List, ListTree, Pin, PinOff } from 'lucide-react'
import micromatch from 'micromatch'
import { useCallback, useMemo, useRef, useState } from 'react'
import HeaderRangeSlider from '@/components/HeaderRangeSlider'
import ColumnsMenu from '@/components/Toolbar.ColumnsMenu'
import RiskSegment from '@/components/Toolbar.RiskSegment'
import SearchBox from '@/components/Toolbar.SearchBox'
import { TreeRow } from '@/components/Tree.Row'
import { SUB_METRIC_COLS } from '@/lib/consts'
import { calculateFolderMetrics, flattenFiles } from '@/lib/metrics'
import { useKeyboardSearch } from '@/lib/useKeyboardSearch'
import { cn } from '@/lib/utils'
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

export default function FileExplorer({ tree, idMap }: { tree: FileNode[]; idMap: Map<string, FileNode> }) {
    const [expandedFolders, setExpandedFolders] = useState<Set<string>>(
        () => new Set(tree.filter((n) => n.type === 'folder').map((n) => n.id)),
    )
    const [viewMode, setViewMode] = useState<'tree' | 'flat'>('tree')
    const [searchMode, setSearchMode] = useState<'glob' | 'normal'>('normal')
    const [query, setQuery] = useState('')
    const searchRef = useRef<HTMLInputElement>(null)
    useKeyboardSearch(searchRef)
    const [isNameColumnPinned, setIsNameColumnPinned] = useState(true)

    const [riskFilter, setRiskFilter] = useState<RiskFilter>('all')
    const [metricConfigs, setMetricConfigs] = useState<MetricConfig[]>([
        { id: 'lineCoverage', label: 'Line Coverage', shortLabel: 'Line', enabled: true },
        { id: 'branchCoverage', label: 'Branch Coverage', shortLabel: 'Branch', enabled: true },
        { id: 'methodCoverage', label: 'Method Coverage', shortLabel: 'Method', enabled: true },
        { id: 'statementCoverage', label: 'Statement Coverage', shortLabel: 'Stmt', enabled: false },
        { id: 'functionCoverage', label: 'Function Coverage', shortLabel: 'Func', enabled: false },
    ])
    const [filterRanges, setFilterRanges] = useState<Record<MetricKey, FilterRange>>({
        lineCoverage: { min: 0, max: 100 },
        branchCoverage: { min: 0, max: 100 },
        methodCoverage: { min: 0, max: 100 },
        statementCoverage: { min: 0, max: 100 },
        functionCoverage: { min: 0, max: 100 },
    })
    const enabledMetrics = metricConfigs.filter((m) => m.enabled)
    const [sortKey, setSortKey] = useState<SortKey>('name')
    const [sortDir, setSortDir] = useState<SortDir>('asc')

    const toggleFolder = (id: string) => {
        setExpandedFolders((prev) => {
            const s = new Set(prev)
            if (s.has(id)) s.delete(id)
            else s.add(id)
            return s
        })
    }

    const toggleMetric = (id: MetricKey) =>
        setMetricConfigs((prev) => prev.map((m) => (m.id === id ? { ...m, enabled: !m.enabled } : m)))

    const updateFilterRange = (id: MetricKey, vals: [number, number]) =>
        setFilterRanges((prev) => ({ ...prev, [id]: { min: vals[0], max: vals[1] } }))

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
                if (node.type === 'file' && node.metrics) {
                    const passesRanges = enabledMetrics.every((cfg) => {
                        const v = node.metrics?.[cfg.id]?.percentage
                        const r = filterRanges[cfg.id]
                        return v !== undefined && v >= r.min && v <= r.max
                    })
                    const passesName = nameMatches(node.name) || nameMatches(node.path)
                    const passesRisk = fileRiskMatches(node)
                    if (passesRanges && passesName && passesRisk) res.push(node)
                } else if (node.type === 'folder' && node.children) {
                    const kept = filterTreeForView(node.children)
                    if (kept.length > 0) res.push({ ...node, children: kept })
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
                    const getValue = (n: FileNode) => {
                        const originalNode = idMap.get(n.id)
                        const m =
                            n.type === 'folder' && originalNode?.children
                                ? calculateFolderMetrics(originalNode.children)
                                : n.metrics
                        return m?.[metric]?.[subMetric] ?? -1
                    }
                    return (getValue(a) - getValue(b)) * dir
                }
                return 0
            })
        },
        [sortDir, sortKey, idMap, viewMode],
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
            setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'))
        } else {
            setSortKey(key)
            setSortDir('asc')
        }
    }

    const metricsForNode = (node: FileNode): Partial<Metrics> | undefined => {
        const originalNode = idMap.get(node.id)
        return node.type === 'folder' && originalNode?.children
            ? calculateFolderMetrics(originalNode.children)
            : node.metrics
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
                                                        onClick={() => setIsNameColumnPinned((p) => !p)}
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
                                                    range={filterRanges[m.id]}
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
