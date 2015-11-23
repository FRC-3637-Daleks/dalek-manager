var deleteFile, selectFile, files;
requirejs(['jquery', 'ko'], function ($, ko) {
    function ViewModel() {
        var self = this;
        self.data = ko.observableArray();
    }

    function updateFileList() {
        $.getJSON('/file/list/' + fileType, function(data){
            files.data(data);
        });
    }

    files = new ViewModel();

    deleteFile = function (file) {
        $.ajax({
            type: 'DELETE',
            url: '/file/' + fileType + '/' + fileName
        }).done(function(){
            updateFileList();
        });
    };

    selectFile = function(file) {

    };

    ko.applyBindings(files);

    $(document).ready(function () {
        console.log('test');
        function close_accordion_section() {
            $('.accordion .accordion-section-title').removeClass('active');
            $('.accordion .accordion-section-content').slideUp(300).removeClass('open');
        }
        console.log('test');
        $('.accordion-section-title').on('click', function (e) {
            console.log('click');
            var currentAttrValue = $(this).find('a').attr('href');

            if ($(e.target).is('.active')) {
                close_accordion_section();
            } else {
                close_accordion_section();
                $(this).addClass('active');
                $('.accordion ' + currentAttrValue).slideDown(300).addClass('open');
            }

            e.preventDefault();
        });
        console.log('test');
        updateFileList();
        window.setInterval(updateFileList, 1000);
    });
});