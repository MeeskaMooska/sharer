import { FormEvent, useState } from "react";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { PageLayout } from "@/components/site/PageLayout";
import { signInWithEmail, type SupportedSchool } from "@/lib/auth";

export const Route = createFileRoute("/signin")({
  head: () => ({
    meta: [
      { title: "Sign In — Sharer" },
      { name: "description", content: "Sign in to access campus sharing pages." },
    ],
  }),
  component: SignInPage,
});

function SignInPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [school, setSchool] = useState<SupportedSchool>("nova");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    signInWithEmail(email || "student@campus.edu", school);
    void navigate({ to: "/items" });
  }

  return (
    <PageLayout>
      <section className="mx-auto max-w-md px-4 py-12">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Sign In</h1>
        <p className="mt-3 text-sm text-muted-foreground">
          No verification for now. Enter any values to continue.
        </p>

        <form onSubmit={handleSubmit} className="mt-8 space-y-4 rounded-xl border border-border bg-background p-5">
          <label className="flex flex-col gap-2 text-sm">
            <span className="font-medium text-foreground">Email</span>
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="you@school.edu"
              className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
              required
            />
          </label>

          <label className="flex flex-col gap-2 text-sm">
            <span className="font-medium text-foreground">School</span>
            <select
              value={school}
              onChange={(event) => setSchool(event.target.value as SupportedSchool)}
              className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
              required
            >
              <option value="nova">NOVA CC</option>
              <option value="uva">UVA</option>
            </select>
          </label>

          <label className="flex flex-col gap-2 text-sm">
            <span className="font-medium text-foreground">Password</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="Any password"
              className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
              required
            />
          </label>

          <button
            type="submit"
            className="inline-flex h-10 w-full items-center justify-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            Sign In
          </button>
        </form>

        <p className="mt-4 text-sm text-muted-foreground">
          New here?{" "}
          <Link to="/signup" className="font-medium text-primary hover:underline">
            Create an account
          </Link>
        </p>
      </section>
    </PageLayout>
  );
}
