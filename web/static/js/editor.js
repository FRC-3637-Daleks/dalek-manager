var editor;

requirejs(['jquery', 'ace/ace'], function($, ace) {

    editor = ace.edit('editor');
    editor.setTheme('ace/theme/monokai');
    editor.getSession().setMode('ace/mode/lua');

    var fileName = document.title;

    $('#save').on('click', function(){
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
    });

    editor.on('input', function() {

        if (editor.session.getUndoManager().isClean()) {
            $('#save').prop("disabled",true);
            document.title = fileName;
            window.onbeforeunload = null;
        }
        else {
            console.log("edit");
            $('#save').prop("disabled",false);
            document.title = "* " + fileName;
            window.onbeforeunload = function(){return "You have unsaved chnages are you sure you want to exit?"};
        }
    });
});