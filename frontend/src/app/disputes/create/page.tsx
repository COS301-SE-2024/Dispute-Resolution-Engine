import dynamic from 'next/dynamic';
const CreateDisputeClient = dynamic(() => import("@/app/disputes/create/CreateDisputeClient"), {
  ssr: false,
});

export default function CreateDispute() {
  return <CreateDisputeClient />;
}