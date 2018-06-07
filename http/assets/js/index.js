// Copies the given text to the clipboard and updates the given element to show
// feedback.
function copy(elem, text) {
    // Copy the text
    const tmp = $("<input>");
    $("body").append(tmp);
    tmp.val(text).select();
    document.execCommand("copy");
    tmp.remove();

    // Update the element and reset after a while.
    const text = elem.text();
    const backgroundColor = elem.css('background-color');

    elem.text('COPIED').css({
        'border-color': '#3c3',
        'background-color': '#3c3',
    }).prop('disabled', true);

    setTimeout(function() {
        elem.text(text).css({
            'border-color': backgroundColor,
            'background-color': backgroundColor,
        }).prop('disabled', false);
    }, 3000);
}