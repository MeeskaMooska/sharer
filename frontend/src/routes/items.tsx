import { createFileRoute } from "@tanstack/react-router";
import { BrowsePage } from "./browse";

export const Route = createFileRoute("/items")({
  head: () => ({
    meta: [
      { title: "Browse Items — Sharer" },
      { name: "description", content: "Browse school supplies, dorm essentials, and more." },
    ],
  }),
  component: BrowsePage,
});
