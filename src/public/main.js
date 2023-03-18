App()

function App(){
    let $name = document.querySelector("#username")
    let $password = document.querySelector("#password")
    let $form = document.querySelector("#form")

    $form.addEventListener("submit", handleSubmit)

    function handleSubmit(event){
        let name = $name.value
        let pass = $password.value
        loginRequest(name, pass)
        event.preventDefault();
    }
}



function loginRequest(name, pass){
    fetch("/api",{
                headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
                },
                method: "POST",
                body: JSON.stringify({name: name, password: pass})
            })
    .then(res => {
        if(!res.ok) throw Error(res.statusText)
        if (res.redirected) {
            window.location.href = res.url;
        }
    })
    .catch(function(res){ console.log(res) })
}