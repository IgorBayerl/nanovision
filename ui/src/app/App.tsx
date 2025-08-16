import SummaryPage from '@/pages/SummaryPage'
import type { SummaryV1 } from '@/types/summary'

export default function AppRoot({ data }: { data: SummaryV1 }) {
    return <SummaryPage data={data} />
}
