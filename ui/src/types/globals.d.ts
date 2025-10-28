import type { SummaryV1 } from '@/types/summary'

declare global {
    interface Window {
        __NANOVISION_SUMMARY__?: SummaryV1
        __NANOVISION_DETAILS__?: DetailsV1
    }
}
