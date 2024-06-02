/** @type {import('next').NextConfig} */
const nextConfig = {
    async redirects() {
        return [
            {
                source: "/",
                destination: "/signup",
                permanent: true
            }
        ]
    }

};

export default nextConfig;
