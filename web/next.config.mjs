import { fileURLToPath } from "node:url";

// English text：`NEXT_EXPORT=1 next build` English text web/out，English text
// nginx web English text。English text(next dev)English text，English text /api English text。
const isExport = process.env.NEXT_EXPORT === "1";
// Vercel demo：English text mock，English text，English text /api English text。
const isMock = process.env.NEXT_PUBLIC_MOCK === "1";

/** @type {import('next').NextConfig} */
const nextConfig = {
  // English text lockfile English text。
  turbopack: { root: fileURLToPath(new URL(".", import.meta.url)) },
  reactCompiler: true,
  // English text IP English text dev English text（HMR），English text。
  // dev English text IPv4 English text /_next/* English text HMR（English text IP English text）。
  // English text：Next English text "*"，English text；"*.*.*.*" English text IPv4。
  allowedDevOrigins: ["*.*.*.*"],
  compiler: {
    removeConsole: process.env.NODE_ENV === "production",
  },
  ...(isExport
    ? {
        // English text：English text Node English text；English text；English text <route>/index.html。
        output: "export",
        images: { unoptimized: true },
        trailingSlash: true,
      }
    : isMock
      ? {
          // Vercel mock demo：English text，English text /api English text。
          images: { unoptimized: true },
        }
      : {
          // English text：English text /api/* English text Go English text（English text :8787，English text AUTOPENTEST_API English text）。
          async rewrites() {
            const backend = process.env.AUTOPENTEST_API ?? "http://localhost:8787";
            return [{ source: "/api/:path*", destination: `${backend}/api/:path*` }];
          },
        }),
};

export default nextConfig;
