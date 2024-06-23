import AddressFields from "@/components/form/address-fields";
import ProfileFields from "@/components/form/profile-fields";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { getProfile } from "@/lib/api/profile";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "DRE - Profile",
};
export default async function Profile() {
  const { data, error } = await getProfile();
  if (error || !data) {
    return <h1>Error</h1>;
  }

  return (
    <main className="md:py-3 h-full">
      <Card variant="page" className="max-h-full flex flex-col">
        <CardHeader className="mb-4 space-y-3">
          <CardTitle>Edit Profile</CardTitle>
          <dl className="grid grid-cols-2 mb-8 ">
            <div>
              <dd className="font-semibold">ID Number</dd>
              <dt className="dark:text-white/50 text-sm">123456789012</dt>
            </div>
            <div>
              <dd className="font-semibold">Date of Birth</dd>
              <dt className="dark:text-white/50 text-sm">{data.birthdate}</dt>
            </div>
          </dl>
        </CardHeader>
        <CardContent className="h-full overflow-y-auto">
          <Tabs defaultValue="profile">
            <TabsList>
              <TabsTrigger value="profile">Profile</TabsTrigger>
              <TabsTrigger value="address">Address</TabsTrigger>
            </TabsList>
            <TabsContent value="profile" forceMount className="data-[state=inactive]:hidden">
              <ProfileFields
                firstName={data.first_name}
                lastName={data.surname}
                email={data.email}
              />
            </TabsContent>
            <TabsContent value="address" forceMount className="data-[state=inactive]:hidden">
              <AddressFields />
            </TabsContent>
          </Tabs>
        </CardContent>
        <CardFooter className="flex justify-between">
          <Button variant="destructive">Delete Account</Button>
          <Button>Save</Button>
        </CardFooter>
      </Card>
    </main>
  );
}
