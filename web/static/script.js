const form = document.getElementById('shorten-form');
const urlInput = document.getElementById('url-input');
const resultDiv = document.getElementById('result');
const shortUrlLink = document.getElementById('short-url');

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const originalUrl = urlInput.value;

    try {
        const response = await fetch('/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: originalUrl }),
        });

        if (!response.ok) {
            throw new Error('Failed to shorten URL');
        }

        const data = await response.json();
        shortUrlLink.href = data.short_url;
        shortUrlLink.textContent = data.short_url;
        resultDiv.classList.remove('hidden');
    } catch (error) {
        console.error(error);
        alert('An error occurred. Please try again.');
    }
});
