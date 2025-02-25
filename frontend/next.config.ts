import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    async rewrites() {
        return [
            {
                source: "/api/v1/expressions",
                destination: process.env.BACKEND_URL + "/api/v1/expressions",
                basePath: false,
            },
            {
                source: "/api/v1/expressions/:id",
                destination: process.env.BACKEND_URL + "/api/v1/expressions/:id",
                basePath: false
            },
            {
                source: "/api/v1/calculate",
                destination: process.env.BACKEND_URL + "/api/v1/calculate",
                basePath: false
            }
        ]
    }
};

export default nextConfig;
