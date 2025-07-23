module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}", "./public/index.html"],
  theme: {
    extend: {
      animation: {
        gradient: "gradient 15s ease infinite",
        float: "float 6s ease-in-out infinite",
        "pulse-glow": "pulse-glow 2s ease-in-out infinite alternate",
      },
      keyframes: {
        gradient: { /* tu definición */ },
        float: { /* tu definición */ },
        "pulse-glow": { /* tu definición */ },
      },
    },
  },
  plugins: [],
};