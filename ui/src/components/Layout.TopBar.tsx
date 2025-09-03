import { ArrowLeft, Github } from 'lucide-react'
import { ThemeSwitch } from '@/components/Theme.Switch'
import { GITHUB_URL } from '@/lib/consts'
import { Button } from '@/ui/button'

export default function TopBar({ title, showBackButton = false }: { title: string; showBackButton?: boolean }) {
    return (
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
                <ThemeSwitch />
            </div>
        </div>
    )
}
