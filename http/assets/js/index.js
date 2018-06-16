$(document).ready(function () {
    clipboard();
});

// Configures the copy buttons.
function clipboard() {
    const clipboard = new ClipboardJS('.copy');

    clipboard.on('success', function (e) {
        // Reset the rest of the copy buttons.
        $('.copy').each(function() {
            $(this).text('COPY').prop('disabled', false).removeClass('copied');
        })

        // Update the element and reset after a while.
        const elem = $(e.trigger);
        elem.text('COPIED').prop('disabled', true).addClass('copied');
        setTimeout(function () {
            elem.text('COPY').prop('disabled', false).removeClass('copied');
        }, 3000);
    })
}

// Checks if the remote stylesheet is loaded, if not, load the local copy.
function load(remote, local) {
    $.each(document.styleSheets, function (i, sheet) {
        if (sheet.href == remote) {
            var rules = sheet.rules ? sheet.rules : sheet.cssRules;
            if (rules.length == 0) {
                document.write(`<link rel="stylesheet" href="${local}" />`);
            }
        }
    });
}