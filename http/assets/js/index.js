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

// Clones the auth element and shows it.
function addAuth() {
    const elem = $('#auth').clone();
    elem.removeAttr('id hidden');
    elem.find('select').attr('name', 'auth[method]');
    elem.find('textarea').attr('name', 'auth[values]');
    elem.find('input').click(function() {
        elem.remove();
    });
    $('.auths').append(elem);
}

// Enables the disabled inputs in the form and submit.
function enableAndSubmit() {
    const elem = $('form.link');
    elem.find('select option').each(function() {
        $(this).prop('disabled', false);
    });
    elem.submit();
}