import type { SummaryV1 } from '@/types/summary'

declare global {
    interface Window {
        __ADLERCOV_SUMMARY__?: SummaryV1
        __ADLERCOV_DETAILS__?: DetailsV1
    }
}
