const API_URL = "https://localhost";


function submit() {
    let u = document.getElementById("url").value;
    if (u === "")
        return;
    fetch(API_URL + "/shorten",
        {
            method: "POST",
            body: JSON.stringify({url: u})
        })
        .then(r => r.json())
        .then(value => {
            document.querySelector("#generated_url").innerHTML = `<a href=${API_URL}/${value.slug} target="_blank">${API_URL}/${value.slug}</a> => ${value.url}`;

            navigator.clipboard.writeText(`${API_URL}/${value.slug}`).then(function () {
                alert("copied the generated url to clipboard")
            }, function () {
                alert("copying the generated url to clipboard failed")
            });


        });

    document.querySelector("#url").textContent = ""
}