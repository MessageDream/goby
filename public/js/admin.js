function active(element) {

}

function changeRole(element) {

}

function forbid(element) {

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
        { data: 'isActive' },
        {
            data: 'role',
            render: function (data, type, row, meta) {
                return data == 1 ? 'admin' : 'member';
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
                if (!row.isActive) {
                    return '<button class="ui button redli inverted tiny admin_user_active" style="margin-bottom:0px;" onclick="active(this)" data-email="' + data + '">Active</button>';
                } else {
                    var roleButton = '<button class="ui button greenli inverted tiny admin_user_role" style="margin-bottom:0px;"  onclick="changeRole(this)" data-role="' + row.role + '" data-email="' + data + '">' + (row.role == 1 ? 'To Member' : 'To Admin') + '</button>';
                    var forbiddenButton = '<button class="ui button blueli inverted tiny admin_user_forbidden" style="margin-bottom:0px;" onclick="forbid(this)" data-forbidden="' + row.isForbidden + '" data-email="' + data + '">' + (row.isForbidden == 1 ? 'Unforbidden' : 'Forbidden') + '</button>';
                    return '<div class="ui bottom aligned">' + roleButton + forbiddenButton + '</div>';
                }
            }
        }
    ]
}).api();