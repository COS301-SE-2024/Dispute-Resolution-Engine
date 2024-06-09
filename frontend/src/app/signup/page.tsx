import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { SignupForm, SignupButton, TextField, SignupField } from "./signup-form";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import CountrySelect from "@/components/form/country-select";
import { Input } from "@/components/ui/input";

export default function Signup() {
  return (
    <main className="md:pt-3 h-full">
      <Card asChild>
        <SignupForm className="mx-auto lg:w-1/2 md:w-3/4 md:h-fit h-full flex flex-col">
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
                  <Input
                    autoComplete="given-name"
                    id="firstName"
                    placeholder="First Name"
                    name="firstName"
                  />
                </SignupField>
                <SignupField name="lastName" label="Last Name">
                  <Input
                    autoComplete="family-name"
                    id="lastName"
                    placeholder="Last Name"
                    name="lastName"
                  />
                </SignupField>
                <SignupField name="email" label="Email">
                  <Input autoComplete="email" id="email" name="email" placeholder="Email" />
                </SignupField>
                <SignupField name="password" label="Password">
                  <Input
                    autoComplete="new-password"
                    id="password"
                    name="password"
                    placeholder="Password"
                    type="password"
                  />
                </SignupField>
                <SignupField name="passwordConfirm" label="Confirm Password">
                  <Input
                    autoComplete="new-password"
                    id="passwordConfirm"
                    name="passwordConfirm"
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
                  <Input
                    autoComplete="off"
                    id="addrProvince"
                    placeholder="Province"
                    type="text"
                    name="addrProvince"
                  />
                </SignupField>
                <SignupField name="addrCity" label="City">
                  <Input
                    autoComplete="c"
                    id="addrCity"
                    placeholder="City"
                    type="text"
                    name="addrCity"
                  />
                </SignupField>
                <SignupField name="addrStreet" label="Street 1">
                  <Input
                    autoComplete="address-line1"
                    id="addrStreet"
                    name="addrStreet"
                    placeholder="Street 1"
                    type="text"
                  />
                </SignupField>
                <SignupField name="addrStreet2" label="Street 2">
                  <Input
                    autoComplete="address-line2"
                    id="addrStreet2"
                    name="addrStreet2"
                    placeholder="Street 2"
                    type="text"
                  />
                </SignupField>
                <SignupField name="addrStreet3" label="Street 3">
                  <Input
                    autoComplete="address-line3"
                    id="addrStreet3"
                    name="addrStreet3"
                    placeholder="Street 3"
                    type="text"
                  />
                </SignupField>
              </TabsContent>
            </Tabs>
          </CardContent>
          <CardFooter className="mt-auto flex md:justify-start justify-end">
            <SignupButton />
            {/* <p role="alert">{state?.data}</p> */}
          </CardFooter>
        </SignupForm>
      </Card>
    </main>
  );
}
