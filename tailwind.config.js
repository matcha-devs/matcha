/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./internal/templates/**/*.{html, go.html}"],
    theme: {
        extend: {},
        colors: {
            primary: '#505c45', //600
            secondary: '#d4f3b7', //000
            accent: {
                100: '#bfdaa4',
                200: '#a9c191',
                300: '#93a87e',
                400: '#7b8d6a',
                500: '#667558',
                700: '#384031',
            },
            white: '#fff',
            black: '#000',
            gray: {
                50: '#f8fafc',
                100: '#f1f5f9',
                200: '#e5e7eb',
                300: '#d1d5db',
                400: '#9ca3af',
                500: '#6b7280',
                600: '#4b5563',
                700: '#374151'
            },
            blue: {
                50: '#eff6ff',
                100: '#dbeafe',
                200: '#93c5fd',
                300: '#93c5fd',
                400: '#60a5fa',
                500: '#3b82f6',
                600: '#2563eb',
                700: '#1d4ed8',
                800: '#1e40af',
                900: '#1e3a8a',
                950: '#172554'
            },
            red: {
                500: '#ef4444',
                600: '#dc2626'
            }
        }
    },
    plugins: [],
}
