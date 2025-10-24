import { ReactNode, useEffect, useRef } from 'react'
import { cn } from '../lib/utils'
import { DESIRE_APP_LOGO } from '../constants/desire'
import { DESIRE_APP_NAME, DESIRE_APP_VERSION, DESIRE_APP_DESCRIPTION, DESIRE_APP_APPLICATION } from '../constants/desire';

interface LogViewerProps {
    logs: string | ReactNode
    height: string
    width: string
    handleResetLogs: () => void
}

const LogViewer = ({ logs, height, width, handleResetLogs }: LogViewerProps) => {
    const logRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (logRef.current) {
            logRef.current.scrollTop = logRef.current.scrollHeight;
        }
    }, [logs]);

    return (
        <div>
            <div className={cn('bg-gray-900 text-green-400 font-mono text-left p-4 overflow-auto text-md rounded-lg z-0 custom-scrollbar log-content', height, width)}
                ref={logRef}
                >
                <pre className='log-content whitespace-pre-wrap break-words'>
                    <br />
                    <p className='text-center bg-gradient-to-b from-yellow-400 to-purple-600 bg-clip-text text-transparent'>
                        {DESIRE_APP_LOGO}
                    </p>
                    <br /><br />
                    <code>{DESIRE_APP_NAME}<br />{DESIRE_APP_VERSION}<br />{DESIRE_APP_DESCRIPTION}<br />{DESIRE_APP_APPLICATION}`</code>
                    <div>{logs}</div>
                </pre>
            </div>
            <button onClick={handleResetLogs} className='absolute right-8 top-8 z-10 text-white bg-gray-800 hover:bg-gray-600 transition-colors bg-opacity-70 px-4 py-1 rounded-md'>Reset</button>

        </div>


    )
}

export default LogViewer