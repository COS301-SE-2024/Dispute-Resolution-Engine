import type { Metadata } from "next";
import LoginForm from "./login-form";

export const metadata: Metadata = {
  title: "Login",
  icons: "/logo.svg",
};

export default function Login() {
  return <LoginForm />;
}
