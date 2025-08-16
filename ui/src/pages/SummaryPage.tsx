import micromatch from 'micromatch'
import { useCallback, useMemo, useRef, useState } from 'react'
import MetricCard from '@/components/MetricCard'
import { ThemeSwitch } from '@/components/Theme.Switch'
import ColumnsMenu from '@/components/Toolbar.ColumnsMenu'
import DensitySelect from '@/components/Toolbar.DensitySelect'
import MetricRangeChip from '@/components/Toolbar.MetricRangeChip'
import RiskSegment from '@/components/Toolbar.RiskSegment'
import SearchBox from '@/components/Toolbar.SearchBox'
import FileTree from '@/components/Tree'
import { averageMetrics, calculateFolderMetrics, flattenFiles } from '@/lib/metrics'
import { buildIdMap, type SortDir, type SortKey } from '@/lib/tree'
import { useKeyboardSearch } from '@/lib/useKeyboardSearch'
import { cn } from '@/lib/utils'
import type { CSSWithVars } from '@/types/css'
import type {
    FileNode,
    FilterRange,
    MetricConfig,
    MetricKey,
    Metrics,
    RiskFilter,
    RiskLevel,
    SummaryV1,
} from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

export default function SummaryPage({ data }: { data: SummaryV1 }) {
    // theme is handled globally; this page just renders
    const tree = data.tree
    const idMap = useMemo(() => buildIdMap(tree), [tree])

    // UI state
    const [expandedFolders, setExpandedFolders] = useState<Set<string>>(
        () => new Set(tree.filter((n) => n.type === 'folder').map((n) => n.id)),
    )
    const [rowDensity, setRowDensity] = useState<'comfortable' | 'compact'>('comfortable')

    const [searchMode, setSearchMode] = useState<'glob' | 'normal'>('normal')

    const [query, setQuery] = useState('')
    const searchRef = useRef<HTMLInputElement>(null) as React.RefObject<HTMLInputElement>
    useKeyboardSearch(searchRef)

    const [riskFilter, setRiskFilter] = useState<RiskFilter>('all')

    const [metricConfigs, setMetricConfigs] = useState<MetricConfig[]>([
        { id: 'lineCoverage', label: 'Line Coverage', shortLabel: 'Line', enabled: true },
        { id: 'branchCoverage', label: 'Branch Coverage', shortLabel: 'Branch', enabled: true },
        { id: 'methodCoverage', label: 'Method Coverage', shortLabel: 'Method', enabled: true },
        {
            id: 'statementCoverage',
            label: 'Statement Coverage',
            shortLabel: 'Stmt',
            enabled: false,
        },
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
                        const v = node.metrics?.[cfg.id as keyof Metrics] as number
                        const r = filterRanges[cfg.id]
                        return v >= r.min && v <= r.max
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

    const sortSiblings = useCallback(
        (nodes: FileNode[]): FileNode[] => {
            const sorted = [...nodes].sort((a, b) => {
                if (a.type !== b.type) return a.type === 'folder' ? -1 : 1
                const dir = sortDir === 'asc' ? 1 : -1
                if (sortKey === 'name') return a.name.localeCompare(b.name) * dir

                const getValue = (n: FileNode) => {
                    if (n.type === 'file') return (n.metrics?.[sortKey as keyof Metrics] as number | undefined) ?? -1
                    const original = idMap.get(n.id)
                    const m = original?.children ? calculateFolderMetrics(original.children) : undefined
                    return m ? (m[sortKey as keyof Metrics] as number) : -1
                }
                const av = getValue(a)
                const bv = getValue(b)
                return (av - bv) * dir
            })
            return sorted.map((n) =>
                n.type === 'folder' && n.children ? { ...n, children: sortSiblings(n.children) } : n,
            )
        },
        [sortDir, sortKey, idMap],
    )

    const sortedView = useMemo(() => sortSiblings(filteredView), [filteredView, sortSiblings])

    const globalAverages = useMemo(() => averageMetrics(flattenFiles(tree)), [tree])
    const onHeaderClick = (key: SortKey) => {
        if (sortKey === key) setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'))
        else {
            setSortKey(key)
            setSortDir('asc')
        }
    }

    const filesInFolderCount = (node: FileNode) =>
        (idMap.get(node.id)?.children || []).filter((c) => c.type === 'file').length

    const metricsForNode = (node: FileNode): Partial<Metrics> | undefined =>
        node.type === 'folder' ? calculateFolderMetrics(idMap.get(node.id)?.children ?? []) : node.metrics

    const riskForMetric = (_node: FileNode, key: MetricKey): RiskLevel => {
        const st = _node.statuses
        return st?.[key] ?? 'safe'
    }

    const headerLeftPadPx = 14 + 32 + (rowDensity === 'compact' ? 8 : 12)

    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="font-bold text-2xl tracking-tight">{data.title || 'Coverage Report'}</h1>
                    <p className="text-muted-foreground text-sm">Backend-driven risk, square UI, column-tied filters</p>
                </div>
                <ThemeSwitch />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {(
                    [
                        { id: 'lineCoverage', label: 'Line Coverage' },
                        { id: 'branchCoverage', label: 'Branch Coverage' },
                        { id: 'methodCoverage', label: 'Method Coverage' },
                        { id: 'statementCoverage', label: 'Statement Coverage' },
                        { id: 'functionCoverage', label: 'Function Coverage' },
                    ] as { id: MetricKey; label: string }[]
                ).map((m) => (
                    <MetricCard
                        key={m.id}
                        label={m.label}
                        value={globalAverages?.[m.id as keyof Metrics] as number | undefined}
                    />
                ))}
            </div>

            <Card className="rounded-md">
                <CardHeader>
                    <div className="flex flex-col gap-3">
                        <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                            <CardTitle className="text-lg">File Coverage Tree</CardTitle>
                            <div className="flex flex-wrap items-center gap-2">
                                <ColumnsMenu metricConfigs={metricConfigs} onToggleMetric={toggleMetric} />
                                <DensitySelect value={rowDensity} onValueChange={setRowDensity} />
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
                    <div className="overflow-x-auto">
                        <table
                            className="w-full"
                            style={
                                {
                                    '--metric-count': String(enabledMetrics.length),
                                    '--metric-col-width': '170px',
                                } as CSSWithVars
                            }
                        >
                            <thead className="sticky top-0 z-10 border-b bg-background/80 backdrop-blur-sm">
                                <tr
                                    className="grid items-start px-2 pt-2 font-semibold text-sm"
                                    style={{
                                        gridTemplateColumns: '1fr repeat(var(--metric-count), var(--metric-col-width))',
                                    }}
                                >
                                    <th className="pb-2 text-left">
                                        <button
                                            type="button"
                                            className="flex items-center gap-1 text-left text-muted-foreground hover:text-foreground"
                                            onClick={() => onHeaderClick('name')}
                                            aria-label="Sort by name"
                                        >
                                            File / Folder
                                        </button>
                                    </th>
                                    {enabledMetrics.map((m) => (
                                        <th
                                            key={m.id}
                                            className="flex flex-col gap-1 pb-2 text-left"
                                            style={{ paddingLeft: headerLeftPadPx }}
                                        >
                                            <button
                                                type="button"
                                                className="flex w-full items-center gap-1 text-left text-muted-foreground hover:text-foreground"
                                                onClick={() => onHeaderClick(m.id)}
                                                aria-label={`Sort by ${m.label}`}
                                                title={`Sort by ${m.label}`}
                                            >
                                                <span className="text-left">{m.shortLabel}</span>
                                                <span className="inline-flex h-3 w-3 items-center justify-center">
                                                    {sortKey === m.id ? (sortDir === 'asc' ? '▲' : '▼') : null}
                                                </span>
                                            </button>
                                            <MetricRangeChip
                                                cfg={m}
                                                value={filterRanges[m.id]}
                                                onChange={(vals) => updateFilterRange(m.id, vals)}
                                            />
                                        </th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody>
                                {sortedView.length > 0 ? (
                                    <FileTree
                                        nodes={sortedView}
                                        expanded={expandedFolders}
                                        toggleFolder={toggleFolder}
                                        enabledMetrics={enabledMetrics}
                                        rowDensity={rowDensity}
                                        metricsForNode={metricsForNode}
                                        riskForMetric={riskForMetric}
                                        filesInFolderCount={filesInFolderCount}
                                    />
                                ) : (
                                    <tr>
                                        <td
                                            colSpan={enabledMetrics.length + 1}
                                            className="py-16 text-center text-muted-foreground"
                                        >
                                            No files match the current filters/search
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
