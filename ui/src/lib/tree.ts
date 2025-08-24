import type { FileNode, Metrics } from '@/types/summary'

export const buildIdMap = (nodes: FileNode[]) => {
    const map = new Map<string, FileNode>()
    const walk = (arr: FileNode[]) =>
        arr.forEach((n) => {
            map.set(n.id, n)
            if (n.children) walk(n.children)
        })
    walk(nodes)
    return map
}

export type SortKey = 'name' | keyof Metrics
export type SortDir = 'asc' | 'desc'

export const filterTreeForView = (
    nodes: FileNode[],
    enabledMetricIds: (keyof Metrics)[],
    ranges: Record<keyof Metrics, { min: number; max: number }>,
    namePredicate: (s: string) => boolean,
    riskPredicate: (f: FileNode) => boolean,
): FileNode[] => {
    const res: FileNode[] = []
    for (const node of nodes) {
        if (node.type === 'file' && node.metrics) {
            const passesRanges = enabledMetricIds.every((k) => {
                const v = node.metrics?.[k]
                const r = ranges[k]
                if (typeof v !== 'number') return false
                return v >= r.min && v <= r.max
            })
            const passesName = namePredicate(node.name) || namePredicate(node.path)
            const passesRisk = riskPredicate(node)
            if (passesRanges && passesName && passesRisk) res.push(node)
        } else if (node.type === 'folder' && node.children) {
            const kept = filterTreeForView(node.children, enabledMetricIds, ranges, namePredicate, riskPredicate)
            if (kept.length > 0) res.push({ ...node, children: kept })
        }
    }
    return res
}
