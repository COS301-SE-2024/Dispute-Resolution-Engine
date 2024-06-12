import AddressFields from "@/components/form/address-fields";
import ProfileFields from "@/components/form/profile-fields";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export default function Profile() {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page">
        <CardHeader>
          <CardTitle>Edit Profile</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="grid grid-cols-2 mb-8">
            <div>
              <dd className="font-semibold">ID Number</dd>
              <dt className="text-gray-500 text-sm">1234567890123</dt>
            </div>
            <div>
              <dd className="font-semibold">Date of Birth</dd>
              <dt className="text-gray-500 text-sm">5 January 2003</dt>
            </div>
          </dl>
          <Tabs defaultValue="profile">
            <TabsList>
              <TabsTrigger value="profile">Profile</TabsTrigger>
              <TabsTrigger value="address">Address</TabsTrigger>
            </TabsList>
            <TabsContent value="profile" forceMount className="data-[state=inactive]:hidden">
              <ProfileFields />
            </TabsContent>
            <TabsContent value="address" forceMount className="data-[state=inactive]:hidden">
              <AddressFields />
            </TabsContent>
          </Tabs>
        </CardContent>
        <CardFooter>
          <Button>Save Changes</Button>
        </CardFooter>
      </Card>
    </main>
  );
}
