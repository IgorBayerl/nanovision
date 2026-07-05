import { ArrowLeft, Github, Maximize2, Minimize2, PanelLeft, PanelRight } from 'lucide-react'
import type { ReactNode } from 'react'
import Footer from '@/components/Layout.Footer'
import { ThemeSwitch } from '@/components/Theme.Switch'
import { useLocalStorage } from '@/hooks/useLocalStorage'
import { GITHUB_URL } from '@/lib/consts'
import { cn } from '@/lib/utils'
import { Button } from '@/ui/button'
import { SIDEBAR_DEFAULT_WIDTH, Sidebar, SidebarInset, SidebarProvider, SidebarTrigger } from '@/ui/sidebar'

interface LayoutProps {
    children: ReactNode
    title: string
    showBackButton?: boolean
    /** Content for the left sidebar (metrics, reports). Omit to hide it. */
    leftSidebar?: ReactNode
    /** Content for the right sidebar (function navigation). Omit to hide it. */
    rightSidebar?: ReactNode
}

export default function Layout({ children, title, showBackButton, leftSidebar, rightSidebar }: LayoutProps) {
    const [isFullWidth, setIsFullWidth] = useLocalStorage('layout-isFullWidth', false)
    const [leftOpen, setLeftOpen] = useLocalStorage('sidebar-left-open', true)
    const [rightOpen, setRightOpen] = useLocalStorage('sidebar-right-open', true)
    const [leftWidth, setLeftWidth] = useLocalStorage('sidebar-left-width', SIDEBAR_DEFAULT_WIDTH)
    const [rightWidth, setRightWidth] = useLocalStorage('sidebar-right-width', SIDEBAR_DEFAULT_WIDTH)

    const toggleFullWidth = () => setIsFullWidth(!isFullWidth)

    return (
        <SidebarProvider className="pt-14">
            {/* Fixed, full-width top bar. Side panels render beneath it (see the
                `pt-14` above and the sidebars' `top-14` offset) so the toolbar
                buttons keep a consistent position regardless of panel state. */}
            <header className="fixed inset-x-0 top-0 z-30 flex h-14 items-center justify-between gap-3 border-border border-b bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/80">
                <div className="flex min-w-0 items-center gap-3">
                    {leftSidebar && (
                        <SidebarTrigger
                            onClick={() => setLeftOpen(!leftOpen)}
                            title={leftOpen ? 'Hide metrics sidebar' : 'Show metrics sidebar'}
                            icon={<PanelLeft className="h-4 w-4" />}
                        />
                    )}
                    {showBackButton && (
                        <a href="./index.html" title="Back to summary">
                            <Button variant="outline" size="sm" className="h-8 w-8 rounded-sm p-0">
                                <ArrowLeft className="h-4 w-4" />
                            </Button>
                        </a>
                    )}
                    <h1 className="truncate font-bold text-2xl tracking-tight">{title || 'Coverage Report'}</h1>
                </div>

                <div className="flex items-center gap-2">
                    <a href={GITHUB_URL} target="_blank" rel="noopener noreferrer">
                        <Button variant="outline" size="sm" className="h-8 w-8 rounded-sm p-0" title="GitHub">
                            <Github className="h-4 w-4" />
                        </Button>
                    </a>
                    <Button
                        variant="outline"
                        size="sm"
                        className="h-8 w-8 rounded-sm p-0"
                        title={isFullWidth ? 'Set constrained width' : 'Set full width'}
                        onClick={toggleFullWidth}
                    >
                        {isFullWidth ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
                    </Button>
                    <ThemeSwitch />
                    {rightSidebar && (
                        <SidebarTrigger
                            onClick={() => setRightOpen(!rightOpen)}
                            title={rightOpen ? 'Hide method navigation' : 'Show method navigation'}
                            icon={<PanelRight className="h-4 w-4" />}
                        />
                    )}
                </div>
            </header>

            {leftSidebar && (
                <Sidebar
                    side="left"
                    open={leftOpen}
                    onOpenChange={setLeftOpen}
                    width={leftWidth}
                    onWidthChange={setLeftWidth}
                >
                    {leftSidebar}
                </Sidebar>
            )}

            <SidebarInset>
                <div className={cn('w-full flex-1 space-y-5 p-6', { 'mx-auto max-w-7xl': !isFullWidth })}>
                    {children}
                    <Footer />
                </div>
            </SidebarInset>

            {rightSidebar && (
                <Sidebar
                    side="right"
                    open={rightOpen}
                    onOpenChange={setRightOpen}
                    width={rightWidth}
                    onWidthChange={setRightWidth}
                >
                    {rightSidebar}
                </Sidebar>
            )}
        </SidebarProvider>
    )
}
