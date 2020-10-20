const encodeFormData=(data)=>{
    return Object.keys(data)
                 .map(key=>encodeURIComponent(key)+"="+encodeURIComponent(data[key]))
                .join("&");
}
function submitForm(ev){
    ev.preventDefault()
    let clientName = document.querySelector("#clientName").value
    let clientURI = document.querySelector("#clientURI").value
    let grantType = document.querySelector("#grantType").value
    let redirectURI = document.querySelector("#redirectURI").value
    let responseType = document.querySelector("#responseType").value
    let scope = document.querySelector("#scope").value
    let authMethod = document.querySelector("#authMethod").value
    const formData={
        "clientName":clientName,
        "clientURI":clientURI,
        "grantType":grantType,
        "redirectURI":redirectURI,
        "responseType":responseType,
        "scope":scope,
        "authMethod":authMethod
   }
    fetch("http://localhost:3001/api/auth/create_client/", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
        body: encodeFormData(formData)
    }).then( resp =>resp.json()
    ).then(data=>{
        if(data){
            window.sessionStorage.setItem("client_id",data.client_id)
            window.sessionStorage.setItem("client_secret",data.client_secret)
        }
    })
        .catch(err=>{

    })
    return 0;
}