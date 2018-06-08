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

// Checks if the remote is loaded, if not the local is loaded.
function load(remote, local) {
    $.each(document.styleSheets, function(i, sheet){
        if (sheet.href == remote) {
          var rules = sheet.rules ? sheet.rules : sheet.cssRules;
          if (rules.length == 0) {
            document.write(`<link rel="stylesheet" href="${local}" />`);
          }
        }
    });
}