import { useId } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import CountrySelect from "./country-select";

export default function AddressFields({
  country,
  province,
  city,
  street,
  street2,
  street3,
}: {
  country: string;
  province: string;
  city: string;
  street: string;
  street2: string;
  street3: string;
}) {
  const countryId = useId();
  const provinceId = useId();
  const cityId = useId();
  const addr1Id = useId();
  const addr2Id = useId();
  const addr3Id = useId();

  return (
    <fieldset className="space-y-3">
      <div>
        <Label htmlFor={countryId}>Country</Label>
        <CountrySelect defaultValue={country} name="addrCountry" />
      </div>
      <div>
        <Label htmlFor={provinceId}>Province</Label>
        <Input
          id={provinceId}
          name="addrProvince"
          defaultValue={province}
          placeholder="Province"
          autoComplete="off"
        />
      </div>
      <div>
        <Label htmlFor={cityId}>City</Label>
        <Input
          id={cityId}
          name="addrCity"
          defaultValue={city}
          placeholder="City"
          type="text"
          autoComplete="off"
        />
      </div>

      <div className="space-y-1">
        <Label htmlFor={addr1Id}>Address Line 1</Label>
        <Input
          id={addr1Id}
          name="addrStreet"
          defaultValue={street}
          autoComplete="address-line1"
          placeholder="Address Line 1"
          type="text"
        />
      </div>
      <div className="space-y-1">
        <Label htmlFor={addr2Id}>Address Line 2</Label>
        <Input
          id={addr2Id}
          name="addrStreet2"
          defaultValue={street2}
          autoComplete="address-line2"
          placeholder="Address Line 2"
          type="text"
        />
      </div>
      <div className="space-y-1">
        <Label htmlFor={addr3Id}>Address Line 3</Label>
        <Input
          id={addr3Id}
          name="addrStreet3"
          defaultValue={street3}
          autoComplete="address-line3"
          placeholder="Address Line 3"
          type="text"
        />
      </div>
    </fieldset>
  );
}
