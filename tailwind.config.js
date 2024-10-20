/** @type {import('tailwindcss').Config} */
import customEmeraldTheme from "./config/custom-emerald-theme.json";

export const content = ["./templates/**/*.html"];
export const theme = {
  extend: {
    fontFamily: {
      poppins: ["Poppins", "sans-serif"],
    },
  },
};
export const plugins = [require("@tailwindcss/typography"), require("daisyui")];
export const daisyui = {
  themes: [customEmeraldTheme],
};
