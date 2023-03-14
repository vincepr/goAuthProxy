App()

function App(){
    console.log("js running")
    let $name = document.querySelector("#username")
    let $password = document.querySelector("#password")
    let $submit = document.querySelector("#submit")
    let $form = document.querySelector("#form")
    $submit.addEventListener("click", handleSubmit)

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
        console.log(res)
        if (res.redirected) {
            window.location.href = res.url;
        }
        return res
    } )
    .then(data => console.log(data))
    .catch(function(res){ console.log(res) })
}