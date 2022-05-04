function get_data() {
    var myArray = [];
    sha512("g")
        .then(x => {
            fetch(window.location.href + "/user_data/" +x)
                .then(response => response.text())
                .then((response) => {

                    const username_str = document.getElementById('inputusername').value
                    const password_str = document.getElementById('inputpassword').value
                    console.log("RESPONSE STRING::", response)
                    const obj = JSON.parse(response)
                    document.getElementById("portal_form").outerHTML = "";
                    console.log(username_str, password_str)
                    myArray = obj.links;
                    var arrayLength = obj.links.length;

                    for (var i = 0; i < arrayLength; i++) {
                        var para = document.createElement("a");
                        para.href = obj.links[i];
                        var node = document.createTextNode(obj.links[i]);
                        para.appendChild(node);
                        var element = document.getElementById("temp");
                        element.appendChild(para);

                        var break_stt = document.createElement("br");
                        element.appendChild(break_stt);

                    }
                })
        })
}
