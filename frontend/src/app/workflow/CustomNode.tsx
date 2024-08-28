import { useCallback } from "react";
import { Handle, Position } from "@xyflow/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

const handleStyle = { left: 10 };

export default function CustomNode(data : any) {
  const onChange = useCallback((evt: any) => {
    console.log(evt.target.value);
  }, []);

  return (
    <>
      <Handle type="target" position={Position.Top} />
      {/* <div>
        <label htmlFor="text">Text:</label>
        <input id="text" name="text" onChange={onChange} className="nodrag" />
      </div> */}
      <Card>
        <CardHeader>
          <CardTitle>Node A</CardTitle>
        </CardHeader>
        <CardContent>
          <Input></Input>
        </CardContent>
      </Card>
      <Handle type="source" position={Position.Bottom} id="a" />
      <Handle type="source" position={Position.Bottom} id="b" style={handleStyle} />
    </>
  );
}
