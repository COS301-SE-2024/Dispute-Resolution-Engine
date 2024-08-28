import { useCallback } from "react";
import { Handle, Position } from "@xyflow/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const handleStyle = { left: 10 };

export default function CustomNode(data: any) {
  const onChange = useCallback((evt: any) => {
    console.log(evt.target.value);
  }, []);
  const arr = ["days", "hours", "minutes", "seconds"];
  const inputs = arr.map((currInput) => {
    return (
      <>
        <Label htmlFor={currInput + "Inp"}>{currInput.charAt(0).toUpperCase() + currInput.slice(1)}</Label>
        <Input id={currInput + "Inp"} type="number" className="h-2 w-16 p-3"></Input>
      </>
    );
  });
  return (
    <div className="bg-opacity-100">
      <Handle type="target" position={Position.Top} />
      <Card className="dark:bg-black">
        <CardHeader className="p-4 text-center">
          <CardTitle contentEditable="true">Node A</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 items-center gap-3">
            {inputs}
          </div>
        </CardContent>
      </Card>
      <Handle type="source" position={Position.Bottom} id="a" />
      <Handle type="source" position={Position.Bottom} id="b" style={handleStyle} />
    </div>
  );
}
