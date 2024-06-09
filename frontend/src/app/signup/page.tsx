import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { SignupForm, SignupButton, TextField, SignupField } from "./signup-form";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import CountrySelect from "@/components/form/country-select";
import { Input } from "@/components/ui/input";

export default function Signup() {
  return (
    <main>
      <Card asChild className="mx-auto md:my-3 lg:w-1/2 md:w-3/4">
        <SignupForm>
          <CardHeader>
            <CardTitle>Create an Account</CardTitle>
          </CardHeader>
          <CardContent asChild>
            <Tabs defaultValue="profile">
              <TabsList>
                <TabsTrigger value="profile">Profile</TabsTrigger>
                <TabsTrigger value="address">Address</TabsTrigger>
              </TabsList>
              <TabsContent value="profile" forceMount className="data-[state=inactive]:hidden">
                <SignupField name="firstName" label="First Name">
                  <Input autoComplete="given-name" id="firstName" placeholder="First Name" />
                </SignupField>
                <SignupField name="lastName" label="Last Name">
                  <Input autoComplete="family-name" id="lastName" placeholder="Last Name" />
                </SignupField>
                <SignupField name="email" label="Email">
                  <Input autoComplete="email" id="email" placeholder="Email" />
                </SignupField>
                <SignupField name="password" label="Password">
                  <Input
                    autoComplete="new-password"
                    id="password"
                    placeholder="Password"
                    type="password"
                  />
                </SignupField>
                <SignupField name="passwordConfirm" label="Confirm Password">
                  <Input
                    autoComplete="new-password"
                    id="passwordConfirm"
                    placeholder="Confirm Password"
                    type="password"
                  />
                </SignupField>
              </TabsContent>
              <TabsContent value="address" forceMount className="data-[state=inactive]:hidden">
                <SignupField name="addrCountry" label="Country">
                  <CountrySelect name="addrCountry" />
                </SignupField>
                <SignupField name="addrProvince" label="Province">
                  <Input autoComplete="off" id="addrProvince" placeholder="Province" type="text" />
                </SignupField>
                <SignupField name="addrCity" label="City">
                  <Input autoComplete="c" id="addrCity" placeholder="City" type="text" />
                </SignupField>
                <SignupField name="addrStreet" label="Street 1">
                  <Input
                    autoComplete="address-line1"
                    id="addrStreet"
                    placeholder="Street 1"
                    type="text"
                  />
                </SignupField>
                <SignupField name="addrStreet2" label="Street 2">
                  <Input
                    autoComplete="address-line2"
                    id="addrStreet2"
                    placeholder="Street 2"
                    type="text"
                  />
                </SignupField>
                <SignupField name="addrStreet3" label="Street 3">
                  <Input
                    autoComplete="address-line3"
                    id="addrStreet3"
                    placeholder="Street 3"
                    type="text"
                  />
                </SignupField>
              </TabsContent>
            </Tabs>
          </CardContent>
          <CardFooter>
            <SignupButton />
            {/* <p role="alert">{state?.data}</p> */}
          </CardFooter>
        </SignupForm>
      </Card>
    </main>
  );
}
