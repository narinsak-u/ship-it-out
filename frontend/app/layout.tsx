import "./globals.css";

export const metadata = { title: "thun-u-der/express", description: "Shipment tracking" };

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <header className="site-header">
          <a className="brand" href="/">
            <span className="brand-mark" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7">
                <path d="m3.5 7 8.5-4 8.5 4v10L12 21l-8.5-4V7Z" />
                <path d="m3.5 7 8.5 4 8.5-4M12 21V11m-4-6 8.5 4" />
              </svg>
            </span>{" "}
            thun-u-der/express
          </a>
          <nav>
            <a href="/">Home</a>
            <a className="active" href="/orders">
              Orders
            </a>
            <a href="/carriers">Carriers</a>
            <span className="signin">Sign in</span>
          </nav>
        </header>
        {children}
        <footer>© 2026 Thun-u-der Express · Logistics control plane</footer>
      </body>
    </html>
  );
}
