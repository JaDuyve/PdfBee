import './App.css';
import FormFillEditor from "./components/FormFillEditor";
import {pdfjs} from "react-pdf";
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';

pdfjs.GlobalWorkerOptions.workerSrc = new URL(
    'pdfjs-dist/build/pdf.worker.min.js',
    import.meta.url,
).toString();

function App() {

    return (
        <>
            <div id="App" className="m-5">
                <FormFillEditor/>
            </div>
        </>
    )
}

export default App
