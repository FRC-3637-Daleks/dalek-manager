var deleteFile, selectFile;
requirejs(['jquery', 'ko'], function ($, ko) {

    function ViewModel() {
        var self = this;
        self.data = ko.observableArray();
    }

    files = new ViewModel();
    var oldList = [];

    function close_accordion_section() {
        $('.accordion .accordion-section-title').removeClass('active');
        $('.accordion .accordion-section-content').slideUp(300).removeClass('open');
    }

    function updateAccordionListener() {
        if (!($(files.data()).not(oldList).length === 0 && $(oldList).not(files.data()).length === 0)) {
            oldList = files.data();
            console.log('listener set');
            $('.accordion-section-title').off('click').on('click', function (e) {
                var currentAttrValue = $(this).find('a').attr('href');
                if ($(e.target).is('.active') || $(e.target).parent().is('.active')) {
                    close_accordion_section();
                } else {
                    close_accordion_section();

                    // Add active class to section title
                    $(this).addClass('active');
                    // Open up the hidden content panel
                    $('.accordion ' + currentAttrValue).slideDown(300).addClass('open');
                }
                e.preventDefault();
            });
        }
    }

    function updateFileList() {
        $.getJSON('/file/list/' + fileType, function (data) {
            files.data(data);
            updateAccordionListener();
        });
    }

    deleteFile = function (file) {
        $.ajax({
            type: 'DELETE',
            url: '/file/' + fileType + '/' + file
        }).done(function(){
            updateFileList();
        });
    };

    selectFile = function (file) {

    };

    ko.applyBindings(files);

    $(document).ready(function () {
        updateFileList();
        window.setInterval(updateFileList, 3000);
    });
});