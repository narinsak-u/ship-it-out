import Link from "next/link";
import "./globals.css";

import AuthControls from "@/features/auth/AuthControls";

export const metadata = { title: "thun-u-der/express", description: "Shipment tracking" };

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className="dark">
      <body className="min-h-screen bg-background text-foreground antialiased">
        <header className="sticky top-0 z-40 border-b border-border bg-background/90 backdrop-blur">
          <div className="mx-auto flex min-h-16 max-w-7xl flex-wrap items-center justify-between gap-x-6 gap-y-2 px-4 py-2 sm:px-6 lg:px-8">
            <Link
              href="/"
              className="flex items-center gap-2 font-mono text-sm font-semibold tracking-wide"
            >
              <span
                className="grid size-8 shrink-0 place-items-center rounded-md bg-primary text-primary-foreground"
                aria-hidden="true"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7">
                  <path d="m3.5 7 8.5-4 8.5 4v10L12 21l-8.5-4V7Z" />
                  <path d="m3.5 7 8.5 4 8.5-4M12 21V11m-4-6 8.5 4" />
                </svg>
              </span>
              <span>thun-u-der/express</span>
            </Link>
            <AuthControls />
          </div>
        </header>
        {children}
        <footer className="mx-auto max-w-7xl px-4 py-8 text-sm text-muted-foreground sm:px-6 lg:px-8">
          © 2026 Thun-u-der Express · Logistics control plane
        </footer>
      </body>
    </html>
  );
}
