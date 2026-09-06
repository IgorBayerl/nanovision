import { useCallback, useMemo } from 'react'
import { useUrlState } from '@/hooks/useUrlState'
import type { ReportSelection } from '@/lib/reportSelection'
import { indicesFromSelection, isReportSelected, selectionFromIndices, toggleReport } from '@/lib/reportSelection'
import type { ReportRef } from '@/lib/validation'

const URL_KEY = 'reports'

function parseSelection(param: string, reportCount: number): ReportSelection {
    if (param.trim() === '') return 0

    const indices = param
        .split(',')
        .map((part) => Number.parseInt(part, 10))
        .filter((index) => Number.isInteger(index) && index >= 0 && index < reportCount)

    return selectionFromIndices(indices)
}

export type ReportSelectionState = {
    reports: ReportRef[]
    selection: ReportSelection
    isSelected: (index: number) => boolean
    isAllSelected: boolean
    isNoneSelected: boolean
    /** Tri-state for the parent checkbox governing the list. */
    groupCheckState: boolean | 'indeterminate'
    toggle: (index: number) => void
    /** Ticks every report, or clears them when they are already all ticked. */
    toggleAll: () => void
    selectOnly: (index: number) => void
    /**
     * Query string that carries this selection to another page, empty when it
     * is the default. The multi-file report is several documents, so links have
     * to spell the selection out; in single-file mode the query survives the
     * hash change on its own.
     */
    linkQuery: string
}

/**
 * Report selection shared by every screen, held in the URL so it survives
 * navigation and can be linked to. Indices address the same reports the
 * coverage masks do.
 */
export function useReportSelection(reports: ReportRef[] | undefined): ReportSelectionState {
    const list = useMemo(() => reports ?? [], [reports])
    const allParam = useMemo(() => list.map((_, index) => index).join(','), [list])

    const [param, setParam] = useUrlState(URL_KEY, allParam)
    const selection = useMemo(() => parseSelection(param, list.length), [param, list.length])

    const commit = useCallback(
        (next: ReportSelection) => setParam(indicesFromSelection(next, list.length).join(',')),
        [setParam, list.length],
    )

    const isAllSelected = param === allParam
    const isNoneSelected = selection === 0

    return {
        reports: list,
        selection,
        isSelected: (index: number) => isReportSelected(selection, index),
        isAllSelected,
        isNoneSelected,
        groupCheckState: isAllSelected ? true : isNoneSelected ? false : 'indeterminate',
        toggle: (index: number) => commit(toggleReport(selection, index)),
        toggleAll: () => commit(isAllSelected ? 0 : parseSelection(allParam, list.length)),
        selectOnly: (index: number) => commit(selectionFromIndices([index])),
        linkQuery: isAllSelected ? '' : `${URL_KEY}=${param}`,
    }
}

/**
 * Adds the current selection to a generated details link. Single-file reports
 * navigate by hash inside one document, so those links are left alone.
 */
export function withReportSelection(targetUrl: string | null | undefined, linkQuery: string): string | undefined {
    if (!targetUrl) return undefined
    if (!linkQuery || targetUrl.startsWith('#')) return targetUrl

    return targetUrl.includes('?') ? `${targetUrl}&${linkQuery}` : `${targetUrl}?${linkQuery}`
}
