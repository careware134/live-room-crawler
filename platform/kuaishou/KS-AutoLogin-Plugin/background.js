// background.js

// Declare a global variable to store the cookie
let websocketCookie = '';

console.log("ks auto login plugin->background>background.......========")
// Listen for completed HTTP requests
chrome.webRequest.onBeforeSendHeaders.addListener(
    function(details) {
        const urlPatterns = [
            'https://live.kuaishou.com/live_api/baseuser/userinfo',
            'https://live.kuaishou.com/live_api/liveroom/websocketinfo'
        ];

        //if (details.url.includes(urlPattern)) {
        if (checkURL(details.url, urlPatterns)) {
            getAllCookiesOfSite('live.kuaishou.com')
        }
    },
    {
        //urls: ['<all_urls>']
        urls: ['*://live.kuaishou.com/*'],
        types: ['xmlhttprequest']
    },
    [ "requestHeaders"]
);

function checkURL(url, patterns) {
    for (let i = 0; i < patterns.length; i++) {
        if (url.includes(patterns[i])) {
            console.log("ks auto login plugin->background::checkURL for url:"+ url + ", match pattern:" + patterns[i])
            return true;
        }
    }
    return false;
}


function getAllCookiesOfSite(site) {
    chrome.cookies.getAll({ domain: site }, function (cookies) {
        websocketCookie = cookies.map(cookie => `${cookie.name}=${cookie.value};`).join(' ');
        console.log("ks auto login plugin->websocketCookie is :" + websocketCookie)
    });
}

// Listen for messages from the content script
// chrome.runtime.onMessage.addListener(
//     function(request, sender, sendResponse) {
//         if (request.action === 'getCookie') {
//             sendResponse(websocketCookie);
//         }
//     }
// );


// background.js
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    if (request.action === 'getCookies') {
        sendResponse({ websocketCookie });
        return true;
    }
});