'use strict'
var alertError = function (error) {
    if (error === 'cancel') {
        return;
    }
    swal({
        type: 'fail',
        title: 'Operation error',
        html: error
    })
};

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

//collaborator
$("#app_add").on("click", function () {
    $("#app_add_modal").modal({
        blurring: true,
        transition: 'fade up',
        closable: true,
        onApprove: function () {
            $('.ui.form_app_add').submit();
            return false;
        }
    }).modal("show");
});

$("#col_add").on("click", function () {
    var app_name = $('#app_name').attr('data-name');
    swal({
        title: 'New collaborator',
        input: 'email',
        text: 'Input an email',
        showCancelButton: true,
        confirmButtonText: 'Submit',
        showLoaderOnConfirm: true,
        preConfirm: function (email) {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/collaborators/' + email,
                    type: 'POST',
                    success: function (data, textStatus) {
                        console.log(data);
                        resolve()
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown)
                    }
                });
            })
        },
        allowOutsideClick: false
    }).then(function (email) {
        swal({
            type: 'success',
            title: 'Operation finished!',
            html: 'Submitted email: ' + email
        }).then(function () {
            document.location.reload();
        });
    }).catch(alertError);
});


$(".ui.col_remove").on("click", function () {

    var app_name = $('#app_name').attr('data-name');
    var email = $(this).attr("data-email");

    swal({
        text: "Are you sure?",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/collaborators/' + email,
                    type: 'DELETE',
                    success: function (data, textStatus) {
                        resolve(email);
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown);
                    }
                });
            })
        }
    }).then(function (email) {
        swal({
            title: 'Removed!',
            text: 'Collaborator of ' + email + 'has been removed',
            type: 'success'
        }).then(function () {
            document.location.reload()
        });
    }).catch(alertError);
});


$(".ui.col_transfer").on("click", function () {

    var app_name = $('#app_name').attr('data-name');
    var email = $(this).attr("data-email");

    swal({
        text: "Are you sure to transfer?",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/transfer/' + email,
                    type: 'POST',
                    success: function (data, textStatus) {
                        resolve(email);
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown);
                    }
                });
            })
        }
    }).then(function (email) {
        swal({
            title: 'Transfer!',
            text: app_name + ' has been transfer to ' + email,
            type: 'success'
        }).then(function () {
            document.location.href = '/web/app/list';
        });
    }).catch(alertError);
});

// deployment
$("#dep_add").on("click", function () {
    var app_name = $('#app_name').attr('data-name');
    swal({
        title: 'New deployment',
        input: 'text',
        text: 'Input a deployment name',
        showCancelButton: true,
        confirmButtonText: 'Submit',
        showLoaderOnConfirm: true,
        preConfirm: function (name) {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/deployments',
                    type: 'POST',
                    contentType: 'application/json;charset=utf-8',
                    data: JSON.stringify({
                        name: name
                    }),
                    success: function (data, textStatus) {
                        resolve()
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown)
                    }
                });
            })
        },
        allowOutsideClick: false
    }).then(function (name) {
        swal({
            type: 'success',
            title: 'Operation finished!',
            html: 'Deployment of: ' + name
        }).then(function () {
            document.location.reload();
        });
    }).catch(alertError);
});

$(".ui.dep_delete").on("click", function () {

    var app_name = $('#app_name').attr('data-name');
    var dep_name = $(this).attr("data-name");

    swal({
        text: "Are you sure?",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/deployments/' + dep_name,
                    type: 'DELETE',
                    success: function (data, textStatus) {
                        resolve(dep_name);
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown);
                    }
                });
            })
        }
    }).then(function (name) {
        swal({
            title: 'Deleted!',
            text: 'Deployment of ' + name + 'has been deleted',
            type: 'success'
        }).then(function () {
            document.location.reload()
        });
    }).catch(alertError);
});

$(".ui.dep_rollback").on("click", function () {

    var app_name = $('#app_name').attr('data-name');
    var dep_name = $(this).attr("data-name");

    swal({
        text: "Are you sure to rollback the last release",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/deployments/' + dep_name + '/rollback',
                    type: 'POST',
                    success: function (data, textStatus) {
                        resolve(dep_name);
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown);
                    }
                });
            })
        }
    }).then(function (name) {
        swal({
            title: 'Rolled Back',
            text: 'The last release of ' + name + 'has been rolled back',
            type: 'success'
        }).then(function () {
            document.location.reload()
        });
    }).catch(alertError);
});

$(".ui.dep_promote").on("click", function () {
    var app_name = $('#app_name').attr('data-name');
    var dep_name = $(this).attr("data-name");
    swal({
        title: 'Promote deployment',
        input: 'text',
        text: 'Input the dest deployment name',
        showCancelButton: true,
        confirmButtonText: 'Submit',
        showLoaderOnConfirm: true,
        preConfirm: function (name) {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/apps/' + app_name + '/deployments/' + dep_name + '/promote/' + name,
                    type: 'POST',
                    success: function (data, textStatus) {
                        resolve()
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(errorThrown)
                    }
                });
            })
        },
        allowOutsideClick: false
    }).then(function (name) {
        swal({
            type: 'success',
            title: 'Operation finished!',
            html: 'Deployment of: ' + name
        }).then(function () {
            document.location.reload();
        });
    }).catch(alertError);
});

$(".ui.dep_new_release").on("click", function () {
    var app_name = $('#app_name').attr('data-name');
    var dep_name = $(this).attr("data-name");
    $("#dep_new_release_modal").modal({
        blurring: true,
        transition: 'fade up',
        closable: true,
        onApprove: function () {
            var formData = new FormData($('#form_dep_release'));
            console.log(formData);
            $.ajax({
                url: '/apps/' + app_name + '/deployments/' + dep_name + '/release',
                type: 'POST',
                data:formData,
                processData: false,  // 不处理数据
                contentType: false,
                success: function (data, textStatus) {

                },
                error: function (XMLHttpRequest, textStatus, errorThrown) {

                }
            });

            return false;
        }
    }).modal("show");
});

