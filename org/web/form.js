async function get_data() {
  const username_str = document.getElementById('inputusername').value
  const password_str = document.getElementById('inputpassword').value
  const promise_hash = await sha512(username_str);
  console.log("hash:", promise_hash);
  const fetch_return = await fetch(window.location.href + "/user_data/" + "newenc");
  const fetch_json = await fetch_return.blob();
  console.log(fetch_json.size)
  const newx = await blobToBase64(fetch_json)
  console.log(newx.slice(37));
  const n = atob(newx.slice(37));
  console.log(n)
  
  var rawData = n;
  key = "5f4dcc3b5aa765d61d8327deb882cf99";

  // Decode the base64 data so we can separate iv and crypt text.
  //var rawData = atob(data);
  var iv = btoa(rawData.substring(0,16));
  var crypttext = btoa(rawData.substring(16));

  // Decrypt...
  var plaintextArray = CryptoJS.AES.decrypt(
    {
      ciphertext: CryptoJS.enc.Base64.parse(crypttext),
      salt: ""
    },
    CryptoJS.enc.Hex.parse(key),
    { iv: CryptoJS.enc.Base64.parse(iv) }
  );

  // Convert hex string to ASCII.
  // See https://stackoverflow.com/questions/11889329/word-array-to-string
  function hex2a(hex) {
    var str = '';
    for (var i = 0; i < hex.length; i += 2)
      str += String.fromCharCode(parseInt(hex.substr(i, 2), 16));
    return str;
  }

  console.log(hex2a(plaintextArray.toString()));
  //    const fetch_json = await fetch_return.text();
  //    const response = JSON.parse(fetch_json);
  //
  //    console.log("RESPONSE STRING::", response)
  //    document.getElementById("login_div").style.display = "none";
  //    document.getElementById("main_div").style.display = "flex";

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
