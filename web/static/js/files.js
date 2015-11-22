var deleteFile, selectFile;
requirejs(['jquery', 'ko'], function ($, ko) {
    $(document).ready(function () {
        function close_accordion_section() {
            $('.accordion .accordion-section-title').removeClass('active');
            $('.accordion .accordion-section-content').slideUp(300).removeClass('open');
        }

        $('.accordion-section-title').on('click', function (e) {
            var currentAttrValue = $(this).find('a').attr('href');

            if ($(e.target).is('.active')) {
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
    });

    deleteFile = function (file) {

    };

    selectFile = function(file) {

    };

    ko.applyBindings(files);
});