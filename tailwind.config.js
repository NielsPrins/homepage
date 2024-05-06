/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./**/*.gohtml"],
    theme: {
        extend: {
            gridTemplateColumns: {
                'auto-fit-minmax-8rem': 'repeat(auto-fit, minmax(8rem, 1fr))',
            },
        },
    },
    plugins: [
        require('@tailwindcss/forms'),
    ],
}

