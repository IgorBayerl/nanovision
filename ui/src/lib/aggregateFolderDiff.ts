import type { DiffStatus, FileNode } from '@/types/summary'

/** Diff kinds that are meaningful to surface on a folder. */
type PropagatedDiff = Extract<DiffStatus, 'added' | 'modified'>

/**
 * Builds a map of folder id -> aggregated diff status by propagating each
 * changed file's status up its ancestor folders (VSCode-style: a folder is
 * decorated when any descendant changed). "modified" wins over "added" so a
 * folder containing any modified file reads as modified.
 */
export function aggregateFolderDiff(nodes: FileNode[]): Map<string, PropagatedDiff> {
    const parentMap = new Map<string, string | undefined>()
    for (const node of nodes) parentMap.set(node.id, node.parentId || undefined)

    const result = new Map<string, PropagatedDiff>()

    const bump = (folderId: string, status: PropagatedDiff) => {
        const current = result.get(folderId)
        // "modified" outranks "added"; once modified it stays modified.
        if (current === 'modified') return
        result.set(folderId, status)
    }

    for (const node of nodes) {
        if (node.type !== 'file') continue
        if (node.diffStatus !== 'added' && node.diffStatus !== 'modified') continue
        const status = node.diffStatus
        let parentId = parentMap.get(node.id)
        while (parentId) {
            bump(parentId, status)
            parentId = parentMap.get(parentId)
        }
    }

    return result
}
