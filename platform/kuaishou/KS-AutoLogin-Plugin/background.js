// background.js

// Declare a global variable to store the cookie
let websocketCookie = '';

console.log("background.......========")
// Listen for completed HTTP requests
chrome.webRequest.onBeforeSendHeaders.addListener(
    function(details) {
        const urlPattern = 'https://live.kuaishou.com/live_api/liveroom/websocketinfo';

        if (details.url.includes(urlPattern)) {
            // Extract the cookie from the request headers
            // const requestHeaders = details.requestHeaders || [];
            // const cookieHeader = requestHeaders.find(header => header.name.toLowerCase() === 'cookie');
            //
            // if (cookieHeader) {
            //     websocketCookie = cookieHeader.value;
            // }

            chrome.cookies.getAll({ url: "https://live.kuaishou.com" }, function(cookies) {
                websocketCookie = cookies.map(cookie => `${cookie.name}=${cookie.value};`).join(' ');
                console.log("websocketCookie is :" + websocketCookie)
            });
        }
    },
    {
        urls: ['<all_urls>']
        //types: ['xmlhttprequest']
    },
    [ "requestHeaders"]
);

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