
$('#key_add').on("click", function () {
    var creator = $(this).attr('data-creator');
    $("#key_add_modal").modal({
        blurring: true,
        transition: 'fade up',
        closable: true,
        onApprove: function () {
            var key_form = document.getElementById("form_key_add");
            var form_data = new FormData(key_form);

            var name = form_data.get('name');
            var date = form_data.get('expire_date');

            var d = new Date(date);
            var td = new Date();
            var ttl = d.getTime() - td.getTime();

            var up_form = new FormData();
            up_form.append('friendlyName', name);
            up_form.append('createdBy', creator);
            up_form.append('description', '');
            up_form.append('ttl', ttl);

            $('#key_confirm').addClass('loading');
            $.ajax({
                url: '/accessKeys',
                type: 'POST',
                data: up_form,
                processData: false,  // 不处理数据
                contentType: false,
                success: function (data, textStatus) {

                    $("#key_add_modal").modal('hide');
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
                    $('#key_confirm').removeClass('loading');
                    $('#form_key_error ul').remove();
                    $('#form_key_add').addClass('error');
                    $('#form_key_error').append('<ul class="list"><li>' + XMLHttpRequest.responseJSON.message + '</li></ul>')
                }
            });

            return false;
        }
    }).modal("show");
    $('#key_calendar').calendar({
        ampm: false,
        type: 'date'
    });
});


$(".ui.key_delete").on("click", function () {

    var key_name = $(this).attr('data-name');
    swal({
        text: "Are you sure to delete this key?",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Yes',
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                $.ajax({
                    url: '/accessKeys/' + key_name,
                    type: 'DELETE',
                    success: function (data, textStatus) {
                        resolve(key_name);
                    },
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        reject(XMLHttpRequest.responseJSON.message);
                    }
                });
            })
        }
    }).then(function (key_name) {
        swal({
            title: 'Deleted!',
            text: 'Key of ' + key_name + 'has been deleted',
            type: 'success'
        }).then(function () {
            document.location.reload()
        });
    }).catch(alertError);
});