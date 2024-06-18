import { Input } from "../ui/input";
import { Label } from "../ui/label";
import CountrySelect from "./country-select";

export default function AddressFields() {
  return (
    <fieldset className="space-y-3">
      <div>
        <Label htmlFor="addrCountry">Country</Label>
        <CountrySelect name="addrCountry" />
      </div>
      <div>
        <Label htmlFor="addrProvince">Province</Label>
        <Input
          autoComplete="off"
          id="addrProvince"
          placeholder="Province"
          type="text"
          name="addrProvince"
        />
      </div>
      <div>
        <Label htmlFor="addrCity">City</Label>
        <Input id="addrCity" placeholder="City" type="text" name="addrCity" />
      </div>

      <div className="space-y-2">
        <Label htmlFor="addrStreet2">Street Address</Label>
        <Input
          autoComplete="address-line1"
          id="addrStreet"
          name="addrStreet"
          placeholder="Address Line 1"
          type="text"
        />
        <Input
          autoComplete="address-line2"
          id="addrStreet2"
          name="addrStreet2"
          placeholder="Address Line 2"
          type="text"
        />

        <Input
          autoComplete="address-line3"
          id="addrStreet3"
          name="addrStreet3"
          placeholder="Address Line 3"
          type="text"
        />
      </div>
    </fieldset>
  );
}
