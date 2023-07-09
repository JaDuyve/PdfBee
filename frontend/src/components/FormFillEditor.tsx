import {GetCurrentFile, OpenFileDialog} from "../../wailsjs/go/main/App";
import {ReactElement, useEffect, useState} from "react";
import {pdfutil} from "../../wailsjs/go/models";
import {EventsOn} from "../../wailsjs/runtime";

export default function FormFillEditor(): ReactElement {
    const [ pdfContent] = usePdfContent();

    function openFileDialog() {
        OpenFileDialog().then((result: pdfutil.Form[]) => setForms(result))
    }

    const [forms, setForms] = useState<pdfutil.Form[]>([]);
    return (
        <>
            <button onClick={openFileDialog}>Open pdf</button>

            {forms.map((formData: pdfutil.Form, i : number) => <Form key={i} form={formData} />)}

            {PdfViewer(pdfContent)}
        </>
    )
}

function Form({form}: {form: pdfutil.Form}): ReactElement {
    return (
        <div>
            <h2>Text fields</h2>

            {form.textFields?.map((textFieldData: pdfutil.TextField, i : number) => <TextField key={i} textField={textFieldData} />)}

        </div>
    )
}

function TextField({textField} : {textField: pdfutil.TextField}): ReactElement {
    return (
        <>
            <label id={textField.id} htmlFor={textField.name}>{textField.name}</label>
            <input id={textField.id} name={textField.name} defaultValue={textField.value}/>
            <br/>
        </>
    )
}

function PdfViewer(pdfContent: string): ReactElement {
    if (pdfContent === "") {
        return (
            <div>
                <p>No pdf present yet</p>
            </div>
        )
    }

    return (
        <div>
            <iframe src={`data:application/pdf;base64,${pdfContent}`} width="100%" height="300px"/>
        </div>
    )
}

function usePdfContent() {
    const [pdfContent, setPdfContent] = useState<string>("");

    const fetchData = async (): Promise<void> => {
        try {
            const response : string = await GetCurrentFile();
            setPdfContent(response);
        } catch (error) {
            console.log("failed to get data")

            // setLoading(false);
        } finally {
            // setLoading(false);
        }
    };


    useEffect(() => {
        EventsOn("current_file_changed", fetchData)
    }, []);

    return [ pdfContent ] as const;
}