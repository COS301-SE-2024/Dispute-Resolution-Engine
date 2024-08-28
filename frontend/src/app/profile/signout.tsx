"use client"

import { Button } from "@/components/ui/button"
import { signout } from "@/lib/actions/auth"

export default function SignOut() {
    return <Button onClick={() => {signout()}}>Sign Out</Button>
}