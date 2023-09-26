//popup.js
// Copy the cookie when the "Copy Cookie" button in the popup is clicked
document.getElementById('copyCookie').addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: 'copyCookie' });
});



function alertMessage(message) {
    const alertDiv = document.createElement("div");
    alertDiv.style.fontFamily = "Arial, sans-serif";
    alertDiv.style.fontSize = "16px";
    alertDiv.style.padding = "10px";
    alertDiv.style.backgroundColor = "lightblue";
    alertDiv.style.border = "1px solid white";
    alertDiv.style.borderRadius = "5px";
    alertDiv.style.position = "fixed";
    alertDiv.style.top = "20px";
    alertDiv.style.left = "50%";
    alertDiv.style.transform = "translateX(-50%)";
    alertDiv.style.width = "30%";
    alertDiv.style.height = "200px";
    alertDiv.style.overflowY = "auto";
    alertDiv.style.zIndex = "9999";

    const bodyDiv = document.createElement("div");
    bodyDiv.style.height = "80px";
    bodyDiv.style.overflowY = "auto";
    bodyDiv.style.marginBottom = "10px";
    alertDiv.appendChild(bodyDiv);

    const messageDiv = document.createElement("div");
    messageDiv.innerHTML = message;
    messageDiv.style.textAlign = "left";
    messageDiv.style.marginBottom = "10px";
    messageDiv.style.margin = "8px 10px";
    bodyDiv.appendChild(messageDiv);

    const footerDiv = document.createElement("div");
    footerDiv.style.position = "absolute";
    footerDiv.style.bottom = "10px";
    footerDiv.style.width = "100%";
    footerDiv.style.height = "70px";
    footerDiv.style.textAlign = "right";
    footerDiv.style.height = "20px";
    alertDiv.appendChild(footerDiv);

    const closeButton = document.createElement("button");
    closeButton.innerText = "Close";
    closeButton.style.padding = "5px 10px";
    closeButton.style.backgroundColor = "whte";
    closeButton.style.color = "white";
    closeButton.style.border = "none";
    closeButton.style.borderRadius = "3px";
    closeButton.style.cursor = "pointer";
    closeButton.addEventListener("click", function() {
        document.body.removeChild(alertDiv);
    });
    footerDiv.appendChild(closeButton);

    closeButton.addEventListener("click", function() {
        document.body.removeChild(alertDiv);
    });
    alertDiv.appendChild(closeButton);

    document.body.appendChild(alertDiv);
}
