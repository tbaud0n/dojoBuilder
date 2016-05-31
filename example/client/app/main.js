define("app/main", [
    "dijit/Dialog",
    "geonef/jig/util/makeDOM",
    'dojo/domReady!'
], function(
    Dialog, makeDOM
) {
    new Dialog({
        style: "min-width:120px;",
        title: "YEAH !!!",
        content: makeDOM(['div', {}, "It works with geonef"])
    }).show();
});