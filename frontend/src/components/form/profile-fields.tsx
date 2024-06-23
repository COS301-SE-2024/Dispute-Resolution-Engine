import { useId } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import GenderSelect from "./gender-select";
import { Gender } from "@/lib/interfaces";
import CountrySelect from "./country-select";

export default function ProfileFields({
  firstName,
  lastName,
  gender,
}: {
  firstName: string;
  lastName: string;
  email: string;
  gender: Gender;
}) {
  const fnameId = useId();
  const lnameId = useId();
  const phoneId = useId();
  const genderId = useId();
  const contryId = useId();

  return (
    <fieldset className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div>
          <Label htmlFor={fnameId}>First Name</Label>
          <Input
            id={fnameId}
            name="firstName"
            autoComplete="given-name"
            placeholder="First Name"
            defaultValue={firstName}
          />
        </div>

        <div>
          <Label htmlFor={lnameId}>Last Name</Label>
          <Input
            id={lnameId}
            name="lastName"
            autoComplete="family-name"
            placeholder="Last Name"
            defaultValue={lastName}
          />
        </div>
      </div>

      <div>
        <Label htmlFor={phoneId}>Phone Number</Label>
        <Input
          id={phoneId}
          name="phoneNumber"
          autoComplete="tel"
          placeholder="Phone Number"
          defaultValue={lastName}
        />
      </div>
      <div>
        <Label htmlFor={genderId}>Gender</Label>
        <GenderSelect name="gender" defaultValue={gender} />
      </div>
      <div>
        <Label htmlFor={contryId}>Nationality</Label>
        <CountrySelect name="country" />
      </div>
    </fieldset>
  );
}
