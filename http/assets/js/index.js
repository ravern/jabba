$(document).ready(function() {
    clipboard();
});

// Configures the copy buttons.
function clipboard() {
    const clipboard = new ClipboardJS('.copy');

    clipboard.on('success', function(e) {
        // Update the element and reset after a while.
        const elem = $(e.trigger);
        const text = elem.text();
        elem.text('COPIED').prop('disabled', true).addClass('copied');
        setTimeout(function() {
            elem.text(text).prop('disabled', false).removeClass('copied');
        }, 3000);
    })
}