import { useId } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import GenderSelect from "./gender-select";
import { Gender } from "@/lib/interfaces";
import CountrySelect from "./country-select";
import { ProfileUpdateField } from "@/app/profile/profile-form";


export default function ProfileFields({
    firstName,
    lastName,
    phone,
    gender,
    country,
}: {
    firstName: string;
    lastName: string;
    phone: string;
    gender: Gender;
    country: string;
}) {
    const fnameId = useId();
    const lnameId = useId();
    const phoneId = useId();
    const genderId = useId();
    const contryId = useId();

    return (
        <fieldset className="space-y-3">
            <div className="grid grid-cols-2 gap-3">
                <ProfileUpdateField name="firstName" label="First Name" id={fnameId}>
                    <Input
                        id={fnameId}
                        name="firstName"
                        autoComplete="given-name"
                        placeholder="First Name"
                        defaultValue={firstName}
                    />
                </ProfileUpdateField>

                <ProfileUpdateField name="surname" label="Last Name" id={lnameId}>
                    <Input
                        id={lnameId}
                        name="surname"
                        autoComplete="family-name"
                        placeholder="Last Name"
                        defaultValue={lastName}
                    />
                </ProfileUpdateField>
            </div>

            <ProfileUpdateField name="phoneNumber" label="Phone Number" id={phoneId}>
                <Input
                    id={phoneId}
                    name="phoneNumber"
                    autoComplete="tel"
                    placeholder="Phone Number"
                    defaultValue={phone}
                />
            </ProfileUpdateField>
            <ProfileUpdateField name="gender" label="Gender" id={genderId}>
                <Label htmlFor={genderId}>Gender</Label>
                <GenderSelect name="gender" defaultValue={gender} />
            </ProfileUpdateField>
            <ProfileUpdateField name="country" label="Nationality" id={contryId}>
                <CountrySelect name="country" defaultValue={country} />
            </ProfileUpdateField>
        </fieldset>
    );
}
