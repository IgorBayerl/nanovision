import { ArrowLeft, Github, Maximize2, Minimize2 } from 'lucide-react'
import Footer from '@/components/Layout.Footer'
import { ThemeSwitch } from '@/components/Theme.Switch'
import { useLocalStorage } from '@/hooks/useLocalStorage'
import { GITHUB_URL } from '@/lib/consts'
import { cn } from '@/lib/utils'
import { Button } from '@/ui/button'

interface LayoutProps {
    children: React.ReactNode
    title: string
    showBackButton?: boolean
}

export default function Layout({ children, title, showBackButton }: LayoutProps) {
    const [isFullWidth, setIsFullWidth] = useLocalStorage('layout-isFullWidth', false)

    const toggleFullWidth = () => setIsFullWidth(!isFullWidth)

    return (
        <div
            className={cn('mx-auto min-h-screen w-full space-y-5 bg-background p-6 text-foreground', {
                'max-w-7xl': !isFullWidth,
            })}
        >
            <div className="flex items-start justify-between">
                <div className="flex items-center gap-3">
                    {showBackButton && (
                        <a href="./index.html" title="Back to summary">
                            <Button variant="outline" size="sm" className="h-8 w-8 rounded-sm p-0">
                                <ArrowLeft className="h-4 w-4" />
                            </Button>
                        </a>
                    )}
                    <h1 className="font-bold text-2xl tracking-tight">{title || 'Coverage Report'}</h1>
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
                </div>
            </div>
            {children}
            <Footer />
        </div>
    )
}
