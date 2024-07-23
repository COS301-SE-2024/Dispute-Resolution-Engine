"use client";
import React, { HTMLAttributes, useEffect, useId, useState } from "react";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import { InfoIcon, SaveIcon, WorkflowIcon } from "lucide-react";

const placeholderText5Paragraph: string =
  "Alternative Dispute Resolution (ADR) provides a way for parties to resolve disputes without involving the judicial system. Traditional litigation processes are costly and time-consuming, varying significantly across different companies and domains. By automating these processes through custom workflow generation and NLP enhancements, the Dispute Resolution Engine aims to drastically increase the speed and cost-effectiveness of ADR.\n" +
  "\n" +
  "The Dispute Resolution Engine (affectionately known as DRE) offers users a convenient way to manage disputes, regardless of their role. The Archive feature allows users to access previously resolved disputes and related analytics. When involved in a dispute, users can easily upload their evidence with a click of a button, and rest assured that they will be notified as events unfold, ensuring they stay informed.\n" +
  "\n" +
  "The star of our show is our flexible workflow management. Users can choose from a selection of pre-made templates, create and customize one to fit their specific case, or ask our LLM-integrated engine to generate one automatically.\n" +
  "\n" +
  "No matter the case, we ensure a swift and smooth experience that will save you time and money.";
export default function Splash() {
  return (
    <div className="container mx-auto px-4 md:px-6 lg:px-8 flex flex-col items-center justify-center">
      <h1 className="text-white font-bold text-center text-5xl mt-10">
        Here to streamline your ADR process
      </h1>

      <div className="mx-auto w-full max-w-sm space-y-2 mt-8">
        <div className="grid grid-cols-3 gap-4">
          <div className="rounded-lg bg-gray-950 p-4 text-center text-gray-50">
            <WorkflowIcon className="mx-auto h-8 w-8" />
            <p className="mt-2 text-sm font-medium">Efficient Workflows</p>
          </div>
          <div className="rounded-lg bg-gray-950 p-4 text-center text-gray-50">
            <SaveIcon className="mx-auto h-8 w-8" />
            <p className="mt-2 text-sm font-medium">Cost Savings</p>
          </div>
          <div className="rounded-lg bg-gray-950 p-4 text-center text-gray-50">
            <InfoIcon className="mx-auto h-8 w-8" />
            <p className="mt-2 text-sm font-medium">Inclusive Summaries</p>
          </div>
        </div>
        <div className="flex justify-center mt-8">
          <Button>Learn More</Button>
        </div>
      </div>
      <div className="flex justify-center">
        <RalphTextGlasses size={400} />
      </div>
      <div className="bg-dre-400 bg-opacity-30 w-fit rounded-3xl mt-9">
        <h1 className="text-white font-bold text-center text-4xl mt-8 pt-7">More About Us</h1>
        <p className="max-w-3xl mt-8 text-justify mx-28 mb-8">Alternative Dispute Resolution (ADR) provides a way for
          parties
          to resolve disputes without involving the judicial system. Traditional litigation processes are costly and
          time-consuming, varying significantly across different companies and domains. By automating these processes
          through custom workflow generation and NLP enhancements, the Dispute Resolution Engine aims to drastically
          increase the speed and cost-effectiveness of ADR.
          <br /><br />
          The Dispute Resolution Engine (affectionately known as DRE) offers users a convenient way to manage disputes,
          regardless of their role. The Archive feature allows users to access previously resolved disputes and related
          analytics. When involved in a dispute, users can easily upload their evidence with a click of a button, and
          rest
          assured that they will be notified as events unfold, ensuring they stay informed.
          <br /><br />
          The star of our show is our flexible workflow management. Users can choose from a selection of pre-made
          templates, create and customize one to fit their specific case, or ask our LLM-integrated engine to generate
          one
          automatically.
          <br /><br />
          No matter the case, we ensure a swift and smooth experience that will save you time and money.
        </p></div>
    </div>
  );
}

function RalphTextGlasses(props: any) {
  const [text, setText] = useState("");

  useEffect(() => {
    const generateRandomText = () => {
      const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
      let result = ""
      const length = 100; // Length of the random text
      for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * characters.length)) + ' '
      }
      return result
    };

    const interval = setInterval(() => {
      setText(generateRandomText())
    }, 200)

    return () => clearInterval(interval)
  }, []);
  const imgSize: number = props.size;
  return (
      <div className="relative w-64 justify-center flex items-center mb-12">
        <Image
          src="/logo-hole.svg"
          alt="Racoon Logo"
          width={imgSize}
          height={imgSize}
          className="absolute top-0 left-0 z-10"
        />
        <div className="w-32 z-0 mt-14 bg-[#03152d]">
          <p className="text-justify text-xs text-[#2f67bd] -tracking-widest">{text}</p>
         </div>
      </div>
    );
}