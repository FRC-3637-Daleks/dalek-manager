/* Team 3637 Dalek Manager - A web base robot configuration manager
 Copyright (C) 2016  Team 3637

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
var editor;
requirejs(['jquery', 'ace/ace', 'mousetrap'], function($, ace, mousetrap) {

    if(/(ports|settings)/.test(fileType)) {
        $('#newGUIFile').removeClass('hidden');
    }

    var nameChange = false;
    var originalName;

    //Load the editor
    editor = ace.edit('editor');
    editor.setTheme('ace/theme/monokai');
    editor.getSession().setMode('ace/mode/' + lang);

    if(readOnly == 'true'){
        editor.setReadOnly(true);
    }

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
        var url = '/file/' + fileType + '/' + fileName;
        url = url.replace(/\/\//g, '/');
        $.ajax({
            type: 'POST',
            url: url,
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

    function updateEditorState() {
        if (editor.session.getUndoManager().isClean() && !nameChange) {
            $('#save').prop("disabled",true);
            document.title = fileName;
            window.onbeforeunload = null;
        }
        else {
            $('#save').prop("disabled",false);
            document.title = "* " + fileName;
            window.onbeforeunload = function(){return "You have unsaved changes are you sure you want to exit?"};
        }
    }

    editor.on('input', function() {
        updateEditorState();
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

    var fileNameBox = $('#fileName');
    if(fileType != '') {
        $('#' + fileType).parent().addClass('active');
        $('#fileNameContainer').find('label').html(fileType + '/')
    }
    if(fileName != '') {
        fileNameBox.attr('value', fileName);
    }
    fileNameBox.on('input', function(e){
        if($(e.target).val() != originalName) {
            nameChange = true;
            fileName = $(e.target).val();
        } else {
            nameChange = false;
            fileName = originalName;
        }
        updateEditorState();
    });
});