import { Input } from "../ui/input";
import { Label } from "../ui/label";

export default function ProfileFields({
  firstName,
  lastName,
  email,
}: {
  firstName: string;
  lastName: string;
  email: string;
}) {
  return (
    <fieldset className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div>
          <Label htmlFor="firstName">First Name</Label>
          <Input
            id="firstName"
            name="firstName"
            autoComplete="given-name"
            placeholder="First Name"
            defaultValue={firstName}
          />
        </div>

        <div>
          <Label htmlFor="firstName">Last Name</Label>
          <Input
            id="lastName"
            name="lastName"
            autoComplete="family-name"
            placeholder="Last Name"
            defaultValue={lastName}
          />
        </div>
      </div>
      <div>
        <Label htmlFor="email">Email</Label>
        <Input
          autoComplete="email"
          id="email"
          name="email"
          placeholder="Email"
          defaultValue={email}
        />
      </div>
    </fieldset>
  );
}
