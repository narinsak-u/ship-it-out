import { Link } from "@tanstack/react-router";
import { Package } from "lucide-react";

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
        <Link to="/" className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-md bg-gradient-accent text-primary-foreground shadow-glow">
            <Package className="h-4 w-4" />
          </div>
          <span className="font-mono text-sm font-semibold tracking-tight">HARBOR/OPS</span>
        </Link>
        <nav className="flex items-center gap-1 font-mono text-sm">
          <Link
            to="/"
            activeOptions={{ exact: true }}
            activeProps={{ className: "bg-secondary text-foreground" }}
            inactiveProps={{ className: "text-muted-foreground hover:text-foreground" }}
            className="rounded-md px-3 py-1.5 transition-colors"
          >
            Home
          </Link>
          <Link
            to="/orders"
            activeProps={{ className: "bg-secondary text-foreground" }}
            inactiveProps={{ className: "text-muted-foreground hover:text-foreground" }}
            className="rounded-md px-3 py-1.5 transition-colors"
          >
            Orders
          </Link>
        </nav>
      </div>
    </header>
  );
}
