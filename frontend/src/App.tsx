import './App.css';
import FormFillEditor from "./components/FormFillEditor";
import {ReactElement, useEffect, useState} from "react";
import {GetCurrentFile} from "../wailsjs/go/main/App";
import {EventsOn} from "../wailsjs/runtime";


function App() {

    return (
        <>
            <div id="App">
                <FormFillEditor/>
            </div>
        </>
    )
}

export default App
