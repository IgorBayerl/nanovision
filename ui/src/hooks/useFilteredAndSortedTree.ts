import micromatch from 'micromatch'
import { useMemo } from 'react'
import { useDebounce } from '@/hooks/useDebounce'
import type {
    DiffFilter,
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
    nodes: FileNode[]
    query: string
    searchMode: 'glob' | 'normal'
    riskFilter: RiskFilter
    diffFilter: DiffFilter
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

/** Sentinel key used for top-level nodes (those with no parentId). */
const ROOT_KEY = ''

const EMPTY_CHILDREN: ReadonlyMap<string, FileNode[]> = new Map()

type Comparator = (a: FileNode, b: FileNode) => number

/**
 * Derives lookup structures from the flat node list in a single linear pass:
 *  - allFiles:     just the file nodes (what every filter runs against)
 *  - parentMap:    node id -> parent id (used to walk ancestors)
 *  - childrenMap:  parent id -> child nodes, preserving the backend's pre-order
 *                  (folders before files, each group sorted by name)
 */
function usePrecomputedNodes(nodes: FileNode[]) {
    return useMemo(() => {
        const allFiles: FileNode[] = []
        const parentMap = new Map<string, string | null>()
        const childrenMap = new Map<string, FileNode[]>()

        for (const node of nodes) {
            const parentId = node.parentId ?? null
            parentMap.set(node.id, parentId)

            const key = parentId ?? ROOT_KEY
            const siblings = childrenMap.get(key)
            if (siblings) siblings.push(node)
            else childrenMap.set(key, [node])

            if (node.type === 'file') allFiles.push(node)
        }

        return { allFiles, parentMap, childrenMap }
    }, [nodes])
}

/**
 * A hook that filters, sorts, and structures the flat node list for display.
 */
export function useFilteredAndSortedTree({
    nodes,
    query,
    searchMode,
    riskFilter,
    diffFilter,
    filterRanges,
    sortKey,
    sortDir,
    viewMode,
    expandedFolders,
    enabledMetrics,
}: HookParams): RenderNode[] {
    const { allFiles, parentMap, childrenMap } = usePrecomputedNodes(nodes)

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

            if (diffFilter === 'changed') {
                if (file.diffStatus !== 'added' && file.diffStatus !== 'modified') {
                    return false
                }
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
    }, [allFiles, query, searchMode, debouncedFilterRanges, riskFilter, diffFilter, enabledMetrics])

    // A single comparator shared by both flat and tree views.
    const comparator = useMemo<Comparator>(() => {
        const dir = sortDir === 'asc' ? 1 : -1

        const sortByName: Comparator = (a, b) => {
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

        return (a, b) => {
            if (sortKey === 'name') return sortByName(a, b)
            if (typeof sortKey === 'object') return sortByMetric(a, b, sortKey)
            return 0
        }
    }, [sortKey, sortDir, viewMode])

    // Flat view: a single sorted list of the filtered files.
    const flatNodes = useMemo<RenderNode[]>(() => {
        if (viewMode !== 'flat') return []
        return [...filteredFiles].sort(comparator).map((node) => ({ ...node, depth: 0 }))
    }, [viewMode, filteredFiles, comparator])

    // Tree view: pre-sort each sibling group once. This is independent of expand/collapse
    // state, so toggling folders never re-sorts. Total cost is O(n log n), not O(n²).
    const sortedChildrenMap = useMemo(() => {
        if (viewMode !== 'tree') return EMPTY_CHILDREN
        const sorted = new Map<string, FileNode[]>()
        for (const [key, siblings] of childrenMap) {
            sorted.set(key, [...siblings].sort(comparator))
        }
        return sorted
    }, [viewMode, childrenMap, comparator])

    return useMemo(() => {
        if (viewMode === 'flat') return flatNodes

        // Mark every filtered file and its ancestor folders as visible.
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

        if (visibleNodeIds.size === 0 && (query.trim().length > 0 || diffFilter === 'changed')) return []

        // Depth-first walk over the pre-sorted sibling groups, honoring expand state.
        const result: RenderNode[] = []
        const walk = (parentKey: string, depth: number) => {
            const siblings = sortedChildrenMap.get(parentKey)
            if (!siblings) return
            for (const node of siblings) {
                if (!visibleNodeIds.has(node.id)) continue
                result.push({ ...node, depth })
                if (node.type === 'folder' && expandedFolders.has(node.id)) {
                    walk(node.id, depth + 1)
                }
            }
        }
        walk(ROOT_KEY, 0)
        return result
    }, [viewMode, flatNodes, filteredFiles, sortedChildrenMap, parentMap, expandedFolders, query, diffFilter])
}
