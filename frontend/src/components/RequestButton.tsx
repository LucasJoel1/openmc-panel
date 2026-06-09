import { useQuery, useQueryClient } from "@tanstack/react-query"
import { Button } from "./ui/button"
import { Check, Loader2, X } from "lucide-react"
import { useEffect } from "react"
import { toast } from "sonner"

const exectuteRequest = async (req: string) => {
    const res = await fetch(req, {
        method: "POST"
    })
    const body = await res.json()
    if (body === null || body.message === null) throw new Error("an unexpected error occured")
    if (!res.ok) throw new Error(body.message)
    return body
}

interface RequestButtonProps {
    req: string
    contents: string
}

export default function RequestButton(props: RequestButtonProps) {
    const queryClient = useQueryClient()
    const { isLoading: isLoading, isSuccess: isSuccess, isError: isError, refetch: runQuery } = useQuery({
        queryKey: ['serverAction', props.req],
        queryFn: () => exectuteRequest(props.req),
        enabled: false
    })

    useEffect(() => {
        if (isSuccess || isError) {
            toast(props.contents + (isSuccess ? " server successfully executed" : " server failed to execute"))
            const timer = setTimeout(() => {
                queryClient.resetQueries({ queryKey: ['serverAction', props.req] })
            }, 4000)

            return () => clearTimeout(timer)
        }

    }, [isSuccess, isError, props.req, queryClient, props.contents])

    return (
        <Button variant="outline" className={"flex-1"} onClick={() => {
            runQuery()
        }}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isSuccess && <Check className="mr-2 h-4 w-4" />}
            {isError && <X className="mr-2 h-4 w-4" />}
            { props.contents }
        </Button>
    )
}