import { Expert } from "@/lib/interfaces/dispute";
import { Button } from "../ui/button";
import { Badge } from "../ui/badge";

export interface ExpertItemProps extends Expert {
}

export default function ExpertItem(props: ExpertItemProps) {
    return <section className="p-2 rounded-lg">
        <div className="flex items-center gap-2">
            <h4 className="text-xl font-bold">{props.full_name}</h4>
            <Badge>{props.role}</Badge>
        </div>
        <div className="flex justify-between flex-wrap gap-3">
            <dl className="grid grid-cols-2">
                <dt className="font-semibold">Email</dt>
                <dd>{props.email}</dd>

                <dt className="font-semibold">Phone</dt>
                <dd>{props.phone}</dd>
            </dl>
            <div className="space-x-2">
                <Button type="submit">Approve</Button>

                <Button variant="destructive">Reject</Button>
            </div>
        </div>
    </section>
}

