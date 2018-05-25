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
        success: function(data, textStatus, xhr){
            if(data.State == "success"){
                User.account = user.account;
                User.token = data.Token;
                SetCookie("user",User.account,false, "/", "", false);
                SetCookie("token",User.token ,false, "/", "", false);
                showManagementView();
            }
        }
    });
}

function SetCookie( name, value, expires, path, domain, secure )
{
// set time, it's in milliseconds
var today = new Date();
today.setTime( today.getTime() );

/*
if the expires variable is set, make the correct
expires time, the current script below will set
it for x number of days, to make it for hours,
delete * 24, for minutes, delete * 60 * 24
*/
if ( expires )
{
expires = expires * 1000 * 60 * 60 * 24;
}
var expires_date = new Date( today.getTime() + (expires) );

document.cookie = name + "=" +escape( value ) +
( ( expires ) ? ";expires=" + expires_date.toGMTString() : "" ) +
( ( path ) ? ";path=" + path : "" ) +
( ( domain ) ? ";domain=" + domain : "" ) +
( ( secure ) ? ";secure" : "" );
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}
