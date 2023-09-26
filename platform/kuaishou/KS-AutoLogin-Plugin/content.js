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
        copyCookieButton.innerText = '复制Cookie';
        copyCookieButton.addEventListener('click', copyCookieToClipboard);
        userInfoElement.parentNode.insertBefore(copyCookieButton, userInfoElement.nextSibling);
    }
}

// Copy the cookie variable value to clipboard
function copyCookieToClipboard() {
    chrome.runtime.sendMessage({ action: 'getCookies' }, function(response) {
        //window.location.reload()
        var websocketCookie = response ? response.websocketCookie : null;
        console.log("ks autologin plugin:: getCookies response is :" + websocketCookie)
        navigator.clipboard.writeText(websocketCookie)
            .then(() => {
                var message = '\n已复制Cookie，请去如影软件中粘贴！';
                //var styledMessage = "<span style='font-size: 20px; font-weight: bold; font-family: Arial, sans-serif;'>" + message + "</span>";
                if (websocketCookie) {
                    alert(message + '\n\nCookie详情:\n' + response.websocketCookie)
                } else {
                    alert("未获取到Cookie！请尝试重新登录或刷新当前页面")
                }
                //alertMessage(styledMessage + '<br>Cookie:<br><br>' + response.websocketCookie)
            })
            .catch(error => {
                alert("获取到Cookie失败！请尝试重新登录或刷新当前页面")
                console.error('ks autologin plugin:: Failed to copy cookies:', error)
            });
    });

}


// Run the automation tasks when the page finishes loading
window.addEventListener('load', () => {
    console.log("ks autologin plugin:: load...")
    autoClickLogin();
    insertCopyCookieButton();
});