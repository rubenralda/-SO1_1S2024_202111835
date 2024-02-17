import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { DatosRam } from "../wailsjs/go/main/App";
//import { DatosRam } from "../wailsjs/go/main/App"

function App() {
    const [resultText, setResultText] = useState();
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => console.log(result);

    function greet() {
        DatosRam().then(updateResultText);
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo" />
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" />
                <button className="btn" onClick={greet}>Greet</button>
            </div>
        </div>
    )
}

export default App
