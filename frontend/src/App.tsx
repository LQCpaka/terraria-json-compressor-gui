import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import LogViewer from './components/LogViewer';
import { GetFirstFile, SelectFile } from '../wailsjs/go/backend/App';

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

    return (
        <div className='m-4 h-full' id="App">
            <LogViewer logs={"Hello\nWorld\nThis is a log viewer component."} height="h-[80vh]" width="w-full" />
            <div className='mt-12 flex flex-row justify-between items-center'>
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
                    <button className='bg-[#0a0a0a] px-4 py-2 rounded-md hover:bg-[#2b2a2a] transition-colors ml-8'>Start JSON Compress</button>
                    <button className='bg-[#0a0a0a] px-4 py-2 rounded-md hover:bg-[#2b2a2a] transition-colors'>Save JSON File</button>
                </div>
            </div>
        </div>
    )
}

export default App
