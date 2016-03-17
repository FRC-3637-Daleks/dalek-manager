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
var pullBin;
requirejs(['jquery'], function ($) {
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