import { TreeRow } from '@/components/Tree.Row'
import type { FileNode, MetricConfig, Metrics } from '@/types/summary'

type RenderNode = FileNode & { depth: number }

interface BodyProps {
    nodes: RenderNode[]
    enabledMetrics: MetricConfig[]
    expandedFolders: Set<string>
    onToggleFolder: (id: string, event: React.MouseEvent | React.KeyboardEvent) => void
    viewMode: 'tree' | 'flat'
    isPinned: boolean
}

const metricsForNode = (node: FileNode): Partial<Metrics> | undefined => {
    return node.metrics
}

export default function FileExplorerBody({
    nodes,
    enabledMetrics,
    expandedFolders,
    onToggleFolder,
    viewMode,
    isPinned,
}: BodyProps) {
    return (
        <div className="w-full">
            {nodes.length > 0 ? (
                nodes.map((node, index) => (
                    <TreeRow
                        key={node.id}
                        node={node}
                        depth={node.depth}
                        enabledMetrics={enabledMetrics}
                        isExpanded={expandedFolders.has(node.id)}
                        onToggleFolder={onToggleFolder}
                        metricsForNode={metricsForNode}
                        viewMode={viewMode}
                        index={index}
                        isPinned={isPinned}
                    />
                ))
            ) : (
                <div className="py-16 text-center text-muted-foreground">No files match the current filters/search</div>
            )}
        </div>
    )
}
