import { useRef } from 'react'
import { api, useApi } from '@/app/api'
import { type Crumb, Crumbs, TopBarActions } from '@/app/components/Chrome'
import { runHref } from '@/app/router'
import Layout from '@/components/Layout'
import { applyDefaultFilters } from '@/lib/applyDefaultFilters'
import SummaryPage from '@/pages/SummaryPage'

export const runCrumbs = (id: number): Crumb[] => [{ label: `run ${id}`, href: runHref(id) }]

/** The pages of a report before its data arrives: the frame, without a spinner. */
export function ReportSkeleton({ crumbs, error }: { crumbs: Crumb[]; error?: Error }) {
    return (
        <Layout title="" nav={<Crumbs items={crumbs} />} actions={<TopBarActions />}>
            {error ? (
                <div className="rounded-md border border-uncovered/40 bg-uncovered/10 p-6 text-sm text-uncovered">
                    {error.message}
                </div>
            ) : (
                <>
                    <div className="h-24 animate-pulse rounded-md bg-muted/50" />
                    <div className="h-[60vh] animate-pulse rounded-md bg-muted/50" />
                </>
            )}
        </Layout>
    )
}

/** The report of a stored run: the same page as the static report, fed by the API. */
export default function RunPage({ id }: { id: number }) {
    const summary = useApi<unknown>(api.summary(id))

    // the report reads its filters from the URL when it mounts, so the run's
    // defaults go there first
    const seeded = useRef(false)
    if (summary.data && !seeded.current) {
        applyDefaultFilters(summary.data)
        seeded.current = true
    }

    const crumbs = runCrumbs(id)
    if (!summary.data) return <ReportSkeleton crumbs={crumbs} error={summary.error} />

    return <SummaryPage data={summary.data} layout={{ nav: <Crumbs items={crumbs} />, actions: <TopBarActions /> }} />
}
