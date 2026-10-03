import { useRef } from 'react'
import { api, useApi } from '@/app/api'
import { Crumbs, TopBarActions } from '@/app/components/Chrome'
import { fileHref } from '@/app/router'
import { applyDefaultFilters } from '@/lib/applyDefaultFilters'
import { summaryHref } from '@/lib/summaryQuery'
import DetailsPage from '@/pages/DetailsPage'
import { ReportSkeleton, runCrumbs } from './RunPage'

export default function FilePage({ id, path }: { id: number; path: string }) {
    const details = useApi<unknown>(api.file(id, path))

    const seeded = useRef(false)
    if (details.data && !seeded.current) {
        applyDefaultFilters(details.data)
        seeded.current = true
    }

    const name = path.slice(path.lastIndexOf('/') + 1)
    const crumbs = [...runCrumbs(id), { label: name, href: fileHref(id, path), title: path }]
    if (!details.data) return <ReportSkeleton crumbs={crumbs} error={details.error} />

    return (
        <DetailsPage
            data={details.data}
            layout={{ nav: <Crumbs items={crumbs} />, actions: <TopBarActions />, backHref: summaryHref() }}
        />
    )
}
