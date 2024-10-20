/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: "/admin",
  assetPrefix: "/admin",
  rewrites() {
    return [{ source: "/admin/_next/:path*", destination: "/_next/:path*" }];
  },
};

export default nextConfig;
