import { useEffect, useRef } from 'react'
import { cn } from '../lib/utils'

interface LogViewerProps {
    logs: string
    height: string
    width: string
}

const LogViewer = ({ logs, height, width }: LogViewerProps) => {
    const logRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (logRef.current) {
            logRef.current.scrollTop = logRef.current.scrollHeight;
        }
    }, [logs]);

    return (
        <div className={cn('bg-gray-900 text-green-400 font-mono text-left p-4 overflow-auto text-sm rounded-lg' , height, width)} ref={logRef}>
            <pre className='whitespace-pre-wrap'>
                <code>{logs}</code>
            </pre>
        </div>
    )
}

export default LogViewer