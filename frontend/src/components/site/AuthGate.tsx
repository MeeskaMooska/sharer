import type { ReactNode } from "react";
import { useEffect, useState } from "react";
import { Navigate } from "@tanstack/react-router";
import { isSignedIn } from "@/lib/auth";

export function AuthGate({ children }: { children: ReactNode }) {
  const [ready, setReady] = useState(false);
  const [signedIn, setSignedIn] = useState(false);

  useEffect(() => {
    setSignedIn(isSignedIn());
    setReady(true);
  }, []);

  if (!ready) {
    return (
      <div className="mx-auto max-w-6xl px-4 py-10 text-sm text-muted-foreground">
        Checking sign in status...
      </div>
    );
  }

  if (!signedIn) {
    return <Navigate to="/signin" replace />;
  }

  return <>{children}</>;
}
