const codesToAllowSwap = [400, 401, 403, 404, 409, 500];

document.body.addEventListener('htmx:beforeSwap', function(evt) {
    if (codesToAllowSwap.includes(evt.detail.xhr.status)) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }
});