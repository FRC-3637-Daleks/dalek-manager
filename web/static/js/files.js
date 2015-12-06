var deleteFile, selectFile, newFile, lang, fileType, fileName, manifest;
requirejs(['jquery', 'ko'], function ($, ko) {
    var validType = /(autonomous|control|ports|settings)/;
    var filePath = window.location.pathname;
    var temp = filePath.split("/").splice(1, filePath.length);
    var typePos = -1;
    for(var i = 0; i < temp.length; i++) {
        if(validType.test(temp[i])) {
            typePos = i;
            break;
        }
    }
    if(typePos >= 0) {
        fileType = temp[typePos];
    } else {
        fileType = '';
    }
    if(temp.length > 1 && (temp.length - 1 > typePos || typePos == -1)) {
        fileName = temp[temp.length - 1];
        var temp2 = fileName.split('.');
        if (temp2.length == 2) {
            lang = temp2[1];
        } else {
            lang = 'text';
        }
        if (lang == 'txt') {
            lang = 'text';
        }
        document.title = fileName;
    }

    //UI setup
    if(fileType != '') {
        $('#' + fileType).parent().addClass('active');
    }

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

    function updateManifest() {
        if (manifest == null || manifest == '') {
            return;
        }
        var boundary = "---------------------------7da24f2e50046";
        var body = '--' + boundary + '\r\n'
            + 'Content-Disposition: form-data; name="file";'
            + 'filename="temp.txt"\r\n'
            + 'Content-type: plain/text\r\n\r\n'
            + JSON.stringify(manifest, null, 4) + '\r\n'
            + '--' + boundary + '--';
        $.ajax({
            type: 'POST',
            url: '/file/manifest.json',
            contentType: "multipart/form-data; boundary=" + boundary,
            data: body
        }).done(function (response) {
            console.log("Manifest save successful: " + response);
            if (manifest == null || manifest == '') {
                $.getJSON('/file/manifest.json', function (data) {
                    manifest = data;
                });
            }
            updateManifestUI();
        }).fail(function (response) {
            console.log("Manifest save failed: " + response);
        });
    }

    function updateManifestUI() {
        if (manifest == null || manifest == '') {
            return;
        }
        function setFile(file) {
            if (file == null || file == '') {
                return;
            }

            temp = file.split('/');
            var fileName = temp[temp.length - 1];
            if (fileName == '') {
                return;
            }
            $('.selected-file').addClass('hidden');
            $('.accordion-section-title-text:contains("' + fileName + '")').parent().find('span').removeClass('hidden');
        }
        switch (fileType) {
            case "autonomous":
                setFile(manifest.runtime.autonomous);
                break;
            case "ports":
                setFile(manifest.runtime.configs.ports);
                break;
            case "controls":
                setFile(manifest.runtime.configs.controls);
                break;
            case "settings":
                setFile(manifest.runtime.configs.settings);
                break;
            case "logs":
                setFile(manifest.runtime.configs.logs);
                break;
            case "binary":
                setFile(manifest.runtime.binary);
                break;
        }
    }

    function updateFileList() {
        if(fileType == null) {
            return;
        }
        $.getJSON('/file/list/' + fileType, function (data) {
            files.data(data);
            updateAccordionListener();
        });
    }

    deleteFile = function (file) {
        $.ajax({
            type: 'DELETE',
            url: '/file/' + fileType + '/' + file
        }).done(function () {
            updateFileList();
        });
    };
    selectFile = function (file) {
        switch (fileType) {
            case "autonomous":
                manifest.runtime.autonomous = fileType + '/' + file;
                break;
            case "ports":
                manifest.runtime.configs.ports = fileType + '/' + file;
                break;
            case "controls":
                manifest.runtime.configs.controls = fileType + '/' + file;
                break;
            case "settings":
                manifest.runtime.configs.settings = fileType + '/' + file;
                break;
            case "logs":
                manifest.runtime.configs.logs = fileType + '/' + file;
                break;
            case "binary":
                manifest.runtime.binary = fileType + '/' + file;
                break;
        }
        updateManifest();
    };
    newFile = function () {
        var name = "untitled.txt",
            pos = 0,
            found = false;
        while (!found) {
            files.data().forEach(function (element) {
                if (element == name) {
                    name = name.substr(0, 8) + ++pos + ".txt";
                } else {
                    found = true;
                }
            });
        }
        window.location = "/editor/" + fileType + "/" + name;
    };
    function handleFileUpload(files) {
        for (var i = 0; i < files.length; i++) {
            var fd = new FormData();
            fd.append('file', files[i]);
            //console.log(fd);
            sendFileToServer(fd, files[i].name);
        }
    }

    function sendFileToServer(formData, fileName) {
        var uploadURL = "/file/" + fileType + '/' + fileName;
        $.ajax({
            url: uploadURL,
            type: "POST",
            contentType: false,
            processData: false,
            cache: false,
            data: formData,
            success: function () {
                window.location = "/editor/" + fileType + '/' + fileName;
            }
        });
    }

    ko.applyBindings(files);
    $(document).ready(function () {
        $('#file').parent().attr('action', '/editor/' + fileType + '/' + fileName);
        updateFileList();
        window.setInterval(updateFileList, 3000);
        function handleFileSelect(evt) {
            var file = evt.target.files[0];
            if (file != null) {
                $('#file').parent().attr('action', '/editor/' + fileType + '/' + file.name).submit();
            }
        }

        document.getElementById('file').addEventListener('change', handleFileSelect, false);
        //$('file').on('change', handleFileSelect);
        $(document).on('dragenter', function (e) {
            e.stopPropagation();
            e.preventDefault();
            $('body').fadeTo('fast', 0.5);
        });
        $(document).on('dragleave ', function (e) {
            e.stopPropagation();
            e.preventDefault();
            $('body').fadeTo('fast', 1);
        });
        $(document).on('dragover', function (e) {
            e.stopPropagation();
            e.preventDefault();
        });
        $(document).on('drop', function (e) {
            e.stopPropagation();
            e.preventDefault();
            var files = e.originalEvent.dataTransfer.files;
            handleFileUpload(files);
        });
        $.getJSON('/file/manifest.json', function (data) {
            manifest = data;
            updateManifestUI();
        });
    });
});