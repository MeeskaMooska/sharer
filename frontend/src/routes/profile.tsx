import { createFileRoute } from "@tanstack/react-router";
import { AuthGate } from "@/components/site/AuthGate";
import { PageLayout } from "@/components/site/PageLayout";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Award, Star } from "lucide-react";

export const Route = createFileRoute("/profile")({
  head: () => ({
    meta: [
      { title: "Profile — Sharer" },
      { name: "description", content: "Your Sharer profile, tier, points, and items." },
    ],
  }),
  component: ProfilePage,
});

function ItemCard({
  title,
  by,
  requested = false,
}: {
  title: string;
  by: string;
  requested?: boolean;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2 text-sm">
        <p className="text-muted-foreground">Shared by {by}</p>
        <p className="text-foreground">
          <span className="font-medium">Status:</span> {requested ? "Requested by students" : "Available"}
        </p>
      </CardContent>
    </Card>
  );
}

function ProfilePage() {
  return (
    <PageLayout>
      <AuthGate>
        <section className="mx-auto max-w-6xl px-4 py-12">
        <h1 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">Profile</h1>
        <p className="mt-3 max-w-2xl text-muted-foreground">
          Your public Sharer identity, item listings, and request activity.
        </p>

        <Card className="mt-8">
          <CardHeader className="flex flex-row items-center gap-5">
            <Avatar className="h-20 w-20">
              <AvatarImage src="" alt="Profile picture" />
              <AvatarFallback className="text-lg">🧑🏾‍🎓</AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <CardTitle className="text-2xl">Jordan Reed</CardTitle>
              <p className="mt-1 text-sm text-muted-foreground">@jordan.reed</p>
              <div className="mt-3 flex flex-wrap gap-2">
                <span className="inline-flex items-center gap-1 rounded-full bg-secondary px-3 py-1 text-xs font-medium text-primary">
                  <Award className="h-3.5 w-3.5" /> Tier: Gold
                </span>
                <span className="inline-flex items-center gap-1 rounded-full bg-secondary px-3 py-1 text-xs font-medium text-primary">
                  <Star className="h-3.5 w-3.5" /> 1,250 pts
                </span>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Sustainability-focused design student sharing practical dorm and study gear.
            </p>
          </CardContent>
        </Card>

        <Card className="mt-8">
          <CardHeader>
            <CardTitle className="text-xl">Post New Item</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <input className="h-10 rounded-md border border-border bg-background px-3 text-sm" placeholder="Item name" />
            <input className="h-10 rounded-md border border-border bg-background px-3 text-sm" placeholder="Category" />
            <input className="h-10 rounded-md border border-border bg-background px-3 text-sm" placeholder="Condition" />
            <input className="h-10 rounded-md border border-border bg-background px-3 text-sm" placeholder="Points value" />
            <button
              type="button"
              className="inline-flex h-10 items-center justify-center rounded-md bg-slate-900 px-4 text-sm font-medium text-white transition-colors hover:bg-slate-800 md:col-span-2"
            >
              Post Item
            </button>
          </CardContent>
        </Card>

        <Card className="mt-8">
          <CardHeader>
            <CardTitle className="text-xl">Request Notifications</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-foreground">
            <p>📩 Maya Patel requested your Mini Fridge</p>
            <p>📩 Ethan Brooks requested your Storage Bin Set</p>
          </CardContent>
        </Card>

        <h2 className="mt-12 text-2xl font-semibold tracking-tight text-foreground">
          My Items — Sharing
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">Items you've shared with campus.</p>
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <ItemCard title="Mini Fridge" by="Jordan Reed" requested />
          <ItemCard title="Storage Bin Set" by="Jordan Reed" requested />
          <ItemCard title="Desk Organizer" by="Jordan Reed" />
        </div>

        <h2 className="mt-12 text-2xl font-semibold tracking-tight text-foreground">
          My Items — Borrowing
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">Items you're currently borrowing.</p>
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <ItemCard title="Calculus Textbook" by="Maya Patel" />
          <ItemCard title="Desk Lamp" by="Nia Johnson" />
        </div>
        </section>
      </AuthGate>
    </PageLayout>
  );
}
