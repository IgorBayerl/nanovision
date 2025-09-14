import micromatch from 'micromatch'
import { useMemo } from 'react'
import { useDebounce } from '@/hooks/useDebounce'
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

/**
 * The delay in milliseconds for debouncing real-time filters.
 * A value of ~150 provides a responsive feel without overwhelming the CPU on large reports.
 */
const DEBOUNCE_DELAY_MS = 1

/**
 * Pre-computes data structures from the raw tree to optimize filtering and lookups.
 */
function usePrecomputedTree(tree: FileNode[]) {
    return useMemo(() => {
        const allFiles: FileNode[] = []
        const parentMap = new Map<string, string | null>()

        function walk(nodes: FileNode[], parentId: string | null) {
            for (const node of nodes) {
                parentMap.set(node.id, parentId)
                if (node.type === 'file') {
                    allFiles.push(node)
                }
                if (node.children) {
                    walk(node.children, node.id)
                }
            }
        }

        walk(tree, null)
        return { allFiles, parentMap }
    }, [tree])
}

/**
 * A hook that filters, sorts, and structures file tree data for display.
 */
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
    const { allFiles, parentMap } = usePrecomputedTree(tree)

    // Debounce the filterRanges prop. The expensive filtering logic will only re-run
    // after the user has stopped dragging the slider for the specified delay.
    const debouncedFilterRanges = useDebounce(filterRanges, DEBOUNCE_DELAY_MS)

    const filteredFiles = useMemo(() => {
        const trimmedQuery = query.trim().toLowerCase()
        const hasQuery = trimmedQuery.length > 0
        const activeRanges = Object.entries(debouncedFilterRanges)

        return allFiles.filter((file) => {
            if (hasQuery) {
                const textToMatch = file.path.toLowerCase()
                const nameMatches =
                    searchMode === 'glob'
                        ? micromatch.isMatch(textToMatch, trimmedQuery, { nocase: true })
                        : textToMatch.includes(trimmedQuery)
                if (!nameMatches) return false
            }

            if (activeRanges.length > 0) {
                const rangeMatches = activeRanges.every(([metricId, range]) => {
                    const percentage = file.metrics?.[metricId]?.percentage
                    if (percentage === undefined) return true
                    return percentage >= range.min && percentage <= range.max
                })
                if (!rangeMatches) return false
            }

            if (riskFilter !== 'all') {
                const visibleStatuses = enabledMetrics
                    .map((cfg) => file.statuses?.[cfg.id])
                    .filter((status): status is RiskLevel => !!status)

                if (visibleStatuses.length === 0) {
                    if (riskFilter !== 'safe') return false
                } else {
                    const hasDanger = visibleStatuses.includes('danger')
                    const hasWarning = visibleStatuses.includes('warning')

                    switch (riskFilter) {
                        case 'danger':
                            if (!hasDanger) return false
                            break
                        case 'warning':
                            if (!hasWarning) return false
                            break
                        case 'safe':
                            if (hasDanger || hasWarning) return false
                            break
                    }
                }
            }

            return true
        })
    }, [allFiles, query, searchMode, debouncedFilterRanges, riskFilter, enabledMetrics])

    const sortedNodes = useMemo(() => {
        const nodesToSort = viewMode === 'flat' ? filteredFiles : tree
        const dir = sortDir === 'asc' ? 1 : -1

        const sortByName = (a: FileNode, b: FileNode) => {
            if (viewMode === 'tree' && a.type !== b.type) return a.type === 'folder' ? -1 : 1
            const nameA = viewMode === 'flat' ? a.path : a.name
            const nameB = viewMode === 'flat' ? b.path : b.name
            return nameA.localeCompare(nameB) * dir
        }

        const sortByMetric = (a: FileNode, b: FileNode, key: { metric: MetricKey; subMetric: string }) => {
            const valA = a.metrics?.[key.metric]?.[key.subMetric as keyof (typeof a.metrics)[string]] ?? -1
            const valB = b.metrics?.[key.metric]?.[key.subMetric as keyof (typeof b.metrics)[string]] ?? -1
            if (valA === valB) return sortByName(a, b)
            return (valA - valB) * dir
        }

        return [...nodesToSort].sort((a, b) => {
            if (sortKey === 'name') return sortByName(a, b)
            if (typeof sortKey === 'object') return sortByMetric(a, b, sortKey)
            return 0
        })
    }, [filteredFiles, tree, viewMode, sortKey, sortDir])

    return useMemo(() => {
        if (viewMode === 'flat') {
            return sortedNodes.map((node) => ({ ...node, depth: 0 }))
        }

        const visibleNodeIds = new Set<string>()
        for (const file of filteredFiles) {
            visibleNodeIds.add(file.id)
            let currentParentId = parentMap.get(file.id)
            while (currentParentId) {
                if (visibleNodeIds.has(currentParentId)) break
                visibleNodeIds.add(currentParentId)
                currentParentId = parentMap.get(currentParentId)
            }
        }

        if (visibleNodeIds.size === 0 && query.trim().length > 0) return []

        const result: RenderNode[] = []
        function buildRenderList(nodes: FileNode[], depth: number) {
            for (const node of nodes) {
                if (!visibleNodeIds.has(node.id)) continue
                result.push({ ...node, depth })
                if (node.type === 'folder' && node.children && expandedFolders.has(node.id)) {
                    const sortedChildren = [...node.children].sort((a, b) => {
                        const aIndex = sortedNodes.findIndex((n) => n.id === a.id)
                        const bIndex = sortedNodes.findIndex((n) => n.id === b.id)
                        return aIndex - bIndex
                    })
                    buildRenderList(sortedChildren, depth + 1)
                }
            }
        }

        buildRenderList(sortedNodes, 0)
        return result
    }, [viewMode, filteredFiles, sortedNodes, parentMap, expandedFolders, query])
}
