var editor;

requirejs(['jquery', 'ace/ace', 'mousetrap'], function($, ace, mousetrap) {

    editor = ace.edit('editor');
    editor.setTheme('ace/theme/monokai');
    editor.getSession().setMode('ace/mode/' + lang);

    var fileName = document.title;
    function confirmExit() {
        return "You have attempted to leave this page. Are you sure?";
    }

    function save() {
        var data = new FormData();
        data.append('file', editor.getValue());
        console.log(data);
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
            console.log(response);
            editor.session.getUndoManager().markClean();
            $('#save').prop("disabled",true);
            document.title = fileName;
            window.onbeforeunload = null;
        }).fail(function(response){
            console.log(response);
        });
    }

    $('#save').on('click', function(){
        save();
    });

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

    editor.on('input', function() {
        if (editor.session.getUndoManager().isClean()) {
            $('#save').prop("disabled",true);
            document.title = fileName;
            window.onbeforeunload = null;
        }
        else {
            $('#save').prop("disabled",false);
            document.title = "* " + fileName;
            window.onbeforeunload = confirmExit;
        }
    });
});