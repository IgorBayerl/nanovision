import { useMemo, useState } from 'react'
import DiffStatusBadge from '@/components/DiffStatusBadge'
import { scrollToLine } from '@/lib/scrollToLine'
import type { Method } from '@/types/summary'
import { Input } from '@/ui/input'
import {
    SidebarContent,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from '@/ui/sidebar'

/**
 * Right-sidebar navigation listing every method/function in the file. Clicking
 * an entry scrolls the source viewer to that method's start line.
 */
export default function FunctionNav({ methods }: { methods: Method[] }) {
    const [query, setQuery] = useState('')

    const filtered = useMemo(() => {
        const q = query.trim().toLowerCase()
        const sorted = [...methods].sort((a, b) => a.startLine - b.startLine)
        if (!q) return sorted
        return sorted.filter((m) => m.name.toLowerCase().includes(q))
    }, [methods, query])

    return (
        <>
            <SidebarHeader>
                <div className="font-semibold text-sm">Methods</div>
                <Input
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="Filter methods…"
                    className="h-8"
                />
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel>{filtered.length} methods</SidebarGroupLabel>
                    <SidebarMenu>
                        {filtered.map((method) => (
                            <SidebarMenuItem key={`${method.name}-${method.startLine}`}>
                                <SidebarMenuButton title={method.name} onClick={() => scrollToLine(method.startLine)}>
                                    <span className="truncate font-mono text-xs">{method.name}</span>
                                    {/* Fixed-width badge + line-number slots keep the
                                        A/M indicators and line numbers aligned in columns. */}
                                    <span className="ml-auto flex shrink-0 items-center gap-1.5">
                                        <span className="flex w-3 justify-center">
                                            <DiffStatusBadge status={method.diffStatus} />
                                        </span>
                                        <span className="w-8 text-right text-muted-foreground text-xs tabular-nums">
                                            {method.startLine}
                                        </span>
                                    </span>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        ))}
                        {filtered.length === 0 && (
                            <div className="px-2 py-4 text-muted-foreground text-xs">No matching methods.</div>
                        )}
                    </SidebarMenu>
                </SidebarGroup>
            </SidebarContent>
        </>
    )
}
