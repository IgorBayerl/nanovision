import type { JSX } from 'react'
import { Fragment } from 'react'
import { TreeRow } from '@/components/Tree.Row'
import type { FileNode, MetricKey, Metrics, RiskLevel } from '@/types/summary'

export default function FileTree({
    nodes,
    expanded,
    toggleFolder,
    enabledMetrics,
    rowDensity,
    metricsForNode,
    riskForMetric,
    filesInFolderCount,
}: {
    nodes: FileNode[]
    expanded: Set<string>
    toggleFolder: (id: string) => void
    enabledMetrics: { id: MetricKey; label: string }[]
    rowDensity: 'comfortable' | 'compact'
    metricsForNode: (n: FileNode) => Partial<Metrics> | undefined
    riskForMetric: (n: FileNode, k: MetricKey) => RiskLevel
    filesInFolderCount: (n: FileNode) => number
}) {
    const render = (arr: FileNode[], depth: number): JSX.Element[] =>
        arr.flatMap((node) => {
            const row = (
                <TreeRow
                    key={node.id}
                    node={node}
                    depth={depth}
                    enabledMetrics={enabledMetrics}
                    rowDensity={rowDensity}
                    isExpanded={expanded.has(node.id)}
                    onToggleFolder={toggleFolder}
                    metricsForNode={metricsForNode}
                    riskForMetric={riskForMetric}
                    filesInFolderCount={filesInFolderCount}
                />
            )
            if (node.type === 'folder' && node.children && expanded.has(node.id)) {
                return [row, ...render(node.children, depth + 1)]
            }
            return [row]
        })

    return <Fragment>{render(nodes, 0)}</Fragment>
}
