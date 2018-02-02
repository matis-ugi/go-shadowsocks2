"use strict";
var Salt = genSalt();
var User = {"account":"","token":"","salt":Salt};

function genSalt() {
    return "1";
}

function login(){
    var user = {
        "account":$("#user").val(),
        "pwd":$("#pwd").val(),
        "salt":Salt
    };
    $.ajax({
        url: 'api/login',
        type: 'POST',
        data: JSON.stringify(user),
        contentType: 'application/json; charset=utf-8',
        dataType: 'json',
        async: false,
        success: function(data) {
            if(data.State == "success"){
                User.account = user.account;
                User.token = data.Token;
                showManagementView();
            }
        }
    });
}
