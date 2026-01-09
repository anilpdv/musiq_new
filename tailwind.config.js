/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/**/*.templ"],
  theme: {
    extend: {
      colors: {
        'neo-bg': '#fffef0',
        'neo-card': '#ffffff',
        'neo-border': '#000000',
        'neo-shadow': '#000000',
        'neo-blue': '#3b82f6',
        'neo-red': '#ef4444',
        'neo-green': '#22c55e',
        'neo-yellow': '#facc15',
        'neo-purple': '#a855f7',
        'neo-orange': '#f97316',
        'neo-pink': '#ec4899',
        'neo-cyan': '#06b6d4',
      },
      boxShadow: {
        'neo': '4px 4px 0px 0px #000',
        'neo-sm': '2px 2px 0px 0px #000',
        'neo-lg': '6px 6px 0px 0px #000',
        'neo-blue': '4px 4px 0px 0px #3b82f6',
        'neo-red': '4px 4px 0px 0px #ef4444',
        'neo-green': '4px 4px 0px 0px #22c55e',
      },
      borderWidth: {
        '3': '3px',
      },
      fontFamily: {
        'heading': ['Space Grotesk', 'system-ui', 'sans-serif'],
        'body': ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
