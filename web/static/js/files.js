var deleteFile, selectFile, lang, fileType, fileName;
requirejs(['jquery', 'ko'], function ($, ko) {
    var filePath = window.location.pathname;
    var temp = filePath.split("/");
    if(temp.length == 3) {
        fileType = '';
        fileName = temp[2];
        temp = temp[2].split('.');
    } else if(temp.length == 4) {
        fileType = temp[2];
        fileName = temp[3];
        temp = temp[3].split('.');
    }
    if(temp.length == 2) {
        lang = temp[1];
    } else {
        lang = 'text';
    }
    if(lang == 'txt') {
        lang = 'text';
    }
    document.title = fileName;

    function ViewModel() {
        var self = this;
        self.data = ko.observableArray();
        self.fileType = fileType;
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