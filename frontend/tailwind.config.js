/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                background: '#0f172a', // Slate 900 (Deep Blue/Blackish)
                primary: '#f97316',   // Orange 500
                secondary: '#4f46e5', // Indigo 600
                accent: '#8b5cf6',    // Violet 500
                surface: '#1e293b',   // Slate 800
            },
            fontFamily: {
                sans: ['Inter', 'sans-serif'],
            },
        },
    },
    plugins: [],
}
