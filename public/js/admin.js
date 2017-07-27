function status(el) {
    var email = $(el).attr("data-email");
    var status = $(el).attr("data-status")
    switch (status) {
        case '0':
            status = '1'
            break;
        case '1':
            status = '2'
            break;
        case '2':
            status = '1'
            break;
    }

    swal({
        text: 'Are you sure to change ' + email + '\'s status?',
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/web/admin/api/users/' + email + '/status',
                    type: 'PATCH',
                    data: { status: status },
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
            title: 'OK!',
            text: 'Status of ' + email + 'has been changed',
            type: 'success'
        }).then(function () {
            admin_users_table.ajax.reload();
        });
    }).catch(alertError);
}

function changeRole(el) {
    var email = $(el).attr("data-email");
    var role = $(el).attr("data-role")
    switch (role) {
        case '0':
            role = '1'
            break;
        case '1':
            role = '0'
            break;
    }

    swal({
        text: 'Are you sure to change ' + email + '\'s role?',
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/web/admin/api/users/' + email + '/role',
                    type: 'PATCH',
                    data: { role: role },
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
            title: 'OK!',
            text: 'Role of ' + email + 'has been changed',
            type: 'success'
        }).then(function () {
            admin_users_table.ajax.reload();
        });
    }).catch(alertError);
}


function fetchAllUsers(params, callback, settings) {
    var start = params.start ? params.start : 0
    var length = params.length ? params.length : 20
    var search = params.search ? params.search.value : ""
    var pageIndex = start / length;
    var returnData = {};
    $.ajax({
        url: '/web/admin/api/users/' + pageIndex + '/' + length + '?email=' + search,
        type: 'GET',
        success: function (data, textStatus) {
            var result = data.users;
            returnData.start = start;
            returnData.recordsTotal = result.totalCount;
            returnData.recordsFiltered = result.totalCount;
            returnData.data = result.data;
            if (callback) {
                callback(returnData);
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            returnData.error = XMLHttpRequest.responseJSON.message;
            if (callback) {
                callback(returnData);
            }
        }
    });
}

var admin_users_table = $('#table_users').dataTable({
    "ordering": false,
    "paginate": true,
    "lengthMenu": [20, 50, 100],
    "processing": true,
    "serverSide": true,
    "ajax": fetchAllUsers,
    "columns": [
        { data: 'userName' },
        { data: 'email' },
        {
            data: 'role',
            render: function (data, type, row, meta) {
                return data == 1 ? 'admin' : 'member';
            }
        },
        {
            data: 'status',
            render: function (data, type, row, meta) {
                switch (data) {
                    case 0:
                        return 'Inactive'
                        break;
                    case 1:
                        return 'Normal'
                    case 2:
                        return 'Forbidden'
                        break;
                    default:
                        return 'Unknown'
                }
            }
        },
        {
            data: 'joinedTime',
            render: function (data, type, row, meta) {
                return moment(new Date(data)).format("YYYY-MM-DD hh:mm:ss");
            }
        },
        {
            data: 'email',
            render: function (data, type, row, meta) {
                if (row.status == 0) {
                    return '<button class="ui button redli inverted tiny admin_user_status" style="margin-bottom:0px;" onclick="status(this)" data-status="' + row.status + '" data-email="' + data + '">Active</button>';
                } else {
                    var roleButton = '<button class="ui button ' + (row.role == 1 ? 'yellowli' : 'greyli') + ' inverted tiny admin_user_role" style="margin-bottom:0px;"  onclick="changeRole(this)" data-role="' + row.role + '" data-email="' + data + '">' + (row.role == 1 ? 'As Member' : 'As Admin') + '</button>';
                    var forbiddenButton = '<button class="ui button '+ (row.status == 1 ? 'greenli':'redli') +' inverted tiny admin_user_status" style="margin-bottom:0px;" onclick="status(this)" data-status="' + row.status + '" data-email="' + data + '">' + (row.status == 2 ? 'Unforbidden' : 'Forbidden') + '</button>';
                    return '<div class="ui bottom aligned">' + roleButton + forbiddenButton + '</div>';
                }
            }
        }
    ]
}).api();

$("#admin_user_add").on("click", function () {
    var app_name = $('#app_name').attr('data-name');
    var dep_name = $(this).attr("data-name");
    $('#admin_user_add_modal').modal({
        blurring: true,
        transition: 'fade up',
        closable: false,
        onApprove: function () {
            $('#admin_add_user_confirm').addClass('loading');
            var user_add_form = document.getElementById("form_admin_user_add");
            var form_data = new FormData(user_add_form);
            $.ajax({
                url: '/web/admin/api/users/add',
                type: 'POST',
                data: form_data,
                processData: false,  // 不处理数据
                contentType: false,
                success: function (data, textStatus) {
                    $("#admin_user_add_modal").modal('hide');
                    admin_users_table.ajax.reload();
                    Lobibox.notify('success', {
                        delay: 2000,
                        msg: 'success'
                    });
                },
                error: function (XMLHttpRequest, textStatus, errorThrown) {
                    $('#admin_add_user_confirm').removeClass('loading');
                    $('#form_admin_user_add_error ul').remove();
                    $('#form_admin_user_add').addClass('error');
                    $('#form_admin_user_add_error').append('<ul class="list"><li>' + XMLHttpRequest.responseJSON.message + '</li></ul>')
                }
            });

            return false;
        }
    }).modal("show");
});