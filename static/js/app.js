
// Using Javascript Modules (https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)
function prompt() {
    let toast = function (c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            position = "top-end",
            title = "",
            footer = ""
        } = c
        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer
        })
    }

    let error = function (c) {
        const {
            msg = "",
            position = "top-end",
            title = "",
            footer = ""
        } = c
        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer
        })
    }

    async function custom(c) {
        const {
            icon ="",
            msg = "",
            title = "",
            showConfirmButton=true
        } = c;
        const { value: result } = await Swal.fire({
            icon :icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton:showConfirmButton,
            willOpen: () => {
                if(c.willOpen !== undefined) {
                    c.willOpen()
                }
            },
            didOpen: () => {
                if(c.didOpen !== undefined) {
                    c.didOpen()
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            }
        })
        // Process the code after the submitting the dates from the popup 
        if(result) {
            if(result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if(c.callback !== undefined) {
                        c.callback(result);
                    }
                }else {
                    c.callback(false)
                }
            }else {
                c.callback(false)
            }
        }
    }


    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom
    }
}