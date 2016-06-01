define([
    "dijit/Dialog",
    "dojo/dom-construct",
    'dojo/domReady!'
], function(
    Dialog, domConstruct
) {
    return {
        showDialog: function(params) {
            var content = "<div>It works in :</div>";
            if (params.buildMode) {
                content += "<h1>BUILD mode ;-)</h1>";
            } else {
                content += "<h1>non-built mode ;-)</h1>";
            }
            new Dialog({
                title: "YEAH !!!",
                content: domConstruct.toDom(content)
            }).show();
        }
    };
});