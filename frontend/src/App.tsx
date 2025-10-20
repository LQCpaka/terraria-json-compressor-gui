import { ReactNode, useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import LogViewer from './components/LogViewer';
import { GetFirstFile, SelectFile, PreviewCSVFile, CompressAndSave, SaveFile } from '../wailsjs/go/backend/App';

interface LogEntry {
    level: string;
    message: string;
}

// Gọi hàm ListFolder
async function loadFolder(path: string) {
    try {
        const files = await GetFirstFile(path);
        console.log('Files:', files);
        return files;
    } catch (error) {
        console.error('Error listing folder:', error);
    }
}

function App() {
    const [folderPath, setFolderPath] = useState('');
    const [LogText, setLogText] = useState<LogEntry[]>([]);
    const [isCompressing, setIsCompressing] = useState(false);

    const handleBrowseClick = async () => {
        try {
            const selectedPath = await SelectFile();
            if (selectedPath) {
                setFolderPath(selectedPath);
                const files = await loadFolder(selectedPath);
                console.log('Files:', files);
            }
        } catch (error) {
            console.error('Error selecting folder:', error);
        }
    };

    const handlePreviewClick = async () => {
        try {
            const result = await PreviewCSVFile(folderPath);
            setLogText(prevItem => [...prevItem, ...result]);
        } catch (error) {
            console.error('Error previewing file:', error);
            setLogText(prevItem => [...prevItem, { level: 'error', message: `Error: ${error}` }]);
        }
    }

    const handleCompressClick = async () => {
        if (!folderPath) {
            setLogText(prevItem => [...prevItem, { level: 'error', message: 'Please select a file first!' }]);
            return;
        }

        setIsCompressing(true);
        try {
            // Compress and save in one step
            const logs = await CompressAndSave(folderPath);
            setLogText(prevItem => [...prevItem, ...logs]);
        } catch (error) {
            console.error('Error compressing file:', error);
            setLogText(prevItem => [...prevItem, { 
                level: 'error', 
                message: `Compression failed: ${error}` 
            }]);
        } finally {
            setIsCompressing(false);
        }
    };

    const renderLog = (log: LogEntry, index: number) => {
        const colorMap: Record<string, string> = {
            error: 'text-red-700',
            warn: 'text-yellow-700',
            info: 'text-slate-400',
            ok: 'text-green-600',
            summary: 'text-blue-600 font-bold'
        };
        
        const colorClass = colorMap[log.level] || 'text-white';
        const levelText = log.level.toUpperCase();
        
        return (
            <div key={index} className={`${colorClass} cursor-text select-text`}>
                [{levelText}] {log.message}
            </div>
        );
    };

    return (
        <div className='flex flex-col h-screen p-4 overflow-hidden' id="App">
            <div className='flex-1 overflow-hidden'>
                <LogViewer logs={
                    <div className="log-content">
                        {LogText.map((log, index) => renderLog(log, index))}
                    </div>
                } height="h-[85vh]" width="w-full" />
            </div>
            <div className='mt-4 flex flex-row justify-between items-center flex-shrink-0'>
                <div className='flex gap-4'>
                    <input
                        type='text'
                        className='rounded-md p-2 w-[30vw] bg-gray-700'
                        value={folderPath}
                        onChange={(e) => setFolderPath(e.target.value)}
                        placeholder='Enter folder path...'
                        readOnly
                    />
                    <button onClick={handleBrowseClick} className='bg-[#0a0a0a] px-4 rounded-md hover:bg-[#2b2a2a] transition-colors'>Browse</button>
                </div>
                <div className='flex gap-4'>
                    <button onClick={handlePreviewClick} className='bg-[#0a0a0a] px-4 py-2 rounded-md hover:bg-[#2b2a2a] transition-colors ml-8'>Preview File</button>
                    <button 
                        onClick={handleCompressClick} 
                        disabled={isCompressing}
                        className='bg-[#0a0a0a] px-4 py-2 rounded-md hover:bg-[#2b2a2a] transition-colors disabled:opacity-50 disabled:cursor-not-allowed'
                    >
                        {isCompressing ? 'Compressing...' : 'Start Compress'}
                    </button>
                </div>
            </div>
        </div>
    )
}

export default App
