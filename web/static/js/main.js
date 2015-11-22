requirejs(['jquery'], function ($) {
    $(document).ready(function () {

        $('#folderViewShrink').find('> span').on('click', function () {
            $('#folderViewShrink').parent().parent().addClass('hidden');
            $('#folderViewExpand').removeClass('hidden')
                .parent().removeClass('col-md-2').addClass('nav-sidebar-collapsed');
            $('#editor').parent().removeClass('col-md-10 col-md-offset-2').addClass('editor-sidebar-collapsed');
            $('.nav-bottom').removeClass('col-md-offset-2').addClass('editor-sidebar-collapsed');
        });

        $('#folderViewExpand').on('click', function () {
            $('#editor').parent().removeClass('editor-sidebar-collapsed').addClass('col-md-10 col-md-offset-2');
            $('#folderViewExpand').addClass('hidden')
                .parent().addClass('col-md-2').removeClass('nav-sidebar-collapsed');
            $('#folderViewShrink').parent().parent().removeClass('hidden');
            $('.nav-bottom').addClass('col-md-offset-2').removeClass('editor-sidebar-collapsed');
        });

        $('#fileViewShrink').find('> span').on('click', function () {
            $('#fileViewShrink').parent().addClass('hidden');
            $('#fileViewExpand').removeClass('hidden')
                .parent().removeClass('col-md-2').addClass('nav-sidebar-collapsed');
            $('#files').addClass('hidden');
        });

        $('#fileViewExpand').on('click', function () {
            $('#fileViewExpand').addClass('hidden')
                .parent().addClass('col-md-2').removeClass('nav-sidebar-collapsed');
            $('#fileViewShrink').parent().removeClass('hidden');
            $('#files').removeClass('hidden');
        });

    });
});