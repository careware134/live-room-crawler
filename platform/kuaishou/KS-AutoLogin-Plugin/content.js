// content.js

// Auto click the element with "login" class
function autoClickLogin() {
    const loginElement = document.querySelector('.login');
    if (loginElement) {
        loginElement.click();
    }
}
// content.js

// Insert "Copy Cookie" button after the element with "user-info" class
function insertCopyCookieButton() {
    const userInfoElement = document.querySelector('.user-info');
    if (userInfoElement) {
        const copyCookieButton = document.createElement('button');
        copyCookieButton.innerText = 'Copy Cookie';
        copyCookieButton.addEventListener('click', copyCookieToClipboard);
        userInfoElement.parentNode.insertBefore(copyCookieButton, userInfoElement.nextSibling);
    }
}

// Copy the cookie variable value to clipboard
function copyCookieToClipboard() {
    chrome.runtime.sendMessage({ action: 'getCookies' }, function(response) {
        console.log("getCookies response is :" + response)
        navigator.clipboard.writeText(response.websocketCookie)
            .then(() => alert('Cookie已复制至剪贴板！请开始直播！\n' +
                'Cookie:\n\n' + response.websocketCookie))
            .catch(error => console.error('Failed to copy cookies:', error));
    });

    //chrome.storage.local.get({ 'cookieValue': websocketCookie };

    navigator.clipboard.writeText(cookieString)
    chrome.cookies
}


// Run the automation tasks when the page finishes loading
window.addEventListener('load', () => {
    console.log("load......=========")
    autoClickLogin();
    insertCopyCookieButton();
});