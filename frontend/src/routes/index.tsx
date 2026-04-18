import { createFileRoute, Link } from "@tanstack/react-router";
import { PageLayout, PlaceholderCard } from "@/components/site/PageLayout";
import { BookOpen, Bed, Recycle, Award } from "lucide-react";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "Sharer — Share School Supplies & Dorm Essentials on Campus" },
      {
        name: "description",
        content:
          "Sharer helps students lend, borrow, and reuse school supplies, dorm items, and everyday essentials across campus. Cut waste, save money, earn points.",
      },
      { property: "og:title", content: "Sharer — Campus Sharing for School & Dorm Items" },
      {
        property: "og:description",
        content:
          "Share textbooks, calculators, mini-fridges, lamps, and more with students across your campus.",
      },
    ],
  }),
  component: HomePage,
});

function SectionHeader({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <div className="mb-6 flex flex-col gap-1">
      <h2 className="text-2xl font-semibold tracking-tight text-foreground md:text-3xl">{title}</h2>
      <p className="text-sm text-muted-foreground">{subtitle}</p>
    </div>
  );
}

function HomePage() {
  return (
    <PageLayout>
      <section className="bg-secondary/50">
        <div className="mx-auto grid max-w-6xl gap-8 px-4 py-16 md:grid-cols-2 md:py-24">
          <div className="flex flex-col justify-center">
            <span className="mb-4 inline-flex w-fit items-center gap-2 rounded-full border border-border bg-background px-3 py-1 text-xs font-medium text-primary">
              <Recycle className="h-3.5 w-3.5" /> Sustainability on campus
            </span>
            <h1 className="text-4xl font-bold leading-tight tracking-tight text-foreground md:text-5xl">
              Share supplies. <span className="text-primary">Stock your dorm.</span>
            </h1>
            <p className="mt-4 max-w-xl text-base text-muted-foreground md:text-lg">
              Sharer connects students to lend and borrow school supplies — textbooks,
              calculators, lab gear — plus dorm essentials like mini-fridges, lamps, fans,
              and kitchenware. Save money, cut waste, and help your campus reuse more.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link
                to="/items"
                className="inline-flex items-center justify-center rounded-md bg-primary px-5 py-2.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-[var(--primary-hover)]"
              >
                Browse Items
              </Link>
              <Link
                to="/points"
                className="inline-flex items-center justify-center rounded-md border border-border bg-background px-5 py-2.5 text-sm font-medium text-foreground transition-colors hover:bg-secondary"
              >
                View Points
              </Link>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {[
              { icon: BookOpen, label: "School Supplies" },
              { icon: Bed, label: "Dorm Essentials" },
              { icon: Recycle, label: "Reuse & Reduce" },
              { icon: Award, label: "Earn Points" },
            ].map(({ icon: Icon, label }, i) => (
              <div
                key={i}
                className="flex flex-col items-start justify-end rounded-xl border border-border bg-background p-5 shadow-sm"
              >
                <Icon className="h-6 w-6 text-primary" />
                <p className="mt-3 text-sm font-medium text-foreground">{label}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-4 py-14">
        <SectionHeader
          title="Featured Top Items"
          subtitle="Coming soon — the highest-value, most in-demand school and dorm items shared by your campus community."
        />
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <PlaceholderCard key={i} label="" />
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-4 pb-16">
        <SectionHeader
          title="Top Tier Members"
          subtitle="Coming soon — meet the students leading the sustainability charge with the highest tier ranks."
        />
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <PlaceholderCard key={i} label="" />
          ))}
        </div>
      </section>
    </PageLayout>
  );
}
