/** @type {import('next').NextConfig} */
const nextConfig = {
    async redirects() {
        return [
            {
                source: "/",
                destination: "/splash",
                permanent: true
            }
        ]
    }

};

export default nextConfig;
