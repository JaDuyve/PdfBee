import {
    GetPdfForm,
    GetPreviewContent,
    OpenFileDialog,
    UpdatePdfForm,
    UpdatePdfFormWithFieldNames
} from "../../wailsjs/go/main/App";
import React, {ReactElement, useEffect, useState} from "react";
import {pdfutil} from "../../wailsjs/go/models";
import {EventsOn} from "../../wailsjs/runtime";

export default function FormFillEditor(): ReactElement {
    const [ pdfContent, form]  = usePdfContent();

    async function openFileDialog() {
        await OpenFileDialog()
        console.log("End OpenFileDialog")
    }

    async function handleFormUpdate() : Promise<void> {
        if (!form) {
            return
        }

        await UpdatePdfForm(form)
    }

    return (
        <div className="flex flex-row">
            <div className="basis-1/5 mr-5 overflow-x-auto h-[60rem]">
                <button onClick={openFileDialog}>Open pdf</button>

                <Form form={form} updateForm={handleFormUpdate} />
            </div>

            <div className="basis-4/5 h-4/5">
                {PdfViewer(pdfContent)}
            </div>
        </div>
    )
}

function Form({form, updateForm}: {form: pdfutil.Form | undefined, updateForm: () => Promise<void>}): ReactElement {
    if (!form) {
        return <></>;
    }

    return (
        <div className="flex flex-col">
            <div className="relative flex flex-row justify-between">
                <h4 className="text-xl font-bold mb-3">
                    Text fields
                </h4>
                <button onClick={(): void => { UpdatePdfFormWithFieldNames().then(() => {
                    console.log("updated form fields form button fill all fields")
                })}}>Fill all fields with name</button>
            </div>

            {form.textFields?.map((textFieldData: pdfutil.TextField, i : number) => <TextField key={i} textField={textFieldData} updateForm={updateForm} />)}
        </div>
    )
}

function TextField({textField, updateForm} : {textField: pdfutil.TextField, updateForm: () => Promise<void>}): ReactElement {
    if (textField.multiline){
        return (
            <div className="py-2">
                <label htmlFor={textField.name} className="block text-sm font-medium text-gray-700">
                    {textField.name}
                </label>

                <textarea
                    id={textField.name}
                    placeholder={textField.name}
                    defaultValue={textField.value}
                    aria-rowcount={5}
                    className="mt-1 w-full rounded-md border-gray-200 shadow-sm md:text-md ms:text-ms"
                    onBlur={updateForm}
                    onChange={(event: React.ChangeEvent<HTMLTextAreaElement>): void => {textField.value = event.target.value}}
                />
            </div>
        )
    }

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
                onBlur={updateForm}
                onChange={(event: React.ChangeEvent<HTMLInputElement>): void => {textField.value = event.target.value}}
            />
        </div>
    )
}

function PdfViewer(pdfContent: string): ReactElement {
    const [numPages, setNumPages] = useState<number>();
    const [pageNumber, setPageNumber] = useState<number>(1);

    function onDocumentLoadSuccess({numPages}: { numPages: number }) {
        setNumPages(numPages)
    }

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
    const [form, setForm] = useState<pdfutil.Form>();

    const fetchPreviewPdfData = async (): Promise<void> => {
        try {
            const response : string = await GetPreviewContent();
            setPdfContent(response);
        } catch (error) {
            console.log("failed to get data", error)
        }
    };

    const fetchFormData = async (): Promise<void> => {
        try {
            const response : pdfutil.Form = await GetPdfForm();
            setForm(response)
        } catch (error) {
            console.log("failed to get form data", error)
        }
    };


    useEffect(() => {
        EventsOn("preview_file_content_updated", fetchPreviewPdfData)
        EventsOn("form_content_updated", fetchFormData)
    }, []);

    return [ pdfContent, form ] as const;
}