const vh = (e)=>{
    const t = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
      , i = "NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm"
      , a = n=>t.indexOf(n)
      , o = n=>a(n) > -1 ? i[a(n)] : n;
    return e.split("").map(o).join("")
}

const b = function(I) {
    switch (I.length) {
    case 4:
        var U = (7 & I.charCodeAt(0)) << 18 | (63 & I.charCodeAt(1)) << 12 | (63 & I.charCodeAt(2)) << 6 | 63 & I.charCodeAt(3)
          , F = U - 65536;
        return String.fromCharCode((F >>> 10) + 55296) + String.fromCharCode((F & 1023) + 56320);
    case 3:
        return String.fromCharCode((15 & I.charCodeAt(0)) << 12 | (63 & I.charCodeAt(1)) << 6 | 63 & I.charCodeAt(2));
    default:
        return String.fromCharCode((31 & I.charCodeAt(0)) << 6 | 63 & I.charCodeAt(1))
    }
}

const w = function(I) {
    const y = /[\xC0-\xDF][\x80-\xBF]|[\xE0-\xEF][\x80-\xBF]{2}|[\xF0-\xF7][\x80-\xBF]{3}/g
    return I.replace(y, b)
}

const E = function(I) {
    buffer = Buffer.from(I,'base64')
    return buffer.toString()
}

const C = function(I) {
    return w(E(I))
}
const k = function(I) {
    return String(I).replace(/[-_]/g, function(U) {
        return U == "-" ? "+" : "/"
    }).replace(/[^A-Za-z0-9\+\/]/g, "")
}
const N = function(I) {
    return C(k(I))
}

const code = vh("5Y2t5nJ977lZ5YvJ55JZ77lO")
const text = N(code)
console.log(text)
