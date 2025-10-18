import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Greet } from "../wailsjs/go/main/App";
import LogViewer from './components/LogViewer';

function App() {
    return (
        <div className='m-4 h-full' id="App">
            <LogViewer logs={"Hello\nWorld\nThis is a log viewer component."} height="h-[80vh]" width="w-full" />
            <div className='mt-12 flex flex-row items-center'>
                <div className='flex gap-4'>
                    <input type='text' className='rounded-md p-2 w-[30vw] bg-gray-700'></input>
                    <button className='bg-[#0a0a0a] px-4 rounded-md hover:bg-[#2b2a2a] transition-colors'>Browse</button>
                </div>
            </div>
        </div>
    )
}

export default App
