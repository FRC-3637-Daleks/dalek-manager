var save;
requirejs(['jquery', 'ko'], function ($, ko) {

    function load() {
        if(files.data().indexOf(fileName) > -1) {
            console.log('Loading: ' + fileName);
            var template, values, AJAX = [];
            AJAX.push($.getJSON('/file/settings/' + manifest.templates.configs.settings));
            AJAX.push($.getJSON('/file/settings/' + fileName));
            $.when.apply($, AJAX).done(function () {
                template = arguments[0][0];
                values = arguments[1][0];
                if (!template.hasOwnProperty("subsystems")) {
                    console.log('Element subsystem not found');
                    return;
                }
                template.subsystems.forEach(function (element) {
                    if (!element.hasOwnProperty("values")) {
                        console.log('Element values not found');
                        return;
                    }
                    var sub = element;
                    element.values.forEach(function (element) {
                        if (values.hasOwnProperty(sub.name)) {
                            if (values[sub.name].hasOwnProperty(element.name)) {
                                element.value = values[sub.name][element.name];
                            }
                        }
                    });
                });
                ko.applyBindings(template, document.getElementById('content'));
            });
        } else {
            console.log('New File');
            $.getJSON('/file/settings/' + manifest.templates.configs.settings, function(template) {
                if (!template.hasOwnProperty("subsystems")) {
                    console.log('Element subsystem not found');
                    return;
                }
                template.subsystems.forEach(function (element) {
                    if (!element.hasOwnProperty("values")) {
                        console.log('Element values not found');
                        return;
                    }
                    var sub = element;
                    element.values.forEach(function (element) {
                        element.value = null;
                    });
                });
                ko.applyBindings(template, document.getElementById('content'));
            });
        }
    }

    postLoad.functions.push(load);

    save = function () {
        var name = $('#fileName').val();
        var json = {};
        var subsystems = $('input').filter(function() {
            return $(this).attr('data-subsystem') != null;
        }).toArray();
        subsystems.forEach(function(element){
            var subsystem = $(element).attr('data-subsystem');
            if(!json.hasOwnProperty(subsystem)) {
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

    $('#save').on('click', function(){
        save();
    });
    if(fileName != '') {
        $('#fileName').attr('value', fileName);
    }
    $('#fileNameContainer').find('label').html(fileType + '/');
});