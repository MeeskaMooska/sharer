import { FormEvent, useState } from "react";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { PageLayout } from "@/components/site/PageLayout";
import { signInWithEmail, type SupportedSchool } from "@/lib/auth";

export const Route = createFileRoute("/signup")({
  head: () => ({
    meta: [
      { title: "Sign Up — Sharer" },
      { name: "description", content: "Sign up to start using Sharer." },
    ],
  }),
  component: SignUpPage,
});

function SignUpPage() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [school, setSchool] = useState<SupportedSchool>("nova");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const fallbackEmail = `${name.trim().toLowerCase().replace(/\s+/g, ".")}@campus.edu`;
    signInWithEmail(email || fallbackEmail || "student@campus.edu", school);
    void navigate({ to: "/items" });
  }

  return (
    <PageLayout>
      <section className="mx-auto max-w-md px-4 py-12">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Sign Up</h1>
        <p className="mt-3 text-sm text-muted-foreground">
          Prototype mode: account is created instantly with any input.
        </p>

        <form onSubmit={handleSubmit} className="mt-8 space-y-4 rounded-xl border border-border bg-background p-5">
          <label className="flex flex-col gap-2 text-sm">
            <span className="font-medium text-foreground">Full name</span>
            <input
              type="text"
              value={name}
              onChange={(event) => setName(event.target.value)}
              placeholder="Student Name"
              className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
              required
            />
          </label>

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
            Create Account
          </button>
        </form>

        <p className="mt-4 text-sm text-muted-foreground">
          Already have an account?{" "}
          <Link to="/signin" className="font-medium text-primary hover:underline">
            Sign in
          </Link>
        </p>
      </section>
    </PageLayout>
  );
}
