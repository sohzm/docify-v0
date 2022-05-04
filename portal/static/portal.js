let url = 'http://localhost:8089/list' //login

fetch(url, {})
    .then(function (a) {
        return a.json(); // call the json method on the response to get JSON
    })
    .then(function (json) {
        window.location.replace(json.links[0]);
    })
