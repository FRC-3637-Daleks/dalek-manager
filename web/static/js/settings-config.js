var save;
requirejs(['jquery', 'ko'], function ($, ko) {

    function load() {
        $.getJSON('/file/settings/template.json', function(data) {
            ko.applyBindings(data, document.getElementById('content'));
        }).fail(function(){
            console.log('Failed to get template.json')
        });
    }

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
    };

    $('#save').on('click', function(){
        save();
    });

    $(document).ready(function (){
        load();
    });
});