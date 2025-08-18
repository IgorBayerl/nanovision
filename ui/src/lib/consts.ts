import type { SortableSubMetricKey } from '@/types/summary'

export const GITHUB_URL = 'https://github.com/IgorBayerl/AdlerCov'
export const WEBSITE_URL = 'https://igorbayerl.github.io/AdlerCov/'
export const DOCS_URL = 'https://igorbayerl.github.io/AdlerCov/'

export const SUB_METRIC_COLS: { id: SortableSubMetricKey; label: string; width: number }[] = [
    { id: 'covered', label: 'Covered', width: 80 },
    { id: 'uncovered', label: 'Uncovered', width: 90 },
    { id: 'coverable', label: 'Coverable', width: 90 },
    { id: 'total', label: 'Total', width: 70 },
    { id: 'percentage', label: 'Percentage %', width: 140 },
]
