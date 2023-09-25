//popup.js
// Copy the cookie when the "Copy Cookie" button in the popup is clicked
document.getElementById('copyCookie').addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: 'copyCookie' });
});