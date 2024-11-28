/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/static/css/**/*.css", "./web/templates/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [
    require("@catppuccin/tailwindcss")({
      defaultFlavour: "mocha",
    }),
  ],
};
