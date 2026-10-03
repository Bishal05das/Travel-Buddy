import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Pin the workspace root to this folder so a lockfile elsewhere on the
  // machine (e.g. in the home directory) isn't picked up instead.
  turbopack: { root: __dirname },
};

export default nextConfig;
