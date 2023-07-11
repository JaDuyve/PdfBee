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
        <div className="flex flex-row">
            <div className="basis-1/5 mr-5 overflow-x-auto h-[60rem]">
                <button onClick={openFileDialog}>Open pdf</button>

                {forms.map((formData: pdfutil.Form, i : number) => <Form key={i} form={formData} />)}
            </div>

            <div className="basis-4/5">
                {PdfViewer(pdfContent)}
            </div>
        </div>
    )
}

function Form({form}: {form: pdfutil.Form}): ReactElement {
    return (
        <div className="flex flex-col">
            <div className="relative flex flex-row justify-between">
                <h4 className="text-xl font-bold mb-3">
                    Text fields
                </h4>
            </div>

            {form.textFields?.map((textFieldData: pdfutil.TextField, i : number) => <TextField key={i} textField={textFieldData} />)}
        </div>
    )
}

function TextField({textField} : {textField: pdfutil.TextField}): ReactElement {
    return (
        <div className="py-2">
            <label htmlFor={textField.name} className="block text-sm font-medium text-gray-700">
                {textField.name}
            </label>

            <input
                type="text"
                id={textField.name}
                placeholder={textField.name}
                defaultValue={textField.value}
                className="mt-1 w-full rounded-md border-gray-200 shadow-sm md:text-md ms:text-ms"
            />
        </div>
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
            <iframe src={`data:application/pdf;base64,${pdfContent}`} className="w-full h-[60rem]"/>
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