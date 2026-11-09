import Link from "next/link";

export default function NotFound() {
  return (
    <main className="mx-auto grid min-h-[60vh] max-w-7xl content-center gap-3 px-4 py-16 sm:px-6 lg:px-8">
      <h1 className="m-0 text-3xl font-semibold tracking-tight">Shipment not found</h1>
      <p className="m-0 text-muted-foreground">404 · Shipment not found.</p>
      <Link
        className="inline-flex min-h-11 w-fit items-center text-primary underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        href="/orders"
      >
        ← Back to orders
      </Link>
    </main>
  );
}
