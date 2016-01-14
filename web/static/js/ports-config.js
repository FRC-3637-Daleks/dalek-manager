var save;
requirejs(['jquery', 'ko'], function ($, ko) {

    function showError(error) {
        $('#error').html(error).removeClass('hidden');
    }

    function load() {
        if (files.data().indexOf(fileName) > -1) {
            console.log('Loading: ' + fileName);
            var template, values, AJAX = [];
            AJAX.push($.getJSON('/file/' + manifest.templates.configs.ports));
            AJAX.push($.getJSON('/file/' + fileType + '/' + fileName));
            $.when.apply($, AJAX).done(function () {
                template = arguments[0][0];
                values = arguments[1][0];
                if (template == null || template == '' || template == {} || values == null || values == '' || values == {}) {
                    console.log('No data loaded');
                    showError('No data loaded');
                    return;
                }
                if (!Array.isArray(template)) {
                    console.log('Invalid template: template is not an array');
                    showError('Invalid template: template is not an array');
                }
                template.forEach(function (element) {
                    var top = element, path = [], depth = 0, setKeys = function (element) {
                        if (element.hasOwnProperty("keys") && Array.isArray(element.keys) && element.keys.length > 0) {
                            depth++;
                            path.push(element.name);
                            element.keys.forEach(setKeys);
                            path.pop();
                        } else {
                            path.push(element.name);
                            var val = values[top.name];
                            path.forEach(function (element) {
                                val = val[element];

                            });
                            path.pop();
                            if (val != null) {
                                element.value = val;
                            }
                        }
                    };
                    if (!(element.hasOwnProperty("name") &&
                        element.hasOwnProperty("min") &&
                        element.hasOwnProperty("max") &&
                        element.hasOwnProperty("keys") &&
                        Array.isArray(element.keys))) {
                        console.log('Element is missing a property (name, min max, keys)');
                        showError('Element is missing a property (name, min max, keys)');
                        return;
                    }
                    element.keys.forEach(setKeys);
                });
                ko.applyBindings(template, document.getElementById('content'));
            }).fail(function() {
                console.log('No data loaded');
                showError('No data loaded');
            });
        } else {
            console.log('New File');
            $.getJSON('/file/' + manifest.templates.configs.ports, function (template) {
                if (!Array.isArray(template)) {
                    console.log('Template is not an array');
                    showError('Template is not an array');
                    return;
                }
                template.forEach(function (element) {
                    var top = element, pos = top.min, setKeys = function (element) {
                        if (element.hasOwnProperty("keys") && Array.isArray(element.keys) && element.keys.length > 0) {
                            element.keys.forEach(setKeys);
                        } else {
                            //console.log(element);
                            if (pos + 1 <= top.max) {
                                element.value = pos++;
                            }
                        }
                    };
                    if (!(element.hasOwnProperty("name") &&
                        element.hasOwnProperty("min") &&
                        element.hasOwnProperty("max") &&
                        element.hasOwnProperty("keys") &&
                        Array.isArray(element.keys))) {
                        console.log('Element is missing a property (name, min max, keys)');
                        showError('Element is missing a property (name, min max, keys)');
                        return;
                    }
                    element.keys.forEach(setKeys);
                });
                ko.applyBindings(template, document.getElementById('content'));
            }).fail(function() {
                console.log('No data loaded');
                showError('No data loaded');
            });
        }
    }

    postLoad.functions.push(load);

    save = function () {
        var name = $('#fileName').val();
        var json = {};
        var subsystems = $('input').filter(function () {
            return $(this).attr('data-subsystem') != null;
        }).toArray();
        subsystems.forEach(function (element) {
            var subsystem = $(element).attr('data-subsystem');
            if (!json.hasOwnProperty(subsystem)) {
                json[subsystem] = {};
            }
            json[subsystem][$(element).attr('name')] = $(element).val();
        });
        console.log(json);
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