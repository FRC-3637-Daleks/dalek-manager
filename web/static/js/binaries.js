var pullBin;
requirejs(['jquery', 'ko'], function ($) {
    pullBin = function() {
        var binName = $('#fileName').val();
        if(binName == null || binName == '') {
            console.log('Invalid bin name: ' + binName);
            return;
        }
        $.ajax({
            type: 'POST',
            url: '/binaries/pull',
            data: {
                fileName: binName
            },
            success: function(response) {
                console.log(response.valueOf())
            }
        });
    };
});