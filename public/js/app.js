'use strict'

$('.ui.form_app_add').form({
    fields: {
        name: {
            identifier: 'name',
            rules: [{
                type: 'empty',
                prompt: 'Please enter the app name'
            }]
        },
        gender: {
            identifier: 'platform',
            rules: [{
                type: 'empty',
                prompt: 'Please select a platform'
            }]
        }
    }
});

$("#app_add").on("click", function () {
    $(".ui.modal.standard").modal({
        blurring: true,
        transition:'fade up',
        closable: true,
        onApprove:function(){
           $('.ui.form_app_add').submit(); 
            return false;
        }
    }).modal("show");
});





