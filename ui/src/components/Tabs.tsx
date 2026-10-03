import { cn } from '@/lib/utils'

export interface Tab<T extends string> {
    id: T
    label: string
    /** Shown after the label, e.g. how many files the tab lists. */
    count?: number
}

interface TabsProps<T extends string> {
    tabs: Tab<T>[]
    active: T
    onChange: (id: T) => void
}

/** The tab bar on top of the report. The side panel stays the same under every tab. */
export default function Tabs<T extends string>({ tabs, active, onChange }: TabsProps<T>) {
    return (
        <div role="tablist" className="flex gap-1 border-border border-b">
            {tabs.map((tab) => {
                const selected = tab.id === active
                return (
                    <button
                        key={tab.id}
                        type="button"
                        role="tab"
                        aria-selected={selected}
                        onClick={() => onChange(tab.id)}
                        className={cn(
                            '-mb-px flex items-center gap-2 border-b-2 px-3 py-2 font-medium text-sm transition-colors',
                            selected
                                ? 'border-primary text-foreground'
                                : 'border-transparent text-muted-foreground hover:text-foreground',
                        )}
                    >
                        {tab.label}
                        {tab.count !== undefined && (
                            <span className="rounded-sm bg-muted px-1.5 text-muted-foreground text-xs tabular-nums">
                                {tab.count}
                            </span>
                        )}
                    </button>
                )
            })}
        </div>
    )
}
