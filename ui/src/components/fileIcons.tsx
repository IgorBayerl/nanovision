import { File } from 'lucide-react'

// The idea is to return specific icons for important file names
export function getFileIcon(_fileName: string) {
    // const ext = fileName.split('.').pop()?.toLowerCase()
    return <File className="h-4 w-4 text-muted-foreground" />
}
