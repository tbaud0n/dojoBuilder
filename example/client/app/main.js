define([
    "dijit/Dialog",
    'dojo/domReady!'
], function(
    Dialog
) {
    new Dialog({
        title: "YEAH !!!",
        content: "<h1>It works ;-)</h1>"
    }).show();
});