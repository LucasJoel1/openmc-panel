import StatusPage from "./pages/StatusPage"

interface PageHandlerProps {
    page: number
}

export default function PageHandler(props: PageHandlerProps) {
    return (
        <div className={"m-4"}>
            <StatusPage visible={props.page === 0} />
        </div>
    )
}