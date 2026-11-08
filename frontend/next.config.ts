import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  typescript: { tsconfigPath: "tsconfig.next.json" },
  env: { NEXT_PUBLIC_CARTO_API_KEY: process.env.VITE_CARTO_API_KEY ?? "" },
};

export default nextConfig;
