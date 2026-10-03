import { ChevronRight } from 'lucide-react'
import { Fragment, type ReactNode } from 'react'
import { GithubIcon } from '@/components/GithubIcon'
import { ThemeSwitch } from '@/components/Theme.Switch'
import { GITHUB_URL } from '@/lib/consts'
import { cn } from '@/lib/utils'
import { Button } from '@/ui/button'

export interface Crumb {
    label: string
    href?: string
    title?: string
}

/** Where the page is: nanovision / project / stream / run. */
export function Crumbs({ items = [] }: { items?: Crumb[] }) {
    return (
        <nav aria-label="Breadcrumb" className="flex min-w-0 items-center gap-1 text-sm">
            <a
                href="/"
                title="Open the newest run"
                className="shrink-0 font-bold text-base tracking-tight hover:text-primary"
            >
                nanovision
            </a>
            {items.map((item, i) => {
                const last = i === items.length - 1
                return (
                    <Fragment key={`${item.href ?? ''}|${item.label}`}>
                        <ChevronRight className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                        {item.href && !last ? (
                            <a
                                href={item.href}
                                title={item.title ?? item.label}
                                className="truncate text-muted-foreground hover:text-primary"
                            >
                                {item.label}
                            </a>
                        ) : (
                            <span
                                title={item.title ?? item.label}
                                className={cn('truncate', last ? 'font-medium' : 'text-muted-foreground')}
                            >
                                {item.label}
                            </span>
                        )}
                    </Fragment>
                )
            })}
        </nav>
    )
}

export function TopBarActions() {
    return (
        <>
            <a href={GITHUB_URL} target="_blank" rel="noopener noreferrer">
                <Button variant="outline" size="sm" className="h-8 w-8 rounded-sm p-0" title="GitHub">
                    <GithubIcon className="h-4 w-4" />
                </Button>
            </a>
            <ThemeSwitch />
        </>
    )
}

/** The frame of the pages that are not a report. */
export function AppShell({ children }: { children: ReactNode }) {
    return (
        <div className="min-h-full text-foreground">
            <header className="fixed inset-x-0 top-0 z-30 flex h-14 items-center justify-between gap-3 border-border border-b bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/80">
                <Crumbs />
                <div className="flex shrink-0 items-center gap-2">
                    <TopBarActions />
                </div>
            </header>
            <main className="mx-auto w-full max-w-6xl space-y-6 px-4 pt-20 pb-12 sm:px-6">{children}</main>
        </div>
    )
}
