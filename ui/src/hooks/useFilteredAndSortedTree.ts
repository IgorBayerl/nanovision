import micromatch from 'micromatch'
import { useCallback, useMemo } from 'react'
import { flattenFiles } from '@/lib/metrics'
import type {
    FileNode,
    FilterRange,
    MetricConfig,
    MetricKey,
    RiskFilter,
    RiskLevel,
    SortDir,
    SortKey,
} from '@/types/summary'

type RenderNode = FileNode & { depth: number }

interface HookParams {
    tree: FileNode[]
    query: string
    searchMode: 'glob' | 'normal'
    riskFilter: RiskFilter
    filterRanges: Record<MetricKey, FilterRange>
    sortKey: SortKey
    sortDir: SortDir
    viewMode: 'tree' | 'flat'
    expandedFolders: Set<string>
    enabledMetrics: MetricConfig[]
}

export function useFilteredAndSortedTree({
    tree,
    query,
    searchMode,
    riskFilter,
    filterRanges,
    sortKey,
    sortDir,
    viewMode,
    expandedFolders,
    enabledMetrics,
}: HookParams): RenderNode[] {
    const nameMatches = useCallback(
        (text: string) => {
            if (!query) return true
            const trimmedQuery = query.trim()
            if (searchMode === 'glob') return micromatch.isMatch(text, trimmedQuery, { nocase: true })
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
                    const getValue = (n: FileNode) => {
                        const m = n.metrics?.[metric]
                        if (!m) return -1
                        return m[subMetric as keyof typeof m] ?? -1
                    }
                    return (getValue(a) - getValue(b)) * dir
                }
                return 0
            })
        },
        [sortDir, sortKey, viewMode],
    )

    return useMemo(() => {
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
}
