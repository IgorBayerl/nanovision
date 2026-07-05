import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'
import { PanelLeft } from 'lucide-react'
import * as React from 'react'
import { useIsMobile } from '@/hooks/useIsMobile'
import { cn } from '@/lib/utils'
import { Button } from '@/ui/button'
import { Separator } from '@/ui/separator'
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/ui/sheet'

// Default width is ~50% wider than the original 17rem. Resizing is clamped
// between MIN and MAX (in px).
export const SIDEBAR_DEFAULT_WIDTH = 408
export const SIDEBAR_MIN_WIDTH = 240
export const SIDEBAR_MAX_WIDTH = 640
const SIDEBAR_WIDTH_MOBILE = '20rem'

type SidebarContextValue = { isMobile: boolean }
const SidebarContext = React.createContext<SidebarContextValue | null>(null)

function useSidebar() {
    const ctx = React.useContext(SidebarContext)
    if (!ctx) throw new Error('useSidebar must be used within a SidebarProvider.')
    return ctx
}

/**
 * Wrapper that lays out sidebars and the main inset in a single flex row.
 * Unlike the stock shadcn provider it does not own a single open-state; each
 * `Sidebar` is independently controlled via `open`/`onOpenChange`, so the left
 * and right sidebars toggle separately.
 */
function SidebarProvider({ className, style, children, ...props }: React.ComponentProps<'div'>) {
    const isMobile = useIsMobile()
    const value = React.useMemo<SidebarContextValue>(() => ({ isMobile }), [isMobile])

    return (
        <SidebarContext.Provider value={value}>
            <div
                data-slot="sidebar-wrapper"
                style={style}
                className={cn('group/sidebar-wrapper flex min-h-svh w-full bg-subtle text-foreground', className)}
                {...props}
            >
                {children}
            </div>
        </SidebarContext.Provider>
    )
}

function Sidebar({
    side = 'left',
    open,
    onOpenChange,
    width = SIDEBAR_DEFAULT_WIDTH,
    onWidthChange,
    className,
    children,
    ...props
}: Omit<React.ComponentProps<'div'>, 'onResize'> & {
    side?: 'left' | 'right'
    open: boolean
    onOpenChange: (open: boolean) => void
    /** Panel width in px (desktop only). */
    width?: number
    /** When provided, a drag handle is rendered so the panel can be resized. */
    onWidthChange?: (width: number) => void
}) {
    const { isMobile } = useSidebar()
    const state = open ? 'expanded' : 'collapsed'

    if (isMobile) {
        return (
            <Sheet open={open} onOpenChange={onOpenChange}>
                <SheetContent
                    side={side}
                    className="w-[var(--sidebar-width)] bg-sidebar p-0 text-sidebar-foreground [&>button]:hidden"
                    style={{ '--sidebar-width': SIDEBAR_WIDTH_MOBILE } as React.CSSProperties}
                >
                    <SheetHeader className="sr-only">
                        <SheetTitle>Sidebar</SheetTitle>
                        <SheetDescription>Sidebar navigation</SheetDescription>
                    </SheetHeader>
                    <div className="flex h-full w-full flex-col">{children}</div>
                </SheetContent>
            </Sheet>
        )
    }

    const startResize = (e: React.MouseEvent) => {
        if (!onWidthChange) return
        e.preventDefault()
        const startX = e.clientX
        const startWidth = width
        const onMove = (ev: MouseEvent) => {
            const delta = side === 'left' ? ev.clientX - startX : startX - ev.clientX
            const next = Math.min(SIDEBAR_MAX_WIDTH, Math.max(SIDEBAR_MIN_WIDTH, startWidth + delta))
            onWidthChange(next)
        }
        const onUp = () => {
            window.removeEventListener('mousemove', onMove)
            window.removeEventListener('mouseup', onUp)
            document.body.style.userSelect = ''
            document.body.style.cursor = ''
        }
        window.addEventListener('mousemove', onMove)
        window.addEventListener('mouseup', onUp)
        document.body.style.userSelect = 'none'
        document.body.style.cursor = 'col-resize'
    }

    const panelWidth = open ? width : 0

    return (
        <div
            className="group peer hidden text-sidebar-foreground md:block"
            data-state={state}
            data-side={side}
            {...props}
        >
            {/* In-flow spacer that pushes the inset; collapses to 0 when closed. */}
            <div className="relative bg-transparent" style={{ width: panelWidth }} />
            {/* Fixed panel that sits off-canvas (width 0) when closed. */}
            <div
                className={cn(
                    // Offset below the fixed 3.5rem top bar so panels sit under it.
                    'fixed top-14 bottom-0 z-10 hidden overflow-hidden md:flex',
                    side === 'left' ? 'left-0 border-r' : 'right-0 border-l',
                    className,
                )}
                style={{ width: panelWidth }}
            >
                <div
                    data-sidebar="sidebar"
                    className="flex h-full flex-col border-sidebar-border bg-sidebar"
                    style={{ width }}
                >
                    {children}
                </div>
                {onWidthChange && open && (
                    <button
                        type="button"
                        aria-label="Resize sidebar"
                        onMouseDown={startResize}
                        className={cn(
                            'absolute inset-y-0 z-20 w-1.5 cursor-col-resize bg-transparent hover:bg-sidebar-ring/40',
                            side === 'left' ? 'right-0' : 'left-0',
                        )}
                    />
                )}
            </div>
        </div>
    )
}

function SidebarInset({ className, ...props }: React.ComponentProps<'main'>) {
    return (
        <main
            data-slot="sidebar-inset"
            className={cn(
                'relative flex min-h-[calc(100svh-3.5rem)] w-full min-w-0 flex-1 flex-col bg-subtle',
                className,
            )}
            {...props}
        />
    )
}

function SidebarTrigger({
    className,
    onClick,
    icon,
    ...props
}: React.ComponentProps<typeof Button> & { icon?: React.ReactNode }) {
    return (
        <Button
            data-slot="sidebar-trigger"
            variant="outline"
            size="sm"
            className={cn('h-8 w-8 rounded-sm p-0', className)}
            onClick={onClick}
            {...props}
        >
            {icon ?? <PanelLeft className="h-4 w-4" />}
        </Button>
    )
}

function SidebarHeader({ className, ...props }: React.ComponentProps<'div'>) {
    return <div data-slot="sidebar-header" className={cn('flex flex-col gap-2 p-3', className)} {...props} />
}

function SidebarFooter({ className, ...props }: React.ComponentProps<'div'>) {
    return <div data-slot="sidebar-footer" className={cn('flex flex-col gap-2 p-3', className)} {...props} />
}

function SidebarContent({ className, ...props }: React.ComponentProps<'div'>) {
    return (
        <div
            data-slot="sidebar-content"
            className={cn('flex min-h-0 flex-1 flex-col gap-2 overflow-auto p-3', className)}
            {...props}
        />
    )
}

function SidebarGroup({ className, ...props }: React.ComponentProps<'div'>) {
    return <div data-slot="sidebar-group" className={cn('flex flex-col gap-2', className)} {...props} />
}

function SidebarGroupLabel({ className, ...props }: React.ComponentProps<'div'>) {
    return (
        <div
            data-slot="sidebar-group-label"
            className={cn('px-1 font-medium text-sidebar-foreground/70 text-xs uppercase tracking-wide', className)}
            {...props}
        />
    )
}

function SidebarSeparator({ className, ...props }: React.ComponentProps<typeof Separator>) {
    return (
        <Separator
            data-slot="sidebar-separator"
            className={cn('mx-1 w-auto bg-sidebar-border', className)}
            {...props}
        />
    )
}

function SidebarMenu({ className, ...props }: React.ComponentProps<'ul'>) {
    return <ul data-slot="sidebar-menu" className={cn('flex w-full min-w-0 flex-col gap-0.5', className)} {...props} />
}

function SidebarMenuItem({ className, ...props }: React.ComponentProps<'li'>) {
    return <li data-slot="sidebar-menu-item" className={cn('group/menu-item relative', className)} {...props} />
}

const sidebarMenuButtonVariants = cva(
    'flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left text-sm outline-none ring-sidebar-ring transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 disabled:pointer-events-none disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground [&>svg]:size-4 [&>svg]:shrink-0',
    {
        variants: {
            size: {
                default: 'h-8',
                sm: 'h-7 text-xs',
            },
        },
        defaultVariants: { size: 'default' },
    },
)

function SidebarMenuButton({
    className,
    asChild = false,
    isActive = false,
    size,
    ...props
}: React.ComponentProps<'button'> &
    VariantProps<typeof sidebarMenuButtonVariants> & { asChild?: boolean; isActive?: boolean }) {
    const Comp = asChild ? Slot : 'button'
    return (
        <Comp
            data-slot="sidebar-menu-button"
            data-active={isActive}
            className={cn(sidebarMenuButtonVariants({ size }), className)}
            {...props}
        />
    )
}

export {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarInset,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarProvider,
    SidebarSeparator,
    SidebarTrigger,
    useSidebar,
}
