"use client";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

type disputeDuration = {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
};

export default function TimerCheckbox(data: any) {
  const [duration, setDuration] = useState<disputeDuration>({
    days: 0,
    hours: 10,
    minutes: 20,
    seconds: 30,
  });
  const [isChecked, setIsChecked] = useState(true);

  const handleInputChange = (
    evt: React.ChangeEvent<HTMLInputElement>,
    unit: keyof disputeDuration
  ) => {
    const value = parseInt(evt.target.value, 10) || 0;
    setDuration((prevDuration) => ({
      ...prevDuration,
      [unit]: value,
    }));
  };

  const handleCheckboxChange = (checked: boolean) => {
    setIsChecked(checked);
  };

  const arr = ["days", "hours", "minutes", "seconds"];
  const inputs = arr.map((currInput) => (
    <div key={currInput}>
      <Label htmlFor={currInput + "Inp"}>
        {currInput.charAt(0).toUpperCase() + currInput.slice(1)}
      </Label>
      <Input
        id={currInput + "Inp"}
        type="number"
        className="h-2 w-16 p-3"
        value={duration[currInput as keyof disputeDuration]}
        onChange={(evt) => handleInputChange(evt, currInput as keyof disputeDuration)}
      />
    </div>
  ));

  const timerText = () => {
    return (
      <div>
        Timer: {duration.days === 0 ? "" : duration.days + "d"}{" "}
        {duration.hours === 0 ? "" : duration.hours + "h"}{" "}
        {duration.minutes === 0 ? "" : duration.minutes + "m"}{" "}
        {duration.seconds === 0 ? "" : duration.seconds + "s"}
      </div>
    );
  };

  return (
    <div className="grid grid-cols-[auto_1fr] gap-3">
      <Checkbox
        id="timerCheckbox"
        checked={isChecked}
        onCheckedChange={handleCheckboxChange}
      />
      <Label htmlFor="timerCheckbox" className={!isChecked ? "pointer-events-none" : ""}>
        {isChecked ? (
          <Accordion type="single" collapsible className="w-full">
            <AccordionItem value="item-1">
              <AccordionTrigger className="pt-0 text-2xl">
                {timerText()}
              </AccordionTrigger>
              <AccordionContent className="grid grid-cols-2 items-center gap-3 pt-3">
                {inputs}
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        ) : (
          <div className="text-gray-600">
            {timerText()}
          </div>
        )}
      </Label>
    </div>
  );
}
