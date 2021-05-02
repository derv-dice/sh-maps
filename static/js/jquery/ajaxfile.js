// Загрузка файла из ответа на ajax запрос и перехват ошибок в случае неудачи
// Использование: Подставить в ajax.success, ajax.error, ajax.xhr функции из этого файла

// Пример использования:
/*
let fd = new FormData
fd.append('file', input.prop('files')[0]);

$.ajax({
    url: '...',
    data: fd,
    processData: false,
    contentType: false,
    type: 'POST',
    success: function (data, textStatus, jqXHR) {
        ajaxDownloadSuccess(data, jqXHR)
        // ...
    },
    error: function (jqXHR) {
        ajaxDownloadError(jqXHR)
        // ...
    },
    xhr: function () {
        return ajaxDownloadXHR()
    },
});
 */

function ajaxDownloadSuccess(data, jqXHR) {
    // Получение имени файла
    var filename = "";
    var disposition = jqXHR.getResponseHeader('Content-Disposition');
    if (disposition && disposition.indexOf('attachment') !== -1) {
        var filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
        var matches = filenameRegex.exec(disposition);
        if (matches != null && matches[1]) filename = matches[1].replace(/['"]/g, '');
    }
    try {
        // выгрузка файла
        var link = document.createElement('a');
        link.href = window.URL.createObjectURL(data);
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    } catch (exc) {
        // Перехват ошибок
        console.log(exc);
    }
}

function ajaxDownloadError(jqXHR) {
    console.log(jqXHR.responseText)
    alert(jqXHR.responseText)
}

function ajaxDownloadXHR() {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 2) {
            if (xhr.status == 200) {
                xhr.responseType = "blob";
            } else {
                xhr.responseType = "text";
            }
        }
    }
    return xhr
}