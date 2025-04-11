/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  darkMode: "selector",
  safelist: [],
  corePlugins: {
    preflight: false,
  },
  plugins: [require("daisyui")],
  daisyui: {
    // Optional: DaisyUI configuration options
    themes: ["light", "dark"], // You can customize the themes here
  },
};
