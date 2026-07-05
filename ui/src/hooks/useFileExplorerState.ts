import { useEffect, useMemo, useRef, useState } from 'react'
import { useKeyboardSearch } from '@/hooks/useKeyboardSearch'
import { useUrlState } from '@/hooks/useUrlState'
import type { DiffFilter, FileNode, FilterRange, MetricKey, RiskFilter, SortDir, SortKey } from '@/types/summary'

const getDefaultEnabledMetrics = (metrics: string[]) => metrics.slice(0, 3)
const EXPANDED_FOLDERS_STORAGE_KEY = 'nanovision-expanded-folders'

// Collect all folder IDs from the flat node list.
const getAllFolderIds = (nodes: FileNode[]): string[] => nodes.filter((n) => n.type === 'folder').map((n) => n.id)

export function useFileExplorerState(nodes: FileNode[], availableMetrics: string[]) {
    const [viewMode, setViewMode] = useUrlState<'tree' | 'flat'>('view', 'tree')
    const [query, setQuery] = useUrlState('q', '')
    const [searchMode, setSearchMode] = useUrlState<'glob' | 'normal'>('qMode', 'normal')
    const [riskFilter, setRiskFilter] = useUrlState<RiskFilter>('risk', 'all')
    const [diffFilter, setDiffFilter] = useUrlState<DiffFilter>('diff', 'all')
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
            return item ? JSON.parse(item) : getAllFolderIds(nodes)
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

    // Map of parent id -> child folder ids, built from the flat node list. Used to
    // resolve descendants for recursive (alt-click) folder expansion.
    const childFolderIdsByParent = useMemo(() => {
        const map = new Map<string, string[]>()
        for (const node of nodes) {
            if (node.type !== 'folder') continue
            const parentId = node.parentId ?? ''
            const siblings = map.get(parentId)
            if (siblings) siblings.push(node.id)
            else map.set(parentId, [node.id])
        }
        return map
    }, [nodes])

    const getDescendantFolderIds = (startId: string): string[] => {
        const ids: string[] = []
        const stack = [...(childFolderIdsByParent.get(startId) ?? [])]
        while (stack.length > 0) {
            const id = stack.pop() as string
            ids.push(id)
            const children = childFolderIdsByParent.get(id)
            if (children) stack.push(...children)
        }
        return ids
    }

    const toggleFolder = (id: string, event: React.MouseEvent | React.KeyboardEvent) => {
        const isRecursive = event.altKey === true

        const newSet = new Set(expandedFolders)
        const shouldExpand = !newSet.has(id)

        if (isRecursive) {
            const allIdsToToggle = [id, ...getDescendantFolderIds(id)]
            for (const folderId of allIdsToToggle) {
                if (shouldExpand) newSet.add(folderId)
                else newSet.delete(folderId)
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
            diffFilter,
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
            setDiffFilter,
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
