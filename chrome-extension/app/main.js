// Copyright 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

var port = null;

var getKeys = function(obj){
    var keys = [];
    for(var key in obj){
        keys.push(key);
    }
    return keys;
}

function appendMessage(text) {
    document.getElementById('response').innerHTML += "<p>" + text + "</p>";
}

function updateUiState() {
    if (port) {
        document.getElementById('connect-button').style.display = 'none';
        document.getElementById('input-text').style.display = 'inline-block';
        document.getElementById('send-message-button').style.display = 'inline-block';
    } else {
        document.getElementById('connect-button').style.display = 'block';
        document.getElementById('input-text').style.display = 'none';
        document.getElementById('send-message-button').style.display = 'none';
    }
}


function sendNativeMessage(query) {
  // message = {"query": document.getElementById('input-text').value.trim()};
    message = {"query": query.trim()};
    port.postMessage(message);
    appendMessage("Sent message: <b>" + JSON.stringify(message) + "</b>");
}

function onNativeMessage(message) {
    appendMessage(message.query);
}

function onDisconnected() {
    appendMessage("Failed to connect: " + chrome.runtime.lastError.message);
    port = null;
    updateUiState();
}

function connect() {
    var hostName = "com.sample.native_msg_golang";
    port = chrome.runtime.connectNative(hostName);
    port.onMessage.addListener(onNativeMessage);
    port.onDisconnect.addListener(onDisconnected);
    appendMessage("Connected to native messaging host <b>" + hostName + "</b>")
    updateUiState();
    chrome.webRequest.onClo
    chrome.webRequest.onResponseStarted.addListener(function(details) {
            sendNativeMessage(details.url)}
        ,  {
            urls: [
                "https://pornhub.org/view_video.php*",
                "https://beeg.com/-*",
                "https://www.pornhub.org/view_video.php*",
                "https://xhamster.com/videos/*",
                "https://www.redtube.com/*",
            ],


        }
    )
}

document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('connect-button').addEventListener('click', connect);
    updateUiState();
});

