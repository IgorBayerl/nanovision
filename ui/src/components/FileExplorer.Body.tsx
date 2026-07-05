import { useWindowVirtualizer } from '@tanstack/react-virtual'
import { useLayoutEffect, useRef, useState } from 'react'
import { TreeRow } from '@/components/Tree.Row'
import type { DiffStatus, FileNode, MetricConfig, Metrics } from '@/types/summary'

type RenderNode = FileNode & { depth: number }

interface BodyProps {
    nodes: RenderNode[]
    enabledMetrics: MetricConfig[]
    expandedFolders: Set<string>
    onToggleFolder: (id: string, event: React.MouseEvent | React.KeyboardEvent) => void
    viewMode: 'tree' | 'flat'
    isPinned: boolean
    /** Folder id -> diff status aggregated from descendant files. */
    folderDiffMap: Map<string, DiffStatus>
}

const metricsForNode = (node: FileNode): Partial<Metrics> | undefined => {
    return node.metrics
}

/** Estimated row height in px; rows are measured on mount so this is just a starting point. */
const ESTIMATED_ROW_HEIGHT = 32

export default function FileExplorerBody({
    nodes,
    enabledMetrics,
    expandedFolders,
    onToggleFolder,
    viewMode,
    isPinned,
    folderDiffMap,
}: BodyProps) {
    const listRef = useRef<HTMLDivElement>(null)

    // The whole page scrolls (there is no inner scroll container), so we virtualize
    // against the window. scrollMargin tells the virtualizer where the list starts
    // relative to the top of the document.
    const [scrollMargin, setScrollMargin] = useState(0)
    useLayoutEffect(() => {
        const el = listRef.current
        if (!el) return
        const update = () => setScrollMargin(el.getBoundingClientRect().top + window.scrollY)
        update()
        window.addEventListener('resize', update)
        return () => window.removeEventListener('resize', update)
    }, [])

    const virtualizer = useWindowVirtualizer({
        count: nodes.length,
        estimateSize: () => ESTIMATED_ROW_HEIGHT,
        overscan: 15,
        scrollMargin,
        getItemKey: (index) => nodes[index].id,
    })

    if (nodes.length === 0) {
        return <div className="py-16 text-center text-muted-foreground">No files match the current filters/search</div>
    }

    return (
        <div ref={listRef} className="w-full">
            <div style={{ height: virtualizer.getTotalSize(), position: 'relative', width: '100%' }}>
                {virtualizer.getVirtualItems().map((virtualRow) => {
                    const node = nodes[virtualRow.index]
                    const diffStatus = node.type === 'folder' ? folderDiffMap.get(node.id) : node.diffStatus
                    return (
                        <div
                            key={virtualRow.key}
                            data-index={virtualRow.index}
                            ref={virtualizer.measureElement}
                            style={{
                                position: 'absolute',
                                top: 0,
                                left: 0,
                                width: '100%',
                                transform: `translateY(${virtualRow.start - scrollMargin}px)`,
                            }}
                        >
                            <TreeRow
                                node={node}
                                depth={node.depth}
                                enabledMetrics={enabledMetrics}
                                isExpanded={expandedFolders.has(node.id)}
                                onToggleFolder={onToggleFolder}
                                metricsForNode={metricsForNode}
                                viewMode={viewMode}
                                index={virtualRow.index}
                                isPinned={isPinned}
                                diffStatus={diffStatus}
                            />
                        </div>
                    )
                })}
            </div>
        </div>
    )
}
