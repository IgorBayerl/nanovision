import { useMemo, useRef, useState } from 'react'
import { useKeyboardSearch } from '@/lib/useKeyboardSearch'
import { useUrlState } from '@/lib/useUrlState'
import type { FileNode, FilterRange, MetricKey, RiskFilter, SortDir, SortKey } from '@/types/summary'

const getDefaultEnabledMetrics = (metrics: string[]) => metrics.slice(0, 3)

export function useFileExplorerState(tree: FileNode[], availableMetrics: string[]) {
    const [viewMode, setViewMode] = useUrlState<'tree' | 'flat'>('view', 'tree')
    const [query, setQuery] = useUrlState('q', '')
    const [searchMode, setSearchMode] = useUrlState<'glob' | 'normal'>('qMode', 'normal')
    const [riskFilter, setRiskFilter] = useUrlState<RiskFilter>('risk', 'all')
    const [isNameColumnPinned, setIsNameColumnPinned] = useUrlState('pinned', true)
    const [sortKey, setSortKey] = useUrlState<SortKey>('sortKey', 'name')
    const [sortDir, setSortDir] = useUrlState<SortDir>('sortDir', 'asc')

    const defaultEnabled = getDefaultEnabledMetrics(availableMetrics)
    const [enabledMetricsParam, setEnabledMetricsParam] = useUrlState('cols', defaultEnabled.join(','))
    const enabledMetrics = useMemo(() => enabledMetricsParam.split(','), [enabledMetricsParam])

    const [filterRanges, setFilterRanges] = useUrlState<Record<MetricKey, FilterRange>>('ranges', {})
    const [expandedFolders, setExpandedFolders] = useState<Set<string>>(
        () => new Set(tree.filter((n) => n.type === 'folder').map((n) => n.id)),
    )
    const searchRef = useRef<HTMLInputElement>(null)
    useKeyboardSearch(searchRef)

    const idToNodeMap = useMemo(() => {
        const map = new Map<string, FileNode>()
        const walk = (nodes: FileNode[]) => {
            for (const node of nodes) {
                map.set(node.id, node)
                if (node.children) walk(node.children)
            }
        }
        walk(tree)
        return map
    }, [tree])

    const getDescendantFolderIds = (startNode: FileNode): string[] => {
        const ids: string[] = []
        const walk = (node: FileNode) => {
            if (node.type === 'folder' && node.children) {
                ids.push(node.id)
                for (const child of node.children) walk(child)
            }
        }
        if (startNode.children) {
            for (const child of startNode.children) walk(child)
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
                        allIdsToToggle.forEach((folderId) => {
                            newSet.add(folderId)
                        })
                    } else {
                        allIdsToToggle.forEach((folderId) => {
                            newSet.delete(folderId)
                        })
                    }
                }
            } else {
                if (shouldExpand) newSet.add(id)
                else newSet.delete(id)
            }
            return newSet
        })
    }

    const toggleMetric = (id: MetricKey) => {
        const current = enabledMetricsParam.split(',').filter(Boolean)
        const newSet = new Set(current)
        if (newSet.has(id)) newSet.delete(id)
        else newSet.add(id)
        const newEnabled = availableMetrics.filter((mId) => newSet.has(mId))
        setEnabledMetricsParam(newEnabled.join(','))
    }

    const updateFilterRange = (id: MetricKey, vals: [number, number]) => {
        const newRanges = { ...filterRanges }
        if (vals[0] === 0 && vals[1] === 100) delete newRanges[id]
        else newRanges[id] = { min: vals[0], max: vals[1] }
        setFilterRanges(newRanges)
    }

    const handleHeaderClick = (key: SortKey) => {
        if (JSON.stringify(sortKey) === JSON.stringify(key)) {
            setSortDir(sortDir === 'asc' ? 'desc' : 'asc')
        } else {
            setSortKey(key)
            setSortDir('asc')
        }
    }

    return {
        state: {
            viewMode,
            query,
            searchMode,
            riskFilter,
            isNameColumnPinned,
            sortKey,
            sortDir,
            enabledMetrics,
            filterRanges,
            expandedFolders,
        },
        setters: {
            setViewMode,
            setQuery,
            setSearchMode,
            setRiskFilter,
            setIsNameColumnPinned,
            setSortKey,
            setSortDir,
            toggleMetric,
            updateFilterRange,
            toggleFolder,
            handleHeaderClick,
        },
        searchRef,
    }
}

export type FileExplorerState = ReturnType<typeof useFileExplorerState>['state']
export type Setters = ReturnType<typeof useFileExplorerState>['setters']
