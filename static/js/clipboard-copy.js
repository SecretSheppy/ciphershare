'use strict'
$(() => {
    $('#copy-to-clipboard').on( 'click', ()=> {
        putTextInClipboard($('#link-to-copy').attr('href'));
    });

});

/**
 * Puts the text into the users clipboard
 * @param valueToCopy the string to put into the clipboard
 */
function putTextInClipboard(valueToCopy) {
    navigator.clipboard.writeText(valueToCopy)
        .then(r => console.log('Text copied to clipboard:', valueToCopy))
        .catch(e => console.log('failed to copy text:', e));
}