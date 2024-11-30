import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  experimental: {
    appDir: true,  // Certifique-se de que a opção appDir está ativada
  },
};

export default nextConfig;
