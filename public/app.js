const API_URL = "https://localhost";
const input = document.querySelector("#url");

function submit() {
    console.log("Submitting...");
    let u = document.getElementById("url").value;
    if (u === "")
        return;
    console.log(u)
    fetch(API_URL + "/shorten",
        {
            method: "POST",
            body: JSON.stringify({url: u})
        })
        .then(r => r.json())
        .then(value => {
            document.querySelector("#generated_url").innerHTML = `<a href=${API_URL}/${value.slug} target="_blank">${API_URL}/${value.slug}</a>`
            console.log(value)

            navigator.clipboard.writeText(`${API_URL}/${value.slug}`).then(function () {
                alert("copied the genereted url to clipboard")
            }, function () {
                alert("copying the generated url to clipboard failed")
            });


        });

    document.querySelector("#generated_url").textContent = ""
}