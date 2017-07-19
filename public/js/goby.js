'use strict';

var csrf;
var suburl;

var alertError = function (error) {
    if (error === 'cancel') {
        return;
    }
    swal({
        type: 'fail',
        title: 'Operation error',
        html: error,
        showCancelButton: false
    })
};

function initInstall() {
    if ($('.install').length == 0) {
        return;
    }

    // Database type change detection.
    $("#db_type").change(function () {
        var sqliteDefault = 'data/gogs.db';

        var dbType = $(this).val();
        if (dbType === "SQLite3") {
            $('#sql_settings').hide();
            $('#pgsql_settings').hide();
            $('#sqlite_settings').show();

            if (dbType === "SQLite3") {
                $('#db_path').val(sqliteDefault);
            }
            return;
        }

        var dbDefaults = {
            "MySQL": "127.0.0.1:3306",
            "PostgreSQL": "127.0.0.1:5432",
            "MSSQL": "127.0.0.1, 1433"
        };

        $('#sqlite_settings').hide();
        $('#sql_settings').show();
        $('#pgsql_settings').toggle(dbType === "PostgreSQL");
        $.each(dbDefaults, function (type, defaultHost) {
            if ($('#db_host').val() == defaultHost) {
                $('#db_host').val(dbDefaults[dbType]);
                return false;
            }
        });
    });

    // storage type change detection.
    $("#storage_type").change(function () {

        var stType = $(this).val();

        if (stType !== "本地") {
            $('#st-remote').show();
            $('#st-local-setting').hide();
            $('#st-oss-setting').toggle(stType === "阿里云OSS");
            $('#st-qn-setting').toggle(stType === "七牛");
        } else {
            $('#st-remote').hide();
            $('#st-oss-setting').hide();
            $('#st-qn-setting').hide();
            $('#st-local-setting').show();
        }

    });

    $('#disable-registration input').change(function () {
        if ($(this).is(':checked')) {
            $('#enable-captcha').checkbox('uncheck');
        }
    });
    $('#enable-captcha input').change(function () {
        if ($(this).is(':checked')) {
            $('#disable-registration').checkbox('uncheck');
        }
    });

}


// For IE
String.prototype.endsWith = function (pattern) {
    var d = this.length - pattern.length;
    return d >= 0 && this.lastIndexOf(pattern) === d;
};

// Adding function to get the cursor position in a text field to jQuery object.
(function ($, undefined) {
    $.fn.getCursorPosition = function () {
        var el = $(this).get(0);
        var pos = 0;
        if ('selectionStart' in el) {
            pos = el.selectionStart;
        } else if ('selection' in document) {
            el.focus();
            var Sel = document.selection.createRange();
            var SelLength = document.selection.createRange().text.length;
            Sel.moveStart('character', -el.value.length);
            pos = Sel.text.length - SelLength;
        }
        return pos;
    }
})(jQuery);


function buttonsClickOnEnter() {
    $('.ui.button').keypress(function (e) {
        if (e.keyCode == 13 || e.keyCode == 32) // enter key or space bar
            $(this).click();
    });
}

function hideWhenLostFocus(body, parent) {
    $(document).click(function (e) {
        var target = e.target;
        if (!$(target).is(body) && !$(target).parents().is(parent)) {
            $(body).hide();
        }
    });
}


$(document).ready(function () {
    csrf = $('meta[name=_csrf]').attr("content");
    suburl = $('meta[name=_suburl]').attr("content");


    buttonsClickOnEnter();

    initInstall();
});



var sideBarIsHide = false;
var ManuelSideBarIsHide = false;
var ManuelSideBarIsState = false;
$(".openbtn").on("click", function () {
    ManuelSideBarIsHide = true;
    if (!ManuelSideBarIsState) {
        resizeSidebar("1");
        ManuelSideBarIsState = true;
    } else {
        resizeSidebar("0");
        ManuelSideBarIsState = false;
    }
});


$(window).resize(function () {
    if (ManuelSideBarIsHide == false) {
        if ($(window).width() <= 767) {
            if (!sideBarIsHide); {
                resizeSidebar("1");
                sideBarIsHide = true;
                $(".colhidden").addClass("displaynone");

            }
        } else {
            if (sideBarIsHide); {
                resizeSidebar("0");
                sideBarIsHide = false;

                $(".colhidden").removeClass("displaynone");

            }
        }
    }
});
var isMobile = window.matchMedia("only screen and (max-width: 768px)");

if (isMobile.matches) {
    resizeSidebar("1");
    $("body")
        .getNiceScroll()
        .remove();
    $(".sidebar")
        .getNiceScroll()
        .remove();

    $(".computer.only").toggleClass("displaynone");
    $(".colhidden").toggleClass("displaynone");
} else {
    $("body").niceScroll({
        cursorcolor: "#3d3b3b",
        cursorwidth: 5,
        cursorborderradius: 0,
        cursorborder: 0,
        scrollspeed: 50,
        autohidemode: true,
        zindex: 9999999
    });
    $(".sidebar").niceScroll({
        cursorcolor: "#3d3b3b",
        cursorwidth: 2,
        cursorborderradius: 0,
        cursorborder: 0,
        scrollspeed: 50,
        autohidemode: true,
        zindex: 9999999
    });

    $(".displaynone .menu").niceScroll({
        cursorcolor: "#3d3b3b",
        cursorwidth: 5,
        cursorborderradius: 0,
        cursorborder: 0,
        scrollspeed: 50,
        autohidemode: true,
        zindex: 9999999
    });
}

function resizeSidebar(op) {

    if (op == "1") {

        $(".ui.sidebar.left").addClass("very thin icon");
        $(".navslide").addClass("marginlefting");
        $(".sidebar.left span").addClass("displaynone");
        $(".sidebar .accordion").addClass("displaynone");
        $(".ui.dropdown.item.displaynone").addClass("displayblock");
        $($(".logo img")[0]).addClass("displaynone");
        $($(".logo img")[1]).removeClass("displaynone");
        $(".hiddenCollapse").addClass("displaynone");


    } else {

        $(".ui.sidebar.left").removeClass("very thin icon");
        $(".navslide").removeClass("marginlefting");
        $(".sidebar.left span").removeClass("displaynone");
        $(".sidebar .accordion").removeClass("displaynone");
        $(".ui.dropdown.item.displaynone").removeClass("displayblock");
        $($(".logo img")[1]).addClass("displaynone");
        $($(".logo img")[0]).removeClass("displaynone");
        $(".hiddenCollapse").removeClass("displaynone");


    }

}

$(".ui.dropdown").dropdown({
    allowCategorySelection: true,
    transition: "fade up"
});
$('.ui.accordion').accordion({
    selector: {}
});


//Sidebar And Navbar Coloring Function (This button on Footer)
function colorize() {
    var a;
    var b;
    var d;
    var z;
    var l;

    if (Cookies.get('sidebarColor') != undefined) {
        if (b == null) {
            b = $(".sidebar").attr("data-color");
        }
        $(".sidemenu").removeClass(b).addClass(Cookies.get('sidebarColor'));
        $(".sidebar").attr("data-color", Cookies.get('sidebarColor'));
    }

    if (Cookies.get('headerColor') != undefined) {
        if (z == null) {
            z = $(".navslide .menu").attr("data-color");
        }
        $(".navslide .menu").removeClass(z).addClass(Cookies.get('headerColor'));
        $(".navslide .menu").attr("data-color", Cookies.get('headerColor'));
    }



    $(".colorlist li a").on("click", function (b) {
        var c = $(this).attr("data-addClass");
        if (l == null) {
            l = $(".navslide .menu").attr("data-color");
        }
        console.log(l);
        $(".navslide .menu").removeClass(l).addClass(c);
        l = c;
        Cookies.set('headerColor', c);
    });
    $(".sidecolor li a").on("click", function (a) {
        var c = $(this).attr("data-addClass");
        // a.preventDefault();
        if (d == null) {
            d = $(".sidebar").attr("data-color");
        }
        $(".sidemenu").removeClass(d).addClass(c);
        $(".accordion").removeClass("inverted").addClass("inverted");
        Cookies.set('sidebarColor', c);
        d = c;
    });
    $(".colorize").popup({
        on: "click"
    });
} (function (i, s, o, g, r, a, m) { i['GoogleAnalyticsObject'] = r; i[r] = i[r] || function () { (i[r].q = i[r].q || []).push(arguments) }, i[r].l = 1 * new Date(); a = s.createElement(o), m = s.getElementsByTagName(o)[0]; a.async = 1; a.src = g; m.parentNode.insertBefore(a, m) })(window, document, 'script', 'https://www.google-analytics.com/analytics.js', 'ga'); ga('create', 'UA-96662612-1', 'auto'); ga('send', 'pageview');
//Sidebar And Navbar Coloring Function (This button on Footer)