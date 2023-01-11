/* https://github.com/L1ghtM4n */

var prettyPrint = true;

function convert () {
    // Results
    let cookies = [];
    // Elements
    let cookiesInput = document.getElementById('inputCookies');
    let cookiesOutput = document.getElementById('outputCookies');
    // Split data
    let lines = cookiesInput.value.split("\n");
    // Iterate
    lines.forEach(function(line, i) {
        var tokens = line.split("\t");
        // We only care for valid cookie def lines
        if (tokens.length == 7) {
            let cookie = {};
            // Trim the tokens
            tokens = tokens.map(function(e) { return e.trim(); });
            // Extract the data
            cookie.domain = tokens[0];
            cookie.httpOnly = tokens[1] === "TRUE";
            cookie.path = tokens[2];
            cookie.secure = tokens[3] === "TRUE";
            // Convert date to a readable format
            let timestamp = tokens[4];
            if (timestamp.length == 17) {
            	timestamp = Math.floor(timestamp / 1000000 - 11644473600);
            }
            cookie.expirationDate = parseInt(timestamp);
            cookie.name = tokens[5];
            cookie.value = tokens[6];
            // Save the cookie.
            cookies.push(cookie);
        }    
    });
    // Done
    if (prettyPrint) {
        cookiesOutput.value = JSON.stringify(cookies, null, 2);
    } else {
        cookiesOutput.value = JSON.stringify(cookies);
    }
    // Focus and Clean
    cookiesOutput.focus();
    cookiesInput.value = '';
}