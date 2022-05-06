document.getElementById("second_list").style.display = "none";
document.getElementById("second_pdf").style.display = "none";

let response;

async function get_data() {
  const username_str = document.getElementById('inputusername').value
  const password_str = document.getElementById('inputpassword').value
  console.log(999)
  const promise_hash = await sha512("salt" + username_str);
  console.log("hash:", promise_hash);
  const fetch_return = await fetch(window.location.href + "/user_data/" + promise_hash);
  const fetch_json = await fetch_return.text();
  response = JSON.parse(fetch_json);

  console.log(response["links"][0][0])
  document.getElementById("first_list").style.display = "none";
  document.getElementById("second_list").style.display = "block";
  document.getElementById("first_pdf").style.display = "none";
  document.getElementById("second_pdf").style.display = "block";

  document.getElementById("forming").style.width = "60%";
  document.getElementById("banner_div").style.width = "40%";

  for (let i = 0; i < response["links"].length -1; i++) {
    if (response["links"][i][3] != "") {
      const subject = document.querySelector('#second_list');
      const string_val = '<div class="list_item" onclick="load_pdf(this)" value="' + i.toString() +'"> <img src="assets/checked.png" width="40px" alt=""> <h1>' + response["links"][i][0] + '</h1></br></br></br></br></br></br></br></br><p> <em> last updated: ' + response["links"][i][3] + '</em></p> </div>'; // sorry for that br*8 :)
      const nor = await subject.insertAdjacentHTML("beforeend", string_val);
    } else {
      const subject = document.querySelector('#second_list');
      const string_val = '<div class="list_item" onclick="load_pdf(this)" value="' + i.toString() +'"> <img src="assets/attention.png" width="40px" alt=""> <h1>' + response["links"][i][0] + '</h1> </div>'
      const nor = await subject.insertAdjacentHTML("beforeend", string_val);

    }
  }

}



async function load_pdf(str) {
  var x = str.getAttribute('value');
  const fetch_returx = await fetch(window.location.href + "/blockchain/" + response["links"][x][1]);
  const fetch_json = await fetch_returx.text();
  console.log(fetch_json)
  newr = JSON.parse(fetch_json);

  const username_s = document.getElementById('inputusername').value
  document.getElementById("embed_pdf").src= newr[username_s][0];
  console.log(newr[username_s][0]);
}

function sha512(str) {
  return crypto.subtle.digest("SHA-512", new TextEncoder("utf-8").encode(str)).then(buf => {
    return Array.prototype.map.call(new Uint8Array(buf), x=>(('00'+x.toString(16)).slice(-2))).join('');
  });
}

const blobToBase64 = blob => {
  const reader = new FileReader();
  reader.readAsDataURL(blob);
  return new Promise(resolve => {
    reader.onloadend = () => {
      resolve(reader.result);
    };
  });
};
