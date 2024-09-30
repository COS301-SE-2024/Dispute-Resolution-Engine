import { z } from "zod";
import ClientPage from "./page-client";
import { getWorkflowDetails } from "@/lib/api/workflow";
import { Workflow } from "@/lib/types";

const querySchema = z.object({
  id: z.coerce
    .number({
      message: "Invalid workflow ID",
    })
    .optional(),
});

export default async function WorkflowDesigner({ searchParams }: { searchParams: unknown }) {
  const params = await querySchema.safeParseAsync(searchParams).then(({ data, error }) => {
    if (error) {
      throw new Error(error.issues[0].message);
    }
    return data!;
  });

  let workflow: Workflow | undefined = undefined;
  if (params.id) {
    workflow = await getWorkflowDetails(params.id);
  }

  return <ClientPage workflow={workflow} />;
}
