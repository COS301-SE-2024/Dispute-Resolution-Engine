import { Result } from "../types";

export async function fetchCountries(): Promise<
  Result<
    {
      code: string;
      label: string;
    }[]
  >
> {
  return {
    status: 200,
    data: [
      {
        code: "za",
        label: "South Africa",
      },
      {
        code: "vn",
        label: "Viet Nam",
      },
      {
        code: "ci",
        label: "China",
      },
    ],
  };
}
