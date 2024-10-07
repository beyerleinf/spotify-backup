/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/ui/public/assets/css/**/*.css",
    "./internal/ui/public/**/*.html",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("@catppuccin/tailwindcss")({
      defaultFlavour: "mocha",
    }),
  ],
};
