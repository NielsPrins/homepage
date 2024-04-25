/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./**/*.gohtml"],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/forms'),
    ],
}

