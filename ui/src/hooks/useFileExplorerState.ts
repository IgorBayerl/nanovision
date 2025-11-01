import { useEffect, useMemo, useRef, useState } from 'react'
import { useKeyboardSearch } from '@/hooks/useKeyboardSearch'
import { useUrlState } from '@/hooks/useUrlState'
import type { FileNode, FilterRange, MetricKey, RiskFilter, SortDir, SortKey } from '@/types/summary'

const getDefaultEnabledMetrics = (metrics: string[]) => metrics.slice(0, 3)
const EXPANDED_FOLDERS_STORAGE_KEY = 'nanovision-expanded-folders'

// Helper function to recursively find all folder IDs in the tree.
const getAllFolderIds = (nodes: FileNode[]): string[] => {
    const ids: string[] = []
    const walk = (arr: FileNode[]) => {
        arr.forEach((n) => {
            if (n.type === 'folder') {
                ids.push(n.id)
            }
            if (n.children) {
                walk(n.children)
            }
        })
    }
    walk(nodes)
    return ids
}

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

    const [expandedFoldersArray, setExpandedFoldersArray] = useState<string[]>(() => {
        try {
            const item = window.localStorage.getItem(EXPANDED_FOLDERS_STORAGE_KEY)
            return item ? JSON.parse(item) : getAllFolderIds(tree)
        } catch (error) {
            console.error('Error reading from localStorage', error)
            return []
        }
    })

    useEffect(() => {
        try {
            window.localStorage.setItem(EXPANDED_FOLDERS_STORAGE_KEY, JSON.stringify(expandedFoldersArray))
        } catch (error) {
            console.error('Error writing to localStorage', error)
        }
    }, [expandedFoldersArray])

    const expandedFolders = useMemo(() => new Set(expandedFoldersArray), [expandedFoldersArray])

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

        const newSet = new Set(expandedFolders)
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
            if (shouldExpand) {
                newSet.add(id)
            } else {
                newSet.delete(id)
            }
        }
        setExpandedFoldersArray(Array.from(newSet))
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
