requirejs(['jquery', 'ace/ace', 'mousetrap'], function($, ace, mousetrap) {

    //UI setup
    if(fileType != '') {
        $('#' + fileType).parent().addClass('active');
    }

    $('#folderViewShrink').find('> span').on('click', function(){
        $('#folderViewShrink').parent().parent().addClass('hidden');
        $('#folderViewExpand').removeClass('hidden')
            .parent().removeClass('col-md-2').addClass('nav-sidebar-collapsed');
        $('#editor').parent().removeClass('col-md-10 col-md-offset-2').addClass('editor-sidebar-collapsed');
        $('.nav-bottom').removeClass('col-md-offset-2').addClass('editor-sidebar-collapsed');

    });

    $('#folderViewExpand').on('click', function(){
        $('#editor').parent().removeClass('editor-sidebar-collapsed').addClass('col-md-10 col-md-offset-2');
        $('#folderViewExpand').addClass('hidden')
            .parent().addClass('col-md-2').removeClass('nav-sidebar-collapsed');
        $('#folderViewShrink').parent().parent().removeClass('hidden');
        $('.nav-bottom').addClass('col-md-offset-2').removeClass('editor-sidebar-collapsed');
    });

    $('#fileViewShrink').find('> span').on('click', function(){
        $('#fileViewShrink').parent().parent().addClass('hidden');
        $('#fileViewExpand').removeClass('hidden')
            .parent().removeClass('col-md-2').addClass('nav-sidebar-collapsed');
    });

    $('#fileViewExpand').on('click', function(){
        $('#fileViewExpand').addClass('hidden')
            .parent().addClass('col-md-2').removeClass('nav-sidebar-collapsed');
        $('#fileViewShrink').parent().parent().removeClass('hidden');
    });

    //Load the editor
    var editor = ace.edit('editor');
    editor.setTheme('ace/theme/monokai');
    editor.getSession().setMode('ace/mode/' + lang);

    //Setup save function
    var fileName = document.title;

    function save() {
        var data = new FormData();
        data.append('file', editor.getValue());
        var boundary = "---------------------------7da24f2e50046";
        var body = '--' + boundary + '\r\n'
            + 'Content-Disposition: form-data; name="file";'
            + 'filename="temp.txt"\r\n'
            + 'Content-type: plain/text\r\n\r\n'
            + editor.getValue() + '\r\n'
            + '--'+ boundary + '--';
        $.ajax({
            type: 'POST',
            url: window.location.pathname,
            contentType: "multipart/form-data; boundary="+boundary,
            data: body
        }).done(function(response){
            console.log("Save successful: " + response);
            editor.session.getUndoManager().markClean();
            $('#save').prop("disabled",true);
            document.title = fileName;
            window.onbeforeunload = null;
        }).fail(function(response){
            console.log("Save failed: " + response);
        });
    }

    //Setup save UI
    $('#save').on('click', function(){
        save();
    });

    editor.on('input', function() {
        if (editor.session.getUndoManager().isClean()) {
            $('#save').prop("disabled",true);
            document.title = fileName;
            window.onbeforeunload = null;
        }
        else {
            $('#save').prop("disabled",false);
            document.title = "* " + fileName;
            window.onbeforeunload = function(){return "You have unsaved changes are you sure you want to exit?"};
        }
    });

    //Bin save keys
    editor.commands.addCommand({
        name: 'save',
        bindKey: {win: 'Ctrl-S', mac: 'Command-S'},
        exec: function(editor) {
            if (!editor.session.getUndoManager().isClean()) {
                save();
            }
        },
        readOnly: true
    });

    mousetrap.bind(['ctrl+s', 'command+s'], function(e){
        if (e.preventDefault) {
            e.preventDefault();
            console.log("Prevented")
        } else {
            // internet explorer
            e.returnValue = false;
        }
        if (!editor.session.getUndoManager().isClean()) {
            save();
        }
        return false;
    });
});