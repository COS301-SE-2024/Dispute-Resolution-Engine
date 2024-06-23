import { useId } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import CountrySelect from "./country-select";

export default function AddressFields() {
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
        <CountrySelect name="addrCountry" />
      </div>
      <div>
        <Label htmlFor={provinceId}>Province</Label>
        <Input
          id={provinceId}
          name="addrProvince"
          type="text"
          placeholder="Province"
          autoComplete="off"
        />
      </div>
      <div>
        <Label htmlFor={cityId}>City</Label>
        <Input id={cityId} placeholder="City" type="text" name="addrCity" />
      </div>

      <div className="space-y-1">
        <Label htmlFor={addr1Id}>Address Line 1</Label>
        <Input
          autoComplete="address-line1"
          id="addrStreet"
          name="addrStreet"
          placeholder="Address Line 1"
          type="text"
        />
      </div>
      <div className="space-y-1">
        <Label htmlFor={addr2Id}>Address Line 2</Label>
        <Input
          id={addr2Id}
          name="addrStreet2"
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
          autoComplete="address-line3"
          placeholder="Address Line 3"
          type="text"
        />
      </div>
    </fieldset>
  );
}
