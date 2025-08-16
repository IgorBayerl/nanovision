import { Moon, Sun } from 'lucide-react'
import { useTheme } from '@/components/Theme.Context'
import { Button } from '@/ui/button'

export function ThemeSwitch() {
    const { mode, setMode } = useTheme()
    return (
        <div className="flex items-center">
            <Button
                variant="outline"
                size="sm"
                className="h-8 w-8 rounded-sm p-0"
                onClick={() => setMode(mode === 'dark' ? 'light' : 'dark')}
                title="Toggle dark mode"
            >
                {mode === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            </Button>
        </div>
    )
}
