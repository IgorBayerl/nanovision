import { useMemo } from 'react'
import FileExplorer from '@/components/FileExplorer'
import Footer from '@/components/Layout.Footer'
import TopBar from '@/components/Layout.TopBar'
import SummaryMetrics from '@/components/SummaryMetrics'
import { averageMetrics, flattenFiles } from '@/lib/metrics'
import { buildIdMap } from '@/lib/tree'
import { cn } from '@/lib/utils'
import type { SummaryV1 } from '@/types/summary'

export default function SummaryPage({ data }: { data: SummaryV1 }) {
    const tree = data.tree
    const idMap = useMemo(() => buildIdMap(tree), [tree])
    const globalAverages = useMemo(() => averageMetrics(flattenFiles(tree)), [tree])

    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <TopBar title={data.title} />
            <SummaryMetrics averages={globalAverages} />
            <FileExplorer tree={tree} idMap={idMap} />
            <Footer />
        </div>
    )
}
