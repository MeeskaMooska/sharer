import type { ReactNode } from "react";
import { Header } from "./Header";
import { Footer } from "./Footer";

export function PageLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col bg-background">
      <Header />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  );
}

export function PagePlaceholder({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children?: ReactNode;
}) {
  return (
    <PageLayout>
      <section className="mx-auto max-w-6xl px-4 py-12">
        <h1 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">{title}</h1>
        <p className="mt-3 max-w-2xl text-muted-foreground">{description}</p>
        <div className="mt-8">{children}</div>
      </section>
    </PageLayout>
  );
}

export function PlaceholderCard({ label }: { label: string }) {
  return (
    <div className="flex h-44 items-center justify-center rounded-xl border border-dashed border-border bg-secondary/30 text-sm font-medium text-muted-foreground">
      {label}
    </div>
  );
}
