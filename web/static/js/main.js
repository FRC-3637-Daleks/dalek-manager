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
var startCheck, stopCheck;
requirejs(['jquery', 'mqtt'], function ($) {
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

        function checkStatus() {
            $.ajax({
                url: "/status",
                type: "HEAD",
                timeout: 1000,
                statusCode: {
                    200: function () {
                        $('#robot-status').html('Running').removeClass('no-status stopped').addClass('running');
                    },
                    204: function () {
                        $('#robot-status').html('Running').removeClass('no-status stopped').addClass('running');
                    },
                    400: function () {
                        $('#robot-status').html('Not Running').removeClass('no-status running').addClass('stopped');
                    },
                    0: function () {
                        $('#robot-status').html('Not Running').removeClass('no-status running').addClass('stopped');
                    }
                }
            });
        }

        var host = mqttHost,
            port = parseInt(mqttPort),
            topic = '#',
            useTLS = false,
            username = null,
            password = null,
            cleanSession = true,
            mqtt,
            reconnectTimeout = 5000,
            serverIntervalId,
            mqttIntervalId;

        function MQTTconnect() {
            if (typeof path == "undefined") {
                path = '/mqtt';
            }
            mqtt = new Paho.MQTT.Client(
                host,
                port,
                path,
                "dalek-manager"
            );
            var options = {
                timeout: 3,
                useSSL: useTLS,
                cleanSession: cleanSession,
                onSuccess: onConnect,
                onFailure: function (message) {
                    //console.log("Connection failed: " + message.errorMessage + "Retrying");
                    $('#mqtt-status').html('Not Running').removeClass('no-status running').addClass('stopped');
                    mqttIntervalId = setTimeout(MQTTconnect, reconnectTimeout);
                }
            };

            mqtt.onConnectionLost = onConnectionLost;

            if (username != null) {
                options.userName = null;
                options.password = null;
            }
            //console.log("Host=" + host + ", port=" + port + ", path=" + path + " TLS = " + useTLS + " username=" + username + " password=" + password);
            mqtt.connect(options);
        }

        function onConnect() {
            //console.log('Connected to ' + host + ':' + port + path);
            mqtt.subscribe(topic, {qos: 0});
            $('#mqtt-status').html('Running').removeClass('no-status stopped').addClass('running');
        }

        function onConnectionLost(response) {
            mqttIntervalId = setTimeout(MQTTconnect, reconnectTimeout);
            //console.log("connection lost: " + responseObject.errorMessage + ". Reconnecting");
            $('#mqtt-status').html('Not Running').removeClass('no-status running').addClass('stopped');

        }

        startCheck = function() {
            checkStatus();
            MQTTconnect();
            serverIntervalId = setInterval(checkStatus, reconnectTimeout);
        };

        stopCheck = function() {
            clearInterval(serverIntervalId);
            clearInterval(mqttIntervalId);
            $('#robot-status').html('No Status').removeClass('running stopped').addClass('no-status');
            $('#mqtt-status').html('No Status').removeClass('running stopped').addClass('no-status');
        };

        $(document).ready(function () {
            //startCheck();
        });
    });
});