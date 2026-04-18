import { createFileRoute } from "@tanstack/react-router";
import { AuthGate } from "@/components/site/AuthGate";
import { PageLayout } from "@/components/site/PageLayout";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export const Route = createFileRoute("/users")({
  head: () => ({
    meta: [
      { title: "Browse Users — Sharer" },
      { name: "description", content: "Browse campus users and their tier information." },
    ],
  }),
  component: UsersPage,
});

const users = [
  { name: "Maya Patel", handle: "@maya.patel", emoji: "👩🏽‍🎓", tier: "Gold", major: "Math", points: 1520 },
  { name: "Ethan Brooks", handle: "@ethan.brooks", emoji: "🧑🏻‍🔬", tier: "Silver", major: "Physics", points: 980 },
  { name: "Nia Johnson", handle: "@nia.johnson", emoji: "👩🏾‍💻", tier: "Bronze", major: "Computer Science", points: 640 },
  { name: "Lucas Kim", handle: "@lucas.kim", emoji: "🧑🏽‍🎓", tier: "Gold", major: "Business", points: 1310 },
  { name: "Avery Chen", handle: "@avery.chen", emoji: "👩🏻‍🏫", tier: "Silver", major: "Education", points: 870 },
  { name: "Jordan Reed", handle: "@jordan.reed", emoji: "🧑🏾‍🎨", tier: "Gold", major: "Design", points: 1250 },
];

function UsersPage() {
  return (
    <PageLayout>
      <AuthGate>
        <section className="mx-auto max-w-6xl px-4 py-12">
          <h1 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">Browse Users</h1>
          <p className="mt-3 max-w-2xl text-muted-foreground">
            Explore student profiles and sustainability tiers across campus.
          </p>

          <div className="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {users.map((user) => (
              <Card key={user.handle}>
                <CardHeader>
                  <CardTitle className="text-lg">
                    <span className="mr-2 text-2xl" aria-hidden="true">
                      {user.emoji}
                    </span>
                    {user.name}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-1 text-sm">
                  <p className="text-muted-foreground">{user.handle}</p>
                  <p><span className="font-medium">Tier:</span> {user.tier}</p>
                  <p><span className="font-medium">Major:</span> {user.major}</p>
                  <p><span className="font-medium">Points:</span> {user.points}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
      </AuthGate>
    </PageLayout>
  );
}
