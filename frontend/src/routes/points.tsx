import { createFileRoute } from "@tanstack/react-router";
import { AuthGate } from "@/components/site/AuthGate";
import { PagePlaceholder, PlaceholderCard } from "@/components/site/PageLayout";

export const Route = createFileRoute("/points")({
  head: () => ({
    meta: [
      { title: "Points — Sharer" },
      { name: "description", content: "Track your sustainability points and tier." },
    ],
  }),
  component: PointsPage,
});

function PointsPage() {
  return (
    <AuthGate>
      <PagePlaceholder
        title="Points & Tiers"
        description="Earn points each time you share or reuse. Your tier, history, and rewards breakdown will live here."
      >
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <PlaceholderCard label="Total Points" />
          <PlaceholderCard label="Current Tier" />
          <PlaceholderCard label="Recent Activity" />
        </div>
      </PagePlaceholder>
    </AuthGate>
  );
}
