export function Footer() {
  return (
    <footer className="mt-16 border-t border-border bg-secondary/40">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-4 py-8 text-sm text-muted-foreground md:flex-row">
        <p>© {new Date().getFullYear()} Sharer — Built for a greener campus.</p>
        <nav className="flex gap-6">
          <a href="#" className="hover:text-primary">About</a>
          <a href="#" className="hover:text-primary">Mission</a>
          <a href="#" className="hover:text-primary">Contact</a>
        </nav>
      </div>
    </footer>
  );
}
