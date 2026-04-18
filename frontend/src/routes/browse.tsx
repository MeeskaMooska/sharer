import { createFileRoute } from "@tanstack/react-router";
import { useMemo, useState } from "react";
import { PageLayout } from "@/components/site/PageLayout";
import { AuthGate } from "@/components/site/AuthGate";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Trophy } from "lucide-react";

export const Route = createFileRoute("/browse")({
  head: () => ({
    meta: [
      { title: "Browse Items — Sharer" },
      { name: "description", content: "Browse school supplies, dorm essentials, and more." },
    ],
  }),
  component: BrowsePage,
});

type SectionName = "School Supplies" | "Dorm Essentials" | "Miscellaneous";
type TierName = "Bronze" | "Silver" | "Gold";

type ItemRecord = {
  id: number;
  name: string;
  section: SectionName;
  publishedOn: string;
  publishedByTier: TierName;
  emoji: string;
  publishedBy: string;
  publisherEmoji: string;
  publisherMajor: string;
};

const allItems: ItemRecord[] = [
  { id: 1, name: "Calculus Textbook", section: "School Supplies", publishedOn: "2026-04-11", publishedByTier: "Gold", emoji: "📘", publishedBy: "Maya Patel", publisherEmoji: "👩🏽‍🎓", publisherMajor: "Math" },
  { id: 2, name: "Physics Lab Kit", section: "School Supplies", publishedOn: "2026-03-30", publishedByTier: "Silver", emoji: "🧪", publishedBy: "Ethan Brooks", publisherEmoji: "🧑🏻‍🔬", publisherMajor: "Physics" },
  { id: 3, name: "Desk Lamp", section: "Dorm Essentials", publishedOn: "2026-04-15", publishedByTier: "Bronze", emoji: "💡", publishedBy: "Nia Johnson", publisherEmoji: "👩🏾‍💻", publisherMajor: "CS" },
  { id: 4, name: "Mini Fridge", section: "Dorm Essentials", publishedOn: "2026-04-02", publishedByTier: "Gold", emoji: "🧊", publishedBy: "Lucas Kim", publisherEmoji: "🧑🏽‍🎓", publisherMajor: "Business" },
  { id: 5, name: "Whiteboard Marker Pack", section: "Miscellaneous", publishedOn: "2026-03-28", publishedByTier: "Silver", emoji: "🖍️", publishedBy: "Avery Chen", publisherEmoji: "👩🏻‍🏫", publisherMajor: "Education" },
  { id: 6, name: "Storage Bin Set", section: "Miscellaneous", publishedOn: "2026-04-10", publishedByTier: "Gold", emoji: "📦", publishedBy: "Jordan Reed", publisherEmoji: "🧑🏾‍🎨", publisherMajor: "Design" },
];

function ItemCard({
  item,
  requested,
  onRequest,
}: {
  item: ItemRecord;
  requested: boolean;
  onRequest: () => void;
}) {
  const sectionStyle =
    item.section === "School Supplies"
      ? "bg-blue-50 border-blue-200"
      : item.section === "Dorm Essentials"
        ? "bg-orange-50 border-orange-200"
        : "bg-purple-50 border-purple-200";

  return (
    <Card className={`border-2 ${sectionStyle}`}>
      <CardHeader className="space-y-2">
        <CardTitle className="text-base">
          <span className="mr-2 text-3xl leading-none align-middle" aria-hidden="true">
            {item.emoji}
          </span>
          {item.name}
        </CardTitle>
        <p className="text-sm text-slate-700">
          From {item.publisherEmoji} {item.publishedBy} ({item.publisherMajor})
        </p>
      </CardHeader>
      <CardContent className="space-y-2 text-sm">
        <p className="text-foreground">
          <span className="font-medium">Section:</span> {item.section}
        </p>
        <p className="text-foreground">
          <span className="font-medium">Date Published:</span> {item.publishedOn}
        </p>
        <p className="text-foreground">
          <span className="font-medium">Published By Tier:</span> {item.publishedByTier}
        </p>
        <button
          type="button"
          onClick={onRequest}
          className="mt-2 inline-flex h-9 items-center justify-center rounded-md bg-slate-900 px-3 text-sm font-medium text-white transition-colors hover:bg-slate-800"
        >
          {requested ? "Requested" : "Request Item"}
        </button>
      </CardContent>
    </Card>
  );
}

export function BrowsePage() {
  const [selectedSection, setSelectedSection] = useState<SectionName | "All">("All");
  const [keyword, setKeyword] = useState("");
  const [publishedOnOrAfter, setPublishedOnOrAfter] = useState("");
  const [showHighTierOnly, setShowHighTierOnly] = useState(false);
  const [sortByDate, setSortByDate] = useState<"newest" | "oldest">("newest");
  const [requestedItemIds, setRequestedItemIds] = useState<number[]>([]);

  function handleRequest(itemId: number) {
    setRequestedItemIds((previous) =>
      previous.includes(itemId) ? previous : [...previous, itemId],
    );
  }

  const filteredItems = useMemo(() => {
    let items = allItems.filter((item) => {
      const sectionMatch = selectedSection === "All" || item.section === selectedSection;
      const keywordMatch = item.name.toLowerCase().includes(keyword.trim().toLowerCase());
      const dateMatch = !publishedOnOrAfter || item.publishedOn >= publishedOnOrAfter;
      const tierMatch = !showHighTierOnly || item.publishedByTier === "Gold";
      return sectionMatch && keywordMatch && dateMatch && tierMatch;
    });

    items = items.sort((a, b) =>
      sortByDate === "newest"
        ? b.publishedOn.localeCompare(a.publishedOn)
        : a.publishedOn.localeCompare(b.publishedOn),
    );

    return items;
  }, [selectedSection, keyword, publishedOnOrAfter, showHighTierOnly, sortByDate]);

  return (
    <PageLayout>
      <AuthGate>
        <section className="mx-auto max-w-6xl px-4 py-12">
        <h1 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">
          Browse Items
        </h1>
        <p className="mt-3 max-w-2xl text-muted-foreground">
          Find campus items by section, keyword, publish date, and member tier.
        </p>

        <div className="mt-8 rounded-xl border border-border bg-secondary/30 p-4 md:p-5">
          <h2 className="text-base font-semibold text-foreground">Filters</h2>
          <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-5">
            <label className="flex flex-col gap-2 text-sm">
              <span className="font-medium text-foreground">Section</span>
              <select
                className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
                value={selectedSection}
                onChange={(event) => setSelectedSection(event.target.value as SectionName | "All")}
              >
                <option value="All">All Sections</option>
                <option value="School Supplies">School Supplies</option>
                <option value="Dorm Essentials">Dorm Essentials</option>
                <option value="Miscellaneous">Miscellaneous</option>
              </select>
            </label>

            <label className="flex flex-col gap-2 text-sm lg:col-span-2">
              <span className="font-medium text-foreground">Search by key name</span>
              <input
                className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
                type="text"
                placeholder="Search item name..."
                value={keyword}
                onChange={(event) => setKeyword(event.target.value)}
              />
            </label>

            <label className="flex flex-col gap-2 text-sm">
              <span className="font-medium text-foreground">Date published</span>
              <input
                className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
                type="date"
                value={publishedOnOrAfter}
                onChange={(event) => setPublishedOnOrAfter(event.target.value)}
              />
            </label>

            <label className="flex flex-col gap-2 text-sm">
              <span className="font-medium text-foreground">Date order</span>
              <select
                className="h-10 rounded-md border border-border bg-background px-3 text-foreground"
                value={sortByDate}
                onChange={(event) => setSortByDate(event.target.value as "newest" | "oldest")}
              >
                <option value="newest">Newest first</option>
                <option value="oldest">Oldest first</option>
              </select>
            </label>
          </div>

          <label className="mt-4 inline-flex items-center gap-2 text-sm text-foreground">
            <input
              type="checkbox"
              checked={showHighTierOnly}
              onChange={(event) => setShowHighTierOnly(event.target.checked)}
            />
            Published by high-tier members only (Gold)
          </label>
        </div>

        <div className="mt-8 flex items-center gap-2 text-sm text-muted-foreground">
          <Trophy className="h-4 w-4 text-primary" />
          {filteredItems.length} item(s) match the current filters
        </div>

        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filteredItems.length > 0 ? (
            filteredItems.map((item) => (
              <ItemCard
                key={item.id}
                item={item}
                requested={requestedItemIds.includes(item.id)}
                onRequest={() => handleRequest(item.id)}
              />
            ))
          ) : (
            <div className="col-span-full rounded-xl border border-dashed border-border bg-secondary/30 p-8 text-center text-sm text-muted-foreground">
              No items found for these filters.
            </div>
          )}
        </div>
        </section>
      </AuthGate>
    </PageLayout>
  );
}
