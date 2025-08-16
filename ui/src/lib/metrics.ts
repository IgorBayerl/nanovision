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
    const sum = f.reduce(
        (acc, x) => {
            if (!x.metrics) return acc
            const m = x.metrics
            acc.lineCoverage += m.lineCoverage
            acc.branchCoverage += m.branchCoverage
            acc.methodCoverage += m.methodCoverage
            acc.statementCoverage += m.statementCoverage
            acc.functionCoverage += m.functionCoverage
            return acc
        },
        {
            lineCoverage: 0,
            branchCoverage: 0,
            methodCoverage: 0,
            statementCoverage: 0,
            functionCoverage: 0,
        },
    )
    return {
        lineCoverage: Math.round(sum.lineCoverage / f.length),
        branchCoverage: Math.round(sum.branchCoverage / f.length),
        methodCoverage: Math.round(sum.methodCoverage / f.length),
        statementCoverage: Math.round(sum.statementCoverage / f.length),
        functionCoverage: Math.round(sum.functionCoverage / f.length),
    }
}

export const calculateFolderMetrics = (children: FileNode[]) => averageMetrics(flattenFiles(children))
