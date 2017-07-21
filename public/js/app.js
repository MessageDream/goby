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
                        reject(XMLHttpRequest.responseJSON.message)
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
                        reject(XMLHttpRequest.responseJSON.message);
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
                        reject(XMLHttpRequest.responseJSON.message);
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
                        reject(XMLHttpRequest.responseJSON.message)
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
                        reject(XMLHttpRequest.responseJSON.message);
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
                        reject(XMLHttpRequest.responseJSON.message);
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
                        reject(XMLHttpRequest.responseJSON.message)
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
        closable: false,
        onApprove: function () {
            $('#release_confirm').addClass('loading');
            var release_form = document.getElementById("form_dep_release");
            var form_data = new FormData(release_form);

            var version = form_data.get('version');
            var rollout = form_data.get('rollout');
            var desc = form_data.get('desc');
            var file = form_data.get('package');
            var packageInfo = {
                appVersion: version,
                isDisabled: false,
                isMandatory: true,
                rollout: parseInt(rollout),
                description: desc
            };
            var up_form = new FormData();
            up_form.append('package', file);
            up_form.append('packageInfo', JSON.stringify(packageInfo));

            $.ajax({
                url: '/apps/' + app_name + '/deployments/' + dep_name + '/release',
                type: 'POST',
                data: up_form,
                processData: false,  // 不处理数据
                contentType: false,
                success: function (data, textStatus) {
                    $("#dep_new_release_modal").modal('hide');
                    Lobibox.notify('success', {
                        delay: 2000,
                        msg: 'success'
                    });
                    var timeOut = setTimeout(function () {
                        clearTimeout(timeOut);
                        document.location.reload();
                    }, 2100);

                },
                error: function (XMLHttpRequest, textStatus, errorThrown) {
                    $('#release_confirm').removeClass('loading');
                    // $("#dep_new_release_modal").modal('hide');
                    $('.ui.error.message.release ul').remove();
                    $('#form_dep_release').addClass('error');
                    $('.ui.error.message.release').append('<ul class="list"><li>' + XMLHttpRequest.responseJSON.message + '</li></ul>')
                    // Lobibox.notify('error', {
                    //     delay: 2000,
                    //     msg:XMLHttpRequest.responseJSON.message
                    // });
                }
            });

            return false;
        }
    }).modal("show");
});



function fetchHistory(data, callback, settings) {
    var app_name = $('#app_name').attr('data-name');
    var dep_name = $('#menu_history').dropdown('get text');

    var returnData = {
        draw: data.draw
    };


    var pkgs = new Promise(function (resolve, reject) {
        $.ajax({
            url: '/apps/' + app_name + '/deployments/' + dep_name + '/history',
            type: 'GET',
            success: function (data, textStatus) {
                resolve(data.history);
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                reject(XMLHttpRequest.responseJSON.message);
            }
        });
    });

    var pkgMetrics = new Promise(function (resolve, reject) {
        $.ajax({
            url: '/apps/' + app_name + '/deployments/' + dep_name + '/metrics',
            type: 'GET',
            success: function (data, textStatus) {
                resolve(data.metrics);
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                reject(XMLHttpRequest.responseJSON.message);
            }
        });
    });

    Promise.all([pkgs, pkgMetrics]).then(function (arr) {
        var histories = arr[0];
        var metrics = arr[1];
        var totalActive = getTotalActiveFromDeploymentMetrics(metrics);
        histories.forEach(function (packageObject) {
            if (metrics[packageObject.label]) {
                packageObject.metrics = {
                    active: metrics[packageObject.label].active,
                    downloaded: metrics[packageObject.label].downloaded,
                    failed: metrics[packageObject.label].failed,
                    installed: metrics[packageObject.label].installed,
                    totalActive: totalActive
                };
            }
        });
        return histories;
    }).then(function (his) {
        returnData.recordsTotal = his.length;
        returnData.recordsFiltered = his.length;
        returnData.data = his;
        if (callback) {
            callback(returnData);
        }
    }).catch(function (err) {
        returnData.error = err;
        if (callback) {
            callback(returnData);
        }
    });
}

function getTotalActiveFromDeploymentMetrics(metrics) {
    var totalActive = 0;
    Object.keys(metrics).forEach((label) => {
        totalActive += metrics[label].active;
    });

    return totalActive;
}

// $(document).ready(function () {
var table = $('#table_history').dataTable({
    "ordering": false,
    "paging": false,
    "searching": false,
    "initComplete": function () {
        // table = this.api();
        // table.draw();
    },
    "ajax": fetchHistory,
    "columns": [
        { data: "label" },
        { data: "appVersion" },
        { data: "isMandatory" },
        { data: "releaseMethod" },
        {
            data: "uploadTime",
            render: function (data, type, row, meta) {
                var date = new Date(data)
                return moment(date).format("YYYY-MM-DD hh:mm:ss");
            }
        },
        { data: "description" },
        {
            data: "metrics",
            render: function (data, type, row, meta) {
                return '<div class="ui list"><div class="item"><i class="attach icon"></i><div class="content"> Active:&nbsp;&nbsp;' + (data.totalActive != 0 ? (data.active / data.totalActive / 100) : 0) + '（' + data.active + '&nbsp;&nbsp;of&nbsp;&nbsp;' + data.totalActive + '）' + '</div></div><div class="item"><i class="download icon"></i><div class="content">Total:&nbsp;&nbsp;' + data.installed + '</div></div><div class="item"><i class="undo icon"></i><div class="content">Rollbacks:&nbsp;&nbsp;' + data.failed + ' </div></div></div>';
            }
        }
    ]
}).api();

$('#menu_history').dropdown({
    onChange: function (value, text, $selectedItem) {
        table.ajax.reload();
    }
});
// });
