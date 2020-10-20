const email = document.querySelector("input#email")
const password = document.querySelector("input#password")

const encodeFormData=(data)=>{
    return Object.keys(data)
        .map(key=>encodeURIComponent(key)+"="+encodeURIComponent(data[key]))
        .join("&");
}
function submitForm(ev){
    let c_id = window.sessionStorage.getItem("client_id")
    let c_secret = window.sessionStorage.getItem("client_secret")
    ev.preventDefault()
    if (email.value=="" && password.value==""){
        alert("Creds not Okay")
        return false
    }
    const req_data={
        "username":email.value,
        "password":password.value,
        "grant_type":"password"
    }
    let headers = new Headers()
    headers.append("Authorization","Basic "+btoa(c_id+":"+c_secret))
    headers.append("Content-Type","application/x-www-form-urlencoded")
    fetch("http://localhost:3001/api/auth/generate_token/",{
        "method":"POST",
        headers:headers,
        body:encodeFormData(req_data)
    }).then(resp=>resp.json())
        .then(data=>{
            window.sessionStorage.setItem("access_token",data.access_token)
            window.sessionStorage.setItem("token_type",data.token_type)
            let token_header = new Headers()
            token_header.append("Authorization",window.sessionStorage.getItem("token_type")+" "+
                window.sessionStorage.getItem("access_token"))
            fetch("http://localhost:3001/api/auth/test",{
                "method":"POST",
                headers:token_header
            })
                .then(resp=>resp.json())
                .then(data=>{
                    console.table(data)
                })
        })
        .catch(err=>{

        })
    return false
}