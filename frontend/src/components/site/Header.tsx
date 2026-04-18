import { Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Leaf } from "lucide-react";
import { applySchoolTheme, getSignedInSchool, isSignedIn, type SupportedSchool } from "@/lib/auth";

const navItems = [
  { to: "/", label: "Home" },
  { to: "/items", label: "Browse Items" },
  { to: "/users", label: "Browse Users" },
  { to: "/points", label: "Points" },
  { to: "/profile", label: "Profile" },
] as const;

export function Header() {
  const [signedIn, setSignedIn] = useState(false);
  const [school, setSchool] = useState<SupportedSchool>("nova");

  useEffect(() => {
    setSignedIn(isSignedIn());
    const selectedSchool = getSignedInSchool();
    setSchool(selectedSchool);
    applySchoolTheme(selectedSchool);
  }, []);

  return (
    <header className="sticky top-0 z-40 w-full border-b border-border bg-background/90 backdrop-blur">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
        <Link to="/" className="flex items-center gap-2 font-semibold text-foreground">
          <span className="flex h-8 w-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <Leaf className="h-4 w-4" />
          </span>
          <span className="inline-flex h-8 items-center rounded-md border border-border bg-background px-2 text-xs font-semibold">
            {school === "uva" ? "UVA" : "NOVA"}
          </span>
          <span>Sharer</span>
        </Link>
        <div className="hidden items-center gap-2 md:flex">
          <nav className="items-center gap-1 md:flex">
            {navItems.map((item) => (
              <Link
                key={item.to}
                to={item.to}
                activeOptions={{ exact: item.to === "/" }}
                className="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-primary"
                activeProps={{ className: "bg-secondary text-primary" }}
              >
                {item.label}
              </Link>
            ))}
          </nav>
          {!signedIn && (
            <>
              <Link
                to="/signin"
                className="inline-flex h-9 items-center justify-center rounded-md border border-border bg-background px-3 text-sm font-medium text-foreground transition-colors hover:bg-secondary"
              >
                Sign In
              </Link>
              <Link
                to="/signup"
                className="inline-flex h-9 items-center justify-center rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
              >
                Sign Up
              </Link>
            </>
          )}
        </div>
        <div className="flex items-center gap-2 md:hidden">
          {!signedIn && (
            <>
              <Link
                to="/signin"
                className="inline-flex h-9 items-center justify-center rounded-md border border-border bg-background px-3 text-sm font-medium text-foreground transition-colors hover:bg-secondary"
              >
                Sign In
              </Link>
              <Link
                to="/signup"
                className="inline-flex h-9 items-center justify-center rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
              >
                Sign Up
              </Link>
            </>
          )}
          <details className="relative">
            <summary className="list-none cursor-pointer rounded-md border border-border px-3 py-2 text-sm">
              Menu
            </summary>
            <div className="absolute right-0 mt-2 flex w-44 flex-col rounded-md border border-border bg-background p-1 shadow-lg">
              {navItems.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  activeOptions={{ exact: item.to === "/" }}
                  className="rounded px-3 py-2 text-sm text-foreground hover:bg-secondary"
                  activeProps={{ className: "bg-secondary text-primary" }}
                >
                  {item.label}
                </Link>
              ))}
            </div>
          </details>
        </div>
      </div>
    </header>
  );
}
