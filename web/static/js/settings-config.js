var save, editor;
requirejs(['jquery', 'jsoneditor'], function ($, jsEditor) {

    function load() {
        if (files.data().indexOf(fileName) > -1) {
            console.log('Loading: ' + fileName);
            var template, values, AJAX = [];
            AJAX.push($.getJSON('/file/' + manifest.templates.configs.settings));
            AJAX.push($.getJSON('/file/' + fileType + '/' + fileName));
            $.when.apply($, AJAX).done(function () {
                template = arguments[0][0];
                values = arguments[1][0];
                editor = new JSONEditor(document.getElementById("config"), {
                    ajax: true,
                    disable_edit_json: true,
                    disable_properties: true,
                    schema: template,
                    theme: 'bootstrap3',
                    iconlib: "bootstrap3"
                });
                editor.setValue(values);
            });
        } else {
            console.log('New File');
            $.getJSON('/file/' + manifest.templates.configs.settings, function (template) {
                editor = new JSONEditor(document.getElementById("config"), {
                    ajax: true,
                    disable_edit_json: true,
                    disable_properties: true,
                    schema: template,
                    theme: 'bootstrap3',
                    iconlib: "bootstrap3"
                });
            });
        }
        editor.on('ready', function () {
            editor.validate();
        });
    }

    postLoad.functions.push(load);

    save = function () {
        var name = $('#fileName').val();
        var json = editor.getValue();
        var boundary = "---------------------------7da24f2e50046";
        var body = '--' + boundary + '\r\n'
            + 'Content-Disposition: form-data; name="file";'
            + 'filename="temp.txt"\r\n'
            + 'Content-type: plain/text\r\n\r\n'
            + JSON.stringify(json, null, 4) + '\r\n'
            + '--' + boundary + '--';
        $.ajax({
            type: 'POST',
            url: '/file/' + fileType + '/' + fileName,
            contentType: "multipart/form-data; boundary=" + boundary,
            data: body
        }).done(function (response) {
            console.log("Save successful: " + response);
        }).fail(function (response) {
            console.log("Save failed: " + response);
        });
    };

    $('#save').on('click', function () {
        save();
    });

    if (fileName != '') {
        $('#fileName').attr('value', fileName);
    }

    $('#fileNameContainer').find('label').html(fileType + '/');
});