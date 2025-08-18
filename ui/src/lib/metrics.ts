import type { FileNode, Metrics } from '@/types/summary'

export const flattenFiles = (nodes: FileNode[]): FileNode[] => {
    const list: FileNode[] = []
    const walk = (arr: FileNode[]) => {
        arr.forEach((n) => {
            if (n.type === 'file') list.push(n)
            if (n.children) walk(n.children)
        })
    }
    walk(nodes)
    return list
}

export const averageMetrics = (files: FileNode[]): Metrics | undefined => {
    const f = files.filter((x) => x.metrics)
    if (f.length === 0) return undefined

    const initial = {
        lineCoverage: { covered: 0, uncovered: 0, coverable: 0, total: 0, percentage: 0 },
        branchCoverage: { covered: 0, uncovered: 0, coverable: 0, total: 0, percentage: 0 },
        methodCoverage: { covered: 0, uncovered: 0, coverable: 0, total: 0, percentage: 0 },
        statementCoverage: { covered: 0, uncovered: 0, coverable: 0, total: 0, percentage: 0 },
        functionCoverage: { covered: 0, uncovered: 0, coverable: 0, total: 0, percentage: 0 },
    }

    const sum = f.reduce((acc, x) => {
        if (!x.metrics) return acc
        for (const key in x.metrics) {
            const metricKey = key as keyof Metrics
            const m = x.metrics[metricKey]
            acc[metricKey].covered += m.covered
            acc[metricKey].uncovered += m.uncovered
            acc[metricKey].coverable += m.coverable
            acc[metricKey].total += m.total
        }
        return acc
    }, initial)

    // Recalculate percentages based on the summed totals
    for (const key in sum) {
        const metricKey = key as keyof Metrics
        const metricSum = sum[metricKey]
        metricSum.percentage =
            metricSum.coverable > 0 ? Math.round((metricSum.covered / metricSum.coverable) * 100) : 100
    }

    return sum
}

export const calculateFolderMetrics = (children: FileNode[]) => averageMetrics(flattenFiles(children))
