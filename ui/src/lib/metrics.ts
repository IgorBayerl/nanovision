import type { FileNode } from '@/types/summary'

/**
 * Flattens a tree of FileNodes into a single array of files.
 * This is still useful for the "flat view" in the file explorer.
 */
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
