import { useMemo } from 'react'
import FileExplorer from '@/components/FileExplorer'
import Footer from '@/components/Layout.Footer'
import TopBar from '@/components/Layout.TopBar'
import SummaryMetrics from '@/components/SummaryMetrics'
import { cn } from '@/lib/utils'
import type { MetadataItem, SummaryV1 } from '@/types/summary'

const NON_METRIC_KEYS = new Set(['files', 'folders'])

const isMetricKey = (key: string): boolean => {
    return !NON_METRIC_KEYS.has(key)
}

export default function SummaryPage({ data }: { data: SummaryV1 }) {
    const { reportInfo, metricKeys } = useMemo(() => {
        let reportInfo: { title: string; items: MetadataItem[] } | undefined

        if (data.metadata) {
            const validItems = data.metadata.filter(
                (item) => item.value !== undefined && (!Array.isArray(item.value) || item.value.length > 0),
            )

            // Only create the info object if there are valid items to show.
            if (validItems.length > 0) {
                reportInfo = {
                    title: 'Report Information',
                    items: validItems,
                }
            }
        }

        const keys = Object.keys(data.totals).filter(isMetricKey)

        return {
            reportInfo,
            metricKeys: keys,
        }
    }, [data])

    return (
        <div className={cn('mx-auto min-h-screen w-full max-w-7xl space-y-5 bg-background p-6 text-foreground')}>
            <TopBar title={data.title} />
            <SummaryMetrics info={reportInfo} metrics={data.totals} metricOrder={metricKeys} />
            <FileExplorer tree={data.tree} availableMetrics={metricKeys} />
            <Footer />
        </div>
    )
}
