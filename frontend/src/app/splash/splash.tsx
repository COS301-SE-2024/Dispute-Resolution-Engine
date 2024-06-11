"use client";
import React, { HTMLAttributes, useEffect, useId, useState } from "react";
import { Button } from "@/components/ui/button";
import Image from "next/image";

const placeholderText5Paragraph: string =
  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Adipiscing diam donec adipiscing tristique. Erat velit scelerisque in dictum non consectetur a. Integer quis auctor elit sed vulputate mi sit. Et odio pellentesque diam volutpat commodo sed egestas egestas. Molestie ac feugiat sed lectus vestibulum. Morbi tristique senectus et netus et. Tincidunt nunc pulvinar sapien et ligula ullamcorper malesuada proin. Bibendum ut tristique et egestas quis. Fames ac turpis egestas integer eget aliquet nibh praesent. Leo in vitae turpis massa. Vel eros donec ac odio tempor orci dapibus ultrices in. Sed tempus urna et pharetra pharetra massa massa ultricies. Scelerisque eleifend donec pretium vulputate sapien nec. Sit amet consectetur adipiscing elit duis. Egestas integer eget aliquet nibh praesent tristique magna sit amet. Commodo odio aenean sed adipiscing diam donec adipiscing tristique.Aenean vel elit scelerisque mauris pellentesque. Odio eu feugiat pretium nibh ipsum consequat nisl. Cursus turpis massa tincidunt dui ut ornare lectus sit. Velit laoreet id donec ultrices tincidunt arcu non. Non enim praesent elementum facilisis. Phasellus egestas tellus rutrum tellus pellentesque. Elit sed vulputate mi sit amet mauris commodo quis. Magna eget est lorem ipsum. Adipiscing tristique risus nec feugiat in fermentum posuere urna. Pellentesque diam volutpat commodo sed egestas egestas fringilla phasellus. Quam elementum pulvinar etiam non quam lacus suspendisse. Auctor augue mauris augue neque gravida in fermentum. Cursus mattis molestie a iaculis at erat pellentesque adipiscing. Lorem ipsum dolor sit amet. Ultrices eros in cursus turpis massa tincidunt dui. Sed lectus vestibulum mattis ullamcorper velit sed. Mi quis hendrerit dolor magna eget est lorem ipsum. Cursus vitae congue mauris rhoncus aenean. Ac turpis egestas sed tempus urna et. Amet dictum sit amet justo donec enim diam.Lacinia at quis risus sed vulputate odio ut enim blandit. Commodo elit at imperdiet dui accumsan sit amet nulla. Placerat vestibulum lectus mauris ultrices eros in cursus turpis. Sed odio morbi quis commodo odio aenean sed adipiscing. Habitasse platea dictumst quisque sagittis purus sit amet volutpat. Enim lobortis scelerisque fermentum dui faucibus in. At volutpat diam ut venenatis tellus in. Amet luctus venenatis lectus magna fringilla. Ultricies lacus sed turpis tincidunt id aliquet risus feugiat. Tortor vitae purus faucibus ornare. Gravida dictum fusce ut placerat orci nulla pellentesque dignissim. Quis auctor elit sed vulputate mi sit amet mauris commodo. Sollicitudin nibh sit amet commodo nulla facilisi. Auctor elit sed vulputate mi. Sed lectus vestibulum mattis ullamcorper velit sed ullamcorper morbi. Nec tincidunt praesent semper feugiat nibh sed.Pretium lectus quam id leo in. Porta nibh venenatis cras sed felis eget velit aliquet. Imperdiet proin fermentum leo vel. Morbi tristique senectus et netus et malesuada. Enim blandit volutpat maecenas volutpat blandit aliquam etiam erat. Felis eget velit aliquet sagittis id consectetur. Id eu nisl nunc mi ipsum faucibus vitae aliquet nec. Ut porttitor leo a diam sollicitudin tempor id eu nisl. Libero volutpat sed cras ornare arcu. Condimentum lacinia quis vel eros donec. Amet mattis vulputate enim nulla aliquet porttitor lacus luctus accumsan. Vitae aliquet nec ullamcorper sit amet risus nullam eget. Ut eu sem integer vitae justo. Eget dolor morbi non arcu risus quis varius quam quisque. Felis imperdiet proin fermentum leo vel orci porta non. Sit amet nulla facilisi morbi tempus iaculis urna id volutpat. Sed egestas egestas fringilla phasellus faucibus. Pharetra pharetra massa massa ultricies mi quis hendrerit dolor magna.Eget gravida cum sociis natoque penatibus. Luctus venenatis lectus magna fringilla urna. Malesuada bibendum arcu vitae elementum curabitur vitae nunc sed velit. Velit euismod in pellentesque massa placerat duis ultricies lacus sed. Nulla aliquet enim tortor at auctor urna nunc id cursus. Faucibus a pellentesque sit amet porttitor. Vivamus at augue eget arcu dictum varius duis at consectetur. Gravida rutrum quisque non tellus orci ac auctor. Consequat id porta nibh venenatis cras sed. Ultrices sagittis orci a scelerisque purus semper eget duis. Eget lorem dolor sed viverra. Ipsum dolor sit amet consectetur. Aliquam etiam erat velit scelerisque in. Tortor at risus viverra adipiscing at in. Nibh praesent tristique magna sit amet. Fermentum odio eu feugiat pretium nibh ipsum.";
export default function Splash() {
  return (
    <div className="container mx-auto px-4 md:px-6 lg:px-8">
      <h1 className="text-white font-bold text-center text-5xl">
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
        <RalphTextGlasses size={200} />
      </div>
      <h1 className="text-white font-bold text-center text-4xl mt-8">More About Us</h1>
      <p className="max-w-3xl mx-auto mt-8 text-justify">{placeholderText5Paragraph}</p>
    </div>
  );
}
function RalphTextGlasses(props: any) {
  // const [text, setText] = useState("");

  // useEffect(() => {
  //   const generateRandomText = () => {
  //     const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  //     let result = "";
  //     const length = 400; // Length of the random text
  //     for (let i = 0; i < length; i++) {
  //       result += characters.charAt(Math.floor(Math.random() * characters.length));
  //     }
  //     return result;
  //   };

  //   const interval = setInterval(() => {
  //     setText(generateRandomText());
  //   }, 200); // Update text every 1000 milliseconds (1 second)

  //   return () => clearInterval(interval); // Clean up the interval on component unmount
  // }, []);
  const text = "jhdsfkjalshdlfkhasdkhgoiasduhg\niasdhfkjsdahoyweoqhgifuhbxbjvsdkaljhfkjsdha\njkfhaskdljnlxzmfhlkajsdhganbslfjdshfkj\nlahsdfkhasdlkjfhasldjfjhdsfkjals\nhdlfkhasdkhgoiasduhg\niasdhfkjsdahoyweoqhgifuhbxbjvsdkaljhfkjs\ndhajkfhaskdljnlxzmfhlkajsdhganbslfj"
  const imgSize: number = props.size;
  return (
      <div className="relative">
        <Image
          src="/ralph.png"
          alt="Racoon Logo"
          width={imgSize}
          height={imgSize}
          className="absolute top-0 left-0 z-10"
        />
        <div className="w-64 z-0">
          <p className="text-justify">{text}</p>
        </div>
      </div>
    );
}

function InfoIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="12" cy="12" r="10" />
      <path d="M12 16v-4" />
      <path d="M12 8h.01" />
    </svg>
  );
}

function SaveIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M15.2 3a2 2 0 0 1 1.4.6l3.8 3.8a2 2 0 0 1 .6 1.4V19a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2z" />
      <path d="M17 21v-7a1 1 0 0 0-1-1H8a1 1 0 0 0-1 1v7" />
      <path d="M7 3v4a1 1 0 0 0 1 1h7" />
    </svg>
  );
}

function WorkflowIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <rect width="8" height="8" x="3" y="3" rx="2" />
      <path d="M7 11v4a2 2 0 0 0 2 2h4" />
      <rect width="8" height="8" x="13" y="13" rx="2" />
    </svg>
  );
}
