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
        <div>
            <div className={cn('bg-gray-900 text-green-400 font-mono text-left p-4 overflow-auto text-sm rounded-lg z-0', height, width)} ref={logRef}>
                <pre className='whitespace-pre-wrap'>
                    <code>{logs}</code>
                </pre>
            </div>
            <button className='absolute right-8 top-8 z-10 text-white bg-gray-800 hover:bg-gray-600 transition-colors bg-opacity-70 px-4 py-1 rounded-md'>Reset</button>

        </div>


    )
}

export default LogViewer