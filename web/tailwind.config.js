/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: "#1978e5",
                "brand-dark": "#0e141b",
                "brand-muted": "#4e7097",
                "brand-bg": "#f6f7f8",
                "border-gray": "#d0dbe7",
                // MD3 design tokens (login page)
                "md-primary": "#006971",
                "md-secondary": "#4646d8",
                "md-surface": "#f8fafb",
                "md-on-surface": "#191c1d",
                "md-on-surface-variant": "#3c494b",
                "md-surface-container-low": "#f2f4f5",
                "md-surface-container-highest": "#e1e3e4",
                "md-outline-variant": "#bbc9cb",
                "md-primary-container": "#3abbc9",
                "md-on-primary-container": "#00474d",
            },
            fontFamily: {
                sans: ["Inter", "sans-serif"],
                headline: ["Space Grotesk", "Inter", "sans-serif"],
            }
        }
    },
    plugins: [],
    // 重点：解决 Tailwind 与 Ant Design 样式冲突（预设样式重置问题）
    corePlugins: {
        preflight: false,
    }
}