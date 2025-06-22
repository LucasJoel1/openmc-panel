import ConsolePage from "./pages/ConsolePage";
import StatusPage from "./pages/StatusPage";

interface PageHandlerProps {
  page: number;
}

export default function PageHandler(props: PageHandlerProps) {
    return (
        <div className="flex-1 min-h-0 p-2">
            <StatusPage visible={props.page === 0} />
            <ConsolePage visible={props.page === 2} isPage={true} />
        </div>
    );
}
