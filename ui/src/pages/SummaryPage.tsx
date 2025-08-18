import FileExplorer from '@/components/FileExplorer'
import Footer from '@/components/Layout.Footer'
import TopBar from '@/components/Layout.TopBar'
import SummaryMetrics from '@/components/SummaryMetrics'
import { cn } from '@/lib/utils'
import type { SummaryV1 } from '@/types/summary'

export default function SummaryPage({ data }: { data: SummaryV1 }) {
    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <TopBar title={data.title} />
            <SummaryMetrics data={data} />
            <FileExplorer tree={data.tree} />
            <Footer />
        </div>
    )
}
