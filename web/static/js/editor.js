var editor;
requirejs(['jquery', 'ace/ace', 'mousetrap'], function($, ace, mousetrap) {

    var nameChange = false;
    var originalName;

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
            url: '/file/' + fileType + '/' + fileName,
            contentType: "multipart/form-data; boundary="+boundary,
            data: body
        }).done(function(response){
            console.log("Save successful: " + response);
            editor.session.getUndoManager().markClean();
            if(nameChange) {
                window.history.pushState(null, null, '/editor/' + fileType + '/' + fileName);
            }
            nameChange = false;
            var temp = fileName.split('.');
            if(temp.length == 2) {
                lang = temp[1];
            } else {
                lang = 'text';
            }
            if(lang == 'txt') {
                lang = 'text';
            }
            editor.getSession().setMode('ace/mode/' + lang);
            updateEditorState()
        }).fail(function(response){
            console.log("Save failed: " + response);
        });
    }

    function updateEditorState() {
        if (editor.getSession().getUndoManager().isClean() && !nameChange) {
            //console.log('yes', editor.session.getUndoManager().isClean(), !nameChange);
            //$('#save').addClass("disabled");
            document.title = fileName;
            window.onbeforeunload = null;
        }
        else {
            //console.log('no', editor.session.getUndoManager().isClean(), !nameChange);
            //$('#save').removeClass("disabled");
            document.title = "* " + fileName;
            window.onbeforeunload = function(){return "You have unsaved changes are you sure you want to exit?"};
        }
    }

    //Load the editor
    editor = ace.edit('editor');
    editor.setTheme('ace/theme/monokai');
    editor.getSession().setMode('ace/mode/' + lang);

    //Bind save keys
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
    editor.getSession().on('change', function(){updateEditorState()});

    mousetrap.bind(['ctrl+s', 'command+s'], function(e){
        if (e.preventDefault) {
            e.preventDefault();
        } else {
            // internet explorer
            e.returnValue = false;
        }
        if (!editor.session.getUndoManager().isClean()) {
            save();
        }
        return false;
    });

    $(document).ready(function(){
        //UI setup
        if(fileType != '') {
            $('#' + fileType).parent().addClass('active');
            $('#fileNameContainer').find('label').html(fileType + '/')
        }
        if(fileName != '') {
            $('#fileName').attr('value', fileName);
        }

        //Setup save UI
        $('#save').on('click', function(){
            save();
        });

        $('#fileName').on('input', function(e){
            if($(e.target).val() != originalName) {
                nameChange = true;
                fileName = $(e.target).val();
            } else {
                nameChange = false;
                fileName = originalName;
            }
            updateEditorState();
        });

        updateEditorState();
    });
});